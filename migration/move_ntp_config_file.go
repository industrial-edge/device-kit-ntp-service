package migration

import (
	"errors"
	"io/fs"
	"log"
	ntpcf "ntpservice/internal/ntpconfigurator"
	. "ntpservice/utils"
	"path/filepath"
)

const oldPath = "/opt/lastntpconfigdate.rec"
const newPath = ntpcf.NtpLastConfigPath
const resourcePermissions = ntpcf.DefaultResourcePermissions

// NtpLastSyncFileMigration
// NTP Last sync time information needs to be persisted on each sync with the ntp server,
// {currently persistence is done via writing content to a plain file}
// Since previously selected location doesn't fit well, file needs to be moved to a proper directory,
type NtpLastSyncFileMigration struct{ fs FileSystemOperations }

func New() NtpLastSyncFileMigration {
	return NtpLastSyncFileMigration{&OsFileSystemOperations{}}
}

func (migration *NtpLastSyncFileMigration) Run() {
	log.Println("Checking Migration requirement for `last time sync info`")

	if migration.isRequired() {
		err := migration.fs.MkdirAll(filepath.Dir(newPath), resourcePermissions)
		if err == nil {
			err = migration.fs.Move(oldPath, newPath)
		}

		if err != nil {
			log.Printf("Migration failed for `lastime sync info`: %s\n", err.Error())
		} else {
			log.Println("Migration is successfully done for `last time sync info`")
		}
	} else {
		log.Println("Skipped: Migration is not required for `last time sync info`")
	}
}

func (migration *NtpLastSyncFileMigration) isRequired() bool {
	isNewFileExist := migration.isFileExist(newPath)
	isOldFileExist := migration.isFileExist(oldPath)

	return !isNewFileExist && isOldFileExist
}

// isFileExist checks if there is a file in a given path, directory with a same name should not be accepted
// Errors will be ignored, implementation should accept file as not exists even if it exists but not reachable
func (migration *NtpLastSyncFileMigration) isFileExist(path string) bool {
	stat, err := migration.fs.Stat(path)

	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			log.Printf("Cannot access the %s file for migration, details: %s", path, err.Error())
		}

		return false
	}

	if stat != nil && stat.IsDir() {
		//There might be a folder name same as with the target file
		log.Printf("Expected file for migration found directory: %s ", path)
		return false
	}

	return true
}
