package ntpconfigurator

import (
	"errors"
	"log"
	v1 "ntpservice/api/siemens_iedge_dmapi_v1"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type tOsUtils struct{}

type tOsUtilsStatus struct{}

type tOsUtilsStatuswithError struct{}

func (o tOsUtils) Commander(command string) ([]byte, error) {
	//out, err := exec.Command(Shell, "-c", command).Output()
	dummyStr := []byte("commander called")
	return dummyStr, nil
}

func (o tOsUtilsStatus) Commander(command string) ([]byte, error) {
	var retval string
	var err error
	if command == "systemctl is-active --quiet ntp" {
		retval = " "
		err = errors.New("\nError system is active")
	} else if command == "systemctl is-active --quiet ntpXYZ" {
		retval = " "
		err = errors.New("\nError system is active")

	} else if command == "ntpq -pn" {
		retval = "      remote           refid      st t when poll reach   delay   offset  jitter" +
			"==============================================================================" +
			"0.debian.pool.n .POOL.          16 p    -   64    0    0.000    0.000   0.000" +
			"*193.30.121.7    131.188.3.223    2 u   35  256  377   64.087   40.219  13.950" +
			"+198.251.86.68   82.64.45.50      2 u   62  256  377   74.710   29.980  21.180"
		err = nil
	} else if command == "touch /etc/ntp.conf" {
		err = errors.New("\nError system is active")
	}

	return []byte(retval), err
}

func prepareNtpConfigurator() *NtpConfigurator {
	var tUt Utils = tOsUtils{}
	tN := NewNtpConfigurator(tUt)

	return tN
}

func Test_WriteConfiguration_WithValidArgument(t *testing.T) {
	tServerList := []string{"99.tr.pool.ntp.org"}

	tN := prepareNtpConfigurator()

	err := exec.Command("bash", "-c", "touch /etc/ntp.conf").Run()
	if err != nil {
		//assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)
	}
	err2 := tN.WriteConfiguration(tServerList)

	assert.Nil(t, err2, "Did not get expected result. Wanted: Nil, got: %q", err)

}

func Test_GetCurrentNtpServers(t *testing.T) {
	tN := prepareNtpConfigurator()
	err := exec.Command("bash", "-c", "touch /etc/ntp.conf").Run()
	if err != nil {
		assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)
	}
	_, err2 := tN.GetCurrentNtpServers()

	assert.Nil(t, err2, "Did not get expected result. Wanted: Nil, got: %q", err2)

}

func Test_GetNtpStatus(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	_, err2 := tN.GetNtpStatus()
	assert.Nil(t, err2, "Did not get expected result. Wanted: Nil, got: %q", err2)

}

func Test_ntpStatusCheckSynced(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	_, err := tN.checkSynced()
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)

}

func Test_parseNtpPeer(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	scanner := "*193.30.121.7    131.188.3.223    2 u   12  256  377   64.087   40.219  13.950"
	address := strings.Fields(scanner)
	PeerDetailsDummy, err := tN.parseNtpPeer(address)
	log.Println("WhenValue-->", PeerDetailsDummy)
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)

}

func Test_getNtpPeerWhenValueNonzero(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	scanner := "*193.30.121.7    131.188.3.223    2 u   35  256  377   64.087   40.219  13.950"
	address := strings.Fields(scanner)
	WhenValue, err := tN.getNtpPeerWhenValue(address)
	log.Println("WhenValue-->", WhenValue)
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)
}

func Test_getNtpPeerWhenValueZero(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	scanner := "*193.30.121.7    131.188.3.223    2 u   -  256  377   64.087   40.219  13.950"
	address := strings.Fields(scanner)
	WhenValue, err := tN.getNtpPeerWhenValue(address)
	log.Println("WhenValue-->", WhenValue)
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)
}

func Test_ntpStatusCheckRunning_WithValid(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	var ntpCheckRunning = "systemctl is-active --quiet ntp"
	_, err := tN.checkRunning(ntpCheckRunning)
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)

}

func Test_ntpStatusCheckRunning_WithError(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	var ntpCheckRunning = "systemctl is-active --quiet ntpXYZ"
	_, err := tN.checkRunning(ntpCheckRunning)
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)

}

func Test_ntpStatusGetSyncedTime_WithValidParam(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	var PeerDetails []*v1.PeerDetails
	var PeerDetailsDummy *v1.PeerDetails
	PeerDetailsDummy = &v1.PeerDetails{} // Sets the pointer
	PeerDetailsDummy.RemoteServer = "*193.30.121.7"
	PeerDetailsDummy.ReferenceID = "131.188.3.223"
	PeerDetailsDummy.Stratum = "2"
	PeerDetailsDummy.Type = "u"
	PeerDetailsDummy.Poll = 35
	PeerDetailsDummy.When = 256
	PeerDetailsDummy.Reach = "377"
	PeerDetailsDummy.Delay = 64.087
	PeerDetailsDummy.Offset = 40.219
	PeerDetailsDummy.Jitter = 13.950
	PeerDetails = append(PeerDetails, PeerDetailsDummy)
	PeerDetails = append(PeerDetails, PeerDetailsDummy)

	IsSynced, LastSyncTime, err := tN.getSyncedTime(PeerDetails)
	log.Println("IsSynced-->", IsSynced)
	log.Println("LastSyncTime-->", LastSyncTime)
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)

}

func Test_ntpStatusGetSyncedTime_WithZeroParam(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	var PeerDetails []*v1.PeerDetails
	var PeerDetailsDummy *v1.PeerDetails
	PeerDetailsDummy = &v1.PeerDetails{} // Sets the pointer
	PeerDetailsDummy.RemoteServer = "*193.30.121.7"
	PeerDetailsDummy.ReferenceID = "131.188.3.223"
	PeerDetailsDummy.Stratum = "2"
	PeerDetailsDummy.Type = "u"
	PeerDetailsDummy.Poll = 35
	PeerDetailsDummy.When = 0
	PeerDetailsDummy.Reach = "377"
	PeerDetailsDummy.Delay = 64.087
	PeerDetailsDummy.Offset = 40.219
	PeerDetailsDummy.Jitter = 13.950
	PeerDetails = append(PeerDetails, PeerDetailsDummy)
	PeerDetails = append(PeerDetails, PeerDetailsDummy)

	IsSynced, LastSyncTime, err := tN.getSyncedTime(PeerDetails)
	log.Println("IsSynced-->", IsSynced)
	log.Println("LastSyncTime-->", LastSyncTime)
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)

}

func Test_ntpStatusCheckLastConfiguredOn(t *testing.T) {
	var tUt Utils = tOsUtilsStatus{}
	tN := NewNtpConfigurator(tUt)
	_, err := tN.checkLastConfiguredOn()
	assert.Nil(t, err, "Did not get expected result. Wanted: Nil, got: %q", err)
}
