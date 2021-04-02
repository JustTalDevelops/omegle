package main

import (
	"github.com/JustTalDevelops/omegleapi"
	"github.com/JustTalDevelops/omegleapi/requests"
	"github.com/sirupsen/logrus"
	"time"
)

// Echo, a simple Omegle bot that echos what the target user said.
// Created by JustTal for Dismegle.

type CustomHandler struct {
	log *logrus.Logger
	omegleapi.NopHandler
}

// OnMessage is called whenever a message is received from the stranger.
func (c CustomHandler) OnMessage(slave *omegleapi.Slave, msg string) {
	time.AfterFunc(1 * time.Second, func() {
		c.log.Info(slave.SendRequest(requests.SendMessage{Content: msg}))
	})
}

func main() {
	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.DebugLevel

	err := omegleapi.NewSlave(&CustomHandler{log: log}, []string{"women"}, omegleapi.SlaveOptions{
		Logger: log,
	}).Start()
	if err != nil {
		panic(err)
	}
}