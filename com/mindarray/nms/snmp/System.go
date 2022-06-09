package snmp

import (
	g "github.com/gosnmp/gosnmp"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strings"
	"time"
)

func System(credentials map[string]interface{}) {
	defer exception.ErrorHandle(credentials)
	var version = g.Version1
	switch credentials["version"] {
	case "version1":
		version = g.Version1
		break
	case "version2":
		version = g.Version2c
		break
	case "version3":
		version = g.Version3
		break
	}

	params := &g.GoSNMP{
		Target:    credentials["ip"].(string),
		Port:      uint16(int(credentials["port"].(float64))),
		Community: credentials["community"].(string),
		Version:   version,
		Timeout:   time.Duration(1) * time.Second,
	}
	err := params.Connect()
	var errors []string
	if err != nil {
		credentials["error"] = err.Error()
		credentials["status"] = "fail"
	} else {
		oid := []string{"1.3.6.1.2.1.1.5.0",
			"1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.1.6.0",
			"1.3.6.1.2.1.1.2.0", "1.3.6.1.2.1.1.3.0"}
		value, _ := params.Get(oid)
		for _, variable := range value.Variables {
			switch variable.Name {
			case ".1.3.6.1.2.1.1.5.0":
				credentials["system_name"] = string(variable.Value.([]byte))
				break
			case ".1.3.6.1.2.1.1.1.0":
				credentials["system.description"] = strings.Replace(string(variable.Value.([]byte)), "\r\n", " ", 9)
				break
			case ".1.3.6.1.2.1.1.6.0":
				if len(variable.Value.([]uint8)) == 0 {
					credentials["system.location"] = "empty"
				} else {
					credentials["system.location"] = string(variable.Value.([]byte))
				}
				break
			case ".1.3.6.1.2.1.1.2.0":
				credentials["system.oid"] = variable.Value
				break
			case ".1.3.6.1.2.1.1.3.0":
				credentials["system.upTime"] = variable.Value
				break
			default:
				errors = append(errors, "unknown interface")
			}

		}

		if len(errors) == 0 {
			credentials["status"] = "success"
		} else {
			credentials["status"] = "fail"
			credentials["error"] = errors
		}
	}
}
