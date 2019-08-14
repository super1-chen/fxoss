package app

type cdsInfo struct {
	SN             string `json:"sn"`
	Company        string `json:"company"`
	Status         string `json:"status"`
	LicenseStartAt string `json:"license_start_at"`
	LicenseEndAt   string `json:"license_end_at"`
	OnlineUser     int64  `json:"online_user"`
	OnlineUserMax  int64  `json:"online_user_max"`
	OnlineUserStr  string
	HitUser        int64  `json:"hit_user"`
	HitUserMax     int64  `json:"hit_user_max"`
	HitUserStr     string `json:"hit_user_str"`
	ServiceKbps    int64  `json:"service_kbps"`
	ServiceKbpsMax int64  `json:"service_kbps_max"`
	ServiceStr     string
	CacheKbps      int64 `json:"cache_kbps"`
	CacheKbpsMax   int64 `json:"cache_kbps_max"`
	CacheStr       string
	MonitorKbps    int64 `json:"monitor_kbps"`
	MonitorKbpsMax int64 `json:"monitor_kbps_max"`
	MonitorStr     string
	Version        string  `json:"version"`
	UpdatedAt      string  `json:"updated_at"`
	Nodes          []*node `json:"nodes"`
}

type labelCdsInfo struct {
	label
	cdsList []*cdsInfo
}

type node struct {
	SN             string `json:"sn"`
	Type           string `json:"type"`
	Status         string `json:"status"`
	HitUser        int64  `json:"hit_user"`
	HitUserMax     int64  `json:"hit_user_max"`
	ServiceKbps    int64  `json:"service_kbps"`
	ServiceKbpsMax int64  `json:"service_kbps_max"`
	CacheKbps      int64  `json:"cache_kbps"`
	CacheKbpsMax   int64  `json:"cache_kbps_max"`
}

type cdsList struct {
	CDS []*cdsInfo `json:"cds"`
}

type cdsDetail struct {
	CDS *cdsInfo `json:"cds"`
}

type label struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Count   int64  `json:"count"`
	CDSList []*cdsInfo
}

type labels struct {
	Labels []*label `json:"labels"`
}

type portInfo struct {
	SSHHost   string `json:"ssh_host"`
	SSHPort   int64  `json:"ssh_port"`
	HTTPUrl   string `json:"http_url"`
	HTTPPort  int64  `json:"http_port"`
	HTTPSURL  string `json:"https_url"`
	HTTPSPort int64  `json:"https_port"`
}

type emailConf struct {
	Address    string `json:"address"`
	Password   string `json:"password"`
	SMTPServer string `json:"smtp_server"`
}

type disk struct {
	Await  string `json:"await"`
	Name   string `json:"name"`
	RS     string `json:"rs"`
	Size   string `json:"size"`
	Status int64  `json:"status"`
	Used   string `json:"usesd"`
	Util   string `json:"util"`
	WS     string `json:"ws"`
}

type disks struct {
	Disk []*disk `json:"disks"`
}

type cdsDiskInfo struct {
	cdsInfo
	disks      []*disk
	Await      string `json:"await"`
	Name       string `json:"name"`
	ReadSpeed  string `json:"rs"`
	Size       string `json:"size"`
	DiskUsed   string `json: "used"`
	Util       string `json:"util"`
	WriteSpeed string `json: "ws"`
}

type diskTypeResult struct {
	domain, sn, company, status, userAndSpeed string
	user, speed, diskType                     int64
}
