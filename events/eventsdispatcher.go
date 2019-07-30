package events

import (
	"log"
	"strconv"

	"../config"
)

var (
	subscribers = make(map[string]map[string]*func([]*interface{}))
)

func SubscribeToEvents(subscriber *func([]*interface{}), eventsName []string) {
	SubscribeNamedToEvents("autoSubscriber", subscriber, eventsName)
}

func Subscribe(subscriber *func([]*interface{}), eventName string) {
	SubscribeNamed("autoSubscriber", subscriber, eventName)
}

func SubscribeNamedToEvents(subscriberName string, subscriber *func([]*interface{}), eventsName []string) {
	for _, eventName := range eventsName {
		SubscribeNamed(subscriberName, subscriber, eventName)
	}
}

func SubscribeNamed(subscriberName string, subscriber *func([]*interface{}), eventName string) {
	eventSubscribers := subscribers[eventName]
	if eventSubscribers == nil {
		eventSubscribers = make(map[string]*func([]*interface{}))
		subscribers[eventName] = eventSubscribers
	}

	newSubscriberName := subscriberName
	for i := 1; eventSubscribers[newSubscriberName] != nil; i++ {
		newSubscriberName = subscriberName + strconv.FormatInt(int64(i), 10)
	}
	eventSubscribers[newSubscriberName] = subscriber
	if config.Debug {
		log.Printf("Added subscriber: %s to event %s\n", newSubscriberName, eventName)
	}
}

func UnsubscribeFromEvents(subscriberName string, eventsName []string) {
	for _, eventName := range eventsName {
		Unsubscribe(subscriberName, eventName)
	}
}

func Unsubscribe(subscriberName string, eventName string) {
	eventSubscribers := subscribers[eventName]
	if eventSubscribers != nil {
		delete(eventSubscribers, subscriberName)
		if config.Debug {
			log.Printf("Removed subscriber: %s to event %s\n", subscriberName, eventName)
		}
	}
}

func RaiseEventAsync(eventName string, args []*interface{}) {
	go RaiseEvent(eventName, args)
}

func RaiseEvent(eventName string, args []*interface{}) {
	if config.Debug {
		log.Printf("Fire event: %s with args: %s\n", eventName, args)
	}
	eventSubscribers := subscribers[eventName]
	if eventSubscribers != nil {
		for subscriberName, eventSubscriber := range eventSubscribers {
			if config.Debug {
				log.Printf("Call subscriber %s\n", subscriberName)
			}
			(*eventSubscriber)(args)
		}
	}
}
