package core

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/labstack/gommon/log"
	"net/http"
	"time"
)

const (
	readTimeOut  = 30 * time.Minute
	writeTimeOut = 30 * time.Minute
)

type WsClient struct {
	WsConn      *websocket.Conn
	HttpRequest *http.Request
	wsManager   *WsManager
	WsSend      chan Message
	Caches      map[string]interface{}
}

func NewWsClient(wsConn *websocket.Conn, resp *http.Request) *WsClient {
	return &WsClient{
		WsConn:      wsConn,
		HttpRequest: resp,
		WsSend:      make(chan Message, 10),
		Caches:      make(map[string]interface{}),
	}
}

func (client *WsClient) GetWSKey() string {
	return client.HttpRequest.Header.Get("Sec-WebSocket-Key")
}

func (client *WsClient) ReadMsg(fn func(c *WsClient, message Message)) {
	defer func() {
		client.wsManager.Unregister <- client
		client.WsConn.Close()
	}()
	client.WsConn.SetReadDeadline(time.Now().Add(readTimeOut))
	client.WsConn.SetPongHandler(func(string) error {
		client.WsConn.SetReadDeadline(time.Now().Add(readTimeOut))
		return nil
	})
	for {
		var msg Message
		err := client.WsConn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Error(err)
			} else if _, ok := err.(*json.SyntaxError); ok {
				client.WsConn.PingHandler()
				continue
			}
			break
		}
		msg.Time = time.Now()
		go fn(client, msg)
	}
}

func (client *WsClient) Remove(fn func(c *WsClient)) {
	fn(client)
}

func (client *WsClient) WriteMsg() {
	ticker := time.NewTicker(writeTimeOut)
	defer func() {
		ticker.Stop()
		client.WsConn.Close()
	}()

	for {
		select {
		case msg, ok := <-client.WsSend:
			if !ok {
				client.WsConn.Close()
				log.Error("close client")
				return
			}
			msg.Time = time.Now()
			client.WsConn.WriteJSON(msg)
		case <-ticker.C:
			client.WsConn.SetWriteDeadline(time.Now().Add(writeTimeOut))
			if err := client.WsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
