package queues

import (
	"bitbucket.org/ignitionrobotics/ign-go"
)

type QueueRepository interface {
	initialize()
	Get(offset, limit *int) ([]interface{}, *ign.ErrMsg)
	Enqueue(entity interface{}) interface{}
	Dequeue() (interface{}, *ign.ErrMsg)
	DequeueOrWait() (interface{}, *ign.ErrMsg)
	MoveToFront(target interface{}) (interface{}, *ign.ErrMsg)
	MoveToBack(target interface{}) (interface{}, *ign.ErrMsg)
	Swap(a interface{}, b interface{}) (interface{}, *ign.ErrMsg)
	Remove(id interface{}) (interface{}, *ign.ErrMsg)
	Count() int
}
