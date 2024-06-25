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

var MessageHandlersList []*MessageHandler = []*MessageHandler{
	{&domain.Journal{}, &infra.LocalRecord{}},
	{&domain.Nop{}, &infra.Nop{}},
}
