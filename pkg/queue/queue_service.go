package queue

import (
	"context"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gitlab.com/ignitionrobotics/web/cloudsim/tools"
	fuel "gitlab.com/ignitionrobotics/web/fuelserver/bundles/users"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

// IService represents a group of methods to expose in the API Rest.
type IService interface {
	GetAll(ctx context.Context, user *fuel.User, page, perPage *int) ([]interface{}, *ign.ErrMsg)
	Count(ctx context.Context, user *fuel.User) (interface{}, *ign.ErrMsg)
	MoveToFront(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg)
	MoveToBack(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg)
	Swap(ctx context.Context, user *fuel.User, groupIDA, groupIDB string) (interface{}, *ign.ErrMsg)
	Remove(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg)
}

// Service is an IService implementation.
type Service struct {
	queue       Queue
	userService users.IService
}

func NewService(queue Queue, userService users.IService) IService {
	var c IService
	c = &Service{
		queue: queue,
		userService: userService,
	}
	return c
}

// GetAll returns a paginated list of elements from the queue.
// If no page or perPage arguments are passed, it sets those value to 0 and 10 respectively.
func (s *Service) GetAll(ctx context.Context, user *fuel.User, page, perPage *int) ([]interface{}, *ign.ErrMsg) {
	if ok := s.userService.IsSystemAdmin(*user.Username); !ok {
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
	return s.queue.Get(&offset, &limit)
}

// Count returns the element count from the queue.
func (s *Service) Count(ctx context.Context, user *fuel.User) (interface{}, *ign.ErrMsg) {
	if ok := s.userService.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.queue.Count(), nil
}

// MoveToFront moves an element by the given groupID to the front of the queue.
func (s *Service) MoveToFront(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := s.userService.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.queue.MoveToFront(groupID)
}

// MoveToBack moves an element by the given groupID to the back of the queue.
func (s *Service) MoveToBack(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := s.userService.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.queue.MoveToBack(groupID)
}

// Swap swaps positions of groupIDs A and B.
func (s *Service) Swap(ctx context.Context, user *fuel.User, groupIDA, groupIDB string) (interface{}, *ign.ErrMsg) {
	if ok := s.userService.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.queue.Swap(groupIDA, groupIDB)
}

// Remove removes an element by the given groupID from the queue.
func (s *Service) Remove(ctx context.Context, user *fuel.User, groupID string) (interface{}, *ign.ErrMsg) {
	if ok := s.userService.IsSystemAdmin(*user.Username); !ok {
		return nil, ign.NewErrorMessage(ign.ErrorUnauthorized)
	}
	return s.queue.Remove(groupID)
}
