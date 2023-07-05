package migration

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io/fs"
	. "ntpservice/utils"
	. "ntpservice/utils/mocks"
	"path/filepath"
	"testing"
)

func Test_MigrationNotRequired_OldPathExistsButDir(t *testing.T) {
	mockFsOp, migration, mockFileInfo := generateMocks()

	mockFileInfo.On("IsDir").Return(true)
	mockFsOp.On("Stat", "/opt/lastntpconfigdate.rec").Return(mockFileInfo, nil)
	mockFsOp.On("Stat", "/etc/iedk/lastntpconfigdate.rec").Return(mockFileInfo, nil)

	migration.Run()

	mockFsOp.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mockFsOp.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationNotRequired_OldFileDoesNotExist(t *testing.T) {
	mockFsOp, migration, mockFileInfo := generateMocks()

	mockFileInfo.On("IsDir").Return(false)
	mockFsOp.On("Stat", "/opt/lastntpconfigdate.rec").Return(mockFileInfo, fs.ErrNotExist)
	mockFsOp.On("Stat", "/etc/iedk/lastntpconfigdate.rec").Return(mockFileInfo, nil)

	migration.Run()

	mockFsOp.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mockFsOp.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationNotRequired_OldFileDoesNotReachable(t *testing.T) {
	mockFsOp, migration, mockFileInfo := generateMocks()

	mockFileInfo.On("IsDir").Return(false)
	mockFsOp.On("Stat", "/opt/lastntpconfigdate.rec").Return(mockFileInfo, fmt.Errorf("cannot Access File"))
	mockFsOp.On("Stat", "/etc/iedk/lastntpconfigdate.rec").Return(mockFileInfo, nil)

	migration.Run()

	mockFsOp.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mockFsOp.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationNotRequired_NewFileAlreadyExist(t *testing.T) {
	mockFsOp, migration, mockFileInfo := generateMocks()

	mockFileInfo.On("IsDir").Return(false)
	mockFsOp.On("Stat", "/opt/lastntpconfigdate.rec").Return(mockFileInfo, nil)
	mockFsOp.On("Stat", "/etc/iedk/lastntpconfigdate.rec").Return(mockFileInfo, nil)

	migration.Run()

	mockFsOp.AssertNotCalled(t, "Move", mock.Anything, mock.Anything)
	mockFsOp.AssertNotCalled(t, "MkdirAll", mock.Anything)
}

func Test_MigrationRequired(t *testing.T) {
	mockFsOp, migration, mockFileInfo := generateMocks()

	mockFileInfo.On("IsDir").Return(false)
	mockFsOp.On("Stat", "/opt/lastntpconfigdate.rec").Return(mockFileInfo, nil)
	mockFsOp.On("Stat", "/etc/iedk/lastntpconfigdate.rec").Return(mockFileInfo, fs.ErrNotExist)
	mockFsOp.On("MkdirAll", filepath.Dir("/etc/iedk/lastntpconfigdate.rec"), fs.FileMode(0666)).Return(nil)
	mockFsOp.On("Move", "/opt/lastntpconfigdate.rec", "/etc/iedk/lastntpconfigdate.rec").Return(nil)

	migration.Run()

	mockFsOp.AssertNumberOfCalls(t, "MkdirAll", 1)
	mockFsOp.AssertNumberOfCalls(t, "Move", 1)
}

func Test_MigrationRequired_CantRun(t *testing.T) {
	mockFsOp, migration, mockFileInfo := generateMocks()

	mockFileInfo.On("IsDir").Return(false)
	mockFsOp.On("Stat", "/opt/lastntpconfigdate.rec").Return(mockFileInfo, nil)
	mockFsOp.On("Stat", "/etc/iedk/lastntpconfigdate.rec").Return(mockFileInfo, fs.ErrNotExist)
	mockFsOp.On("MkdirAll", filepath.Dir("/etc/iedk/lastntpconfigdate.rec"), fs.FileMode(0666)).Return(fmt.Errorf("can't create folder(s) for new file path"))

	migration.Run()

	mockFsOp.AssertNumberOfCalls(t, "MkdirAll", 1)
	mockFsOp.AssertNotCalled(t, "Move", "/opt/lastntpconfigdate.rec", "/etc/iedk/lastntpconfigdate.rec")
}

func generateMocks() (*MockFileSystem, NtpLastSyncFileMigration, *MockFileInfo) {
	mockFsOp := new(MockFileSystem)
	return mockFsOp, NtpLastSyncFileMigration{mockFsOp}, new(MockFileInfo)
}

func Test_OsPackageIntegratedToMigration(t *testing.T) {
	migration := New()

	_ = migration.fs.(*OsFileSystemOperations)

	assert.NoError(t, nil)
}
