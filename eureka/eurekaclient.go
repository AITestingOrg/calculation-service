package eureka

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func PostToEureka() {
	var localIpAddr string = GetLocalIpAddress()
	jsonRequest := RequestBody{
		Instance{
			HostName:         localIpAddr,
			App:              "calculationservice",
			IpAddr:           localIpAddr,
			VipAddress:       "calculationservice",
			SecureVipAddress: "calculationservice",
			Status:           "UP",
			Port:             Port{8080, true},
			HomePageUrl:      "http://" + localIpAddr + ":8080",
			StatusPageUrl:    "http://" + localIpAddr + ":8080/api/v1/status",
			HealthCheckUrl:   "http://" + localIpAddr + ":8080/api/v1/status",
			DataCenterInfo:   DataCenter{"com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo", "MyOwn"},
			Metadata:         MetaData{InstanceId: ""},
		},
	}

	jsonParsed, err := json.Marshal(jsonRequest)
	if err != nil {
		log.Printf("ERROR parsing JSON from struct")
		panic(err)
	}

	log.Printf("Registering with Eureka...")
	eureka := os.Getenv("EUREKA_SERVER")
	if eureka == "" {
		eureka = "discoveryservice"
	}
	url := fmt.Sprintf("http://%s:8761/eureka/apps/calculationservice", eureka)
	json := []byte(jsonParsed)

	request, _ := http.NewRequest("POST", url, bytes.NewBuffer(json))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		panic(responseErr)
	}
	if response.StatusCode == 204 {
		log.Printf("Registered with Eureka!")
	} else {
		log.Printf("Response Status: %s", response.Status)
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		log.Printf("Response Body: %s", string(body))
		errorMessage := fmt.Sprintf("ERROR: Recieved response with status code: \"%d\"  but expected response with status 204. \nResponse Body: %s", response.StatusCode, string(body))
		panic(errorMessage)
	}
}

func CheckEurekaService(eurekaUp bool) bool {
	duration := time.Duration(15) * time.Second
	time.Sleep(duration)

	eureka := os.Getenv("EUREKA_SERVER")
	if eureka == "" {
		eureka = "discoveryservice"
	}
	url := fmt.Sprintf("http://%s:8761/eureka/", eureka)

	request, _ := http.NewRequest("GET", url, nil)
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		log.Printf("Response error")
		return false
	}
	if response.Status != "204 No Content" {
		log.Printf("Success, Eureka was found")
		return true
	}
	return false
}

func StartHeartbeat() {
	log.Printf("Initializing heartbeat for every 30 seconds")
	ticker := time.NewTicker(time.Second * 30)
	go func() {
		for range ticker.C {
			log.Printf("Sending heartbeat to Eureka...")
			eureka := os.Getenv("EUREKA_SERVER")
			if eureka == "" {
				eureka = "discoveryservice"
			}
			url := fmt.Sprintf("http://%s:8761/eureka/apps/calculationservice/%s", eureka, GetLocalIpAddress())
			request, _ := http.NewRequest("PUT", url, nil)
			request.Header.Add("Content-Type", "application/json")
			client := &http.Client{}
			client.Do(request)
			log.Printf("Heartbeat sent to Eureka")
		}
	}()
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
