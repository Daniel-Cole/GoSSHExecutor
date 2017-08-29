package main

import
(
	"os"
	"fmt"
	"time"
	"io/ioutil"

	"github.com/daniel-cole/GoSSHExecutor/log"
	"github.com/daniel-cole/GoSSHExecutor/sshclient"
	"github.com/daniel-cole/GoSSHExecutor/sshexecprsr"

	"github.com/alexflint/go-arg"
	"golang.org/x/crypto/ssh"
	"bufio"
	"strings"
	"github.com/howeyc/gopass"
)

type args struct {
	CommandFile string		`arg:"required,help:The path to the file which contains the commands to execute on the target hosts."`
	Username 	string		`arg:"required,help:The username which will be used to connect to the target hosts."`
	SSHKey 		string		`arg:"required,help:The path to your SSH key which is used to connect to the target hosts."`
	SSHKeyPass 	string		`arg:"help:The password for the SSH key."`
	Port 		string		`arg:"help:The port that is used to connect to the target host. Defaults to 22."`
	Concurrent 	bool		`arg:"help:Execute Commands Concurrently. Concurrency is enabled by default. To disable use --concurrent=false"`
	Halt		bool		`arg:"help:If this option is specified the program will terminate after displaying confirmation information."`
	Timeout 	int32		`arg:"help:This is the maximum time that the program will run for in seconds. Defaults to 60."`
	TargetHosts []string 	`arg:"positional,help:The target hosts to execute the specified commands against. This can either be a CIDR block or a list of hosts."`
}

func main(){
	log.Init(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	args := args{}

	//set defaults
	args.Concurrent = true
	args.Port = "22"
	args.Timeout = 60

	arg.MustParse(&args)

	hosts, err := sshexecprsr.ParseHosts(args.TargetHosts)
	if err != nil {
		log.LogFatal("failed to parse hosts", err)
	}

	fmt.Println("Hosts to run commands on\n-------------------------------------")
	for index, host := range hosts {
		fmt.Printf("Host %d: %s\n", index, host)
	}
	fmt.Println("-------------------------------------")

	promptContinue("Please confirm that the hosts above are correct and type 'y' to continue; 'n' to abort:\n")


	commands := sshexecprsr.ParseCommands(args.CommandFile)

	fmt.Println("Commands to be executed: \n-------------------------------------")
	for index, cmd := range commands {
		fmt.Printf("%d: %s\n", index, cmd)
	}
	fmt.Println("-------------------------------------")


	promptContinue("Please confirm that the commands above are correct and type 'y' to continue; 'n' to abort:\n")

	if args.Halt {
		log.LogInfo.Println("Terminating program prior to execution as halt is set to true.")
		os.Exit(0)
	}

	if args.SSHKeyPass == "" {
		args.SSHKeyPass = promptPassword()
	}

	clientConfig := sshclient.CreateClientConfig(args.Username, args.SSHKey, args.SSHKeyPass)

	results := make(chan string)
	timeout := time.After(time.Duration(args.Timeout) * time.Second)

	for _, host := range hosts {
		if args.Concurrent {
			go func(host string, commands []string, port string, clientConfig *ssh.ClientConfig) {
				results <- sshclient.ExecuteCommands(host, commands, port, clientConfig)
			}(host, commands, args.Port, clientConfig)
		} else {
			results <- sshclient.ExecuteCommands(host, commands, args.Port, clientConfig)
		}
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <- results:
			fmt.Printf(fmt.Sprintf("\n%s\n", res))
		case <- timeout:
			fmt.Println("Timed out waiting for results")
			os.Exit(0)
		}
	}

}

func promptContinue(message string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf(message)
	text, _ := reader.ReadString('\n')
	switch strings.TrimRight(text, "\n") {
	case "y":
	case "n":
		log.LogFatal("aborting program due to 'n' character received during confirmation check.", nil)
	default:
		fmt.Println("Input not recognised.")
		promptContinue(message)
	}
}

func promptPassword() string {
	fmt.Println("Please enter your password (If you have specified an ssh key this will be your ssh key password):")
	password, err := gopass.GetPasswdMasked()
	if err != nil {
		log.LogFatal("failed to get get password input", err)
	}
	return string(password)
}





