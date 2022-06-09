package winrm

import (
	"github.com/masterzen/winrm"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strings"
)

func System(credentials map[string]interface{}) {

	defer exception.ErrorHandle(credentials)

	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)

	client, err := winrm.NewClient(endpoint, username, password)

	if err != nil {
		credentials["error"] = err.Error()
		credentials["status"] = "fail"

	} else {
		clients, er := client.CreateShell()
		if er != nil {
			userReadableError := strings.Contains(er.Error(), "connection refused")
			if userReadableError {
				credentials["error"] = "wrong ip or port ( connection refused )"
			} else {
				credentials["error"] = "wrong username or password ( unable to authenticate )"
			}
			credentials["status"] = "fail"
		} else {

			a := "aa"

			output := ""

			ac := "(Get-WmiObject win32_operatingsystem).name;(Get-WMIObject win32_operatingsystem).version;whoami;(Get-WMIObject win32_operatingsystem).LastBootUpTime;"

			output, _, _, err = client.RunPSWithString(ac, a)

			res1 := strings.Split(output, "\n")

			credentials["system.os.name"] = strings.Replace(strings.Split(res1[0], "\r")[0], "\\", ": ", 9)
			credentials["system.os.version"] = strings.Split(res1[1], "\r")[0]
			credentials["system.user.name"] = strings.Replace(strings.Split(res1[2], "\r")[0], "\\", ": ", 2)
			credentials["system.up.time"] = strings.Split(res1[3], "\r")[0]

			credentials["status"] = "success"
			clients.Close()
		}
	}
}
