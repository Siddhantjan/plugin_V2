package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strings"
	"time"
)

func Process(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	sshHost := credentials["ip"].(string)
	sshPort := int(credentials["port"].(float64))
	sshUser := credentials["username"].(string)
	sshPassword := credentials["password"].(string)

	config := &ssh.ClientConfig{
		Timeout:         10 * time.Second,
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
			credentials["error"] = err.Error()
			credentials["status"] = "fail"

		} else {
			cmd := "ps aux | awk ' {print $1 \" \"$2 \" \"  $3 \" \"  $4 \" \" $11}'"
			combo, err := session.CombinedOutput(cmd)
			if err != nil {
				credentials["status"] = "fail"
				credentials["error"] = err.Error()

			}
			output := string(combo)
			splitByLine := strings.Split(output, "\n")
			var processes []map[string]interface{}
			for index := 1; index < len(splitByLine)-1; index++ {
				processValue := make(map[string]interface{})
				res := strings.Split(splitByLine[index], " ")
				processValue["process.user"] = res[0]
				processValue["process.id"] = res[1]
				processValue["process.memory.percentage"] = res[3]
				processValue["process.command"] = res[4]
				processes = append(processes, processValue)
			}

			credentials["processes"] = processes
			credentials["status"] = "success"

		}

	}
}
