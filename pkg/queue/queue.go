package queue

import "gitlab.com/ignitionrobotics/web/ign-go"

type IQueue interface {
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

type Queue struct {
	dao *ign.Queue
}

func New() *IQueue {
	var q IQueue
	q = initialize()
	return &q
}

func initialize() *Queue {
	q := Queue{}
	q.dao = ign.NewQueue()
	return &q
}

// Get returns the entire launch queue.
// If `offset` and `limit` are not nil, it will return up to `limit` results from the provided `offset`.
func (q *Queue) Get(offset, limit *int) ([]interface{}, *ign.ErrMsg) {
	if offset == nil || limit == nil {
		return q.dao.GetElements()
	}
	return q.dao.GetFilteredElements(*offset, *limit)
}

// Remove removes a groupID from the queue.
func (q *Queue) Remove(id interface{}) (interface{}, *ign.ErrMsg) {
	groupID, ok := id.(string)

	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if err := q.dao.Remove(groupID); err != nil {
		return nil, err
	}

	return groupID, nil
}

// Enqueue enqueues a groupID on the queue.
// Returns the groupID that was pushed.
func (q *Queue) Enqueue(entity interface{}) interface{} {
	groupID, ok := entity.(string)

	if !ok {
		return nil
	}

	q.dao.Enqueue(groupID)
	return entity
}

// Dequeue returns the next groupID from the queue.
func (q *Queue) Dequeue() (interface{}, *ign.ErrMsg) {
	return q.dao.Dequeue()
}

// DequeueOrWait returns the next groupID from the queue or waits until there is one available.
func (q *Queue) DequeueOrWait() (interface{}, *ign.ErrMsg) {
	return q.dao.DequeueOrWaitForNextElement()
}

// MoveToFront moves a target groupID to the front of the queue.
func (q *Queue) MoveToFront(target interface{}) (interface{}, *ign.ErrMsg) {
	groupID, ok := target.(string)

	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if err := q.dao.MoveToFront(groupID); err != nil {
		return nil, err
	}
	return target, nil
}

// MoveToBack moves a target element to the front of the queue.
func (q *Queue) MoveToBack(target interface{}) (interface{}, *ign.ErrMsg) {
	groupID, ok := target.(string)

	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if err := q.dao.MoveToBack(groupID); err != nil {
		return nil, err
	}
	return target, nil
}

// Swap switch places between groupID A and groupID B.
func (q *Queue) Swap(a interface{}, b interface{}) (interface{}, *ign.ErrMsg) {
	var groupIDA string
	var groupIDB string
	var ok bool

	if groupIDA, ok = a.(string); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if groupIDB, ok = b.(string); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	err := q.dao.Swap(groupIDA, groupIDB)
	if err != nil {
		return nil, err
	}

	res := SwapOutput{
		ItemA: Item{
			GroupID: groupIDA,
		},
		ItemB: Item{
			GroupID: groupIDB,
		},
	}
	return res, nil
}

// Count returns the length of the underlying queue's slice
func (q *Queue) Count() int {
	return q.dao.GetLen()
}
