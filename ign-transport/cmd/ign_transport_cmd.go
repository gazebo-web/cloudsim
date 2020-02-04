package main

import (
	igntran "bitbucket.org/ignitionrobotics/web-cloudsim/ign-transport"
	"fmt"
)

func main() {
	// An ign-transport node with explicit Partition set
	partition := "part"
	node1, _ := igntran.NewIgnTransportNode(&partition)
	defer node1.Free()

	// Another ign-transport node, without explicit Partition
	node2, _ := igntran.NewIgnTransportNode(nil)
	defer node2.Free()

	fmt.Println("Example: About to register 2 ign-transport callbacks")
	igncb := func(msg []byte, msgType string) {
		fmt.Println("Node1: hi friend / ", msg)
	}
	node1.IgnTransportSubscribe("/foo", igncb)
	defer node1.IgnTransportUnsubscribe("/foo")

	igncb2 := func(msg []byte, msgType string) {
		fmt.Printf("Node2: Another callback. Got msg: [%s]. Type: %s\n", msg, msgType)
	}
	node2.IgnTransportSubscribe("/foo", igncb2)

	node2.IgnTransportPublishStringMsg("/foo", "should only be seen by Node2")
	node1.IgnTransportPublishStringMsg("/foo", "This one by Node1")

	// Trick to wait for an "enter" key press before exiting...
	fmt.Scanln()
}
