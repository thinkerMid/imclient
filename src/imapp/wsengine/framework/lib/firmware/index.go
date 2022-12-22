package firmware

import (
	"bytes"
	"github.com/andybalholm/brotli"
	"io/ioutil"
	"sync"
	"ws/framework/plugin/json"
	"ws/framework/utils"
)

type firmware struct {
	OSVersion   string `json:"version"`
	BuildNumber string `json:"buildid"`
}

type device struct {
	Name       string     `json:"name"`
	Identifier string     `json:"identifier"`
	Firmwares  []firmware `json:"firmwares"`
}

//func fetch() {
//	client := netpoll.HTTP(
//		networkConstant.ConnectionConfig{},
//		&tls.Config{
//			ServerName: "api.ipsw.me",
//		},
//	)
//
//	deviceList := make([]device, 0)
//
//	err := httpApi.DoAndBind(
//		client, &deviceList,
//		httpApi.Url("https://api.ipsw.me/v4/devices"),
//		httpApi.Method(hertzConst.MethodGet),
//	)
//
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	// 已收集的darwin系统支持的版本
//	var osVersionList = []string{
//		"13.0", "13.1", "13.1.1", "13.1.2", "13.1.3", "13.2", "13.2.2", "13.2.3", "13.3", "13.3.1", "13.4", "13.4.5", "13.4.6", "13.4.8",
//		"14.0", "14.0.1", "14.0.2", "14.2", "14.3", "14.4", "14.4.1", "14.4.2", "14.5", "14.5.1", "14.6", "14.7", "14.7.1", "14.8", "14.8.1",
//		"15.0", "15.0.1", "15.0.2", "15.1", "15.1.1", "15.2", "15.2.1", "15.3", "15.3.1", "15.4", "15.4.1", "15.5", "15.6",
//		"16.0",
//	}
//
//	filter := make([]device, 0)
//
//	for _, v := range deviceList {
//		if !strings.Contains(v.Identifier, "iPhone") {
//			continue
//		}
//
//		identifierCodeSplit := strings.Split(v.Identifier, ",")
//
//		// 手机产品代号 小于8的不要
//		code := strings.ReplaceAll(identifierCodeSplit[0], "iPhone", "")
//		codeNumber, _ := strconv.Atoi(code)
//		if codeNumber < 8 {
//			continue
//		}
//
//		deviceFirmwares := device{}
//
//		// 设备固件
//		err = httpApi.DoAndBind(
//			client, &deviceFirmwares,
//			httpApi.Url(fmt.Sprintf("https://api.ipsw.me/v4/device/%s?type=ipsw", v.Identifier)),
//			httpApi.Method(hertzConst.MethodGet),
//		)
//
//		if err != nil {
//			fmt.Println(err.Error())
//			return
//		}
//
//		firmwares := make([]firmware, 0)
//
//		// 固件版本过滤
//		for _, f := range deviceFirmwares.Firmwares {
//			versionSplit := strings.Split(f.OSVersion, ".")
//			version, _ := strconv.Atoi(versionSplit[0])
//
//			// 小于13不要
//			if version < 13 {
//				continue
//			}
//
//			find := false
//
//			// 不再支持的列表里不要
//			for _, osVersion := range osVersionList {
//				if f.OSVersion == osVersion {
//					find = true
//					break
//				}
//			}
//
//			if find {
//				firmwares = append(firmwares, f)
//			}
//		}
//
//		deviceFirmwares.Firmwares = firmwares
//
//		// 移除带括号的产品名字
//		phoneNameSplit := strings.Split(deviceFirmwares.Name, " (")
//		deviceFirmwares.Name = phoneNameSplit[0]
//
//		filter = append(filter, deviceFirmwares)
//
//		time.Sleep(time.Second * 1)
//	}
//
//	b, _ := json.Marshal(filter)
//	ioutil.WriteFile("./firmware.json", b, 0666)
//
//	buf := bytes.NewBuffer(make([]byte, 0))
//	w := brotli.NewWriterLevel(buf, brotli.BestCompression)
//	w.Write(b)
//	w.Flush()
//
//	fmt.Println(buf.Bytes())
//}

