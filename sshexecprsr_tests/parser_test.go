package sshexecprsr_tests

import (
	"testing"
	"reflect"

	"github.com/daniel-cole/GoSSHExecutor/sshexecprsr"
)

//testing doesn't cover IPv6 :(

func TestValidHostParse1(t *testing.T){

	hosts := []string{"198.168.24.0", "10.100.22.5", "172.22.0.15"}

	parsedHosts, _ := sshexecprsr.ParseHosts(hosts)

	if !reflect.DeepEqual(parsedHosts, hosts) {
		t.Error("parsed hosts do not equal the expected hosts:", parsedHosts, hosts)
	}

}

func TestValidCIDRParse1(t *testing.T){
	hosts := []string{"10.100.82.0/29", "10.100.70.5"}

	expectedHosts := []string{
		"10.100.82.0", "10.100.82.1", "10.100.82.2",
		"10.100.82.3", "10.100.82.4", "10.100.82.5",
		"10.100.82.6", "10.100.82.7", "10.100.70.5"}

	parsedHosts, err := sshexecprsr.ParseHosts(hosts)
	if err != nil {
		t.Error("failed to parse hosts", hosts)
		return
	}

	if !reflect.DeepEqual(parsedHosts, expectedHosts){
		t.Error("parsed hosts do not equal the expected hosts:", parsedHosts, expectedHosts)
	}
}

func TestInvalidHostParse1(t *testing.T){

	hosts := []string{"bobandalice", "172.16.8.14"}

	parsedHosts, err := sshexecprsr.ParseHosts(hosts)
	if err != nil {
		return
	}

	if parsedHosts != nil {
		t.Error("parsed hosts should return nil", parsedHosts)
	}

}

func TestInvalidHostParse2(t *testing.T){

	hosts := []string{"172.16.8.14", "bob.and.alice.party"}

	parsedHosts, err := sshexecprsr.ParseHosts(hosts)
	if err != nil {
		return
	}

	if parsedHosts != nil {
		t.Error("parsed hosts should return nil", parsedHosts)
	}

}

func TestInvalidCidrParse1(t *testing.T){

	hosts := []string{"10.100.82.0/82"}

	parsedHosts, err := sshexecprsr.ParseHosts(hosts)
	if err != nil {
		return
	}

	if parsedHosts != nil {
		t.Error("parsed hosts should return nil", parsedHosts)
	}

}

func TestInvalidCidrParse2(t *testing.T){

	hosts := []string{"10.100.82.0/68", "172.16.0.18"}

	parsedHosts, err := sshexecprsr.ParseHosts(hosts)
	if err != nil {
		return
	}

	if parsedHosts != nil {
		t.Error("parsed hosts should return nil", parsedHosts)
	}

}