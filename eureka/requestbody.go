package eureka

type RequestBody struct {
	Instance Instance `json:"instance"`
}

type Instance struct {
	HostName         string     `json:"hostName"`
	App              string     `json:"app"`
	IpAddr           string     `json:"ipAddr"`
	VipAddress       string     `json:"vipAddress"`
	SecureVipAddress string     `json:"secureVipAddress"`
	Status           string     `json:"status"`
	Port             Port       `json:"port"`
	SecurePort		 Port		`json:"securePort"`
	HomePageUrl      string     `json:"homePageUrl"`
	StatusPageUrl    string     `json:"statusPageUrl"`
	HealthCheckUrl   string     `json:"healthCheckUrl"`
	DataCenterInfo   DataCenter `json:"dataCenterInfo"`
	Metadata         MetaData   `json:"metadata"`
}

type Port struct {
	PortNum int `json:"$"`
	Enabled bool `json:"@enabled"`
}

type DataCenter struct {
	Class string `json:"@class"`
	Name string `json:"name"`
}

type MetaData struct {
	InstanceId string `json:"instanceId"`
}
