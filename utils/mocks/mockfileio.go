package mocks

import (
	"bytes"
)

// MockFileIO Basic Mock File with the functionality of both io.Reader and io.Writer
type MockFileIO struct {
	content *bytes.Buffer
}

func NewEmptyMockFile() *MockFileIO {
	return &MockFileIO{bytes.NewBuffer([]byte{})}
}

func NewMockFile(content string) *MockFileIO {
	return &MockFileIO{bytes.NewBufferString(content)}
}

func (mock *MockFileIO) Read(p []byte) (n int, err error) {
	return mock.content.Read(p)
}

func (mock *MockFileIO) Write(p []byte) (n int, err error) {
	return mock.content.Write(p)
}

func (mock *MockFileIO) Close() error {
	return nil
}
