package cmd

import (
	"github.com/bonzonkim/url-shortener/handlers"
	"github.com/bonzonkim/url-shortener/network"
	"github.com/bonzonkim/url-shortener/storage"
)

type Cmd struct {
	Network		*network.Network
	Handlers	*handlers.Handlers
	Store		*storage.RedisStore
}

func NewCmd() *Cmd {
	c := &Cmd{}

	c.Store = storage.NewRedisStore()
	c.Handlers = handlers.NewHandlers(c.Store)
	c.Network = network.NewNetwork(c.Handlers)
	c.Network.ServerStart(":8080")

	return c
}
