package unblock

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	hertzConst "github.com/cloudwego/hertz/pkg/protocol/consts"
	"reflect"
	"strings"
	"time"
	containerInterface "ws/framework/application/container/abstract_interface"
	"ws/framework/lib/msisdn"
	networkOperators "ws/framework/lib/network_operators"
	unblockQueryContent "ws/framework/lib/unblock_query_content"
	"ws/framework/plugin/json"
	httpApi "ws/framework/plugin/network/http_api"
	"ws/framework/utils"
)

// UnblockSupportInfo .
type UnblockSupportInfo struct {
	DebugInfo                   string `json:"Debug info"`
	CCode                       string `json:"CCode"`
	Hash                        string `json:"Hash"`
	Version                     string `json:"Version"`
	LC                          string `json:"LC"`
	LG                          string `json:"LG"`
	Context                     string `json:"Context"`
	Carrier                     string `json:"Carrier"`
	Manufacturer                string `json:"Manufacturer"`
	Model                       string `json:"Model"`
	OS                          string `json:"OS"`
	UserAgent                   string `json:"UserAgent"`
	SocketConn                  string `json:"Socket Conn"`
	Connection                  string `json:"Connection"`
	LastVoIPCallBlocked         string `json:"Last VoIP call blocked"`
	NetworkType                 string `json:"Network Type"`
	Datacenter                  string `json:"Datacenter"`
	RadioMCCMNC                 string `json:"Radio MCC-MNC"`
	SIMMCCMNC                   string `json:"SIM MCC-MNC"`
	FreeSpaceBuiltIn            string `json:"Free Space Built-In"`
	FreeSpaceRemovable          string `json:"Free Space Removable"`
	FAQResultsReturned          string `json:"FAQ Results Returned"`
	FAQResultsRead              string `json:"FAQ Results Read"`
	CachedConnectionFAQ         string `json:"Cached Connection FAQ"`
	CachedConnectionFAQTimeRead string `json:"Cached Connection FAQ Time Read"`
	DeviceISO8601               string `json:"Device ISO8601"`
	SmbCount                    string `json:"Smb count"`
	EntCount                    string `json:"Ent count"`
	Interface                   string `json:"Interface"`
	DA                          string `json:"DA"`
	DBCorrupted                 string `json:"DB corrupted"`
	Backup                      string `json:"Backup"`
	BUError                     string `json:"BUError"`
	VideoCalls                  string `json:"Video Calls"`
	Payments                    string `json:"Payments"`
	VoiceOver                   string `json:"VoiceOver"`
	LargerText                  string `json:"Larger Text"`
	SO                          string `json:"SO"`
	DeviceID                    string `json:"DeviceID"`
	MDEnabled                   string `json:"MDEnabled"`
	HasMdCompanion              string `json:"HasMdCompanion"`
	Anid                        string `json:"anid"`
}

// file .
type file struct {
	Name    string `json:"name"`
	Content []byte `json:"content"`
}

// UnblockContentResponse .
type UnblockContentResponse struct {
	EmailTo      string `json:"emailTo"`
	EmailSubject string `json:"emailSubject"`
	EmailContent string `json:"emailContent"`
	AttachFile   file   `json:"attachFile"`
}

// Do .
func Do(appIocContainer containerInterface.IAppIocContainer) (*UnblockContentResponse, error) {
	query := unblockQueryContent.New()

	supportInfo, err := makeSupportInfo(appIocContainer, query)
	if err != nil {
		return nil, err
	}

	fileBody, err := createTarFileAndCompressGzip(appIocContainer, supportInfo)
	if err != nil {
		return nil, err
	}

	response := UnblockContentResponse{}
	response.EmailTo = "iphone@support.whatsapp.com"
	response.EmailSubject = "Question about WhatsApp for iPhone"
	response.EmailContent = fillFieldsValue(query, supportInfo)
	response.AttachFile.Name = "logs.tar.gz"
	response.AttachFile.Content = fileBody

	return &response, nil
}

