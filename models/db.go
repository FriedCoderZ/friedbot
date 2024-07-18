package models

import (
	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

type Service struct {
	*xorm.Engine
}

func NewService(path string) *Service {
	engine, err := xorm.NewEngine("sqlite3", path)
	if err != nil {
		panic(err)
	}
	return &Service{engine}
}

func (svc *Service) Start() {
	err := svc.Sync(
		new(User),
		new(Message),
	)
	if err != nil {
		panic(err)
	}
}
