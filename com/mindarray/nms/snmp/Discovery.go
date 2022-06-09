package snmp

import (
	g "github.com/gosnmp/gosnmp"
	exception "pluginCopy/com/mindarray/nms/exceptionHandler"
	"strings"
	"time"
)

func Discovery(credentials map[string]interface{}) {
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
		Timeout:   time.Duration(2) * time.Second,
	}
	err := params.Connect()
	if err != nil {
		credentials["status"] = "fail"
		credentials["error"] = err.Error()

	} else {
		credentials["status"] = "success"
	}
	sysName, err := params.Get([]string{".1.3.6.1.2.1.1.5.0"})

	if err != nil {
		credentials["error"] = err.Error()
		credentials["status"] = "fail"
		userReadableError := strings.Contains(err.Error(), "connection refused")
		if userReadableError {
			credentials["error"] = " connection refused "
		}

	} else {
		credentials["status"] = "success"
		for _, value := range sysName.Variables {
			credentials["host"] = string(value.Value.([]byte))
		}

		walkOid := "1.3.6.1.2.1.2.2.1.1"
		index := "1.3.6.1.2.1.2.2.1.1."
		description := "1.3.6.1.2.1.2.2.1.2."
		name := "1.3.6.1.2.1.31.1.1.1.1."
		operationalStatus := "1.3.6.1.2.1.2.2.1.8."
		alias := "1.3.6.1.2.1.31.1.1.1.18."

		var walkOidArray []string
		walk := params.Walk(walkOid, func(pdu g.SnmpPDU) error {
			switch pdu.Type {
			case g.IPAddress:
				result := pdu.Value
				walkOidArray = append(walkOidArray, result.(string))
				break
			case g.Integer:
				result := g.ToBigInt(pdu.Value)
				walkOidArray = append(walkOidArray, result.String())
				break
			case g.OctetString:
				result := pdu.Value.([]byte)
				walkOidArray = append(walkOidArray, string(result))
				break
			case g.Gauge32:
				result := pdu.Value
				walkOidArray = append(walkOidArray, result.(string))
			default:
				result := pdu.Value
				walkOidArray = append(walkOidArray, result.(string))
			}
			return nil
		},
		)
		if walk != nil {
			credentials["error"] = walk.Error()

		}

		var oids []string
		for i := 0; i < len(walkOidArray); i++ {
			oids = append(oids, index+walkOidArray[i])
			oids = append(oids, description+walkOidArray[i])
			oids = append(oids, name+walkOidArray[i])
			oids = append(oids, operationalStatus+walkOidArray[i])
			oids = append(oids, alias+walkOidArray[i])
		}
		var startIndex = 0
		var endIndex = 40
		var resultArray []interface{}
		for {
			if len(resultArray) == len(oids) {
				break
			}
			output, error := params.Get(oids[startIndex:endIndex])
			if error != nil {
				credentials["error"] = error.Error()

			}
			for _, variable := range output.Variables {

				resultArray = append(resultArray, ConvertData(variable))
			}
			startIndex = endIndex
			endIndex = endIndex + 40

			if endIndex > len(oids) {
				endIndex = len(oids)
			}
		}
		var interfaces []map[string]interface{}
		for i := 0; i < len(resultArray); i = i + 5 {
			interfaceValue := make(map[string]interface{})
			interfaceValue["interface.index"] = resultArray[i].(int)
			interfaceValue["interface.description"] = resultArray[i+1]
			interfaceValue["interface.name"] = resultArray[i+2]
			if resultArray[i+3] == 1 {
				interfaceValue["interface.operational.status"] = "up"
			} else {
				interfaceValue["interface.operational.status"] = "down"
			}
			if resultArray[i+4] == "" {
				interfaceValue["interface.alias.name"] = "empty"
			} else {
				interfaceValue["interface.alias.name"] = resultArray[i+4]
			}
			interfaces = append(interfaces, interfaceValue)
		}

		credentials["interfaces"] = interfaces

		credentials["status"] = "success"

	}
}

func ConvertData(pdu g.SnmpPDU) interface{} {
	if pdu.Value == " " {
		return pdu.Value
	}
	switch pdu.Type {
	case g.IPAddress:
		return pdu.Value
	case g.Integer:
		return pdu.Value
	case g.OctetString:
		return strings.Trim(string(pdu.Value.([]byte)), "\"")
	default:
		return pdu.Value
	}

}
