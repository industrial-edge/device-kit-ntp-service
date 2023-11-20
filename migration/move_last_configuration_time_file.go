package migration

import (
	"errors"
	"io/fs"
	"log"
	ntpcf "ntpservice/internal/ntpconfigurator"
	. "ntpservice/utils/files"
	"path/filepath"
)

const oldPath = "/opt/lastntpconfigdate.rec"
const newPath = ntpcf.NtpLastConfigPath
const resourcePermissions = ntpcf.DefaultResourcePermissions

// LastConfigurationTimeOfNTPClientMigration
// Last Configuration date-time for NTP Client information needs to be persisted after each configuration done via rpc method `SetNtpServer` call
// {with the current implementation persistence is done via writing content to a plain file}
// Since previously selected location doesn't fit well, file needs to be moved to a proper directory,
// Migration is required to keep previously set last configuration date-time information during update
type LastConfigurationTimeOfNTPClientMigration struct {
	FileSystemOperations
	FileUtil
}

func NewLastConfigurationTimeOfNTPClientMigration() LastConfigurationTimeOfNTPClientMigration {
	fileSystem := &OsFileSystemOperations{}
	fileUtils := &OsFileUtils{FileSystemOperations: fileSystem}

	return LastConfigurationTimeOfNTPClientMigration{fileSystem, fileUtils}
}

func (migration *LastConfigurationTimeOfNTPClientMigration) Start() {
	log.Println("Checking Migration requirement for `last configuration time of ntp client`")

	if migration.isRequired() {
		err := migration.MkdirAll(filepath.Dir(newPath), resourcePermissions)
		if err == nil {
			err = migration.Move(oldPath, newPath)
		}

		if err != nil {
			log.Printf("Migration failed for `last configuration time of ntp client`: %s\n", err.Error())
		} else {
			log.Println("Migration is successfully done for `last configuration time of ntp client`")
		}
	} else {
		log.Println("Skipped: Migration is not required for `last configuration time of ntp client`")
	}
}

func (migration *LastConfigurationTimeOfNTPClientMigration) isRequired() bool {
	isNewFileExist := migration.isFileExist(newPath)
	isOldFileExist := migration.isFileExist(oldPath)

	return !isNewFileExist && isOldFileExist
}

func (migration *LastConfigurationTimeOfNTPClientMigration) isFileExist(path string) bool {
	if ok, err := migration.IsFileExist(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Printf("Cannot access the %s file for `last configuration time of ntp client` migration, details: %s", path, err.Error())
		}
		return false
	} else {
		return ok
	}
}
