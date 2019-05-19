package fxoss

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/super1-chen/fxoss/conf"
	"github.com/super1-chen/fxoss/utils"
)

type oss struct {
	User, Password, Host, SshUser, SshPassword string
	Conf                                       *conf.Config
	HttpClient                                 *http.Client
}

var (
	UserKey    = "FXOSS_USER"
	HostKey    = "FXOSS_HOST"
	PWDKey     = "FXOSS_PWD"
	SSHUserKey = "FXOSS_SSH_USER"
	SSHPWDKey  = "FXOSS_SSH_PWD"
	ConfDirKey = "FXOSS_DIR"
	Retry      = 3
)

func New(now time.Time) (oss *oss, err error) {

	user, err := envLookUp(UserKey)
	if err != nil {
		return
	}

	host, err := envLookUp(HostKey)
	if err != nil {
		return
	}

	pwd, err := envLookUp(PWDKey)
	if err != nil {
		return
	}

	sshUser, err := envLookUp(SSHUserKey)
	if err != nil {
		return
	}

	sshPwd, err := envLookUp(SSHPWDKey)
	if err != nil {
		return
	}
	confDir := os.Getenv(ConfDirKey)

	if confDir == "" {
		confDir = "/tmp"
	}

	tokenPath := path.Join(confDir, conf.TokenJson)

	oss.User = user
	oss.Host = host
	oss.Password = pwd
	oss.SshUser = sshUser
	oss.SshPassword = sshPwd
	oss.HttpClient = &http.Client{Timeout: time.Minute} // set default client timeout 60s

	if _, err := os.Stat(tokenPath); os.IsNotExist(err) {
		if err = oss.UpdateToken(tokenPath); err != nil {
			return nil, err
		}
	}

	conf, err := conf.NewConfig(tokenPath)

	if err != nil {
		return nil, err
	}

	if !conf.IsValid(oss.Host, now) {
		oss.UpdateToken(tokenPath)

	} else {
		oss.Conf = conf
	}

	return

}

// ShowCDSList shows all cds list info
func (oss *oss) ShowCDSList(option string, long bool) error {
	// api doc: https://doc.fxdata.cn/jenkins/cloud/doc-api/build/#list-cds75

	api := "/v1/cds"
	errorMsg := "get cds list from api failed"
	successMsg := "get cds list from api successfully"
	data := new(CDSList)

	var headers []string
	var content [][]string

	b, err := oss.get(api)

	if err != nil {
		utils.ErrorPrintln(errorMsg, false)
		return fmt.Errorf("%s, %v", errorMsg, err)
	}

	if err = json.Unmarshal(b, &data); err != nil {
		utils.ErrorPrintln("decode cds list failed", false)
		return fmt.Errorf("decode cds list failed, %v", err)
	}

	utils.SuccessPrintln(successMsg)
	if len(data.CDS) == 0 {
		utils.ColorPrintln("CDS list is empty", utils.Yellow)
		return nil
	}
	if long {
		headers = []string{
			"#",
			"company",
			"sn",
			"status",
			"license_start",
			"license_end",
			"online_user(max)",
			"hit_user(max)",
			"service_kbps(max)",
			"cache_kbps(max)",
			"monitor_kbps(max)",
			"version",
			"updated_at",
		}

		for index, cds := range data.CDS {
			index += 1
			cds.OnlineUserStr = utils.FormatItem(cds.OnlineUser, cds.OnlineUserMax)
			cds.HitUserStr = utils.FormatItem(cds.HitUser, cds.HitUserMax)
			cds.ServiceStr = utils.FormatItem(cds.ServiceKbps, cds.ServiceKbpsMax)
			cds.CacheStr = utils.FormatItem(cds.CacheKbps, cds.CacheKbpsMax)
			cds.MonitorStr = utils.FormatItem(cds.MonitorKbps, cds.MonitorKbpsMax)
			content = append(content,
				[]string{
					strconv.Itoa(index),
					cds.Company,
					cds.SN,
					cds.Status,
					cds.LicenseStartAt,
					cds.LicenseEndAt,
					cds.OnlineUserStr,
					cds.HitUserStr,
					cds.ServiceStr,
					cds.CacheStr,
					cds.MonitorStr,
					cds.Version,
					cds.UpdatedAt,
				})
		}

	} else {
		headers = []string{"#", "company", "sn", "status", "version", "update_at"}

		for index, cds := range data.CDS {
			index += 1
			content = append(content, []string{
				strconv.Itoa(index),
				cds.Company,
				cds.SN,
				cds.Status,
				cds.Version,
				cds.UpdatedAt,
			})

		}
	}
	utils.PrintTable(headers, content)
	return nil
}

