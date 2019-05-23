package fxoss

import (
	"bytes"
	"encoding/json"
	"strings"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/super1-chen/fxoss/utils"
	"github.com/super1-chen/fxoss/logger"
)

// OSS server struct

type OSS struct {
	User, Password, Host, SSHUser, SSHPassword string
	HTTPClient                                 *http.Client
	logger                                     *log.Logger
	config                                     config
}

var (
	userKey      = "FXOSS_USER"
	hostKey      = "FXOSS_HOST"
	pwdKey       = "FXOSS_PWD"
	sshUserKey   = "FXOSS_SSH_USER"
	sshPwdKey    = "FXOSS_SSH_PWD"
	confDirKey   = "FXOSS_DIR"
	defaultRetry = 3
	tokenJson    = "token.json"
)

// NewOssServer create a new oss server for command line tools
func NewOssServer(now time.Time, verbose bool, config config) (*OSS, error) {

	user, err := envLookUp(userKey)
	if err != nil {
		return nil, err
	}

	host, err := envLookUp(hostKey)
	if err != nil {
		return nil, err
	}

	pwd, err := envLookUp(pwdKey)
	if err != nil {
		return nil, err
	}

	sshUser, err := envLookUp(sshUserKey)
	if err != nil {
		return nil, err
	}

	sshPwd, err := envLookUp(sshPwdKey)
	if err != nil {
		return nil, err
	}
	confDir := os.Getenv(confDirKey)

	if confDir == "" {
		confDir = "/tmp"
	}
	tokenPath := path.Join(confDir, tokenJson)
	o := &OSS{
		User:        user,
		Host:        host,
		Password:    pwd,
		SSHUser:     sshUser,
		SSHPassword: sshPwd,
		HTTPClient:  &http.Client{Timeout: time.Minute},
		logger:      logger.Mylogger(verbose),
		config:      config,
	}

	if _, err = os.Stat(tokenPath); os.IsNotExist(err) {
		o.logger.Printf("update token")
		if err = o.updateToken(confDir); err != nil {
			return nil, err
		}
	}

	o.logger.Printf("config %+v", config)
	if err != nil {
		fmt.Printf("%v", err)
		return nil, err
	}
	o.logger.Printf("%s", now)

	if !o.config.IsValid(o.Host, now) {
		o.logger.Printf("config is invalid update...")
		o.updateToken(confDir)

	} else {
		o.logger.Printf("config is valid skip update...")
	}

	return o, nil

}

