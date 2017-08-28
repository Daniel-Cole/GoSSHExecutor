package sshexecprsr

import (
	"os"
	"bufio"
	"net"

	"github.com/daniel-cole/GoSSHExecutor/log"
)

//parses a list of IP addresses or CIDR blocks specified in the given hosts parameter
//and returns a list of all the IP addresses
func ParseHosts(hosts []string) ([]string, error) {
	var parsedHosts []string
	for _, host:= range hosts {
		ip := net.ParseIP(host)
		if ip == nil {
			ip, ipnet, err := net.ParseCIDR(host)
			if err != nil {
				return nil, err
			}

			for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
				parsedHosts = append(parsedHosts, ip.String())
			}
		} else {
			parsedHosts = append(parsedHosts, host)
		}
	}

	return parsedHosts, nil
}

func ParseCommands(commandFile string) []string {
	file, err := os.Open(commandFile)

	if err != nil {
		log.LogFatal("Failed to read in specified command file. ", err)
	}
	defer file.Close()

	commands := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		commands = append(commands,scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.LogFatal("Failed to read commands in command file", err)
	}

	return commands
}

func inc(ip net.IP) {
	for j := len(ip)-1; j>=0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

