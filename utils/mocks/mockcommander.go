package mocks

import "github.com/stretchr/testify/mock"

type MockCommander struct {
	mock.Mock
}

func (mock *MockCommander) Commander(command string) ([]byte, error) {
	args := mock.Called(command)
	return args.Get(0).([]byte), args.Error(1)
}
