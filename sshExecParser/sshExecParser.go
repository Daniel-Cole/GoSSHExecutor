package sshExecParser

import (
	"regexp"
	"os"
	"bufio"
	"github.com/daniel-cole/GoSSHExecutor/log"
	"fmt"
)

func ParseHosts(hosts []string) []string {
	//TODO: Allow for host name to be specified and check target exists before running
	host_regex := "^[0-9]+\\.[0-9]+\\.[0-9]+\\.[0-9]+$"
	for _, host:= range hosts {
		validHost, _ := regexp.MatchString(host_regex, host)
		if !validHost {
			log.LogFatal(fmt.Sprintf("Failed to parse host: %s", host), nil)
		}
	}
	return hosts
	//TODO: Allow for a range of hosts to be specified.
	//for _, host := range hosts {
	//	host_range, _ := regexp.MatchString(fmt.Sprintf("%s-%s", host_regex, host_regex), host)
	//	if(host_range){
	//
	//	}
	//}
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

