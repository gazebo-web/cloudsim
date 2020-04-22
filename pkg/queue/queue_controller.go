package queue

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// IController represents a group of methods to expose in the API Rest.
type IController interface {
	GetAll(ctx context.Context, user *fuel.User, page, perPage *int) ([]interface{}, *ign.ErrMsg)
	Count(ctx context.Context, user *fuel.User) (interface{}, *ign.ErrMsg)
	MoveToFront(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg)
	MoveToBack(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg)
	Swap(ctx context.Context, user *fuel.User, groupIDA, groupIDB string) (interface{}, *ign.ErrMsg)
	Remove(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg)
}

// Controller is an IController implementation.
type Controller struct {
	services services
}

// services is a group of services used by the Controller.
type services struct {
	user  *users.Service
	queue IQueue
}

// GetAll returns a paginated list of elements from the queue.
// If no page or perPage arguments are passed, it sets those value to 0 and 10 respectively.
func (c *Controller) GetAll(ctx context.Context, user *fuel.User, page, perPage *int) ([]interface{}, *ign.ErrMsg) {
	if ok := c.services.user.Accessor.IsSystemAdmin(*user.Name); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	if page == nil {
		page = tools.Intptr(0)
	}
	if perPage == nil {
		perPage = tools.Intptr(10)
	}
	offset := *page * *perPage
	limit := *perPage
	return c.services.queue.Get(&offset, &limit)
}

// Count returns the element count from the queue.
func (c *Controller) Count(ctx context.Context, user *fuel.User) (interface{}, *ign.ErrMsg) {
	if ok := c.services.user.Accessor.IsSystemAdmin(*user.Name); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return c.services.queue.Count(), nil
}

// MoveToFront moves an element by the given groupID to the front of the queue.
func (c *Controller) MoveToFront(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := c.services.user.Accessor.IsSystemAdmin(*user.Name); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return c.services.queue.MoveToFront(groupID)
}

// MoveToBack moves an element by the given groupID to the back of the queue.
func (c *Controller) MoveToBack(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := c.services.user.Accessor.IsSystemAdmin(*user.Name); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return c.services.queue.MoveToBack(groupID)
}

// Swap swaps positions of groupIDs A and B.
func (c *Controller) Swap(ctx context.Context, user *fuel.User, groupIDA, groupIDB string) (interface{}, *ign.ErrMsg) {
	if ok := c.services.user.Accessor.IsSystemAdmin(*user.Name); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return c.services.queue.Swap(groupIDA, groupIDB)
}

// Remove removes an element by the given groupID from the queue.
func (c *Controller) Remove(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := c.services.user.Accessor.IsSystemAdmin(*user.Name); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return c.services.queue.Remove(groupID)
}
