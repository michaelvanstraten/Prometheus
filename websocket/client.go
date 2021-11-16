package websocket

import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type Client struct {
	functions map[string]func(*[]byte, *Client)
	connection *websocket.Conn
}

func (c *Client) Handle() {
	for {
		_, data, err := c.connection.ReadMessage()
		if err != nil {
			return
		}
		var newMessage = Message{}
		if err := json.Unmarshal(data, &newMessage); err == nil {
			if handler, ok := c.functions[newMessage.Reciever]; ok {
				go handler(&newMessage.Data, c)
			}
		} else {
			var err = c.connection.WriteMessage(websocket.TextMessage, []byte("Message not recognised"))
			println(string(data))
			if err != nil {
				println(err.Error())
			}
		}
	}
}