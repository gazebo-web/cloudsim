package simulations

import (
	"github.com/go-playground/form"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/application"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/interfaces"
	"gitlab.com/ignitionrobotics/web/cloudsim/pkg/users"
	"gopkg.in/go-playground/validator.v9"
	"testing"
)

func TestIntegration_StartSimulation(t *testing.T) {
	var app interfaces.IApplication
	app = application.NewMock()

	db, err := gorm.Open("sqlite", "test.db")
	if err != nil {
		t.Fail()
	}

	repository := NewRepository(db, app.Name())
	service := NewService(repository)

	input := NewControllerInput{
		SimulationService: service,
		UserService:       users.NewUserServiceMock(),
		FormDecoder:       form.NewDecoder(),
		Validator:         validator.New(),
	}

	controller := NewController(input)


}