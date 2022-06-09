package winrm

import (
	"github.com/masterzen/winrm"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strconv"
	"strings"
)

func Disk(credentials map[string]interface{}) {
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
			ac := "Get-WmiObject win32_logicaldisk |Foreach-Object {$_.DeviceId,$_.Freespace,$_.Size}"
			output, _, _, err = client.RunPSWithString(ac, a)
			res := strings.Split(output, "\r\n")
			var disks []map[string]interface{}
			var usedBytes int64
			var totalBytes int64
			for index := 0; index < len(res); index = index + 3 {
				disk := make(map[string]interface{})
				disk["disk.name"] = strings.Split(res[index], ":")[0]

				if (index+2) > len(res) || res[index+1] == "" {
					disk["disk.free.bytes"] = 0
					disk["disk.total.bytes"] = 0
					disk["disk.available.bytes"] = 0
					disk["disk.used.percent"] = 0
					disk["disk.free.percent"] = 0
					disks = append(disks, disk)
					break
				}

				bytes, _ := strconv.ParseInt(res[index+1], 10, 64)
				usedBytes = usedBytes + bytes
				disk["disk.available.bytes"], _ = strconv.ParseInt(res[index+1], 10, 64)
				bytes, _ = strconv.ParseInt(res[index+2], 10, 64)
				totalBytes = totalBytes + bytes

				disk["disk.total.bytes"] = bytes
				disk["disk.used.bytes"] = (disk["disk.total.bytes"]).(int64) - (disk["disk.available.bytes"]).(int64)
				disk["disk.used.percent"] = ((float64((disk["disk.total.bytes"]).(int64)) - float64((disk["disk.used.bytes"]).(int64))) / float64((disk["disk.total.bytes"].(int64)))) * 100
				disk["disk.free.percent"] = 100 - disk["disk.used.percent"].(float64)
				disks = append(disks, disk)

			}
			credentials["disk.total.bytes"] = totalBytes
			credentials["disk.used.byes"] = usedBytes
			credentials["disk.available.bytes"] = totalBytes - usedBytes
			credentials["disk.used.percent"] = ((float64(totalBytes) - float64(usedBytes)) / float64(totalBytes)) * 100
			credentials["disk.available.percent"] = 100.00 - (credentials["disk.used.percent"]).(float64)
			credentials["disks"] = disks
			credentials["status"] = "success"
			clients.Close()
		}
	}
}
