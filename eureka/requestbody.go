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
	Port             int        `json:"port"`
	HomePageUrl      string     `json:"homePageUrl"`
	StatusPageUrl    string     `json:"statusPageUrl"`
	HealthCheckUrl   string     `json:"healthCheckUrl"`
	DataCenterInfo   DataCenter `json:"dataCenterInfo"`
	Metadata         MetaData   `json:"metadata"`
}

type DataCenter struct {
	Name string `json:"name"`
}

type MetaData struct {
	InstanceId string `json:"instanceId"`
}
