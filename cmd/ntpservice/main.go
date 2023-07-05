/*
 * Copyright (c) 2021 Siemens AG
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
	lastSyncTimeMigration := migration.New()
	lastSyncTimeMigration.Run()

	ntpServiceApp := ntpservice.CreateServiceApp()
	ntpServiceApp.StartApp()
	if err := ntpServiceApp.StartGRPC(os.Args); err != nil {
		log.Printf("Cannot start gRPC server! : %s \n", err)
		return
	}
}
