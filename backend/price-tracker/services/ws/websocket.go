package ws

import "github.com/gorilla/websocket"

type WebSocketFetcher interface {
	// Connect returns the connection to websocket
	Connect() (*websocket.Conn, error)
	// Fetch keeps running every few seconds to retrieve latest price from given connection
	Fetch(conn *websocket.Conn) error
}
