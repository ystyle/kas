package core

import (
	"encoding/json"
	"time"
)

type Message struct {
	Time    time.Time
	Type    string // messageType
	DriveID string
	Data    interface{} // content
}

func NewMessage(Type string, message interface{}) Message {
	return Message{
		Type: Type,
		Data: message,
	}
}

func (msg *Message) GetString() string {
	return msg.Data.(string)
}

func (msg *Message) GetInt() int {
	return msg.Data.(int)
}

func (msg *Message) GetInt64() int64 {
	return msg.Data.(int64)
}

func (msg *Message) GetFloat() float64 {
	return msg.Data.(float64)
}

func (msg *Message) GetBool() bool {
	return msg.Data.(bool)
}

func (msg *Message) JsonParse(v interface{}) error {
	buff, err := json.Marshal(msg.Data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(buff, v)
	return err
}
