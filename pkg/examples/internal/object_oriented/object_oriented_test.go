package object_oriented

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestObjectOriented(t *testing.T) {

	t.Run("Composition instead of inheritance", func(t *testing.T) {
		human := Human{1, "John"}
		doctor := Doctor{"doctor's diploma", human}

		assert.Equal(t, "John", doctor.Name())
		assert.Equal(t, "Dr. John", doctor.Title())
		assert.Equal(t, "John", human.Title())
	})

	t.Run("Interface", func(t *testing.T) {

		t.Run("Doctor heals patient", func(t *testing.T) {
			doctor := Doctor{Human: Human{name: "Dr. John"}}
			Healing(doctor)
			assert.True(t, doctor.Heal())
		})

		t.Run("Sorcerer heals patient", func(t *testing.T) {
			magician := Sorcerer{Human: Human{name: "Merlin"}}
			Healing(magician)
			assert.True(t, magician.Heal())
		})
	})

	t.Run("Equality", func(t *testing.T) {
		human := Human{1, "John"}
		doctor := Doctor{"doctor's diploma", human}

		assert.False(t, human.Equal(doctor))
		assert.True(t, human.Equal(doctor.Human))
		assert.True(t, doctor.Equal(human))
	})

	t.Run("Type assertion", func(t *testing.T) {
		var human interface{} = Human{1, "John"}
		var doctor interface{} = Doctor{"doctor's diploma", Human{1, "John"}}

		// Interface -> Human, return false because it misses the Heal method
		_, ok := human.(Healer)
		assert.False(t, ok)

		// Interface -> Healer
		healer, ok := doctor.(Healer)
		assert.True(t, ok)

		// Healer -> interface -> Doctor
		_, ok = interface{}(healer).(Doctor)
		assert.True(t, ok)
	})

	// In Go, a nil slice, map, channel, or pointer is not considered equal to a nil interface.
	// This is because the type information is still present even if the value is nil.
	t.Run("nil assertion", func(t *testing.T) {
		t.Run("nil slice", func(t *testing.T) {
			var nilArray []int = nil

			assert.NotEqual(t, nil, nilArray)
			assert.Nil(t, nilArray)
			t.Logf("Type of nilArray: %s, Value: %v", reflect.TypeOf(nilArray), nilArray)
		})

		t.Run("nil map", func(t *testing.T) {
			var nilMap map[string]int = nil

			assert.NotEqual(t, nil, nilMap)
			assert.Nil(t, nilMap)
			t.Logf("Type of nilMap: %s, Value: %v", reflect.TypeOf(nilMap), nilMap)
		})

		t.Run("nil channel", func(t *testing.T) {
			var nilChannel chan int = nil

			assert.NotEqual(t, nil, nilChannel)
			assert.Nil(t, nilChannel)
			t.Logf("Type of nilChannel: %s, Value: %v", reflect.TypeOf(nilChannel), nilChannel)
		})

		t.Run("nil pointer", func(t *testing.T) {
			var nilPointer *int = nil

			assert.NotEqual(t, nil, nilPointer)
			assert.Nil(t, nilPointer)
			t.Logf("Type of nilPointer: %s, Value: %v", reflect.TypeOf(nilPointer), nilPointer)
		})

		t.Run("nil interface", func(t *testing.T) {
			var nilInterface interface{} = nil

			assert.Equal(t, nil, nilInterface)

			t.Run("nil healer", func(t *testing.T) {
				var nilHealer Healer = nil

				assert.Equal(t, nil, nilHealer)
			})
		})
	})
}

func TestEvent(t *testing.T) {
	testCases := []struct {
		name         string
		event        Event
		expectedType EventType
	}{
		{"EventCreated", Event{ID: "1", Type: EventCreated}, EventCreated},
		{"EventDeleted", Event{ID: "2", Type: EventDeleted}, EventDeleted},
		{"EventUpdated", Event{ID: "3", Type: EventUpdated}, EventUpdated},
		{"EventType is not a hard type", Event{ID: "3", Type: "random"}, "random"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedType, tc.event.Type)
		})
	}

	t.Run("deserialization gives the string not matching the type", func(t *testing.T) {
		invalidJSON := `{"id": "1", "type": "read"}`

		var event Event
		err := json.Unmarshal([]byte(invalidJSON), &event)

		assert.NoError(t, err)
		assert.Equal(t, "read", string(event.Type))
	})
}
