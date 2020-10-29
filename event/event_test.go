package event

import (
	"fmt"
	"testing"
	"time"
)

func TestEventBus(t *testing.T) {

	aListener := NewListener()
	bListener := NewListener()

	aListener.SetEventType("testEvent")
	aListener.SetCallback(func(ae AEvent) {
		fmt.Println("A called")
	})

	bListener.SetEventType("testEvent")
	bListener.SetCallback(func(ae AEvent) {
		fmt.Println("B called")
	})

	_ = Manager.Run()

	Bus <- AEvent{Type: "testEvent"}

	time.Sleep(5 * time.Second)
}