var deviceProfileBody = []byte{11, 33, 71, 128, 156, 7, 182, 141, 124, 235, 100, 113, 52, 48, 237, 119, 41, 178, 181, 23, 77, 246, 185, 230, 43, 219, 137, 45, 101, 247, 97, 228, 29, 235, 246, 147, 5, 84, 192, 205, 162, 211, 213, 199, 255, 255, 171, 233, 0, 22, 234, 107, 236, 71, 61, 176, 132, 86, 90, 86, 76, 35, 187, 255, 181, 253, 198, 26, 106, 161, 177, 78, 20, 104, 39, 96, 217, 182, 71, 181, 59, 40, 175, 14, 3, 191, 84, 99, 171, 117, 160, 93, 244, 42, 21, 204, 247, 175, 130, 171, 188, 181, 207, 26, 241, 220, 33, 55, 1, 215, 252, 221, 126, 62, 126, 21, 91, 111, 86, 20, 193, 14, 94, 48, 70, 205, 215, 135, 10, 203, 181, 37, 80, 166, 33, 202, 152, 57, 98, 156, 30, 16, 70, 77, 180, 244, 56, 48, 40, 220, 210, 181, 73, 163, 209, 226, 32, 65, 117, 179, 91, 204, 224, 226, 36, 13, 230, 20, 152, 62, 212, 93, 57, 185, 211, 90, 146, 224, 68, 211, 18, 144, 158, 66, 15, 25, 145, 43, 194, 142, 139, 3, 66, 110, 73, 66, 26, 188, 81, 206, 204, 65, 183, 245, 112, 226, 68, 148, 48, 96, 208, 90, 134, 20, 39, 142, 81, 29, 9, 147, 3, 6, 6, 35, 37, 204, 224, 78, 145, 233, 27, 241, 236, 246, 114, 84, 19, 139, 3, 58, 82, 200, 65, 211, 199, 68, 51, 203, 122, 228, 197, 176, 154, 162, 128, 25, 76, 221, 100, 100, 15, 192, 211, 52, 84, 226, 96, 88, 246, 148, 70, 100, 151, 145, 134, 166, 209, 43, 13, 147, 41, 54, 137, 148, 172, 49, 139, 40, 70, 184, 129, 14, 180, 245, 10, 20, 17, 234, 60, 81, 129, 82, 4, 8, 167, 66, 53, 137, 163, 24, 255, 186, 35, 136, 85, 168, 134, 64, 97, 52, 40, 129, 248, 32, 250, 32, 148, 18, 237, 35, 95, 66, 161, 122, 217, 6, 250, 48, 212, 17, 52, 144, 57, 224, 86, 90, 189, 50, 28, 132, 186, 165, 195, 199, 205, 11, 23, 84, 239, 97, 179, 109, 160, 53, 40, 199, 34, 231, 205, 237, 191, 85, 162, 0, 81, 122, 192, 69, 137, 222, 136, 199, 209, 7, 92, 238, 189, 219, 20, 244, 119, 82, 91, 18, 151, 161, 162, 14, 20, 97, 63, 223, 142, 165, 85, 40, 222, 33, 29, 6, 160, 10, 197, 71, 193, 11, 20, 109, 165, 15, 184, 220, 34, 109, 231, 230, 248, 248, 202, 147, 48, 168, 73, 219, 86, 251, 100, 192, 219, 130, 220, 210, 101, 9, 31, 244, 137, 207, 155, 206, 14, 224, 35, 168, 146, 181, 62, 247, 251, 75, 162, 85, 127, 183, 160, 156, 237, 145, 176, 131, 43, 155, 183, 5, 185, 75, 112, 191, 100, 237, 23, 173, 249, 187, 5, 45, 118, 111, 132, 126, 229, 95, 215, 29, 113, 241, 3, 246, 69, 6, 30, 42, 226, 197, 230, 6, 45, 0, 226, 2, 42, 110, 207, 10, 217, 226, 114, 107, 54, 36, 119, 158, 46, 158, 35, 55, 82, 197, 12}
var deviceList []device
var once sync.Once

func deviceData() []device {
	once.Do(func() {
		rd := brotli.NewReader(bytes.NewBuffer(deviceProfileBody))
		rb, _ := ioutil.ReadAll(rd)

		if err := json.Unmarshal(rb, &deviceList); err != nil {
			panic(err.Error())
		}

		deviceProfileBody = nil
	})

	return deviceList
}

// NewAppleFirmware .
func NewAppleFirmware() Apple {
	data := deviceData()
	deviceSize := len(data) - 1

	chooseDeviceIdx := utils.RandInt64(0, int64(deviceSize))
	device := data[chooseDeviceIdx]

	firmwareSize := len(device.Firmwares) - 1

	chooseFirmwareIdx := utils.RandInt64(0, int64(firmwareSize))
	firmware := device.Firmwares[chooseFirmwareIdx]

	return Apple{
		device.Name,
		firmware.OSVersion,
		firmware.BuildNumber,
	}
}

// Apple .
type Apple struct {
	production  string
	osVersion   string
	buildNumber string
}

// GetOSVersion .
func (a Apple) GetOSVersion() string {
	return a.osVersion
}

// GetManufacturer .
func (a Apple) GetManufacturer() string {
	return "Apple"
}

// GetProduction .
func (a Apple) GetProduction() string {
	return a.production
}

// GetBuildNumber .
func (a Apple) GetBuildNumber() string {
	return a.buildNumber
}
