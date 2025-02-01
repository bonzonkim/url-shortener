package network

import (
	"github.com/bonzonkim/url-shortener/handlers"
	"github.com/gin-gonic/gin"
)

type Network struct {
	engine		*gin.Engine
	router		*Router
}

func NewNetwork(h *handlers.Handlers) *Network {
	n := &Network{
		engine: gin.New(),
	}
	n.router = NewRouter(n, h)
	return n
}

func (n *Network) ServerStart(port string) error {
	return n.engine.Run(port)
}
