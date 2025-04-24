/*
 * Copyright © Siemens 2023 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package mocks

import "github.com/stretchr/testify/mock"

type MockCommander struct {
	mock.Mock
}

func (mock *MockCommander) Commander(command string) ([]byte, error) {
	args := mock.Called(command)
	return args.Get(0).([]byte), args.Error(1)
}
