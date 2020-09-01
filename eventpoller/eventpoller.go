package eventpoller

import (
	"github.com/fbelisk/aeolus/internal"
)

type eventpoller struct {
	p internal.Poller
	clients map[int]*client.Client
}



