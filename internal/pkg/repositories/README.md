# Repositories package
The goal of this document is to explain how to use the generic repository interface. As it is right now, it supports the following ORMs:
- GORM

This package groups a set of interfaces to `Create`, `Find`, `Update` and `Delete` entities from a certain database. It makes use of the `Entity` interface to define tables that will have the information provided with the repository.

## GORM
We've created a gorm implementation that satisfies the interfaces described before. In the following section you'll find a couple of examples to help you understand how to use this specific implementation in your codebase.

### Setup
First of all, you need to set up a model and implement the Entity interface.

```golang
type Car struct {
    gorm.Model
    Color string
    Owner string
}

func (Car) TableName() string {
	return "cars"
}

func (Car) SingularName() string {
	return "Car"
}

func (Car) PluralName() string {
	return "Cars"
}
```

After that, you'll need to initialize gorm's repository implementation. `NewRepository()` accepts 3 arguments: A `*gorm.DB` instance, an implementation for the `ign.Logger` interface, and a pointer to an entity of type `Car` that we created before.
```golang
func main() {
    db, err := gorm.Open(...)
    if err != nil {
        os.Exit(1)
    }
    
    carLogger := ign.NewLoggerNoRollbar("cars", ign.VerbosityDebug) // Use your own.
    
    repository := repositories.NewRepository(db, carLogger, &Car{})
}
```

And that's it! Now you're ready to start using this implementation.

### Repository
In this section we're going take a look at the different methods to perform CRUD operations with our brand-new entity.

#### Create car
```golang
func CreateCars(repository repositories.Repository) {
	cars := []*Car{
		{
			Color:  "Red",
			Owner:  "OwnerA",
		},
		{
			Color:  "Red",
			Owner:  "OwnerB",
		},
		{
			Color:  "Blue",
			Owner:  "OwnerC",
		},
		{
			Color:  "Blue",
			Owner:  "OwnerD",
		},
		{
			Color:  "Green",
			Owner:  "OwnerE",
		},
	}

	var entities []domain.Entity
	for _, car := range cars {
		entities = append(entities, car)
	}
	output, err := repository.Create(entities)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Cars:", output)
}
```

### Get one car
To get a specific car, you need to use filters:

```golang
func GetCar(repository repositories.Repository) {
	var car Car

	f := repositories.NewGormFilter("owner = ?", "OwnerA")

	err := repository.FindOne(&car, f)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Car:", car)
}
```

### Get cars
```golang
func GetCars(repository repositories.Repository) {
	var cars []Car

	f := repositories.NewGormFilter("color = ? OR color = ?", "Red", "Blue")

	err := repository.Find(&cars, nil, nil, f)
	if err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Cars:", cars)
}
```

### Update cars
```golang
func UpdateCars(repository repositories.Repository) {
	f := repositories.NewGormFilter("color = ?", "Blue")
	data := map[string]interface{}{ "owner": "LightBlue" }
	err := repository.Update(data, f)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```

### Delete cars
```golang
func DeleteCars(repository repositories.Repository) {
	f := repositories.NewGormFilter("color = ?", "Red")
	err := repository.Delete(f)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
```