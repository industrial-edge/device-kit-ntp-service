/*
 * Copyright © Siemens 2020 - 2025. ALL RIGHTS RESERVED.
 * Licensed under the MIT license
 * See LICENSE file in the top-level directory
 */

package ntpconfigurator

import (
	"bufio"
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
const ntpSecConfigPath = "/etc/ntpsec/ntp.conf"
const NtpLastConfigPath = "/etc/iedk/lastntpconfigdate.rec"
const DefaultResourcePermissions = 0666
const ntpSecCheckRunning = "/usr/bin/systemctl is-active --quiet ntpsec"
const ntpCheckPeers = "ntpq -pn"
const StartNtpSecService = "/usr/bin/systemctl start ntpsec.service"
const StopNtpSecService = "/usr/bin/systemctl stop ntpsec.service"

// UpdateSystemTimeCmd If the servers are not reachable `ntpd -gq` will never end,
// this will block ntpservice indefinitely, `timeout` used to prevent this behavior
const UpdateSystemTimeCmd = "timeout 20 ntpd -gq"
const CommanderError = "Command(): Error command failed!"

// NewNtpConfigurator It returns a value of type *NtpConfigurator.
func NewNtpConfigurator(utVal Utils) *NtpConfigurator {
	var ntpconfigurator = NtpConfigurator{Ut: utVal}
	return &ntpconfigurator
}

// ReplaceCurrentNtpServersOrPools Deletes the lines starting with pool and server prefixes in /etc/ntpsec/ntp.conf file and all blank lines in the file.
func (n *NtpConfigurator) ReplaceCurrentNtpServersOrPools(serverList []string) {
	file, err := os.Open(ntpSecConfigPath)
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

	// Changes are rewritten to /etc/ntpsec/ntp.conf file.
	err = os.WriteFile(ntpSecConfigPath, []byte(output), 0644)
	if err != nil {
		log.Fatalln(err)
	}
}

// WriteConfiguration The configurations sent by the client are tested and written to /etc/ntpsec/ntp.conf file. Then the ntp service is restarted.
func (n *NtpConfigurator) WriteConfiguration(serverList []string) error {
	n.ReplaceCurrentNtpServersOrPools(serverList)

	err := UpdateSystemTime(n.Ut)
	if err != nil {
		log.Printf("Could not update system time! : %s", err.Error())
	} else {
		log.Printf("Command(`%s`) execution is successfull and system time updated", UpdateSystemTimeCmd)
	}
	return nil
}

func UpdateSystemTime(cmdUtils Utils) error {
	if _, err := cmdUtils.Commander(StopNtpSecService); err != nil {
		log.Println(CommanderError, StopNtpSecService, err)
		return err
	}
	if _, err := cmdUtils.Commander(UpdateSystemTimeCmd); err != nil {
		log.Println(CommanderError, UpdateSystemTimeCmd, err)
		return err
	}
	if _, err := cmdUtils.Commander(StartNtpSecService); err != nil {
		log.Println(CommanderError, StartNtpSecService, err)
		return err
	}
	return nil
}

// GetCurrentNtpServers Lines starting with the server prefix in the /etc/ntpsec/ntp.conf file are sent to the client.
func (n *NtpConfigurator) GetCurrentNtpServers() ([]string, error) {
	var ntpServers []string
	// The contents of /etc/ntpsec/ntp.conf file in the device are read.
	input, err := os.ReadFile(ntpSecConfigPath)
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

	status.IsNtpServiceRunning, err = n.checkRunning(ntpSecCheckRunning)
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
		log.Println(CommanderError, command, err)
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
	data, err = os.ReadFile(NtpLastConfigPath)
	if err != nil {
		log.Println("ReadFile returns error :", err.Error())
		data = []byte(" ")
	}
	LastConfigurationTime = string(data)
	log.Println("LastConfigurationTime-->", LastConfigurationTime)
	return LastConfigurationTime, nil
}
