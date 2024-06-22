package service

import (
	"github.com/auifzysr/kaburasutegi/domain"
	"github.com/auifzysr/kaburasutegi/repository"
)

type Service struct {
	domain.MessageHandler
	repository.Recorder
}
