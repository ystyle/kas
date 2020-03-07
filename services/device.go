package services

import (
	"fmt"
	"github.com/labstack/gommon/log"
	uuid "github.com/satori/go.uuid"
	"github.com/ystyle/kas/core"
	"github.com/ystyle/kas/model"
	"time"
)

func Register(client *core.WsClient, message core.Message) {
	uid := uuid.NewV4().String()
	err := model.DB().Save(&model.Drive{
		ID:        uid,
		CreatedAt: time.Now(),
		Last:      time.Now(),
	})
	if err != nil {
		log.Error(err)
		client.WsSend <- core.NewMessage("register:faild", fmt.Errorf("注册失败%w", err).Error())
		return
	}
	client.WsSend <- core.NewMessage("register:success", "注册成功")
}