// ShowCDSDetail show all
func (oss *oss) ShowCDSDetail(sn string) error {

	var cdsContent [][]string
	var nodeContent [][]string
	cdsHeaders :=  []string{
		"company", "sn", "status", "license_start",
		"license_end", "online_user(max)", "hit_user(max)",
		"service_kbps(max)", "cache_kbps(max)",
		"monitor_kbps(max)", "version", "updated_at"}

	data, err := oss.getCDSDetail(sn)
	if err != nil {
		return err
	}

	// data is empty
	if data.CDS.SN == "" {
		utils.ColorPrintln(fmt.Sprintf("CDS information is empty with sn: '%q'", sn), utils.Yellow)
		return nil
	}

	cds := data.CDS
	cds.OnlineUserStr = utils.FormatItem(cds.OnlineUser, cds.OnlineUserMax)
	cds.HitUserStr = utils.FormatItem(cds.HitUser, cds.HitUserMax)
	cds.ServiceStr = utils.FormatItem(cds.ServiceKbps, cds.ServiceKbpsMax)
	cds.CacheStr = utils.FormatItem(cds.CacheKbps, cds.CacheKbpsMax)
	cds.MonitorStr = utils.FormatItem(cds.MonitorKbps, cds.MonitorKbpsMax)

	cdsContent = append(cdsContent,
		[]string{
			cds.Company,
			cds.SN,
			cds.Status,
			cds.LicenseStartAt,
			cds.LicenseEndAt,
			cds.OnlineUserStr,
			cds.HitUserStr,
			cds.ServiceStr,
			cds.CacheStr,
			cds.MonitorStr,
			cds.Version,
			cds.UpdatedAt,
		})
	utils.PrintTable(cdsHeaders, cdsContent)

	if len(cds.Nodes) == 0 {
		return nil
	}
	nodeHeaders := []string{"#", "sn", "type", "status", "hit_user(max)", "cache_kbps(max)", "service_kbps(max)"}
	for index, node := range cds.Nodes {
		index += 1
		nodeContent = append(nodeContent, []string{
			strconv.Itoa(index),
			node.SN,
			node.Type,
			node.Status,
			utils.FormatItem(node.HitUser, node.HitUserMax),
			utils.FormatItem(node.CacheKbps, node.CacheKbpsMax),
			utils.FormatItem(node.ServiceKbps, node.ServiceKbpsMax),
		})
	}

	utils.PrintTable(nodeHeaders, nodeContent)

	return nil
}

// LoginCDS uses ssh to login CDS server via ssh-tunnel or frpc-tunnel
func (oss *oss) LoginCDS(sn string, retry int, f bool) error {
	// here use white-box test method.
	if retry == 0 {
		retry = Retry
	}
	company, host, port, err := oss.getSSHInfo(sn, f)

	if err != nil {
		return fmt.Errorf("get ssh host port info failed: %v", err)
	}

	c, err := oss.sshClient(host, port, Retry)
	if err != nil {
		return err
	}

	defer c.Close()

	err = utils.RunTerminal(c)
	if err != nil {
		return err
	}

	utils.SuccessPrintln(fmt.Sprintf("*** LOG OUT FROM %s ******", company))
	return nil
}

// ShowCDSPort shows cds port information by specified sn
func (oss *oss) ShowCDSPort(sn string) error {

	var company, sshHost, sshPort string
	port, err := oss.getCDSPort(sn)

	if err != nil {
		return err
	}

	sshHost = port.SshHost
	sshPort = strconv.Itoa(int(port.SshPort))

	detail, err := oss.getCDSDetail(sn)
	if err == nil {
		company = detail.CDS.Company
	}

	headers := []string{"company", "ssh_host", "ssh_token"}
	content :=[][]string{{company, sshHost, sshPort}}
	utils.PrintTable(headers, content)

	return nil
}

