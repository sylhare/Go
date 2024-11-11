package services

import (
	"github.com/stretchr/testify/assert"
	"src/pkg/examples/mocks"
	"testing"
)

// To show how to use the mock
func TestServiceProcess(t *testing.T) {
	mockService := new(mocks.Service)

	mockService.On("Process").Return("mocked result")

	result := mockService.Process()

	assert.Equal(t, "mocked result", result)
	mockService.AssertExpectations(t)
}
