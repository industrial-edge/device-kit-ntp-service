/*
 * Copyright (c) 2021 Siemens AG
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package main

import (
	ntpservice "ntpservice/app"
	"os"
)

func main() {

	ntpServiceApp := ntpservice.CreateServiceApp()
	ntpServiceApp.StartApp()
	ntpServiceApp.StartGRPC(os.Args)
}
