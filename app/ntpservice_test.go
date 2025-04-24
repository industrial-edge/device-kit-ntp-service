/*
 * Copyright © Siemens 2020 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package app

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	v1 "ntpservice/api/siemens_iedge_dmapi_v1"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"time"
)

type tConfigurator struct {
}

func (c tConfigurator) WriteConfiguration(serverList []string) error {
	return errors.New("Failed WriteConfiguratiob")
}

func Test_VerifyArgsForStartGRPC_WithLessArgs(t *testing.T) {
	//Create App to use
	tApp := CreateServiceApp()

	//Test with 0 argument
	tArgs := []string{}
	tErr := tApp.StartGRPC(tArgs)

	assert.NotNil(t, tErr, "Did not get expected result. Wanted: %q, got: %q", "parameter not supported!", tErr)

	//Test with 1 argument
	tArgs = []string{"ntpserver"}
	tErr = tApp.StartGRPC(tArgs)

	assert.NotNil(t, tErr, "Did not get expected result. Wanted: %q, got: %q", "parameter not supported!", tErr)

	//Test with 2 arguments
	tArgs = []string{"ntpserver", "unix"}
	tErr = tApp.StartGRPC(tArgs)

	assert.NotNil(t, tErr, "Did not get expected result. Wanted: %q, got: %q", "parameter not supported!", tErr.Error())
}

func Test_VerifyArgsForStartGRPC_WithInappropriateArgs(t *testing.T) {
	//Create App to use
	tApp := CreateServiceApp()

	tApp.StartApp()

	tArgs := []string{"ntpserver", "dummy", "11111"}
	tErr := tApp.StartGRPC(tArgs)

	//Kill the goroutine
	tApp.done <- true

	//Connection failure is expected with dummy sock
	assert.Equal(t, "parameter not supported: dummy", tErr.Error(), "Did not get expected result. Wanted: %q, got: %q", "parameter not supported: dummy", tErr.Error())
}

func Test_VerifyArgsForStartGRPC_WithDummySocketForUnix(t *testing.T) {
	//Create App to use
	tApp := CreateServiceApp()

	//wait until system is up and goroutines are running
	tApp.StartApp()
	time.Sleep(time.Second * 2)

	tArgs := []string{"ntpserver", "unix", "/dummy/unix/path.sock"}
	tErr := tApp.StartGRPC(tArgs)

	//Kill the goroutine
	tApp.done <- true

	//Connection failure is expected with dummy sock
	assert.NotNil(t, tErr, "Did not get expected result. got: %q", tErr)
}

func Test_SetNtpServerFailure(t *testing.T) {
	//Prepare Test Data
	var dummyctx context.Context

	//Create App to use
	tApp := CreateServiceApp()

	//inject new configurator
	tApp.configurator = tConfigurator{}

	tApp.StartApp()

	serverList := v1.Ntp{}
	serverList.NtpServer = []string{}

	_, err := tApp.serverInstance.SetNtpServer(dummyctx, &serverList)

	//Kill the goroutine
	tApp.done <- true

	t.Log("Error:", err.Error())
	assert.Contains(t, err.Error(), "Failed to Set", "Did not get expected result. expected: %q got: %q", "Failed to Set", err.Error())
}

func Test_chownSocketFailure(t *testing.T) {
	//Fail the function with Non existing file path
	err := chownSocket("Non/existing/Path", "root", "root")

	assert.NotNil(t, err, "Did not get expected result. got: %q", err)
}

func Test_GetNtpServerFailure(t *testing.T) {
	//Prepare Test Data
	var dummyctx context.Context

	//Create App to use
	tApp := CreateServiceApp()

	//inject new configurator
	tApp.configurator = tConfigurator{}

	tApp.StartApp()
	err := exec.Command("bash", "-c", "mkdir -p /tmp/ntpsec && touch /tmp/ntpsec/ntp.conf").Run()
	t.Log(err)
	if err != nil {
		assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)
	}
	_, err2 := tApp.serverInstance.GetNtpServer(dummyctx, &emptypb.Empty{})

	//Kill the goroutine
	tApp.done <- true

	assert.Nil(t, err2, "Did not get expected result for GetNtpServer() method. Wanted: Nil, got: %q", err2)
}
