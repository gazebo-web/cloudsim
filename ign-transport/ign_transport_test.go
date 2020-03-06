package igntransport

import (
	"fmt"
	"testing"
	"time"
)

func TestIgnPublish(t *testing.T) {
	ign, _ := NewIgnTransportNode(nil)
	defer ign.Free()

	ign.IgnTransportPublishStringMsg("/foo", "SOME TEXT")

	ign.IgnTransportPublishStringMsg("/bar", "00000")

	ign.IgnTransportPublishStringMsg("/bar", "11111")

	ign.IgnTransportPublishStringMsg("/bar", "222222")
	ign.IgnTransportPublishStringMsg("/foo", "3333")
	ign.IgnTransportPublishStringMsg("/bar", "4444")
}

func TestIgnSubscribe(t *testing.T) {
	ign, _ := NewIgnTransportNode(nil)
	defer ign.Free()

	igncb := func(msg []byte, msgType string) {
		fmt.Println("callback FOO / ", msg)
	}
	ign.IgnTransportSubscribe("/foo", igncb)
	defer ign.IgnTransportUnsubscribe("/foo")

	igncbbar := func(msg []byte, msgType string) {
		fmt.Println("callback BAR / ", msg)
	}
	ign.IgnTransportSubscribe("/bar", igncbbar)

	time.Sleep(5 * time.Second)
}
