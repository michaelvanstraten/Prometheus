package websocket

import (
	"sync"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
)

type Message struct {
	Reciever string
	Data []byte
}

type WebsocketFunctions struct {
	Autherized func(*http.Request) bool
	Functions map[string]func(*[]byte, *Client)
}

type Websocket struct {
	upgrader websocket.Upgrader
	clients map[uuid.UUID]Client
	mux sync.Mutex
	functions []WebsocketFunctions
}

func New(Functions ...WebsocketFunctions) *Websocket {
	var NewWebsocket = &Websocket{}
	NewWebsocket.upgrader = websocket.Upgrader{
		ReadBufferSize: 512,
		WriteBufferSize: 512,
	}
	NewWebsocket.clients = make(map[uuid.UUID]Client)
	NewWebsocket.functions = Functions
	return NewWebsocket
}

func (w *Websocket) AddClient(Client *Client) uuid.UUID {
	var ClientID = uuid.New()
	w.mux.Lock()
	for _, ok := w.clients[ClientID]; ok; {
		ClientID = uuid.New()
	}
	w.clients[ClientID] = *Client
	w.mux.Unlock()
	return ClientID
}

func (w *Websocket) RemoveClient(ClientID uuid.UUID) {
	if Client, ok := w.clients[ClientID]; ok {
		Client.connection.Close()
	}
	w.mux.Lock()
	delete(w.clients, ClientID)
	w.mux.Unlock()
}

func (w *Websocket) Listener(W http.ResponseWriter, R *http.Request, _ httprouter.Params) {
	if conn, err := w.upgrader.Upgrade(W, R, nil); err == nil {
		var newClient = Client{}
		newClient.connection = conn
		newClient.functions = make(map[string]func(*[]byte, *Client))
		for _, functions := range w.functions {
			if functions.Autherized(R) {
				for Destination, function := range functions.Functions {
					newClient.functions[Destination] = function
				}
			}
		}
		go w.AddClient(&newClient)
		go newClient.Handle()
	}
}