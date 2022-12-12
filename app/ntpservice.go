/*
 * Copyright (c) 2022 Siemens AG
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package app

import (
	"context"
	"errors"
	"os/user"
	"strconv"
	"time"

	"log"
	"net"
	v1 "ntpservice/api/siemens_iedge_dmapi_v1"
	ntpcf "ntpservice/internal/ntpconfigurator"
	"os"

	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DeviceModelService interface {
	StartGRPC(args []string)
	StartApp()
}

type ntpServer struct {
	v1.UnimplementedNtpServiceServer
	channelWr       chan []string
	errWr           chan error
	ntpConfigurator *ntpcf.NtpConfigurator
}

type MainApp struct {
	serverInstance *ntpServer
	configurator   configuratorApi
	done           chan bool
}

type configuratorApi interface {
	WriteConfiguration(serverList []string) error
}

func CreateServiceApp() *MainApp {
	app := MainApp{}
	ut := ntpcf.OsUtils{}
	vt := ntpcf.NewNtpConfigurator(ut)
	app.serverInstance = &ntpServer{
		channelWr:       make(chan []string),
		errWr:           make(chan error),
		ntpConfigurator: vt,
	}
	app.done = make(chan bool)

	app.configurator = ntpcf.NewNtpConfigurator(ut)

	return &app
}

func chownSocket(address string, userName string, groupName string) error {
	us, err1 := user.Lookup(userName)
	group, err2 := user.LookupGroup(groupName)
	if err1 == nil && err2 == nil {
		uid, _ := strconv.Atoi(us.Uid)
		gid, _ := strconv.Atoi(group.Gid)
		err3 := os.Chmod(address, os.FileMode.Perm(0660))
		err4 := os.Chown(address, uid, gid)
		if err3 != nil || err4 != nil {
			return errors.New("file permissions failed")
		} else {
			log.Println(uid, " : ", gid)
			return nil
		}

	} else {
		return errors.New("file permissions failed")
	}
}

// StartGRPC Necessary actions are taken to start the GRPC server.
func (app *MainApp) StartGRPC(args []string) error {
	const message string = "ERROR: Could not start monitor with bad arguments! \n " +
		"Sample usage:\n  ./ntpservice unix /tmp/devicemodel/ntp.socket \n" +
		"  ./ntpservice tcp localhost:50006"

	if len(args) != 3 {
		log.Println(message)
		return errors.New("parameter not supported")
	}
	typeOfConnection := args[1]
	address := args[2]
	if typeOfConnection != "unix" && typeOfConnection != "tcp" {
		log.Println(message)
		return errors.New("parameter not supported: " + typeOfConnection)
	}

	if typeOfConnection == "unix" {

		if err := os.RemoveAll(os.Args[2]); err != nil {
			return errors.New("socket could not removed: " + typeOfConnection)
		}
	}

	lis, err := net.Listen(typeOfConnection, address)

	if err != nil {
		log.Println("Failed to listen: ", err.Error())
		return errors.New("Failed to listen: " + err.Error())

	}
	err = chownSocket(address, "root", "docker")
	if err != nil {
		return err
	}

	log.Print("Started listening on : ", typeOfConnection, " - ", address)
	s := grpc.NewServer()

	v1.RegisterNtpServiceServer(s, app.serverInstance)
	if err := s.Serve(lis); err != nil {
		log.Printf("Failed to serve: %v", err)
		return errors.New("Failed to serve: " + err.Error())
	}

	return nil
}

// StartApp When a request is received by the client, the processes start here.
func (app *MainApp) StartApp() {
	var serverList []string
	// waits for the new ServerList
	go func() {

		for {
			select {
			case <-app.done:
				log.Println("app done!")
				return
			case serverList = <-app.serverInstance.channelWr:
				err := app.configurator.WriteConfiguration(serverList)
				app.serverInstance.errWr <- err
			}
		}
	}()
}

// GRPC method implementations ################################################################################
// ############################################################################################################

//Implementation of RPC method given v1 proto file

// SetNtpServer This method applies the ntp configurations sent by the client
func (n ntpServer) SetNtpServer(ctx context.Context, serverList *v1.Ntp) (*emptypb.Empty, error) {
	log.Println("SetNtpServer() enter")
	log.Println("Values passed by the client to the SetNtpServer() method: ", serverList)
	//pass the server list for WriteConfiguration
	n.channelWr <- serverList.NtpServer
	defer log.Println("SetNtpServer() leave")
	//err result
	if err := <-n.errWr; err != nil {
		log.Println("SetNtpServer() Failed to Set")
		return &emptypb.Empty{}, status.New(codes.Unknown, "Failed to Set").Err()
	}

	//Save ntp last setting time
	currentTime := time.Now()
	ntpSettingTime := currentTime.Format("2006.01.02 15:04:05")
	log.Println("Ntp Last Setting Time : " + ntpSettingTime)

	f, err := os.OpenFile("/opt/lastntpconfigdate.rec", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		log.Println("Error Opening file : ")
	} else {
		_, err = f.WriteString(ntpSettingTime)
	}
	defer f.Close()

	return &emptypb.Empty{}, status.New(codes.OK, "fine").Err()
}

// GetNtpServer ntp configurations in the device are sent to the client.
func (n ntpServer) GetNtpServer(ctx context.Context, e *emptypb.Empty) (serverList *v1.Ntp, err error) {
	log.Println("GetNtpServer() enter")
	valueWithServerPrefix, err := n.ntpConfigurator.GetCurrentNtpServers()
	if err != nil {
		log.Fatalln("GetNtpServer() Failed to GetCurrentNtpServers()")
	}
	serverList = &v1.Ntp{NtpServer: valueWithServerPrefix}
	log.Println("Server list sent to client:", serverList)
	log.Println("GetNtpServer() leave")
	return serverList, status.New(codes.OK, "fine").Err()
}

// GetStatus check ntp peers and synced behaviours with setting date time.
func (n ntpServer) GetStatus(ctx context.Context, e *emptypb.Empty) (status *v1.Status, err error) {
	log.Println("GetStatus() enter")
	status, err = n.ntpConfigurator.GetNtpStatus()
	return status, err
}
