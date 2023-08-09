package core

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/util"
	"github.com/ystyle/kas/util/env"
	"golang.org/x/net/context"
)

type Service func(client *WsClient, message Message)

type WsManager struct {
	MaxConnect int
	clients    map[*WsClient]bool
	services   map[string]Service
	Register   chan *WsClient
	Unregister chan *WsClient
	ctx        context.Context
}

var wm = &WsManager{
	services:   make(map[string]Service),
	Register:   make(chan *WsClient),
	Unregister: make(chan *WsClient),
	clients:    make(map[*WsClient]bool),
}

func init() {
	if wm.MaxConnect == 0 {
		wm.MaxConnect = env.GetInt("MAX_CONNECT", 100)
	}
}

func GetWsManager() *WsManager {
	return wm
}

func (m *WsManager) Run() {
	for {
		select {
		case client := <-m.Register:
			client.wsManager = m
			m.clients[client] = true
			process(client)
		case client := <-m.Unregister:
			if _, ok := m.clients[client]; ok {
				client.cancel()
				delete(m.clients, client)
				close(client.WsSend)
			}
		}
	}
}

func process(client *WsClient) {
	log.Info("WSClient key: ", client.GetWSKey())
	if len(client.wsManager.clients) >= client.wsManager.MaxConnect {
		client.WsSend <- NewMessage("Error", "MaxConnect")
		return
	}
	go client.WriteMsg()
	go client.ReadMsg(func(c *WsClient, message Message) {
		services := client.wsManager.services
		if service, ok := services[message.Type]; ok {
			service(c, message)
		} else {
			log.Warn("Message type: %s not have any service", message.Type)
		}
	})
}

func (m *WsManager) GetClients() map[*WsClient]bool {
	return m.clients
}

func (m *WsManager) RegisterService(messageType string, service Service) {
	if temp, ok := m.services[messageType]; ok {
		pre := util.GetFunctionName(temp)
		curr := util.GetFunctionName(service)
		log.Warn(fmt.Sprintf("%s has the same service: %s, and it will be covered by %s.\n", messageType, pre, curr))
	}
	m.services[messageType] = service
}
