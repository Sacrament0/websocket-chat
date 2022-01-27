# WebSocket Chat App

This is a simple chat app based on websockets. It allows users to chat back and forth, and shows the list of connected users, updated in real time.

[![ezgif.com-gif-maker-1e8df4435a13804b5.gif](https://s10.gifyu.com/images/ezgif.com-gif-maker-1e8df4435a13804b5.gif)](https://gifyu.com/image/Sbftb)

This application uses:
 - [Gorilla WebSocket](https://github.com/gorilla/websocket)
 - [Jet package](https://github.com/CloudyKit/jet) - for easier work with templates
 - [ReconnectingWebSocket](https://github.com/joewalnes/reconnecting-websocket) - a small JavaScript library that decorates the WebSocket API to provide a WebSocket connection that will automatically reconnect if the connection is dropped
 - [Notie](https://cdnjs.com/libraries/notie) - for decorating alert notifications

The idea of this app:

 1. Server renders home page. Client(browser) connects to this page. JS inside the client reconnects to server using ip:port/ws address. WebSocket connection has been formed.
 2. Client is listening to user's actions. If action is detected, client sends message to server using JSON. 
 3. Server is serving client messages. If there is a message, server serves it and sends response to client using JSON.
 4. Client is serving server response. If there is a response, client rederes smth or do smth. 
 
This application is made according to [Trevor Sawler's](https://github.com/tsawler) lessons