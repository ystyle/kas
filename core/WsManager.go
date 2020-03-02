package core

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"github.com/ystyle/kas/util"
	"github.com/ystyle/kas/util/env"
)

type Service func(client *WsClient, message Message)

type WsManager struct {
	clients    []*WsClient
	MaxConnect int
	services   map[string]Service
}

var wm = &WsManager{
	services: make(map[string]Service),
}

func init() {
	if wm.MaxConnect == 0 {
		wm.MaxConnect = env.GetInt("MAX_CONNECT", 100)
	}
}

func GetWsManager() *WsManager {
	return wm
}

func (m *WsManager) Add(client *WsClient) {
	log.Info("WSClient key: ", client.GetWSKey())
	if len(m.clients) >= m.MaxConnect {
		client.WsSend <- NewMessage("Error", "MaxConnect")
		return
	}
	client.wsManager = m
	m.clients = append(m.clients, client)
	go client.WriteMsg()
	client.ReadMsg(func(c *WsClient, message Message) {
		if service, ok := m.services[message.Type]; ok {
			service(c, message)
		} else {
			log.Warn("Message type: %s not have any service", message.Type)
		}
	})
}

func (m *WsManager) GetClients() []*WsClient {
	return m.clients
}

func (m *WsManager) Remove(client *WsClient) {
	for i, wsClient := range m.clients {
		if wsClient == client {
			m.clients = append(m.clients[:i], m.clients[i+1:]...)
			break
		}
	}
}

func (m *WsManager) RegisterService(messageType string, service Service) {
	if temp, ok := m.services[messageType]; ok {
		pre := util.GetFunctionName(temp)
		curr := util.GetFunctionName(service)
		log.Warn(fmt.Sprintf("%s has the same service: %s, and it will be covered by %s.\n", messageType, pre, curr))
	}
	m.services[messageType] = service
}
