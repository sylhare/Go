package internal

import "fmt"

type Human struct {
	ID   int    // public, `human.ID` can be called because it starts with a capital letter
	name string // private outside the package because it starts with a lowercase letter
}

func (h Human) GetName() string {
	return h.name
}

type Healer interface {
	Heal() bool
}

// Guarantee that Doctor implements Healer.
var _ Healer = Doctor{}       // Verify that T implements I.
var _ Healer = (*Doctor)(nil) // Verify that *T implements I.

type Doctor struct {
	License string
	Human
}

func (d Doctor) Heal() bool {
	return true
}

type Magician struct {
	Human
}

func (m Magician) Heal() bool {
	return true
}

func HealPatient(h Healer) {
	if h.Heal() {
		fmt.Println("Patient healed successfully.")
	} else {
		fmt.Println("Failed to heal the patient.")
	}
}

// Guarantee that Human and Doctor implements Equaler.
var _ Equaler = Human{}        // Verify that T implements I.
var _ Equaler = Doctor{}       // Verify that T implements I.
var _ Equaler = (*Human)(nil)  // Verify that *T implements I.
var _ Equaler = (*Doctor)(nil) // Verify that *T implements I.

type Equaler interface {
	Equal(Equaler) bool
}

// When facing interface{} type in Go, you can convert it using .(Type) for type assertion.
func (h Human) Equal(e Equaler) bool {
	otherPerson, ok := e.(Human)
	if !ok {
		return false
	}
	return h.name == otherPerson.name
}

type EventType string

// Looks like an enum but does not have any type safety.
const (
	EventCreated EventType = "created"
	EventDeleted EventType = "deleted"
	EventUpdated EventType = "updated"
)

type Event struct {
	ID   string    `json:"id"`
	Type EventType `json:"type"`
}
