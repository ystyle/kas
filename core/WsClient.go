package core

import (
	"github.com/labstack/gommon/log"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
)

type WsClient struct {
	WsConn      *websocket.Conn
	WsSend      chan Message
	wsManager   *WsManager
	HttpRequest *http.Request
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
			websocket.JSON.Send(client.WsConn, msg)
		}
	}
}
