package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strings"
	"time"
)

func System(credentials map[string]interface{}) {
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

			terminalCommand := "uname -a | awk  '{ print $1 \" \" $2  \" \" $4 \" \"$6 \" \" $7 \" \" $8 \" \"$9 }'"
			combo, er := session.CombinedOutput(terminalCommand)
			output := string(combo)

			res := strings.Split(output, "\n")
			systemValue := strings.Split(res[0], " ")

			credentials["system.os.name"] = systemValue[0]
			credentials["system.user.name"] = systemValue[1]
			credentials["system.os.version"] = systemValue[2]
			credentials["system.up.time"] = systemValue[3] + " " + systemValue[4] + " " + systemValue[5] + " " + systemValue[6]

			session.Close()

			session, err = sshClient.NewSession()

			if err != nil {

				credentials["error"] = er.Error()
				credentials["status"] = "fail"

			} else {

				credentials["status"] = "success"

			}
			runningProcess := " vmstat | awk '{print $1 \" \" $2 \" \"  $12}'"

			combo, er = session.CombinedOutput(runningProcess)

			output = string(combo)

			res = strings.Split(output, "\n")

			processValue := strings.Split(res[2], " ")

			credentials["system.running.process"] = processValue[0]
			credentials["system.blocking.process"] = processValue[1]
			credentials["system.context.switching"] = processValue[2]

		}
	}
}
