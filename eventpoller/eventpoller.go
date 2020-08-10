package eventpoller

import (
	"github.com/fbelisk/aeolus/internal"
	"jarvis/client"
)

type eventpoller struct {
	p internal.Poller
	clients map[int]*client.Client
}



