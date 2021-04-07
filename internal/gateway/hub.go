package gateway

import (
	"sync"
)

type BasicHub struct {
	subscribers map[string][]Client
	sync.RWMutex
}

type Subscriber func(string) error
type Unsubscriber func(string) error

func (h *BasicHub) BaseSubscribe(topic string, client Client, subscriber Subscriber) error {
	h.Lock()
	defer h.Unlock()

	if h.subscribers == nil {
		h.subscribers = make(map[string][]Client)
	}

	if _, ok := h.subscribers[topic]; ok {
		h.subscribers[topic] = append(h.subscribers[topic], client)
	} else {
		subscriber(topic)
		h.subscribers[topic] = []Client{client}
	}

	return nil
}

func (h *BasicHub) BaseUnsubscribe(topic string, client Client, unsubscriber Unsubscriber) error {
	h.Lock()
	defer h.Unlock()

	clients, ok := h.subscribers[topic]
	if !ok {
		return nil
	}
	for i := 0; i < len(clients); i++ {
		if clients[i] == client {
			h.subscribers[topic] = append(clients[:i], clients[i+1:]...)
		}
	}

	if len(h.subscribers[topic]) == 0 {
		unsubscriber(topic)
		delete(h.subscribers, topic)
	}
	return nil
}

func (h *BasicHub) Broadcast(topic string, message []byte) error {
	h.RLock()
	defer h.RUnlock()
	clients, ok := h.subscribers[topic]
	if !ok {
		return nil
	}
	for i := 0; i < len(clients); i++ {
		clients[i].Send(message)
	}
	return nil
}

func (h *BasicHub) BaseClose() error {
	h.Lock()
	defer h.Unlock()
	for _, clients := range h.subscribers {
		for _, client := range clients {
			client.Close()
		}
	}
	h.subscribers = make(map[string][]Client)
	return nil
}
