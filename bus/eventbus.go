package bus

import (
	"github.com/asaskevich/EventBus"
	"github.com/novikovSU/rocketchat-desktop-native/config"
	"log"
)

var (
	Bus = EventBus.New()
)

func Publish(topic string, args ...interface{}) {
	if config.Debug {
		log.Printf("Fire event: %s %s\n", topic, args)
	}
	Bus.Publish(topic, args...)
}

func SubscribeAsync(topic string, fn interface{}) {
	err := Bus.SubscribeAsync(topic, fn, false)
	if err != nil {
		log.Panicf("Invalid argument %s. It must be a function!", fn)
	}
}
