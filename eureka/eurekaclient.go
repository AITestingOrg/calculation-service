package eureka

import (
	"net/http"
	"bytes"
	"log"
	"net"
	"encoding/json"
)

func PostToEureka() {
	var localIpAddr string = GetLocalIpAddress()
	jsonRequest := RequestBody {
		Instance {
			HostName: "localhost",
			App: "calculationservice",
			IpAddr: localIpAddr,
			VipAddress: "calculationservice",
			SecureVipAddress: "calculationservice",
			Status: "UP",
			Port: 8000,
			StatusPageUrl: "http://" + localIpAddr + ":8000/api/v1/status",
			DataCenterInfo: DataCenter { Name: "MyOwn" },
			Metadata: MetaData { InstanceId: "" },
		},
	}

	jsonParsed, err := json.Marshal(jsonRequest)
	if err != nil {
		log.Printf("ERROR parsing JSON from struct")
		panic(err)
	}

	log.Printf("Registering with Eureka...")
	url := "http://discovery-service:8761/eureka/apps/calculationservice"
	json := []byte(jsonParsed)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(json))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		panic(responseErr)
	}
	if response.Status != "204 No Content" {
		log.Printf(response.Status)
		panic("ERROR: 204 Response Not Returned")
	}
	log.Printf("Registered with Eureka!")
}

func GetLocalIpAddress() string {
	addrs, err := net.InterfaceAddrs()
    if err != nil {
        panic(err)
    }
    for _, address := range addrs {
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
			}
        }
    }
    return "localhost"
}