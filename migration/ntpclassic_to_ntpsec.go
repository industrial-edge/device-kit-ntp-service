/*
 * Copyright © Siemens 2023 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package migration

import (
	"bufio"
	"fmt"
	"io"
	"log"
	ntpcf "ntpservice/internal/ntpconfigurator"
	. "ntpservice/utils/files"
	"os"
	"path/filepath"
	"strings"
)

const NTPClassicConfPath = "/etc/ntp.conf"
const NTPSecConfPath = "/etc/ntpsec/ntp.conf"
const ntpSecBackupConfPath = "/etc/ntpsec/ntp.conf.backup"
const ntpSecMigrationFilePath = "/etc/iedk/ntp/migration/ntpsec.migration"
const iedkMigrationTag = "#iedk-migration"
const ntpSecVersion = "1.2.x"

var ntpMigrationCommands = []string{"server", "tos", "pool"}

type NTPClassicToNTPSecMigration struct {
	FileSystemOperations
	FileUtil
	ntpcf.Utils
}

func NewNTPClassicToNTPSecMigration() NTPClassicToNTPSecMigration {
	fileSystem := OsFileSystemOperations{}
	fileUtil := OsFileUtils{FileSystemOperations: &fileSystem}

	return NTPClassicToNTPSecMigration{
		&fileSystem,
		&fileUtil,
		&ntpcf.OsUtils{}}
}

func (migration *NTPClassicToNTPSecMigration) isRequired() (bool, error) {
	log.Println("Checking migration requirement for `ntp-classic to ntp-sec`.")
	migrationFileExists, err1 := migration.IsFileExist(ntpSecMigrationFilePath)
	oldNtpConfExists, err2 := migration.IsFileExist(NTPClassicConfPath)

	if err1 != nil {
		return false, err1
	} else if err2 != nil {
		return false, err2
	} else {
		return !migrationFileExists && oldNtpConfExists, nil
	}
}

func (migration *NTPClassicToNTPSecMigration) Start() error {
	if isRequired, err := migration.isRequired(); err != nil {
		log.Printf("Error: Cannot check is ntp-classic to ntp-sec migration is required, reason: %s, related path: %s", err.Error(), ntpSecMigrationFilePath)

		return err
	} else if isRequired {
		log.Printf("`ntp-classic to ntp-sec` migration is required!, Running the Migration")

		return migration.run()
	} else {
		if migrationFileExists, err := migration.IsFileExist(ntpSecMigrationFilePath); !migrationFileExists {
			if err := migration.commentOutNTPSecDefaultsConfigurations(); err != nil {
				return fmt.Errorf("Cannot disable defaults of ntpsec configuration: %s\n", err.Error())
			}
			if err != nil {
				return fmt.Errorf("Cannot read migration file (%s), err: %s\n", ntpSecMigrationFilePath, err.Error())
			} else if err := migration.createNtpSecMigrationFile(false); err != nil {
				return fmt.Errorf("Creating file  %s failed: %s\n", ntpSecMigrationFilePath, err.Error())
			}
			log.Printf("Skipping `ntp-classic to ntp-sec` migration since it is not required !")
		} else {
			log.Printf("Skipping `ntp-classic to ntp-sec` migration since it has already been done!")
		}

	}

	return nil
}

func (migration *NTPClassicToNTPSecMigration) run() error {
	if err := migration.backupDefaultNTPSecConf(); err != nil {
		return err
	} else if err := migration.migrateNtpSecToNtpClassic(); err != nil {
		migration.rollback()
		return err
	} else {
		migration.finalize()

		log.Println("Migration is successfully done for `ntp-classic to ntp-sec`.")
		return nil
	}
}

func (migration *NTPClassicToNTPSecMigration) migrateNtpSecToNtpClassic() error {
	if err := migration.commentOutNTPSecDefaultsConfigurations(); err != nil {
		return fmt.Errorf("Disabling ntp-sec default configurations failed: %s\n", err.Error())
	}

	if err := migration.copyNTPClassicCommandsToNTPSec(); err != nil {
		return fmt.Errorf("Copying ntp-classic to ntp-sec commands failed: %s\n", err.Error())
	}

	if err := migration.createNtpSecMigrationFile(true); err != nil {
		return fmt.Errorf("Creating file  %s failed: %s\n", ntpSecMigrationFilePath, err.Error())
	}

	return nil
}

func (migration *NTPClassicToNTPSecMigration) copyNTPClassicCommandsToNTPSec() error {
	oldConfig, err := migration.Open(NTPClassicConfPath)
	if err != nil {
		log.Printf("File %s not found. Migration ends.", NTPClassicConfPath)
		return nil
	}
	defer oldConfig.Close()

	newConfig, err := migration.OpenFile(NTPSecConfPath, os.O_APPEND|os.O_WRONLY, ntpcf.DefaultResourcePermissions)
	if err != nil {
		return err
	}
	defer newConfig.Close()

	previousCommands, err := migration.fetchNTPClassicCommands(oldConfig)
	if err != nil {
		return fmt.Errorf("Fetching ntp-classic commands from %s failed: %s\n", NTPClassicConfPath, err.Error())
	}

	if err := migration.appendCommandsToConfigFile(previousCommands, newConfig); err != nil {
		return fmt.Errorf("Injecting commands failed: %s\n", err.Error())
	}
	return nil
}

func (migration *NTPClassicToNTPSecMigration) createNtpSecMigrationFile(isUpgraded bool) error {
	if err := migration.MkdirAll(filepath.Dir(ntpSecMigrationFilePath), resourcePermissions); err != nil {
		return err
	}
	if migrationFile, err := migration.Create(ntpSecMigrationFilePath); err != nil {
		return err
	} else {
		migrationFileWriter := bufio.NewWriter(migrationFile)
		content := ntpSecVersion
		if !isUpgraded {
			content += " #not-upgraded"
		}

		if _, err := migrationFileWriter.WriteString(content); err != nil {
			return err
		}
		if err = migrationFileWriter.Flush(); err != nil {
			return err
		}
		defer migrationFile.Close()
	}
	return nil
}

func (migration *NTPClassicToNTPSecMigration) commentOutNTPSecDefaultsConfigurations() error {
	if file, err := migration.Open(NTPSecConfPath); err == nil {
		scanner := bufio.NewScanner(file)
		configBuilder := strings.Builder{}

		migration.findAndCommentOut(scanner, &configBuilder)

		if err = migration.CreateOrUpdateFile(NTPSecConfPath, configBuilder.String()); err != nil {
			return fmt.Errorf("Cannot update file content %s, error:%s\n", NTPSecConfPath, err.Error())
		}
		defer file.Close()
	} else {
		return fmt.Errorf("%s not found. Migration Failed: %s", NTPSecConfPath, err.Error())
	}

	return nil
}

func (migration *NTPClassicToNTPSecMigration) findAndCommentOut(scanner *bufio.Scanner, configBuilder *strings.Builder) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		for _, command := range ntpMigrationCommands {
			if strings.HasPrefix(line, command) {
				line = fmt.Sprintf("#%s %s", line, iedkMigrationTag)
			}
		}
		configBuilder.WriteString(line + "\n")
	}
}

func (migration *NTPClassicToNTPSecMigration) appendCommandsToConfigFile(commands []string, file io.Writer) error {
	writer := bufio.NewWriter(file)

	if _, err := writer.WriteString(strings.Join(commands, "\n")); err != nil {
		return err
	}

	if err := writer.Flush(); err != nil {
		return err
	}

	return nil
}

func (migration *NTPClassicToNTPSecMigration) fetchNTPClassicCommands(reader io.Reader) ([]string, error) {
	var ntpCommandList []string

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		for _, command := range ntpMigrationCommands {
			if strings.HasPrefix(line, command) {
				ntpCommandList = append(ntpCommandList, fmt.Sprintf("%s #iedk-migration", line))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return ntpCommandList, nil
}

func (migration *NTPClassicToNTPSecMigration) backupDefaultNTPSecConf() error {
	if err := migration.Copy(NTPSecConfPath, ntpSecBackupConfPath); err != nil {
		return fmt.Errorf("cannot back up ntpsec configuration file: %w", err)
	} else {
		log.Printf("Backed up of ntpsec conf: %s for migration in to: %s", ntpSecBackupConfPath, ntpSecBackupConfPath)
		return nil
	}
}

func (migration *NTPClassicToNTPSecMigration) rollback() {
	if err := migration.Move(ntpSecBackupConfPath, NTPSecConfPath); err != nil {
		log.Printf("Error: Cannot rollback the changes made during migration, err: %s", err.Error())
	} else {
		log.Printf("Rollbacked ntpsec conf file %s  via backup:  %s", NTPSecConfPath, ntpSecBackupConfPath)
	}
}

func (migration *NTPClassicToNTPSecMigration) finalize() {
	if err := migration.Remove(ntpSecBackupConfPath); err != nil {
		log.Printf("Warning: Cannot remove ntp sec config backup file: %s, reason: %s", ntpSecBackupConfPath, err.Error())
	}

	if err := ntpcf.UpdateSystemTime(migration); err != nil {
		log.Printf("Warning: Couldn't update system time: %s ", err.Error())
	}
}
