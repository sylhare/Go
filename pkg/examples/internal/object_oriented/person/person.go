package person

import (
	"math/rand"
	"src/pkg/examples/internal/object_oriented"
)

type Human struct {
	ID   int
	name string
}

func New(name string) Human {
	return Human{
		ID:   rand.New(rand.NewSource(2)).Int(),
		name: name,
	}
}

func (h Human) Name() string {
	return h.name
}

func (h Human) Title() string {
	return h.name
}

func (h Human) Profile() object_oriented.HealthProfile {
	return object_oriented.HealthProfile{}
}

// var _ hospital.Patient = Human{} // creates circular dependency
