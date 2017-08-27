package sshclient

import (
	"fmt"
	"bytes"
	"io/ioutil"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"github.com/daniel-cole/GoSSHExecutor/log"
)

func CreateClientConfig(username string, sshKey string, sshKeyPass string) *ssh.ClientConfig{
	authMethod := getPrivateKeyFile(sshKey, []byte(sshKeyPass))

	return &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func ExecuteCommands(hostname string, commands []string, port string, clientConfig *ssh.ClientConfig) string {

	log.LogInfo.Printf("Running commands on host %s over port %s as user: %s\n", hostname, port, clientConfig.User)
	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", hostname, port), clientConfig)
	if err != nil {
		log.LogFatal("failed to create ssh connection", err)
	}

	session, _ := conn.NewSession()
	defer session.Close()

	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf

	for _, cmd := range commands {
		session.Run(cmd)
	}

	return fmt.Sprintf("%s -> %s", hostname, stdoutBuf.String())
}

func decryptKey(key []byte, password []byte) []byte {
	block, rest := pem.Decode(key)
	if len(rest) > 0 {
		log.LogFatal("Unexpected data found in private key file", nil)
	}

	if x509.IsEncryptedPEMBlock(block) {
		der, err := x509.DecryptPEMBlock(block, password)
		if err != nil {
			log.LogFatal("Decrypt failed", err)
		}
		log.LogInfo.Println("Successfully decrypted private key with the provided password")
		return pem.EncodeToMemory(&pem.Block{Type: block.Type, Bytes: der})
	}
	return nil
}

func getPrivateKeyFile(file string, password []byte) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		log.LogFatal("Failed to read private SSH key", err)
	}

	key, err := ssh.ParsePrivateKey(decryptKey(buffer, password))
	if err != nil {
		log.LogFatal("Failed to parse private SSH key", err)
	}

	return ssh.PublicKeys(key)
}