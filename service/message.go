package service

import (
	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/infra"
	"github.com/auifzysr/kaburasutegi/repository"
)

type MessageHandler struct {
	messageBuilder domain.MessageBuilder
	recorder       repository.Recorder
}

var MessageHandlersList []*MessageHandler

func init() {
	makeMessageHandlersList()
}

func makeMessageHandlersList() {
	MessageHandlersList = append(MessageHandlersList, &MessageHandler{&domain.Journal{}, &infra.LocalRecord{}})
	MessageHandlersList = append(MessageHandlersList, &MessageHandler{&domain.Nop{}, &infra.Nop{}})
}
