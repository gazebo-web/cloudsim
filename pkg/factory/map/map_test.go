package factorymap

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/suite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/factory"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

const (
	type1 = "type1"
	type2 = "type2"
	type3 = "type3"
)

// Tester is a test interface. The factory is expected to return values that implement Tester
type Tester interface {
	Do() interface{}
}

// Dependency
type Dep struct{}

// Do doubles the passed value
func (d *Dep) Op(value int) int {
	return value * 2
}

// Type 1
type Type1 struct {
	value int
	dep   *Dep
}

// Do passes the value to its dependency and returns the processed value.
func (t *Type1) Do() interface{} {
	return t.dep.Op(t.value)
}

type Type1Config struct {
	Value int
}

type Type1Dependencies struct {
	*Dep `validate:"required"`
}

func (td *Type1Dependencies) Validate() error {
	return validator.New().Struct(td)
}

func Type1NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Cast the passed config into the expected type
	typeConfig := &Type1Config{}
	if err := mapstructure.Decode(config, typeConfig); err != nil {
		return err
	}

	// Get dependencies
	typeDependencies := &Type1Dependencies{}
	// This call will both populate typeDependencies and check that all the required dependencies have been passed
	if err := dependencies.ToStruct(typeDependencies); err != nil {
		return err
	}

	// Create and set the object
	return factory.SetValue(out, &Type1{
		value: typeConfig.Value,
		dep:   typeDependencies.Dep,
	})
}

// Type 2
type Type2 struct {
	value string
}

// Do returns a copy of the string containing the value string twice.
func (t *Type2) Do() interface{} {
	return fmt.Sprint(t.value, t.value)
}

type Type2Config struct {
	Value string
}

func Type2NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Cast the passed config into the expected type
	typeConfig := &Type2Config{}
	if err := mapstructure.Decode(config, typeConfig); err != nil {
		return err
	}

	// Create and set the object
	return factory.SetValue(out, &Type2{
		value: typeConfig.Value,
	})
}

// Type 3
type Type3 struct {
	value bool
}

// Do flips the boolean value.
func (t *Type3) Do() interface{} {
	return !t.value
}

type Type3Config struct {
	Value bool
}

func Type3NewFunc(config interface{}, dependencies factory.Dependencies, out interface{}) error {
	// Cast the passed config into the expected type
	typeConfig := &Type3Config{}
	if err := mapstructure.Decode(config, typeConfig); err != nil {
		return err
	}

	// Create and set the object
	return factory.SetValue(out, &Type3{
		value: typeConfig.Value,
	})
}

func TestFactoryMapSuite(t *testing.T) {
	suite.Run(t, new(factoryMapSuite))
}

type factoryMapSuite struct {
	suite.Suite
	m Map
}

func (s *factoryMapSuite) SetupSuite() {
	// Prepare factory
	s.m = NewMap().(Map)
	s.Nil(s.m.Register(type1, Type1NewFunc))
	s.Nil(s.m.Register(type2, Type2NewFunc))
	s.Nil(s.m.Register(type3, Type3NewFunc))
}

func (s *factoryMapSuite) TestNewType1() {
	// Prepare factory config
	value := 1
	cfg := &factory.Config{
		Type: type1,
		Config: factory.ConfigValues{
			"Value": value,
		},
	}

	// Configure dependencies
	dep := factory.NewDependencies()
	dep.Set("dep", &Dep{})

	var out Tester
	err := s.m.New(cfg, dep, &out)
	s.Nil(err)
	s.NotNil(out)

	// Type 1 Do() doubles a number
	s.Equal(2, out.Do())
}

func (s *factoryMapSuite) TestNewType1MissingDependencies() {
	// Prepare factory config
	value := 1
	cfg := &factory.Config{
		Type: type1,
		Config: factory.ConfigValues{
			"Value": value,
		},
	}

	// Prepare dependencies object without required dependencies
	dep := factory.NewDependencies()

	var out Tester
	s.Error(s.m.New(cfg, dep, &out))
}

func (s *factoryMapSuite) TestNewType2() {
	// Prepare factory config
	value := "test"
	cfg := &factory.Config{
		Type: type2,
		Config: factory.ConfigValues{
			"Value": value,
		},
	}

	var out Tester
	s.Nil(s.m.New(cfg, nil, &out))
	s.NotNil(out)

	// Type 2 Do() returns a copy of the string containing the value string twice.
	s.Equal("testtest", out.Do())
}

func (s *factoryMapSuite) TestNewType3() {
	// Prepare factory config
	value := false
	cfg := &factory.Config{
		Type: type3,
		Config: factory.ConfigValues{
			"Value": value,
		},
	}

	var out Tester
	s.Nil(s.m.New(cfg, nil, &out))
	s.NotNil(out)

	// Type 3 Do() flips its boolean value
	s.Equal(true, out.Do())
}

func (s *factoryMapSuite) TestRegister() {
	// Register Type 3 as another type
	type4 := "type4"
	s.Nil(s.m.Register(type4, Type3NewFunc))

	// Prepare factory config
	value := false
	cfg := &factory.Config{
		Type: type4,
		Config: factory.ConfigValues{
			"Value": value,
		},
	}

	var out Tester
	s.Nil(s.m.New(cfg, nil, &out))
	s.NotNil(out)

	// Type 4 Do() flips its boolean value
	s.Equal(true, out.Do())
}

func (s *factoryMapSuite) TestRegisterExisting() {
	// Attempt to register Type 1
	s.Equal(factory.ErrFactoryTypeAlreadyExists, s.m.Register(type1, Type1NewFunc))
}