// 创建tar文件再gzip压缩
func createTarFileAndCompressGzip(appIocContainer containerInterface.IAppIocContainer, supportInfo *UnblockSupportInfo) ([]byte, error) {
	device := appIocContainer.ResolveDeviceService().Context()
	configuration := appIocContainer.ResolveWhatsappConfiguration()

	tarFileBody := bytes.NewBuffer(make([]byte, 0))
	tarWriter := tar.NewWriter(tarFileBody)

	// debuginfo.json
	debugJson, _ := json.MarshalIndent(supportInfo, "", "\t")
	err := writeToTar(tarWriter, "debuginfo.json", debugJson)
	if err != nil {
		return nil, fmt.Errorf("生成文件内容失败")
	}

	// whatsapp-%s-WhatsApp-0-launch.log
	launchLogName := fmt.Sprintf("whatsapp-%s-WhatsApp-0-launch.log", time.Now().Format("2006-01-02-15-04-05-.000"))
	launchLog := bytes.NewBuffer(make([]byte, 0))
	launchLog.WriteString("\n")
	commitHash := configuration.CommitHash
	waVersionString := configuration.VersionString
	launchLog.WriteString(fmt.Sprintf("Device: %s | System: iOS %s (%s) | WhatsApp version: %s | Hash: %s |  | launchID: %s\n", device.Device, device.OsVersion, device.BuildNumber, waVersionString, commitHash, utils.GenUUID4()))
	launchLog.WriteString("\n")

	err = writeToTar(tarWriter, launchLogName, launchLog.Bytes())
	if err != nil {
		return nil, fmt.Errorf("生成文件内容失败")
	}

	_ = tarWriter.Close()

	// gzip
	gzipFile := bytes.NewBuffer(make([]byte, 0))
	gzipWriter := gzip.NewWriter(gzipFile)
	gzipWriter.OS = 3 // unix
	gzipWriter.ModTime = time.Now()

	_, _ = gzipWriter.Write(tarFileBody.Bytes())
	_ = gzipWriter.Close()

	return gzipFile.Bytes(), nil
}

// 填充supportInfo结构体数据
func fillFieldsValue(query string, supportInfo *UnblockSupportInfo) string {
	var stringBuilder strings.Builder

	stringBuilder.WriteString("\n" + query + "\n\n")
	stringBuilder.WriteString("--[[Support Info]]--\n")

	reflectClass := reflect.ValueOf(supportInfo).Elem()
	reflectClassChildType := reflectClass.Type()
	fieldsCount := reflectClassChildType.NumField()

	for i := 0; i < fieldsCount; i++ {
		field := reflectClassChildType.Field(i)
		fieldValue := reflectClass.Field(i)

		stringBuilder.WriteString(field.Tag.Get("json"))
		stringBuilder.WriteString(": ")
		stringBuilder.WriteString(fieldValue.String())
		stringBuilder.WriteString("\n")
	}

	stringBuilder.WriteString("\n\n\n\n\n\n")
	stringBuilder.WriteString("Sent from my iPhone")

	return stringBuilder.String()
}

