/*
 * Copyright (c) 2021 Siemens AG
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package ntpconfigurator

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	v1 "ntpservice/api/siemens_iedge_dmapi_v1"
)

// Utils struct
type Utils interface {
	Commander(command string) ([]byte, error)
}

// OsUtils struct
type OsUtils struct{}

// Commander do shell command.
func (o OsUtils) Commander(command string) ([]byte, error) {
	out, err := exec.Command(shell, "-c", command).Output()
	return out, err
}

// NtpConfigurator struct
type NtpConfigurator struct {
	Ut Utils
}

const shell = "bash"
const restartNtpService = "/usr/bin/systemctl restart ntp.service"
const ntpDate = "ntpdate "
const ntpConfigPath = "/etc/ntp.conf"
const ntplastconfigPath = "/etc/lastntpconfigdate.rec"
const ntpCheckRunning = "/usr/bin/systemctl is-active --quiet ntp"
const ntpCheckPeers = "ntpq -pn"

// NewNtpConfigurator It returns a value of type *NtpConfigurator.
func NewNtpConfigurator(utVal Utils) *NtpConfigurator {
	var ntpconfigurator = NtpConfigurator{Ut: utVal}
	return &ntpconfigurator
}

// ReplaceCurrentNtpServersOrPools Deletes the lines starting with pool and server prefixes in /etc/ntp.conf file and all blank lines in the file.
func (n *NtpConfigurator) ReplaceCurrentNtpServersOrPools(serverList []string) {
	file, err := os.Open(ntpConfigPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	builder := strings.Builder{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "pool") && !strings.HasPrefix(line, "server") {
			builder.WriteString(line)
			builder.WriteString("\n")
		}
	}
	for _, val := range serverList {
		builder.WriteString("server " + val + "\n")
	}
	output := builder.String()

	// Changes are rewritten to /etc/ntp.conf file.
	err = ioutil.WriteFile(ntpConfigPath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// WriteConfiguration The configurations sent by the client are tested and written to /etc/ntp.conf file. Then the ntp service is restarted.
func (n *NtpConfigurator) WriteConfiguration(serverList []string) error {
	var command string
	var out []byte
	var err error

	// The ntp configurations entered by the user are controlled by the ntpdate command.
	for _, val := range serverList {
		val = strings.Split(val, " ")[0]
		command = ntpDate + "-u " + val
		out, err = n.Ut.Commander(command)
		if err != nil {
			log.Println("Command(): Error command failed!", command, err)
			continue
		} else {
			log.Println("Command(): ", command, "-> out:", string(out))
			break
		}
	}
	n.ReplaceCurrentNtpServersOrPools(serverList)
	// After the changes, the ntp service is restarted.
	out, err = n.Ut.Commander(restartNtpService)
	return nil
}

// GetCurrentNtpServers Lines starting with the server prefix in the /etc/ntp.conf file are sent to the client.
func (n *NtpConfigurator) GetCurrentNtpServers() ([]string, error) {
	var ntpServers []string
	// The contents of /etc/ntp.conf file in the device are read.
	input, err := ioutil.ReadFile(ntpConfigPath)
	if err != nil {
		log.Fatalln(err)
	}
	lines := strings.Split(string(input), "\n")
	for i, line := range lines {
		// only lines with a server prefix are taken.
		if strings.HasPrefix(line, "server") {
			ntpServers = append(ntpServers, strings.TrimLeft(lines[i], "server "))
		}
	}
	return ntpServers, err
}

// GetNtpStatus is used for checking Ntp running, peers and last configuration times.
func (n *NtpConfigurator) GetNtpStatus() (*v1.Status, error) {

	status := &v1.Status{}
	var err error

	status.IsNtpServiceRunning, err = n.checkRunning(ntpCheckRunning)
	status.PeerDetails, err = n.checkSynced()
	status.IsSynced, status.LastSyncTime, err = n.getSyncedTime(status.PeerDetails)
	status.LastConfigurationTime, err = n.checkLastConfiguredOn()
	return status, err
}

// ntpStatusCheckRunning Check ntp service is running or not with command 'systemctl is-active --quiet ntp'
func (n *NtpConfigurator) checkRunning(ntpCheckRunning string) (bool, error) {
	var IsNtpServiceRunning = true
	command := ntpCheckRunning
	_, err := n.Ut.Commander(command)
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			IsNtpServiceRunning = false
			log.Println("systemctl exit code is ", exitError.ExitCode())
		}
	}
	log.Println("IsNtpServiceRunning-->", IsNtpServiceRunning)
	return IsNtpServiceRunning, nil
}

// ntpStatusCheckSynced Check all ntp remote server peerings conditions with sync time with command 'ntpq -pn'.
func (n *NtpConfigurator) checkSynced() ([]*v1.PeerDetails, error) {

	var out []byte
	var err error
	var PeerDetails []*v1.PeerDetails
	command := ntpCheckPeers
	out, err = n.Ut.Commander(command)
	if err != nil {
		log.Println("Command(): Error command failed!", command, err)
		return PeerDetails, nil
	}
	log.Println("Command(): ", command, "-> out:\n", string(out))
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	scanner.Split(bufio.ScanLines)
	totalLineCount := 0
	for scanner.Scan() {
		totalLineCount++
	}
	peersCount := totalLineCount - 2
	if peersCount > 0 {
		log.Println("Number of peers--> ", peersCount)
		scanner := bufio.NewScanner(strings.NewReader(string(out)))
		var PeerDetailsDummy *v1.PeerDetails
		PeerDetailsDummy = &v1.PeerDetails{} // Sets the pointer
		index := 0
		for scanner.Scan() {
			index++
			if index < 3 {
				continue
			}
			address := strings.Fields(scanner.Text())
			PeerDetailsDummy, _ = n.parseNtpPeer(address)
			PeerDetails = append(PeerDetails, PeerDetailsDummy)

		}
	} else {
		peersCount = 0
	}
	return PeerDetails, nil
}

// parseNtpPeer adjust all peer line elements
func (n *NtpConfigurator) parseNtpPeer(peer []string) (*v1.PeerDetails, error) {
	PeerDetailsDummy := &v1.PeerDetails{}
	var intValue int64
	var floatValue float64
	if len(peer) > 2 {
		PeerDetailsDummy.RemoteServer = peer[0]
		PeerDetailsDummy.ReferenceID = peer[1]
		PeerDetailsDummy.Stratum = peer[2]
		PeerDetailsDummy.Type = peer[3]
		PeerDetailsDummy.When, _ = n.getNtpPeerWhenValue(peer)
		intValue, _ = strconv.ParseInt(peer[5], 10, 32)
		PeerDetailsDummy.Poll = int32(intValue)
		PeerDetailsDummy.Reach = peer[6]
		floatValue, _ = strconv.ParseFloat(peer[7], 8)
		PeerDetailsDummy.Delay = float32(floatValue)
		floatValue, _ = strconv.ParseFloat(peer[8], 8)
		PeerDetailsDummy.Offset = float32(floatValue)
		floatValue, _ = strconv.ParseFloat(peer[9], 8)
		PeerDetailsDummy.Jitter = float32(floatValue)
	}
	return PeerDetailsDummy, nil
}

// ntpStatusGetSyncedTime check when parameters to set synced and lastsynced time.
func (n *NtpConfigurator) getSyncedTime(PeerDetails []*v1.PeerDetails) (bool, string, error) {

	var IsSynced = false
	var LastSyncTime = " "
	if len(PeerDetails) > 0 {
		for i := 0; i < len(PeerDetails); i++ {
			if PeerDetails[i].When > 0 && PeerDetails[i].RemoteServer[:1] == "*" {
				IsSynced = true
				t := time.Now()
				newT := t.Add(time.Duration(-PeerDetails[i].When) * time.Second)
				LastSyncTime = newT.String()
			}
		}
	}
	log.Println("IsSynced-->", IsSynced)
	log.Println("LastSyncTime-->", LastSyncTime)
	return IsSynced, LastSyncTime, nil
}

// GetNtpPeerWhenValue , check character '*' in peer line.
func (n *NtpConfigurator) getNtpPeerWhenValue(PeerLine []string) (int32, error) {
	var intValue int64
	var whenValue int32
	whenValue = 0
	if PeerLine[4] != " -" {
		intValue, _ = strconv.ParseInt(PeerLine[4], 10, 32)
		whenValue = int32(intValue)
	}
	return whenValue, nil
}

// ntpStatusCheckLastConfiguredOn get first ntp setting time from file.
func (n *NtpConfigurator) checkLastConfiguredOn() (string, error) {
	var LastConfigurationTime string
	var data []byte
	var err error
	data, err = ioutil.ReadFile(ntplastconfigPath)
	if err != nil {
		log.Println("ReadFile returns error :", err.Error())
		data = []byte(" ")
	}
	LastConfigurationTime = string(data)
	log.Println("LastConfigurationTime-->", LastConfigurationTime)
	return LastConfigurationTime, nil
}
