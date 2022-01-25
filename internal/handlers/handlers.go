package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/gorilla/websocket"
)

// channel for getting data from client
var wsChan = make(chan WsPayload)

// contains list of user connections
var clients = make(map[WebSocketConnection]string)

// set with settings for page loading
var views = jet.NewSet(
	// specifies where templates are stored
	jet.NewOSFileSystemLoader("./html"),
	// we can change app and running it at the same time
	jet.InDevelopmentMode(),
)

// set with settings for websocket
var upgradeConnection = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handles Home page
func Home(w http.ResponseWriter, r *http.Request) {
	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}
}

// *websocket.Conn type keeper
type WebSocketConnection struct {
	*websocket.Conn
}

// Structure for sending JSON response to client
type WsJsonResponse struct {
	Action         string   `json:"action"`
	Message        string   `json:"message"`
	MessageType    string   `json:"message_type"`
	ConnectedUsers []string `json:"connected_users"`
}

// Structure for getting JSON response from client
type WsPayload struct {
	Action   string              `json:"action"`
	Username string              `json:"username"`
	Message  string              `json:"message"`
	Conn     WebSocketConnection `json:"-"`
}

// Upgrades connection to websockets
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	//getting connection
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client connected to endpoint")

	// contains response to client
	var response WsJsonResponse
	response.Message = `<em><small>Connected to server</small></em>`

	// connection with WebSocketConnection type
	conn := WebSocketConnection{Conn: ws}
	//adding connection to map
	clients[conn] = ""

	// sends response to client
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println(err)
	}
	// listens response from client in loop
	go ListenForWs(&conn)

}

// gets data from client
func ListenForWs(conn *WebSocketConnection) {
	// recover in case of error
	defer func() {
		if r := recover(); r != nil {
			log.Println("Error", fmt.Sprintf("%v", r))
		}
	}()
	// contains response from client
	var payload WsPayload
	// getting client response in loop
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {
			//do nothing if it is empty
		} else {
			//store data if we have smth
			payload.Conn = *conn
			// send to channel
			wsChan <- payload
		}
	}
}

// Serves messages from client
func ListenToWsChannel() {

	var response WsJsonResponse

	for {
		// reading from channel in loop
		e := <-wsChan

		switch e.Action {
		case "username":
			// giving name to current connection
			clients[e.Conn] = e.Username
			//getting the list of users
			users := getUserList()
			//forming response
			response.Action = "list_users"
			response.ConnectedUsers = users
			//sending response to all users
			broadcastToAll(response)
		case "left":
			response.Action = "list_users"
			// delete current connection from list
			delete(clients, e.Conn)
			// forming a new list of users
			users := getUserList()
			response.ConnectedUsers = users
			// sending response
			broadcastToAll(response)
		case "broadcast":
			// sending incoming message from user to all users
			response.Action = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.Username, e.Message)
			broadcastToAll(response)
		}
	}
}

// gets a list of users
func getUserList() []string {
	var userList []string
	for _, x := range clients {
		if x != "" {
			userList = append(userList, x)
		}
	}
	sort.Strings(userList)
	return userList
}

// sends data to all clients
func broadcastToAll(response WsJsonResponse) {
	// sending response to all connections in list of users
	for client := range clients {
		err := client.WriteJSON(response)
		if err != nil {
			// close connection in case of error
			log.Println("websocket err")
			_ = client.Close()
			delete(clients, client)
		}
	}
}

// renders page
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	// gets a template
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}
	// renders template
	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