// ShowCDSList shows all cds list info
func (oss *OSS) ShowCDSList(option string, long bool) error {
	// api doc: https://doc.fxdata.cn/jenkins/cloud/doc-api/build/#list-cds75

	api := "/v1/cds"
	errorMsg := "get cds list from api failed"
	successMsg := "get cds list from api successfully"
	data := new(cdsList)
	var cdsList []*cdsInfo

	var headers []string
	var content [][]string

	b, err := oss.get(api)

	if err != nil {
		utils.ErrorPrintln(errorMsg, false)
		return fmt.Errorf("%s, %v", errorMsg, err)
	}

	if err = json.Unmarshal(b, &data); err != nil {
		oss.logger.Printf("decode list failed %v", err)
		utils.ErrorPrintln("decode cds list failed", false)
		return fmt.Errorf("decode cds list failed, %v", err)
	}

	utils.SuccessPrintln(successMsg)
	if len(data.CDS) == 0 {
		oss.logger.Printf("decode list failed %v", err)
		utils.ColorPrintln("CDS list is empty", utils.Yellow)
		return nil
	}

	if option == "" {
		cdsList = data.CDS
	} else {
		for _, cds := range data.CDS {
			if strings.Contains(cds.SN, option) || strings.Contains(cds.Company, option) {
				cdsList = append(cdsList, cds)
			}
		}
	}
	if len(cdsList) == 0 {
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

		for index, cds := range cdsList {
			index ++
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

		for index, cds := range cdsList {
			index ++
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

// ShowCDSDetail show all cds detail information
func (oss *OSS) ShowCDSDetail(sn string) error {

	var cdsContent [][]string
	var nodeContent [][]string
	cdsHeaders := []string{
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
		oss.logger.Printf("cds nodes is empty return.")
		utils.ColorPrintln(fmt.Sprintf("Nodes list of CDS %q is empty", sn), utils.Yellow)
		return nil
	}

	utils.SuccessPrintln(fmt.Sprintf("CDS %q Nodes list", sn))

	nodeHeaders := []string{"#", "sn", "type", "status", "hit_user(max)", "cache_kbps(max)", "service_kbps(max)"}
	for index, node := range cds.Nodes {
		index ++
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
func (oss *OSS) LoginCDS(sn string, retry int, f bool) error {
	// here use white-box test method.
	if retry == 0 {
		retry = defaultRetry
	}
	company, host, port, err := oss.getSSHInfo(sn, f)
	oss.logger.Printf("%s %s:%d", company, host, port)
	if err != nil {
		return fmt.Errorf("get ssh host port info failed: %v", err)
	}

	c, err := oss.sshClient(host, port, retry)
	if err != nil {
		return err
	}

	defer c.Close()

	err = utils.RunTerminal(company, c)
	if err != nil {
		return err
	}

	utils.SuccessPrintln(fmt.Sprintf("*** LOG OUT FROM %s ******", company))
	return nil
}

// ShowCDSPort shows cds port information by specified sn
func (oss *OSS) ShowCDSPort(sn string) error {

	var company, sshHost, sshPort string
	port, err := oss.getCDSPort(sn)

	if err != nil {
		return err
	}

	sshHost = port.SSHHost
	sshPort = strconv.Itoa(int(port.SshPort))

	detail, err := oss.getCDSDetail(sn)
	if err == nil {
		company = detail.CDS.Company
	}

	headers := []string{"company", "ssh_host", "ssh_port"}
	content := [][]string{{company, sshHost, sshPort}}
	utils.PrintTable(headers, content)

	return nil
}

// getCDSPort gets cds port info
func (oss *OSS) getCDSPort(sn string) (*portInfo, error) {

	api := fmt.Sprintf("/v1/icaches/%s/ports", sn)
	port := new(portInfo)

	b, err := oss.get(api)
	if err != nil {
		return nil, fmt.Errorf("get cds port info failed, %v", err)
	}

	if err = json.Unmarshal(b, &port); err != nil {
		err = fmt.Errorf("unmarshal cds port info failed, %v", err)
	}
	return port, nil
}

// getCDSDetail gets CDS detail infomation
func (oss *OSS) getCDSDetail(sn string) (detail *cdsDetail, err error) {

	api := fmt.Sprintf("/v1/cds/%s", sn)
	errorMsg := fmt.Sprintf("GET cds detail information with %q failed", sn)
	successMsg := fmt.Sprintf("GET cds detail information with %q success", sn)

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

func (oss *OSS) getSSHInfo(sn string, f bool) (company, sshHost string, sshPort int, err error) {
	detail, err := oss.getCDSDetail(sn)
	if err == nil {
		company = detail.CDS.Company
	}
	if f {
		sshHost = "OSS.fxdata.cn"
		port, err := utils.SN2Port(sn)
		if err != nil {
			return "", "", 0, err
		}
		p, err := strconv.Atoi(port)

		if err != nil {
			return "", "", 0, err
		}
		sshPort = p

	} else {
		port, err := oss.getCDSPort(sn)
		if err != nil {
			return "", "", 0, err
		}
		sshHost = port.SSHHost
		sshPort = int(port.SshPort)
	}

	return company, sshHost, sshPort, nil
}

// gets information by the given api
func (oss *OSS) get(api string) ([]byte, error) {

	url := fmt.Sprintf("%s%s", oss.Host, api)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, fmt.Errorf("create new get request %s failed %v", api, err)
	}

	req.Header.Set("X-auth-token", oss.config.GetToken())
	req.Header.Set("Content-Type", "application/json")

	resp, err := oss.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("reqest get %s failed %v", api, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request get %s status %s", api, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)
}

func (oss *OSS) post(api string, body io.Reader, needToken bool) ([]byte, error) {
	url := fmt.Sprintf("%s%s", oss.Host, api)
	fmt.Println(url)
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, fmt.Errorf("create new post request %s failed %v", api, err)
	}

	if needToken && oss.config != nil {
		req.Header.Set("X-auth-token", oss.config.GetToken())
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := oss.HTTPClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("request post %s failed %v", api, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request post %s status %s", api, resp.Status)
	}
	return ioutil.ReadAll(resp.Body)

}

// updateToken will download a new token from OSS and make a new config of itself
func (oss *OSS) updateToken(folderName string) error {

	b, err := oss.getNewToken()
	if err != nil {
		return fmt.Errorf("get new token failed %s", err)
	}

	r := bytes.NewReader(b)

	if err = oss.config.Update(r); err != nil {
		log.Printf("update config failed %v", err)
		return err
	}

	if err = oss.config.SetHost(oss.Host); err != nil {
		log.Printf("set host failed %v", err)
	}


	if err = oss.config.Save(folderName, tokenJson); err != nil {
		log.Printf("save token failed %s ", path.Join(folderName, tokenJson))
	}

	return nil
}

func (oss *OSS) getNewToken() ([]byte, error) {
	api := "/v1/auth/tokens"

	data := map[string]string{"username": oss.User, "password": oss.Password}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	body := bytes.NewReader(b)

	return oss.post(api, body, false)
}

func (oss *OSS) sshClient(host string, port, retry int) (*ssh.Client, error) {

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
		User: oss.SSHUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(oss.SSHPassword),
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
	value, ok := os.LookupEnv(key);
	if !ok {
		// todo add logger in here
		return value, fmt.Errorf(errPrefix, key)
	}
	return value, nil
}
