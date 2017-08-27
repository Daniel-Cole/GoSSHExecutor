package main

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"

	"github.com/daniel-cole/GoSSHExecutor/log"
	"github.com/daniel-cole/GoSSHExecutor/sshclient"
	"github.com/daniel-cole/GoSSHExecutor/sshExecParser"

	"github.com/alexflint/go-arg"
	"golang.org/x/crypto/ssh"
)

type args struct {
	CommandFile string		`arg:"required,help:The path to the file which contains the commands to execute on the target hosts."`
	Username 	string		`arg:"required,help:The username which will be used to connect to the target hosts."`
	SSHKey 		string		`arg:"required,help:The path to your SSH key which is used to connect to the target hosts."`
	SSHKeyPass 	string		`arg:"help:The password for the SSH key."`
	Port 		string		`arg:"help:The port that is used to connect to the target host. Defaults to 22."`
	Concurrent 	bool		`arg:"help:Execute Commands Concurrently. Concurrency is enabled by default. To disable use --concurrent=false"`
	Halt		bool		`arg:"help:If this option is specified the program will terminate after displaying confirmation information."`
	TargetHosts []string 	`arg:"positional,help:The target hosts to execute the specified commands against. This can either be a CIDR block or a list of hosts."`
}

func main(){
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	args := args{}

	//set defaults
	args.Concurrent = true
	args.Port = "22"

	arg.MustParse(&args)

	hosts := sshExecParser.ParseHosts(args.TargetHosts)

	fmt.Println("Hosts to run commands on\n-------------------------------------")
	for index, host := range hosts {
		fmt.Printf("%d: %s\n", index, host)
	}
	fmt.Println("-------------------------------------")


	commands := sshExecParser.ParseCommands(args.CommandFile)

	fmt.Println("Commands to be executed: \n-------------------------------------")
	for index, cmd := range commands {
		fmt.Printf("%d: %s\n", index, cmd)
	}
	fmt.Println("-------------------------------------")


	if args.Halt {
		log.LogInfo.Println("Terminating program prior to execution as halt is set to true.")
		os.Exit(0)
	}

	clientConfig := sshclient.CreateClientConfig(args.Username, args.SSHKey, args.SSHKeyPass)

	results := make(chan string)
	timeout := time.After(5 * time.Second)

	for _, host := range hosts {
		go func(host string, commands []string, port string, clientConfig *ssh.ClientConfig) {
			results <- sshclient.ExecuteCommands(host, commands, port, clientConfig)
		}(host, commands, args.Port, clientConfig)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <- results:
			fmt.Print(res)
		case <- timeout:
			fmt.Println("Timed out waiting for results")
			return
		}
	}

}





