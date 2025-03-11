package healer

import (
	"github.com/stretchr/testify/assert"
	"src/pkg/examples/internal/object_oriented/person"
	"testing"
)

func TestHealerName(t *testing.T) {
	doctor := Doctor{Human: person.New("John")}
	assert.Equal(t, "Dr. John", doctor.Name())
	assert.NotNil(t, doctor.ID)
}

func TestMagician(t *testing.T) {
	wizard := Wizard{Human: person.New("Gandalf")}
	sorcerer := Sorcerer{Human: person.New("Merlin")}
	magician := Magician{Wizard: wizard, Sorcerer: sorcerer}
	assert.True(t, magician.Heal())
}
