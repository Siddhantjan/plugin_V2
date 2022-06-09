package winrm

import (
	"github.com/masterzen/winrm"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strings"
)

func Discovery(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	var errorOccurred []string
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)
	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	client, err := winrm.NewClient(endpoint, username, password)
	if err != nil {
		credentials["status"] = "fail"

		errorOccurred = append(errorOccurred, err.Error())
	}
	_, err2 := client.CreateShell()
	if err2 != nil {
		userReadableError := strings.Contains(err2.Error(), "connection refused")
		if userReadableError {
			errorOccurred = append(errorOccurred, "wrong ip or port ( connection refused )")
		} else {
			errorOccurred = append(errorOccurred, "wrong username or password ( unable to authenticate )")
		}
		credentials["status"] = "fail"
	}
	if len(errorOccurred) == 0 {
		credentials["status"] = "success"
		a := "aa"
		output := ""
		cmd := "hostname"
		output, _, _, err = client.RunPSWithString(cmd, a)
		if err != nil {
			credentials["status"] = "fail"
			credentials["error"] = err.Error()

		} else {
			credentials["host"] = strings.Split(output, "\r\n")[0]
		}
	} else {
		credentials["status"] = "fail"
		credentials["error"] = errorOccurred
	}
}
