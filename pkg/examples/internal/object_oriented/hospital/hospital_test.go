package hospital

import (
	"src/pkg/examples/internal/object_oriented/healer"
	"src/pkg/examples/internal/object_oriented/person"
	"testing"
)

func TestTreat(t *testing.T) {
	h := Hospital{
		healers: []healer.Healer{
			healer.Doctor{},
		},
	}

	h.Treat(person.New("John"))
	healPatient(h.healers[0])
}