// gets cds port info
//
//{
//"http_port": 57045,
//"http_url": "http://rhelp.fxdata.cn:57045",\n
//"https_port": 34755,\n
//"https_url": "https://rhelp.fxdata.cn:34755",\n
//"ssh_host": "rhelp.fxdata.cn",\n
//"ssh_port": 36079\n
//}
func (oss *oss) getCDSPort(sn string)(*PortInfo, error) {
	api := fmt.Sprintf("/v1/cds/%s", sn)
	port := new(PortInfo)

	b, err := oss.get(api)
	if err != nil {
		return nil, fmt.Errorf("get cds port info failed, %v", err)
	}

	if err = json.Unmarshal(b, &port); err != nil {
		err = fmt.Errorf("unmarshal cds port info failed, %v", err)
	}
	return port, nil
}

//gets CDS detail infomation
func (oss *oss) getCDSDetail(sn string)(detail *CdsDetail, err error){

	api := fmt.Sprintf("/v1/cds/%s", sn)
	errorMsg := fmt.Sprintf("GET cds detail information with %s failed", sn)
	successMsg := fmt.Sprintf("GET cds detail information with %s success", sn)

	b, err := oss.get(api)

	if err != nil {
		utils.ErrorPrintln(errorMsg, false)
		return
	}

	if err = json.Unmarshal(b, &detail); err != nil {
		utils.ErrorPrintln("decode cds list failed", false)
		return nil, fmt.Errorf("decode cds list failed, %v", err)
	}

	utils.SuccessPrintln(successMsg)
	return
}

func (oss *oss) getSSHInfo(sn string, f bool)(company, sshHost string, sshPort int, err error) {
	detail, err := oss.getCDSDetail(sn)
	if err == nil {
		company = detail.CDS.Company
	}
	if f {
		sshHost = "oss.fxdata.cn"
		port, err := utils.SN2Port(sn)
		if err != nil {
			return
		}
		p, err := strconv.Atoi(port)

		if err != nil {
			return
		}
		sshPort = p

	} else {
		port, err := oss.getCDSPort(sn)
		if err != nil {
			return
		}
		sshHost = port.SshHost
		sshPort = int(port.SshPort)
		return
	}
	return
}

// gets information by the given api
func (oss *oss) get(api string) ([]byte, error) {

	url := fmt.Sprintf("%s%s", oss.Host, api)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("create new get request %s failed %v", api, err)
	}

	req.Header.Set("Content-Type", oss.Conf.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := oss.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("reqest get %s failed %v", api, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request get %s status %s", api, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

func (oss *oss) post(api string, body io.Reader, needToken bool) ([]byte, error) {
	url := fmt.Sprintf("%s%s", oss.Host, api)
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, fmt.Errorf("create new post request %s failed %v", api, err)
	}

	if needToken && oss.Conf != nil {
		req.Header.Set("Content-Type", oss.Conf.Token)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oss.HttpClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("reqest post %s failed %v", api, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request post %s status %s", api, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)

}

// UpdateToken will download a new token from oss and make a new config of itself
func (oss *oss) UpdateToken(folderName string) error {
	c := new(conf.Config)
	b, err := oss.GetNewToken()
	if err != nil {
		return fmt.Errorf("get new token failed %s", err)
	}
	r := bytes.NewReader(b)
	c.Update(r)
	c.Host = oss.Host
	if err := c.Save(folderName, conf.TokenJson); err != nil {
		log.Printf("save token failed %s ", path.Join(folderName, conf.TokenJson))
	}
	oss.Conf = c
	return nil
}

func (oss *oss) GetNewToken() ([]byte, error) {
	api := "/v1/auth/tokens"

	data := map[string]string{"username": oss.User, "password": oss.Password,}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(b)

	return oss.post(api, body, false)
}

func (oss *oss) sshClient(host string, port, retry int)(*ssh.Client, error) {


	Cb := func(user, instruction string, questions []string, echos []bool) (answers []string, err error) {
		answers = make([]string, len(questions))
		for i, question := range questions {
			fmt.Fprintf(os.Stdout, "%s\r\n", question)
			bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return nil, err
			}
			password := string(bytePassword)
			answers[i] = password
		}
		return answers, nil
	}
	sshConfig := &ssh.ClientConfig{
		User: oss.SshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(oss.SshPassword),
			ssh.RetryableAuthMethod(ssh.KeyboardInteractive(Cb), retry),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Minute,
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)

	if err != nil {
		return nil, fmt.Errorf("ssh dail: connection failed %s", err)
	}
	return client, nil

}

func envLookUp(key string) (value string, err error) {
	errPrefix := "missing env: %s"
	if value, ok := os.LookupEnv(key); !ok {
		// todo add logger in here
		return value, fmt.Errorf(errPrefix, key)
	}
	return value, nil
}
