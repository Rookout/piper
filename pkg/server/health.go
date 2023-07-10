package server

type Health interface {
	check(msg healthCheck) error
	recovery(msg healthCheck) error
	failure(msg healthCheck) error
}

type healthCheck struct {
	id string
}

type HealthChecker struct {
	Check    chan *healthCheck
	Recovery chan *healthCheck
	Failure  chan *healthCheck
}

func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		Check:    make(chan *healthCheck),
		Recovery: make(chan *healthCheck),
		Failure:  make(chan *healthCheck),
	}
}

func (h *HealthChecker) check(msg *healthCheck) error {
	h.Check <- msg
	return nil
}

func (h *HealthChecker) recovery(msg *healthCheck) error {
	h.Recovery <- msg
	return nil
}

func (h *HealthChecker) failure(msg *healthCheck) error {
	h.Failure <- msg
	return nil
}
