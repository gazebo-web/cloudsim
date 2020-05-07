package queues

import (
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// Deprecated: LaunchQueueRepository represents the Launch Queue Data Access Object
type LaunchQueueRepository struct {
	queue *ign.Queue
}

// Deprecated: NewLaunchQueueRepository returns a new LaunchQueueRepository instance.
func NewLaunchQueueRepository() (lq *LaunchQueueRepository) {
	lq = &LaunchQueueRepository{}
	lq.initialize()
	return
}

// QueueRepository Implementation

// Deprecated: Get returns the entire launch queue.
// If `offset` and `limit` are not nil, it will return up to `limit` results from the provided `offset`.
func (lq *LaunchQueueRepository) Get(offset, limit *int) ([]interface{}, *ign.ErrMsg) {
	if offset == nil || limit == nil {
		return lq.queue.GetElements()
	}
	return lq.queue.GetFilteredElements(*offset, *limit)
}

// Deprecated: Remove removes a groupID from the queue.
func (lq *LaunchQueueRepository) Remove(id interface{}) (interface{}, *ign.ErrMsg) {
	groupID, ok := id.(string)

	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if err := lq.queue.Remove(groupID); err != nil {
		return nil, err
	}

	return groupID, nil
}

// Deprecated: initialize initializes the queue data structure.
func (lq *LaunchQueueRepository) initialize() {
	lq.queue = ign.NewQueue()
}

// Deprecated: Enqueue enqueues a groupID on the queue.
// Deprecated: Returns the groupID that was pushed.
func (lq *LaunchQueueRepository) Enqueue(entity interface{}) interface{} {
	groupID, ok := entity.(string)

	if !ok {
		return nil
	}

	lq.queue.Enqueue(groupID)
	return entity
}

// Deprecated: Dequeue returns the next groupID from the queue.
func (lq *LaunchQueueRepository) Dequeue() (interface{}, *ign.ErrMsg) {
	return lq.queue.Dequeue()
}

// Deprecated: DequeueOrWait returns the next groupID from the queue or waits until there is one available.
func (lq *LaunchQueueRepository) DequeueOrWait() (interface{}, *ign.ErrMsg) {
	return lq.queue.DequeueOrWaitForNextElement()
}

// Deprecated: MoveToFront moves a target groupID to the front of the queue.
func (lq *LaunchQueueRepository) MoveToFront(target interface{}) (interface{}, *ign.ErrMsg) {
	groupID, ok := target.(string)

	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if err := lq.queue.MoveToFront(groupID); err != nil {
		return nil, err
	}
	return target, nil
}

// Deprecated: MoveToBack moves a target element to the front of the queue.
func (lq *LaunchQueueRepository) MoveToBack(target interface{}) (interface{}, *ign.ErrMsg) {
	groupID, ok := target.(string)

	if !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if err := lq.queue.MoveToBack(groupID); err != nil {
		return nil, err
	}
	return target, nil
}

// Deprecated: Swap switch places between groupID A and groupID B.
func (lq *LaunchQueueRepository) Swap(a interface{}, b interface{}) (interface{}, *ign.ErrMsg) {
	var groupIDA string
	var groupIDB string
	var ok bool

	if groupIDA, ok = a.(string); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	if groupIDB, ok = b.(string); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorCastingID)
	}

	err := lq.queue.Swap(groupIDA, groupIDB)
	if err != nil {
		return nil, err
	}

	res := SwapResponse{
		a: QueueItemResponse{
			GroupID: groupIDA,
		},
		b: QueueItemResponse{
			GroupID: groupIDB,
		},
	}
	return res, nil
}

// Deprecated: Count returns the length of the underlying queue's slice
func (lq *LaunchQueueRepository) Count() int {
	return lq.queue.GetLen()
}
