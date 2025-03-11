package healer

import (
	"src/pkg/examples/internal/object_oriented/person"
)

type Healer interface {
	Heal() bool
}

// Guarantee that Doctor implements Healer.
var _ Healer = Doctor{}       // Verify that T implements I.
var _ Healer = (*Doctor)(nil) // Verify that *T implements I.

type Licence struct {
	Number string
}

type Doctor struct {
	Licence Licence
	person.Human
}

func (d Doctor) Name() string {
	return "Dr. " + d.Human.Name() // d.name only works if in the same package as Human
}

func (d Doctor) Heal() bool {
	return true
}

type Sorcerer struct {
	person.Human
}

func (m Sorcerer) Heal() bool {
	return true
}

type Wizard struct {
	person.Human
}

type Magician struct {
	Wizard
	Sorcerer
}

type Butcher struct {
	person.Human
}

type Surgeon interface {
	Name() string
	Heal() bool
}

type surgeon struct {
	Doctor
	Butcher
}

func NewSurgeon(name string) Surgeon {
	human := person.New(name)
	return surgeon{
		Doctor: Doctor{
			Licence: Licence{Number: "1234"},
			Human:   human,
		},
		Butcher: Butcher{Human: human},
	}
}
