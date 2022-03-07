package api_client

type TargetsItArmy struct {
	Online  []string `json:"online"`
	Offline []string `json:"offline"`
}

type Targets struct {
	LastUpdateMillis int64      `json:"lastUpdateMillis"`
	Statuses         []Statuses `json:"statuses"`
}
type Data struct {
	Ch []string `json:"CH"`
	Ru []string `json:"RU"`
	Au []string `json:"AU"`
	It []string `json:"IT"`
	At []string `json:"AT"`
	Lt []string `json:"LT"`
	Pt []string `json:"PT"`
	Ir []string `json:"IR"`
	Tr []string `json:"TR"`
	Fr []string `json:"FR"`
	Ca []string `json:"CA"`
	De []string `json:"DE"`
	Kz []string `json:"KZ"`
	Md []string `json:"MD"`
	Nl []string `json:"NL"`
	Us []string `json:"US"`
	Ua []string `json:"UA"`
	Hk []string `json:"HK"`
}
type Ping struct {
	Data             Data  `json:"data"`
	LastUpdateMillis int64 `json:"lastUpdateMillis"`
}
type Ports struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
	Reason   string `json:"reason"`
	State    string `json:"state"`
	Service  string `json:"service"`
}

type Info struct {
	LastModified     string      `json:"last_modified"`
	NetworkRange     string      `json:"network_range"`
	MonthlyPageviews interface{} `json:"monthly_pageviews"`
	Server           string      `json:"server"`
	Route            string      `json:"route"`
	NameServer       string      `json:"name_server"`
	HostedBy         string      `json:"hosted_by"`
	Onip             []string    `json:"onip"`
	PaidTill         string      `json:"paid_till"`
}

type Statuses struct {
	Ping       Ping     `json:"ping"`
	IP         string   `json:"ip"`
	Ips        []string `json:"ips"`
	Ports      []Ports  `json:"ports"`
	Priority   bool     `json:"priority"`
	Status     string   `json:"status"`
	Info       Info     `json:"info,omitempty"`
	URL        string   `json:"url"`
	Cloudflare bool     `json:"cloudflare?"`
}
