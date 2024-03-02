package main

import (
	"sync"
)

type Item struct {
	Topic string
	Data  Message
}

type Subscriber struct {
	Topic   string
	Channel chan Item
}

type PubSub struct {
	Mutex       sync.RWMutex
	Subscribers []Subscriber
}

func NewPubSub() *PubSub {
	return &PubSub{}
}

func (ps *PubSub) Subscribe(topic string) <-chan Item {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()

	channel := make(chan Item, 1)
	ps.Subscribers = append(ps.Subscribers, Subscriber{Topic: topic, Channel: channel})

	return channel
}

func (ps *PubSub) Publish(topic string, data Message) {
	ps.Mutex.RLock()
	defer ps.Mutex.RUnlock()

	for _, sub := range ps.Subscribers {
		if sub.Topic == topic {
			sub.Channel <- Item{Topic: topic, Data: data}
		}
	}
}

func (ps *PubSub) Unsubscribe(topic string, subscriberChan <-chan Item) {
	ps.Mutex.Lock()
	defer ps.Mutex.Unlock()

	for i, subscriber := range ps.Subscribers {
		if subscriber.Topic == topic && subscriber.Channel == subscriberChan {
			ps.Subscribers = append(ps.Subscribers[:i], ps.Subscribers[i+1:]...) // remove subscriber | append( subs[0 to i excluding], subs[i+1 to end] )
			break
		}
	}
}
