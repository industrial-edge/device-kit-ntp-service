/*
 * Copyright © Siemens 2023 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package migration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/fs"
	. "ntpservice/utils/files"
	. "ntpservice/utils/mocks"
	"path/filepath"
	"testing"
)

func Test_MigrationNotRequired_OldPathExistsButDir(t *testing.T) {
	mocks, migration := generateMocks()

	mocks.fu.On("IsFileExist", "/opt/lastntpconfigdate.rec").Return(false, nil)
	mocks.fu.On("IsFileExist", "/etc/iedk/lastntpconfigdate.rec").Return(false, nil)

	migration.Start()

	mocks.fs.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mocks.fs.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationNotRequired_OldFileDoesNotExist(t *testing.T) {
	mocks, migration := generateMocks()

	mocks.fu.On("IsFileExist", "/opt/lastntpconfigdate.rec").Return(false, fs.ErrNotExist)
	mocks.fu.On("IsFileExist", "/etc/iedk/lastntpconfigdate.rec").Return(false, nil)

	migration.Start()

	mocks.fs.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mocks.fs.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationNotRequired_OldFileDoesNotReachable(t *testing.T) {
	mocks, migration := generateMocks()

	mocks.fu.On("IsFileExist", "/opt/lastntpconfigdate.rec").Return(false, fmt.Errorf("cannot Access File"))
	mocks.fu.On("IsFileExist", "/etc/iedk/lastntpconfigdate.rec").Return(false, nil)

	migration.Start()

	mocks.fs.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mocks.fs.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationNotRequired_NewFileAlreadyExist(t *testing.T) {
	mocks, migration := generateMocks()

	mocks.fu.On("IsFileExist", "/opt/lastntpconfigdate.rec").Return(true, nil)
	mocks.fu.On("IsFileExist", "/etc/iedk/lastntpconfigdate.rec").Return(true, nil)

	migration.Start()

	mocks.fs.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mocks.fs.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationRequired(t *testing.T) {
	mocks, migration := generateMocks()

	mocks.fu.On("IsFileExist", "/opt/lastntpconfigdate.rec").Return(true, nil)
	mocks.fu.On("IsFileExist", "/etc/iedk/lastntpconfigdate.rec").Return(false, fs.ErrNotExist)
	mocks.fs.On("MkdirAll", filepath.Dir("/etc/iedk/lastntpconfigdate.rec"), fs.FileMode(0666)).Return(nil)
	mocks.fs.On("Move", "/opt/lastntpconfigdate.rec", "/etc/iedk/lastntpconfigdate.rec").Return(nil)

	migration.Start()

	mocks.fs.AssertNumberOfCalls(t, "MkdirAll", 1)
	mocks.fs.AssertNumberOfCalls(t, "Move", 1)
}

func Test_MigrationRequired_CantRun(t *testing.T) {
	mocks, migration := generateMocks()

	mocks.fu.On("IsFileExist", "/opt/lastntpconfigdate.rec").Return(true, nil)
	mocks.fu.On("IsFileExist", "/etc/iedk/lastntpconfigdate.rec").Return(false, fs.ErrNotExist)
	mocks.fs.On("MkdirAll", filepath.Dir("/etc/iedk/lastntpconfigdate.rec"), fs.FileMode(0666)).Return(fmt.Errorf("can't create folder(s) for new file path"))

	migration.Start()

	mocks.fs.AssertNumberOfCalls(t, "MkdirAll", 1)
	mocks.fs.AssertNotCalled(t, "Move", "/opt/lastntpconfigdate.rec", "/etc/iedk/lastntpconfigdate.rec")
}

func Test_OsPackageIntegratedToLastConfigurationTimeOfNTPClientMigration(t *testing.T) {
	migration := NewLastConfigurationTimeOfNTPClientMigration()

	fileUtils := migration.FileUtil.(*OsFileUtils)
	_ = migration.FileSystemOperations.(*OsFileSystemOperations)
	_ = fileUtils.FileSystemOperations.(*OsFileSystemOperations)

	assert.NoError(t, nil)
}

type mocks struct {
	fs *MockFileSystem
	fu *MockFileUtil
}

func generateMocks() (mocks, LastConfigurationTimeOfNTPClientMigration) {
	mockFsOp := new(MockFileSystem)
	mockFileUtil := new(MockFileUtil)

	migration := LastConfigurationTimeOfNTPClientMigration{mockFsOp, mockFileUtil}
	mocks := mocks{mockFsOp, mockFileUtil}

	return mocks, migration
}
