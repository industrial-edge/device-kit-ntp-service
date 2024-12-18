/*
 * Copyright (c) Siemens 2021
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */
package main

import (
	"log"
	ntpservice "ntpservice/app"
	"ntpservice/migration"
	"os"
)

func main() {
	runMigrations()

	ntpServiceApp := ntpservice.CreateServiceApp()
	ntpServiceApp.StartApp()
	if err := ntpServiceApp.StartGRPC(os.Args); err != nil {
		log.Printf("Cannot start gRPC server! : %s \n", err)
		return
	}
}

func runMigrations() {
	runLastConfiguredTimeOfNTPClientMigration()
	runNtpClassicToNtpSecMigration()
}

func runNtpClassicToNtpSecMigration() {
	ntpClassicToNTPSecMigration := migration.NewNTPClassicToNTPSecMigration()
	if err := ntpClassicToNTPSecMigration.Start(); err != nil {
		log.Printf("Migration failed, %s", err.Error())
		os.Exit(1)
	}
}

func runLastConfiguredTimeOfNTPClientMigration() {
	lastConfTimeMigration := migration.NewLastConfigurationTimeOfNTPClientMigration()
	lastConfTimeMigration.Start()
}
