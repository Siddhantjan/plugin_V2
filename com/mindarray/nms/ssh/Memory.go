package ssh

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strconv"
	"strings"
	"time"
)

func Memory(credentials map[string]interface{}) {
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

			cmd := "free -b | awk  '{if ($1 != \"total\") print $1 \" \" $2 \" \" $3 \" \" $4 \" \"$7}'"
			combo, _ := session.CombinedOutput(cmd)
			output := string(combo)
			res := strings.Split(output, "\n")

			memoryValue := strings.Split(res[0], " ")
			totalBytes, _ := strconv.ParseInt(memoryValue[1], 10, 64)
			credentials["memory.total.bytes"] = totalBytes

			usedBytes, _ := strconv.ParseInt(memoryValue[2], 10, 64)
			credentials["memory.used.bytes"] = usedBytes
			credentials["memory.free.bytes"], _ = strconv.ParseInt(memoryValue[3], 10, 64)
			credentials["memory.available.bytes"], _ = strconv.ParseInt(memoryValue[4], 10, 64)

			swapValue := strings.Split(res[1], " ")
			credentials["memory.swap.total.bytes"], _ = strconv.ParseInt(swapValue[1], 10, 64)
			credentials["memory.swap.used.bytes"], _ = strconv.ParseInt(swapValue[2], 10, 64)
			credentials["memory.swap.free.bytes"], _ = strconv.ParseInt(swapValue[3], 10, 64)
			usedPercent := float64(totalBytes-usedBytes) / float64(totalBytes)

			credentials["memory.used.percent"] = usedPercent
			credentials["memory.available.percent"] = 100 - usedPercent

			credentials["ip"] = credentials["ip"]
			credentials["metric.group"] = credentials["metric.group"]
			credentials["status"] = "success"
		}
	}
}
