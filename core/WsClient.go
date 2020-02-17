package core

import (
	"github.com/labstack/gommon/log"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
	"time"
)

type WsClient struct {
	WsConn      *websocket.Conn
	HttpRequest *http.Request
	wsManager   *WsManager
	WsSend      chan Message
	Caches      map[string]interface{}
}

func NewWsClient(wsConn *websocket.Conn) *WsClient {
	return &WsClient{
		WsConn:      wsConn,
		HttpRequest: wsConn.Request(),
		WsSend:      make(chan Message, 10),
		Caches:      make(map[string]interface{}),
	}
}

func (client *WsClient) GetWSKey() string {
	return client.HttpRequest.Header.Get("Sec-WebSocket-Key")
}

func (client *WsClient) ReadMsg(fn func(c *WsClient, message Message)) {
	for {
		var msg Message
		err := websocket.JSON.Receive(client.WsConn, &msg)
		if err != nil {
			if !strings.Contains(err.Error(), "EOF") {
				log.Error(err)
			}
			client.wsManager.Remove(client)
			break
		}
		msg.Time = time.Now()
		fn(client, msg)
	}
}

func (client *WsClient) Remove(fn func(c *WsClient)) {
	fn(client)
}

func (client *WsClient) WriteMsg() {
	for {
		select {
		case msg, ok := <-client.WsSend:
			if !ok {
				client.WsConn.WriteClose(500)
				log.Error("close client")
				return
			}
			msg.Time = time.Now()
			websocket.JSON.Send(client.WsConn, msg)
		}
	}
}
