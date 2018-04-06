package utils

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Instance struct {
	IpAddress struct {
		InnerXML string `xml:",innerxml"`
	} `xml:"instance>ipAddr"`
}

func GetIpAddress() string {
	eureka := os.Getenv("EUREKA_SERVER")
	if eureka == "" {
		eureka = "discovery-service"
	}
	url := fmt.Sprintf("http://%s:8761/eureka/apps/gmapsadapter", eureka)

	var maxAttempts int = 5
	retryGET(maxAttempts, url)

	request, _ := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	response, responseErr := client.Do(request)
	if responseErr != nil {
		panic(responseErr)
	}
	log.Println(response)
	log.Printf("Reading XML body...")
	body, _ := ioutil.ReadAll(response.Body)

	var instance Instance
	unmarshallError := xml.Unmarshal(body, &instance)
	if unmarshallError != nil {
		panic(unmarshallError)
	}

	log.Printf("Received Ip Address!")

	return instance.IpAddress.InnerXML
}

//Retry GET requests to specified url according to maxAttempts
func retryGET(maxAttempts int, url string) {
	log.Printf("Attempting to connect to " + url)
	var response *http.Response
	for i := 0; i < maxAttempts; i++ {
		request, _ := http.NewRequest("GET", url, nil)
		client := &http.Client{}
		response, _ = client.Do(request)
		if response.Status == "200 OK" {
			log.Printf("Success response returned, continuing")
			break
		}
	}
	if response.Status != "200 OK" {
		log.Printf("Response was never successful, please increase attempts")
	}
}
