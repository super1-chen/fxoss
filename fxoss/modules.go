package fxoss


type CDSInfo struct {
	SN string `json:"sn"`
	Company string `json:"company"`
	Status string `json:"status"`
	LicenseStartAt string `json:"license_start_at"`
	LicenseEndAt string `json:"license_end_at"`
	OnlineUser int64 `json:"online_user"`
	OnlineUserMax int64 `json:"online_user_max"`
	OnlineUserStr string
	HitUser int64 `json:"hit_user"`
	HitUserMax int64 `json:"hit_user_max"`
	HitUserStr string `json:"hit_user_str"`
	ServiceKbps int64 `json:"service_kbps"`
	ServiceKbpsMax int64 `json:"service_kbps_max"`
	ServiceStr string
	CacheKbps int64 `json:"cache_kbps"`
	CacheKbpsMax int64 `json:"cache_kbps_max"`
	CacheStr string
	MonitorKbps int64 `json:"monitor_kbps"`
	MonitorKbpsMax int64 `json:"monitor_kbps_max"`
	MonitorStr string
	Version string `json:"version"`
	UpdatedAt string `json:"updated_at"`
	Nodes []*Node `json:"nodes"`
}

type Node struct {
	SN string `json:"sn"`
	Type string `json:"type"`
	Status string `json:"status"`
	HitUser int64 `json:"hit_user"`
	HitUserMax int64 `json:"hit_user_max"`
	ServiceKbps int64 `json:"service_kbps"`
	ServiceKbpsMax int64 `json:"service_kbps_max"`
	CacheKbps int64 `json:"cache_kbps"`
	CacheKbpsMax int64 `json:"cache_kbps_max"`
}

type CDSList struct {
	CDS []*CDSInfo `json:"cds"`
}

type CdsDetail struct {
	CDS *CDSInfo `json:"cds"`
}

type PortInfo struct {
	SshHost string `json:"ssh_host"`
	SshPort int64 `json:"ssh_port"`
	HttpUrl string `json:"http_url"`
	HttpPort int64 `json:"http_port"`
	HttpsUrl string `json:"https_url"`
	HttpsPort int64 `json:"https_port"`
}