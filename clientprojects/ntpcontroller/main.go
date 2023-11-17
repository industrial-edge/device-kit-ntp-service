/*
 * Copyright (c) 2021 Siemens AG
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package main

import (
	"context"
	"log"
	"net"
	ntpapi "ntpcontroller/api"
	"os"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	"google.golang.org/grpc"
)

func createNtpFromOsArgs(args []string) *ntpapi.Ntp {
	ntpData := &ntpapi.Ntp{}
	//ntpData.NtpServer = []string{"0.tr.pool.ntp.org", "1.tr.pool.ntp.org", "2.tr.pool.ntp.org"}
	ntpData.NtpServer = args

	return ntpData
}

func main() {
	// Set up a connection to the server.
	if len(os.Args) < 3 || os.Args[1] != "unix" && os.Args[1] != "tcp" {
		log.Println("You must give an argument when running golang file." +
			"[uds /var/run/devicemodel/edge.sock OR tcp localhost:50006]\nUsage:" +
			" go run main.go uds /var/run/devicemodel/edge.sock OR go run main.go tcp localhost:50006")
		return
	}

	var ntpData *ntpapi.Ntp
	args := os.Args[3:]

	ntpData = createNtpFromOsArgs(args)

	//ntpData.NtpServer = []string{}
	log.Println("ntpserver parameters: ", ntpData.NtpServer)

	var conn *grpc.ClientConn
	var err error
	if os.Args[1] == "tcp" {
		conn, err = grpc.Dial(os.Args[2], grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		} else {
			log.Print("connected to " + os.Args[2])
		}
	}
	if os.Args[1] == "unix" {
		conn, err = grpc.Dial(
			os.Args[2],
			grpc.WithInsecure(),
			grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
				return net.DialTimeout("unix", addr, timeout)
			}))
		if err != nil {
			log.Fatalf("Did not connect: %v", err)
		} else {
			log.Print("connected to " + os.Args[2])
		}
	}

	defer conn.Close()

	sysClient2 := ntpapi.NewNtpServiceClient(conn)
	response2, err2 := sysClient2.GetNtpServer(context.Background(), &emptypb.Empty{})
	log.Println("err2", err2)
	log.Print("response: ", response2)
	testJson, _ := protojson.Marshal(response2)
	log.Println(string(testJson))
	if err2 != nil {
		log.Printf("GET-action could not performed %v", err)
	} else {
		log.Println("GET-message sent : ", response2)
		log.Println("GET-DONE!")
	}

	response3, err3 := sysClient2.GetStatus(context.Background(), &emptypb.Empty{})
	log.Println("err2", err3)
	log.Print("response: ", response3)
	testJson3, _ := protojson.Marshal(response3)
	log.Println(string(testJson3))
	if err3 != nil {
		log.Printf("GET-action could not performed %v", err3)
	} else {
		log.Println("GET-message sent : ", response3)
		log.Println("GET-DONE!")
	}

}
