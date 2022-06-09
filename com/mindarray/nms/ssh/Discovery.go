package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strings"
	"time"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip"].(string)
	sshPort := int(credentials["port"].(float64))
	sshUser := credentials["username"].(string)
	sshPassword := credentials["password"].(string)

	config := &ssh.ClientConfig{
		Timeout:         6 * time.Second,
		User:            sshUser,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Config: ssh.Config{Ciphers: []string{
			"aes128-ctr", "aes192-ctr", "aes256-ctr",
		}},
	}
	config.Auth = []ssh.AuthMethod{ssh.Password(sshPassword)}
	addr := fmt.Sprintf("%s:%d", sshHost, sshPort)
	sshClient, er := ssh.Dial("tcp", addr, config)

	if er != nil {

		userReadableError := strings.Contains(er.Error(), "connection refused")
		if userReadableError {
			credentials["error"] = "wrong ip or port ( connection refused )"
		} else {
			credentials["error"] = "wrong username or password ( unable to authenticate )"
		}
		credentials["status"] = "fail"

	} else {
		session, err := sshClient.NewSession()
		if err != nil {
			credentials["status"] = "fail"
			credentials["error"] = err.Error()
		} else {
			credentials["status"] = "success"
		}
		cmd := "hostname"
		combo, err := session.CombinedOutput(cmd)
		output := string(combo)
		if err != nil {
			credentials["status"] = "fail"
			credentials["error"] = er.Error()
		} else {
			credentials["status"] = "success"
			credentials["host"] = strings.Split(output, "\n")[0]
		}
	}
}
