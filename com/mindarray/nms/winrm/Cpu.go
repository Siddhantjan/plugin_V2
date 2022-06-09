package winrm

import (
	"github.com/masterzen/winrm"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"regexp"
	"strings"
)

func Cpu(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	host := (credentials["ip"]).(string)
	port := int(credentials["port"].(float64))
	username := credentials["username"].(string)
	password := credentials["password"].(string)

	endpoint := winrm.NewEndpoint(host, port, false, false, nil, nil, nil, 0)
	//result := make(map[string]interface{})
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
			output := ""
			a := "aa"
			ac := "(Get-Counter '\\Processor(*)\\% Idle Time','\\Processor(*)\\% Processor Time','\\Processor(*)\\% user time' -ErrorAction SilentlyContinue).countersamples | Format-List  -Property Path,Cookedvalue;"
			output, _, _, err = client.RunPSWithString(ac, a)
			re := regexp.MustCompile("Path\\s*\\:\\s\\\\+[\\w\\-#]+\\\\(\\w*\\([\\w\\-#]+\\))\\\\%?\\s*(\\w*\\s*\\w*)\\s*\\w*\\s*:\\s*([\\d\\.]+)")
			value := re.FindAllStringSubmatch(output, -1)
			var counters = 3
			var cores []map[string]interface{}
			size := len(value) / counters

			for index := 0; index < len(value)/counters; index++ {
				count := 0
				core := make(map[string]interface{})
				res := strings.Split(value[index][1], "(")
				if strings.Split(res[1], ")")[0] == "_total" {
					credentials["system.cpu.idle.percent"] = value[index][3]
					credentials["system.cpu.process.percent"] = value[count+size][3]
					credentials["system.cpu.user.percent"] = value[count+size+size][3]
				} else {
					core["core.name"] = value[index][2]
					core["core.idle.percent"] = value[index][3]
					core["core.process.percent"] = value[count+size][3]
					core["core.user.percent"] = value[count+size+size][3]

				}
				cores = append(cores, core)
			}
			credentials["cores"] = cores
			//result["ip"] = credentials["ip"]
			//result["metric.group"] = credentials["metric.group"]
			credentials["status"] = "success"
			//data, _ := json.Marshal(result)
			//fmt.Print(string(data))
			clients.Close()

		}

	}
}
