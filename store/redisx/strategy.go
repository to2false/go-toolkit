package redisx

import (
	"sync/atomic"
)

var (
	roundRobinCount uint64
)

func RoundRobin(clients []Node) Node {
	i := (roundRobinCount) % uint64(len(clients))
	atomic.AddUint64(&roundRobinCount, 1)

	return clients[i]
}
