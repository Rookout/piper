package server

import (
	"github.com/rookout/piper/pkg/clients"
	"golang.org/x/net/context"
)

type healthCheck struct {
	id  int64
	ctx *context.Context
}

type HealthChecker struct {
	clients *clients.Clients
	Check   chan *healthCheck
	Failure chan *healthCheck
}

func NewHealthChecker(clients *clients.Clients) *HealthChecker {
	return &HealthChecker{
		clients: clients,
		Check:   make(chan *healthCheck),
		Failure: make(chan *healthCheck),
	}
}

func (h *HealthChecker) New(msg *healthCheck) error {
	h.Check <- msg
	return nil
}

func (h *HealthChecker) Fail(msg *healthCheck) error {
	h.Failure <- msg
	return nil
}

func (h *HealthChecker) Handle(healthCheck *healthCheck) {
	//log.Printf("Health check started of: %d\n", healthCheck.id)
	//failedHooks, err := h.clients.GitProvider.PingHooks(healthCheck.ctx)
	//if err != nil {
	//	log.Printf("error in pinging hooks %s", err)
	//}

}
