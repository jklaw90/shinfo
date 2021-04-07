package gateway

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

type ChatSSEClient struct {
	ctx     context.Context
	topic   string
	hub     Hub
	rw      http.ResponseWriter
	flusher http.Flusher
	sender  chan []byte
	done    chan bool
}

func NewSSEClient(ctx context.Context, topic string, hub Hub, rw http.ResponseWriter, flusher http.Flusher) *ChatSSEClient {
	sender := make(chan []byte, 10)

	c := &ChatSSEClient{
		ctx:     ctx,
		topic:   topic,
		hub:     hub,
		sender:  sender,
		rw:      rw,
		flusher: flusher,
	}
	return c
}

func (c *ChatSSEClient) Run() {
	defer close(c.sender)

	c.rw.Header().Set("Content-Type", "text/event-stream")
	c.rw.Header().Set("Cache-Control", "no-cache")
	c.rw.Header().Set("Connection", "keep-alive")

	tk := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-tk.C:
			fmt.Fprint(c.rw, ": ping\n\n")
			c.flusher.Flush()
		case <-c.ctx.Done():
			c.hub.Unsubscribe(c.topic, c)
			return
		case <-c.done:
			return
		case msg := <-c.sender:
			fmt.Println(msg)
		}
	}
}

func (c *ChatSSEClient) Send(message []byte) error {
	c.sender <- message
	return nil
}

func (c *ChatSSEClient) Close() error {
	c.done <- true
	return nil
}
