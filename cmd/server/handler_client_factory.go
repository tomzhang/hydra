package server

import (
	"github.com/julienschmidt/httprouter"
	"github.com/ory-am/fosite"
	"github.com/ory-am/hydra/client"
	"github.com/ory-am/hydra/config"
	"github.com/ory-am/hydra/herodot"
	"golang.org/x/net/context"
)

func newClientManager(c *config.Config) client.Manager {
	ctx := c.Context()

	switch con := ctx.Connection.(type) {
	case *config.MemoryConnection:
		return &client.MemoryManager{
			Clients: map[string]*fosite.DefaultClient{},
			Hasher:  ctx.Hasher,
		}
	case *config.RethinkDBConnection:
		con.CreateTableIfNotExists("hydra_policies")
		m := &client.RethinkManager{Session: con.GetSession()}
		m.ColdStart()
		m.Watch(context.Background())
		return m
	default:
		panic("Unknown connection type.")
	}
}

func newClientHandler(c *config.Config, router *httprouter.Router, manager client.Manager) *client.Handler {
	ctx := c.Context()
	h := &client.Handler{
		H: &herodot.JSON{},
		W: ctx.Warden, Manager: manager,
	}

	h.SetRoutes(router)
	return h
}