// 写入到tar文件
func writeToTar(tarWriter *tar.Writer, name string, body []byte) error {
	header := new(tar.Header)
	header.Name = name
	header.Size = int64(len(body))
	header.Mode = 420
	header.Format = tar.FormatUnknown
	header.ModTime = time.Now()

	err := tarWriter.WriteHeader(header)
	if err != nil {
		return err
	}

	_, err = tarWriter.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func makeSupportInfo(appIocContainer containerInterface.IAppIocContainer, queryContent string) (*UnblockSupportInfo, error) {
	networkOperator, err := searchNetworkOperator(appIocContainer.ResolveJID().User)
	if err != nil {
		return nil, err
	}

	device := appIocContainer.ResolveDeviceService().Context()
	configuration := appIocContainer.ResolveWhatsappConfiguration()
	uuid := utils.GenUUID4()

	client := appIocContainer.ResolveHttpClient()

	//1.先查询网页 如果回包内容是[] 下面有个参数填0 如果回包[]里面有内容，下面的参数填3
	body, err := httpApi.Do(
		client,
		httpApi.Url("https://faq.whatsapp.com/client_search.php"),
		httpApi.Method(hertzConst.MethodGet),
		httpApi.QueryParams("query", queryContent),
		httpApi.QueryParams("lg", device.Language),
		httpApi.QueryParams("lc", device.Country),
		httpApi.QueryParams("platform", "iphone"),
		httpApi.QueryParams("anid", uuid),
		httpApi.UserAgent(appIocContainer.ResolveDeviceService().PrivateStatsAgent()),
	)

	if err != nil {
		return nil, fmt.Errorf("网络异常")
	}

	supportInfo := UnblockSupportInfo{}
	supportInfo.DebugInfo = "unregistered"
	supportInfo.CCode = device.Area
	supportInfo.Version = configuration.VersionString
	supportInfo.LC = device.Country
	supportInfo.LG = device.Language
	supportInfo.Context = "blocked +" + device.Area + device.Phone
	supportInfo.Carrier = networkOperator
	supportInfo.Manufacturer = device.Manufacturer
	supportInfo.Model = device.Device
	supportInfo.OS = device.OsVersion
	supportInfo.UserAgent = appIocContainer.ResolveDeviceService().DeviceAgent()
	supportInfo.SocketConn = "DN"
	supportInfo.Connection = "none"
	supportInfo.LastVoIPCallBlocked = "no"
	supportInfo.NetworkType = "Unknown"
	supportInfo.Datacenter = "Unknown"
	supportInfo.RadioMCCMNC = "N/A"
	supportInfo.SIMMCCMNC = device.Mcc + "-" + device.Mnc
	supportInfo.Hash = configuration.CommitHash
	mb := utils.RandInt64(0, 20000) + 10000
	k := mb*1024*1024 + utils.RandInt64(0, 1023)
	freeSpace := fmt.Sprintf("%d (%d MB)", k, mb)
	supportInfo.FreeSpaceBuiltIn = freeSpace

	if len(body) > 0 {
		supportInfo.FAQResultsReturned = "3" //根据上面的网页请求返回值来 返回空是0 ，返回有网页数据是3
	} else {
		supportInfo.FAQResultsReturned = "0" //根据上面的网页请求返回值来 返回空是0 ，返回有网页数据是3
	}

	supportInfo.FreeSpaceRemovable = "Not Present"
	supportInfo.FAQResultsRead = "n/a"
	supportInfo.CachedConnectionFAQ = "no"
	supportInfo.CachedConnectionFAQTimeRead = "n/a"
	supportInfo.DeviceISO8601 = time.Now().Format("2006-01-02 15:04:05.000-0700")
	supportInfo.SmbCount = "0"
	supportInfo.EntCount = "0"
	supportInfo.Interface = "WiFi/None"
	supportInfo.DA = "0"
	supportInfo.DBCorrupted = "false"
	supportInfo.Backup = "off"
	supportInfo.BUError = "none"
	supportInfo.VideoCalls = "false"
	supportInfo.Payments = "false"
	supportInfo.VoiceOver = "false"
	supportInfo.LargerText = "false"
	supportInfo.SO = "C"
	supportInfo.DeviceID = "0"
	supportInfo.MDEnabled = "false"
	supportInfo.HasMdCompanion = "false"
	supportInfo.Anid = uuid

	return &supportInfo, nil
}

func searchNetworkOperator(phoneNumber string) (networkOperator string, err error) {
	imsi, err := msisdn.Parse(phoneNumber, true)

	if err != nil {
		return
	}

	networkOperator = networkOperators.FindOperator(imsi.CountryName)

	if len(networkOperator) == 0 {
		err = fmt.Errorf("找不到运营商")
		return
	}

	return
}
