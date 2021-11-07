package lb

import (
	"errors"
	"sync/atomic"
)

// ErrNonExistServer for non existent server
var ErrNonExistServer = errors.New("non existent server")

// RoundRobin is an interface for representing round-robin balancing.
type RoundRobin interface {
	Next() string
}

type serviceRoundrobin struct {
	backends []string
	next     uint32
}

// New returns RoundRobin implementation(*roundrobin).
func New(backends ...string) (RoundRobin, error) {
	if len(backends) == 0 {
		return nil, ErrNonExistServer
	}

	return &serviceRoundrobin{
		backends: backends,
	}, nil
}

// Next returns next address
func (s *serviceRoundrobin) Next() string {
	n := atomic.AddUint32(&s.next, 1)
	return s.backends[(int(n)-1)%len(s.backends)]
}
