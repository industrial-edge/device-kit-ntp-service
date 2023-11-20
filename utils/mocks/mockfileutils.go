package mocks

import "github.com/stretchr/testify/mock"

type MockFileUtil struct {
	mock.Mock
}

func (mock *MockFileUtil) IsFileExist(path string) (bool, error) {
	args := mock.Called(path)

	return args.Bool(0), args.Error(1)
}

func (mock *MockFileUtil) Copy(source string, target string) error {
	args := mock.Called(source, target)
	return args.Error(0)
}

func (mock *MockFileUtil) CreateOrUpdateFile(path string, content string) error {
	args := mock.Called(path, content)
	return args.Error(0)
}
