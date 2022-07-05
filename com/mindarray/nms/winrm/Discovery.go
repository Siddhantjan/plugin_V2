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
			//	errorOccurred = append(errorOccurred, "wrong username or password ( unable to authenticate )")
			errorOccurred = append(errorOccurred, err2.Error())
		}
		credentials["status"] = "fail"
	}
	if len(errorOccurred) == 0 {
		credentials["status"] = "success"
		a := "aa"
		output := ""
		//cmd := "hostname"
		cmd := "(Get-Counter '\\Processor(_Total)\\% Processor Time'," +
			"'\\Processor(_Total)\\% User Time'," +
			"'\\Processor(_Total)\\% Interrupt Time'," +
			"'\\Processor(_Total)\\% Idle Time'," +
			"'\\System\\Processor Queue Length'," +
			"'\\Memory\\Available MBytes'," +
			"'\\Memory\\Committed Bytes'," +
			"'\\Memory\\Commit Limit'," +
			"'\\LogicalDisk(*)\\Disk Write Bytes/sec'," +
			"'\\LogicalDisk(*)\\Disk Read Bytes/sec'," +
			"'\\LogicalDisk(*)\\Disk Reads/sec'," +
			"'\\LogicalDisk(*)\\Disk Writes/sec'," +
			"'\\LogicalDisk(*)\\% Disk Read Time'," +
			"'\\LogicalDisk(*)\\% Disk Write Time'," +
			"'\\LogicalDisk(*)\\% Disk Time'," +
			"'\\LogicalDisk(*)\\Avg. Disk Queue Length'," +
			"'\\Network Interface(*)\\Bytes Received/sec'," +
			"'\\Network Interface(*)\\Output Queue Length'," +
			"'\\Network Interface(*)\\Bytes Sent/sec'," +
			"'\\Network Interface(*)\\Bytes Total/sec'," +
			"'\\Network Interface(*)\\Packets Received/sec'," +
			"'\\Network Interface(*)\\Packets Sent/sec'," +
			"'\\System\\System Up Time','\\System\\Processes'," +
			"'\\System\\Threads','\\Memory\\Pages/sec'," +
			"'\\Memory\\Page Faults/sec'," +
			"'\\Processor(_Total)\\Interrupts/sec'," +
			"'\\System\\Context Switches/sec'," +
			"'\\Memory\\Free & Zero Page List Bytes'," +
			"'\\Memory\\Pool Paged Bytes'," +
			"'\\Memory\\Pool Nonpaged Bytes' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue"
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
