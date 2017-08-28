[![Build Status](https://travis-ci.org/daniel-cole/GoSSHExecutor.svg?branch=master)](https://travis-ci.org/Daniel-Cole/GoSSHExecutor)

# GoSSHExecutor
This is a simple Go application to execute commands on multiple remote hosts
using ssh key authentication

# Arguments
| Argument        | Flag            | Description  | Allowed Values
| ------------- |:---------------------|:-----|:-----------|
| Command File      | -\-commandfile | The path to the file which contains a list of the commands to execute against the target hosts | string |
| Username      | -\-username      |   The specified username to connect to the target hosts as | string |
| SSH key |  -\-sshkey       |    The SSH key for the user that will connect to the target hosts | string | 
| SSH key pass  | -\-sshkeypass   | The password for the SSH key | string |
| Port          | -\-port           | Port number to connect to the target hosts using SSH | string |
| Concurrent          |  -\-concurrent        |  Runs commands concurrently on different hosts.  Defaults to enabled. | boolean
| Halt  | -\-halt  | If this is set to true then the program will terminate prior to execution | boolean
| Target hosts | -\-targethosts | A list of hosts which are IP addresses or CIDR blocks. i.e. 192.168.0.15 or 10.100.1.0/27 | []string

## Under Construction
