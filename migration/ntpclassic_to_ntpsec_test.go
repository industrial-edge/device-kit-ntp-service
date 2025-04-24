/*
 * Copyright © Siemens 2023 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package migration

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"io/fs"
	"log"
	"ntpservice/internal/ntpconfigurator"
	"ntpservice/utils/files"
	. "ntpservice/utils/mocks"
	"os"
	"path/filepath"
	"testing"
)

const defaultNtpSecCommands = "" +
	"tos maxclock 9\n" +
	"server 1.tr.pool.ntp.org\n" +
	"pool 2.tr.pool.ntp.org\n" +
	"restrict default kod nomodify nopeer noquery limited"

const ntpClassicCommands = "" +
	"tos maxclock 11\n" +
	"tos minclock 1 minsane 1\n" +
	"server 0.tr.pool.ntp.org\n" +
	"    server 2.tr.pool.ntp.org\n" +
	"pool 1.tr.pool.ntp.org\n" +
	"restrict default kod nomodify nopeer noquery limited"

func Test_MigrationSkipped_IfMigrationNotRequired(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()

	mocks.fu.On("IsFileExist", ntpSecMigrationFilePath).Return(true, nil)
	mocks.fu.On("IsFileExist", NTPClassicConfPath).Return(false, nil)

	err := migration.Start()
	assert.NoError(t, err)

	mocks.fs.AssertNotCalled(t, "Create", mock.Anything)
	mocks.cmd.AssertNotCalled(t, "Commander")
}
func Test_MigrationDoneWithoutUpgrading_IfMigrationNotRequired(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()

	ntpMigrationContent := NewEmptyMockFile()

	mocks.fu.On("IsFileExist", ntpSecMigrationFilePath).Return(false, nil)
	mocks.fu.On("IsFileExist", NTPClassicConfPath).Return(false, nil)
	mocks.fs.On("MkdirAll", filepath.Dir(ntpSecMigrationFilePath), fs.FileMode(0666)).Return(nil)
	mocks.fs.On("Create", ntpSecMigrationFilePath).Return(ntpMigrationContent, nil)
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()
	actual, err := io.ReadAll(ntpMigrationContent)
	assert.NoError(t, err)

	mocks.fs.AssertCalled(t, "Create", ntpSecMigrationFilePath)
	mocks.cmd.AssertNotCalled(t, "Commander")
	assert.Equal(t, "1.2.x #not-upgraded", string(actual))
}

func Test_NTPSecDefaultConfWillBeDisabled_NTPClassicConfNotExists(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	defaultNtpSecCommands := "server 2.tr.pool.ntp.org\n" + defaultNtpSecCommands
	ntpSecDefaultConf := NewMockFile(defaultNtpSecCommands)

	mocks.fu.On("IsFileExist", ntpSecMigrationFilePath).Return(false, nil)
	mocks.fu.On("IsFileExist", NTPClassicConfPath).Return(false, nil)
	mocks.fs.On("Open", NTPSecConfPath).Return(ntpSecDefaultConf, nil).Once()
	mocks.fs.On("Create", NTPSecConfPath).Return(ntpSecDefaultConf, nil).Once()
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()
	assert.NoError(t, err)

	const expectDisabledCommands = "" +
		"#server 2.tr.pool.ntp.org #iedk-migration\n" +
		"#tos maxclock 9 #iedk-migration\n" +
		"#server 1.tr.pool.ntp.org #iedk-migration\n" +
		"#pool 2.tr.pool.ntp.org #iedk-migration\n" +
		"restrict default kod nomodify nopeer noquery limited\n"

	mocks.fu.AssertCalled(t, "CreateOrUpdateFile", NTPSecConfPath, expectDisabledCommands)
}

func Test_MigrationRuns_IfMigrationRequired(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	mocks.fu.On("IsFileExist", ntpSecMigrationFilePath).Return(false, nil)
	mocks.fu.On("IsFileExist", NTPClassicConfPath).Return(true, nil)
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()
	assert.NoError(t, err)

	mocks.fs.AssertCalled(t, "Open", NTPSecConfPath)
	mocks.fs.AssertCalled(t, "OpenFile", NTPSecConfPath, os.O_APPEND|os.O_WRONLY, fs.FileMode(0666))
	mocks.fs.AssertCalled(t, "Create", ntpSecMigrationFilePath)
	mocks.fs.AssertCalled(t, "Remove", ntpSecBackupConfPath)
	mocks.cmd.AssertCalled(t, "Commander", ntpconfigurator.UpdateSystemTimeCmd)
}

func Test_CommandsCopied_WhenMigrationRuns(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	ntpSecConf := NewEmptyMockFile()
	mocks.fs.On("Open", NTPClassicConfPath).Return(NewMockFile(ntpClassicCommands), nil).Once()
	mocks.fs.On("OpenFile", NTPSecConfPath, os.O_APPEND|os.O_WRONLY, fs.FileMode(0666)).Return(ntpSecConf, nil).Once()
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()

	actual, err := io.ReadAll(ntpSecConf)
	assert.NoError(t, err)

	const expectedCommands = "" +
		"tos maxclock 11 #iedk-migration\n" +
		"tos minclock 1 minsane 1 #iedk-migration\n" +
		"server 0.tr.pool.ntp.org #iedk-migration\n" +
		"server 2.tr.pool.ntp.org #iedk-migration\n" +
		"pool 1.tr.pool.ntp.org #iedk-migration"
	assert.Equal(t, expectedCommands, string(actual))
}

func Test_MigrationDisablesNTPSecDefaults(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	ntpSecConf := NewMockFile(defaultNtpSecCommands)
	mocks.fs.On("Open", NTPClassicConfPath).Return(NewMockFile(ntpClassicCommands), nil).Once()
	mocks.fs.On("Open", NTPSecConfPath).Return(ntpSecConf, nil).Once()
	mocks.fs.On("Create", NTPSecConfPath).Return(ntpSecConf, nil).Once()
	mocks.fs.On("OpenFile", NTPSecConfPath, os.O_APPEND|os.O_WRONLY, fs.FileMode(0666)).Return(ntpSecConf, nil).Once()
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()
	assert.NoError(t, err)

	const expectDisabledCommands = "" +
		"#tos maxclock 9 #iedk-migration\n" +
		"#server 1.tr.pool.ntp.org #iedk-migration\n" +
		"#pool 2.tr.pool.ntp.org #iedk-migration\n" +
		"restrict default kod nomodify nopeer noquery limited\n"
	mocks.fu.AssertCalled(t, "CreateOrUpdateFile", NTPSecConfPath, expectDisabledCommands)
}

func Test_MigrationSkipped_MigrationFileError(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	mocks.fu.On("IsFileExist", ntpSecMigrationFilePath).Return(false, errors.New("an error occurred"))
	mocks.fu.On("IsFileExist", NTPClassicConfPath).Return(true, errors.New("an error occurred"))

	err := migration.Start()

	mocks.fs.AssertNotCalled(t, "Create", mock.Anything)
	mocks.fs.AssertNotCalled(t, "MkdirAll", filepath.Dir(ntpSecMigrationFilePath), fs.FileMode(0666))
	mocks.cmd.AssertNotCalled(t, "Commander", ntpconfigurator.StartNtpSecService)
	assert.ErrorContains(t, err, "an error occurred")
}

func Test_MigrationSkipped_OpeningNewNTPConfError(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	mocks.fs.On("Open", NTPSecConfPath).Return(NewEmptyMockFile(), errors.New(" error occurred: cannot open new ntp.conf"))
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()

	mocks.fs.AssertNotCalled(t, "Create", NTPSecConfPath)
	mocks.cmd.AssertNotCalled(t, "Commander", ntpconfigurator.StartNtpSecService)
	mocks.fs.AssertCalled(t, "Move", ntpSecBackupConfPath, NTPSecConfPath)
	assert.ErrorContains(t, err, " error occurred: cannot open new ntp.conf")
}

func Test_MigrationSkipped_CopyNTPClassicCommandsError(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	mocks.fs.On("OpenFile", NTPSecConfPath, os.O_APPEND|os.O_WRONLY, fs.FileMode(0666)).Return(NewEmptyMockFile(), errors.New("error occurred while copying NTP classic commands"))
	mocks.fs.On("Move", ntpSecBackupConfPath, NTPSecConfPath).Return(nil)
	mocks.fs.On("CreateOrUpdateFile", NTPSecConfPath, "").Return(nil).Once()
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()

	mocks.fs.AssertNotCalled(t, "Create", ntpSecMigrationFilePath)
	mocks.cmd.AssertNotCalled(t, "Commander", ntpconfigurator.StartNtpSecService)
	assert.ErrorContains(t, err, "error occurred while copying NTP classic commands")
}

func Test_MigrationSkipped_RollbackRuns_WritingCommandsToNtpSecError(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	mocks.fs.On("OpenFile", NTPSecConfPath, os.O_APPEND|os.O_WRONLY, fs.FileMode(0666)).Return(NewEmptyMockFile(), errors.New("error occurred"))
	mocks.fu.On("CreateOrUpdateFile", NTPSecConfPath, "").Return(nil).Once()
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()

	mocks.fs.AssertNotCalled(t, "Create", ntpSecMigrationFilePath)
	mocks.fs.AssertCalled(t, "Move", ntpSecBackupConfPath, NTPSecConfPath)
	mocks.cmd.AssertNotCalled(t, "Commander", ntpconfigurator.StartNtpSecService)
	assert.ErrorContains(t, err, "error occurred")
}

func Test_MigrationSkipped_RollbackNotRuns_NtpClassicConfNotFound(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	mocks.fs.On("Open", NTPClassicConfPath).Return(NewEmptyMockFile(), errors.New("error occurred"))
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()

	mocks.fs.AssertCalled(t, "Create", ntpSecMigrationFilePath)
	mocks.fs.AssertNotCalled(t, "Move", ntpSecBackupConfPath, NTPSecConfPath)
	assert.NoError(t, err)
}

func Test_MigrationSkipped_CreatingNtpSecMigrationFileError(t *testing.T) {
	mocks, migration := generateMockNtpSecMigration()
	mocks.fs.On("Create", ntpSecMigrationFilePath).Return(NewEmptyMockFile(), errors.New("error occurred when creating ntp-Sec.migration file"))
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()

	mocks.cmd.AssertNotCalled(t, "Commander", ntpconfigurator.StartNtpSecService)
	assert.ErrorContains(t, err, "error occurred when creating ntp-Sec.migration file")
}

func Test_OsPackageIntegratedToNptClassicToNtpSecMigration(t *testing.T) {
	migration := NewNTPClassicToNTPSecMigration()

	_ = migration.FileSystemOperations.(*files.OsFileSystemOperations)
	_ = migration.FileUtil.(*files.OsFileUtils)
	_ = migration.Utils.(*ntpconfigurator.OsUtils)

	assert.NoError(t, nil)
}

func Test_UserInformed_IfUpdateSystemTimeFails(t *testing.T) {
	var logOutput = bytes.Buffer{}
	log.SetOutput(&logOutput)
	defer log.SetOutput(os.Stdout)

	mocks, migration := generateMockNtpSecMigration()
	mocks.fs.On("Remove", ntpSecBackupConfPath).Return(fmt.Errorf("cannot remove backup file")).Once()
	mocks.cmd.On("Commander", ntpconfigurator.UpdateSystemTimeCmd).Return([]byte{}, fmt.Errorf("cannot update system time"))
	mockRemainingMethodsWithSuccessfulResults(mocks)

	err := migration.Start()

	assert.NoError(t, err)
	assert.Contains(t, logOutput.String(), "Warning: Couldn't update system time: cannot update system time")
	assert.Contains(t, logOutput.String(), "Warning: Cannot remove ntp sec config backup file: /etc/ntpsec/ntp.conf.backup, reason: cannot remove backup file")
}

func Test_ValuesOfConstantsAreValid(t *testing.T) {
	assert.Equal(t, NTPClassicConfPath, "/etc/ntp.conf")
	assert.Equal(t, NTPSecConfPath, "/etc/ntpsec/ntp.conf")
	assert.Equal(t, ntpSecBackupConfPath, "/etc/ntpsec/ntp.conf.backup")
	assert.Equal(t, ntpSecMigrationFilePath, "/etc/iedk/ntp/migration/ntpsec.migration")
	assert.Equal(t, iedkMigrationTag, "#iedk-migration")
	assert.Equal(t, ntpSecVersion, "1.2.x")
}

type NtpToNtpSecMocks struct {
	fs  *MockFileSystem
	fu  *MockFileUtil
	cmd *MockCommander
}

func generateMockNtpSecMigration() (*NtpToNtpSecMocks, NTPClassicToNTPSecMigration) {
	mockFsOp := new(MockFileSystem)
	mockCommanderOp := new(MockCommander)
	mockFuOp := new(MockFileUtil)
	migration := NTPClassicToNTPSecMigration{mockFsOp, mockFuOp, mockCommanderOp}

	return &NtpToNtpSecMocks{fs: mockFsOp, fu: mockFuOp, cmd: mockCommanderOp}, migration
}

// Important!
// If one of the methods below need to be mocked then it must be mocked before calling this method.
// The mock definitions in here will be ignored if the exact mocks were created before calling this function
func mockRemainingMethodsWithSuccessfulResults(m *NtpToNtpSecMocks) {
	//#
	//File System mock
	//mocks for old ntp conf file
	m.fs.On("Open", NTPClassicConfPath).Return(NewEmptyMockFile(), nil).Once()
	//mocks for new ntp conf file
	m.fs.On("Open", NTPSecConfPath).Return(NewEmptyMockFile(), nil).Once()
	m.fs.On("OpenFile", NTPSecConfPath, os.O_APPEND|os.O_WRONLY, fs.FileMode(0666)).Return(NewEmptyMockFile(), nil).Once()
	m.fs.On("Create", NTPSecConfPath).Return(NewEmptyMockFile(), nil).Once()
	//mocks for migration file
	m.fs.On("Create", ntpSecMigrationFilePath).Return(NewEmptyMockFile(), nil).Once()
	m.fs.On("MkdirAll", filepath.Dir(ntpSecMigrationFilePath), fs.FileMode(0666)).Return(nil).Once()
	//mocks for backup file
	m.fs.On("Remove", ntpSecBackupConfPath).Return(nil).Once()
	//mocks for move
	m.fs.On("Move", ntpSecBackupConfPath, NTPSecConfPath).Return(nil)

	//#
	//File Util mock
	m.fu.On("IsFileExist", NTPClassicConfPath).Return(true, nil)
	m.fu.On("IsFileExist", ntpSecMigrationFilePath).Return(false, nil).Once()
	m.fu.On("Copy", NTPSecConfPath, ntpSecBackupConfPath).Return(nil).Once()
	m.fu.On("CreateOrUpdateFile", NTPSecConfPath, mock.Anything).Return(nil).Once()

	//#
	//CMD executions mock
	m.cmd.On("Commander", ntpconfigurator.StopNtpSecService).Return([]byte("NTPD Stopped!"), nil)
	m.cmd.On("Commander", ntpconfigurator.UpdateSystemTimeCmd).Return([]byte("System Time Updated!"), nil)
	m.cmd.On("Commander", ntpconfigurator.StartNtpSecService).Return([]byte("NTPD Started!"), nil)
}
