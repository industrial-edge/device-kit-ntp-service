/*
 * Copyright © Siemens 2023 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package mocks

import (
	"github.com/stretchr/testify/mock"
	"io/fs"
	. "ntpservice/utils/files"
	"time"
)

type MockFileSystem struct {
	mock.Mock
}

type MockFileInfo struct {
	mock.Mock
}

func (mock *MockFileSystem) Open(path string) (FileIO, error) {
	args := mock.Called(path)

	return args.Get(0).(*MockFileIO), args.Error(1)
}

func (mock *MockFileSystem) OpenFile(path string, flag int, perm fs.FileMode) (FileIO, error) {
	args := mock.Called(path, flag, perm)

	return args.Get(0).(*MockFileIO), args.Error(1)
}

func (mock *MockFileSystem) Create(name string) (FileIO, error) {
	args := mock.Called(name)

	return args.Get(0).(*MockFileIO), args.Error(1)
}

func (mock *MockFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	return mock.Called(path, perm).Error(0)
}

func (mock *MockFileSystem) Stat(name string) (fs.FileInfo, error) {
	args := mock.Called(name)
	return args.Get(0).(*MockFileInfo), args.Error(1)
}

func (mock *MockFileSystem) Remove(path string) error {
	return mock.Called(path).Error(0)
}

func (mock *MockFileSystem) Move(source string, target string) error {
	return mock.Called(source, target).Error(0)
}

func (mock *MockFileInfo) Name() string {
	return mock.Called().String(0)
}

func (mock *MockFileInfo) Size() int64 {
	return int64(mock.Called().Int(0))
}

func (mock *MockFileInfo) Mode() fs.FileMode {
	return fs.FileMode(mock.Called().Int(0))
}

func (mock *MockFileInfo) ModTime() time.Time {
	return mock.Called().Get(0).(time.Time)
}

func (mock *MockFileInfo) IsDir() bool {
	return mock.Called().Bool(0)
}

func (mock *MockFileInfo) Sys() any {
	return mock.Called().Error(0)
}
