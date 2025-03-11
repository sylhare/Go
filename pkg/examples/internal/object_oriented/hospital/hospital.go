package hospital

import (
	"fmt"
	"src/pkg/examples/internal/object_oriented"
	"src/pkg/examples/internal/object_oriented/healer"
	"src/pkg/examples/internal/object_oriented/person"
)

type Hospital struct {
	healers []healer.Healer
}

type HealthProfile struct {
	Age    int
	Weight int
}

type Patient interface {
	Name() string
	Profile() object_oriented.HealthProfile
}

func (h Hospital) Treat(patient Patient) {
	fmt.Printf("Treating patient %s...\n", patient.Name())
	// loop through healers and treat the person
}

func healPatient(h healer.Healer) {
	if h.Heal() {
		fmt.Println("Patient healed. 💪")
	} else {
		fmt.Println("Failed. 😵")
	}
}

var _ Patient = person.Human{}
