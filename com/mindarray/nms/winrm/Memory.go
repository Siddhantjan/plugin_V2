package winrm

import (
	"github.com/masterzen/winrm"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strconv"
	"strings"
)

func Memory(credentials map[string]interface{}) {
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
			ac := "Get-WmiObject win32_OperatingSystem |%{\"{0} {1} {2} {3}\" -f $_.totalvisiblememorysize, $_.freephysicalmemory, $_.totalvirtualmemorysize, $_.freevirtualmemory}"
			output, _, _, err = client.RunPSWithString(ac, a)
			res1 := strings.Split(output, " ")

			totalSpaceMemory, _ := strconv.ParseInt(strings.TrimSpace(res1[0]), 10, 64)
			totalSpaceVirtual, _ := strconv.ParseInt(strings.TrimSpace(res1[2]), 10, 64)
			freeSpaceMemory, _ := strconv.ParseInt(strings.TrimSpace(res1[1]), 10, 64)
			freeSpaceVirtual, _ := strconv.ParseInt(strings.TrimSpace(res1[3]), 10, 64)
			totalSpace := float64(totalSpaceMemory + totalSpaceVirtual)
			freeSpace := float64(freeSpaceVirtual + freeSpaceMemory)
			percent := float64(freeSpace/totalSpace) * 100

			credentials["memory.total.bytes"] = totalSpaceMemory * 1000
			credentials["memory.free.bytes"] = freeSpaceMemory * 1000
			credentials["memory.used.bytes"] = (totalSpaceMemory - freeSpaceMemory) * 1000
			credentials["memory.virtual.total.bytes"] = totalSpaceVirtual * 1000
			credentials["memory.virtual.free.bytes"] = freeSpaceVirtual * 1000
			credentials["memory.virtual.used.bytes"] = (totalSpaceVirtual - freeSpaceVirtual) * 1000
			credentials["memory.used.percent"] = percent
			credentials["memory.available.percent"] = 100.0 - percent
			credentials["status"] = "success"
			clients.Close()

		}
	}
}
