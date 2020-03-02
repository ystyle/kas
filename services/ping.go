package services

import "github.com/ystyle/kas/core"

func Ping(client *core.WsClient, message core.Message) {
	client.WsSend <- core.NewMessage("pong", "pong")
}
