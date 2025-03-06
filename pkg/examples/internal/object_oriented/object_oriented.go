package object_oriented

import "fmt"

type Human struct {
	ID   int    // public, `human.ID` can be called because it starts with a capital letter
	name string // private outside the package because it starts with a lowercase letter
}

func (h Human) Name() string {
	return h.name
}

func (h Human) Title() string {
	return h.name
}

type Healer interface {
	Heal() bool
}

// Guarantee that Doctor implements Healer.
var _ Healer = Doctor{}       // Verify that T implements I.
var _ Healer = (*Doctor)(nil) // Verify that *T implements I.

type DoctorInterface interface {
	Healer
	Name() string
}

var _ DoctorInterface = Doctor{}

type Doctor struct {
	License string
	Human
}

func (d Doctor) Title() string {
	return "Dr. " + d.name
}

func (d Doctor) Heal() bool {
	return true
}

type Sorcerer struct {
	Human
}

func (m Sorcerer) Heal() bool {
	return true
}

func Healing(h Healer) {
	if h.Heal() {
		fmt.Println("Healed. 💪")
	} else {
		fmt.Println("Failed. 😵")
	}
}

type HealthProfile struct {
	Age    int
	Weight int
}

// Guarantee that Human and Doctor implements Equaler.
var _ Equaler = Human{} // Verify that (receiver T) implements I.
var _ Equaler = Doctor{}
var _ Equaler = (*Human)(nil) // Verify that (receiver *T) implements I.
var _ Equaler = (*Doctor)(nil)

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
