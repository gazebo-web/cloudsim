package users

import (
	"gitlab.com/ignitionrobotics/web/fuelserver/permissions"
	"gitlab.com/ignitionrobotics/web/ign-go"
)

type Service struct {
	Accessor IUserAccessor
}

func NewService(p *permissions.Permissions, sysAdmin string) (*Service, error) {
	s := Service{}
	dbConfig, err := newDbConfig()
	if err != nil {
		return nil, err
	}
	db, err := ign.InitDbWithCfg(dbConfig)
	if err != nil {
		return nil, err
	}

	ua, err := NewUserAccessor(p, db, sysAdmin)
	if err != nil {
		return nil, err
	}
	s.Accessor = ua
	return &s, nil
}