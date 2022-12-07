package apple

import (
	"sync"
)

/*
iOS APP user-agent信息伪造
darwin: https://www.theiphonewiki.com/wiki/Kernel#iOS.2FiPadOS
CFNetwork版本号与darwin一一对应，手动网页搜索收集而来。
*/

// DarwinSystemManager .
type DarwinSystemManager struct{}

type darwinSystem struct {
	DarwinVersion    string
	CFNetworkVersion string
	IOSVersions      []string
}

var darwinSystemMapOnce sync.Once
var darwinSystemList []darwinSystem

// DarwinSystemManagerInstance .
func DarwinSystemManagerInstance() DarwinSystemManager {
	darwinSystemMapOnce.Do(func() {
		darwinSystemList = []darwinSystem{
			{"19.0.0", "1107.1",
				[]string{
					"13.0",
					"13.1", "13.1.1", "13.1.2", "13.1.3",
					"13.2", "13.2.2", "13.2.3",
				}},
			{"19.2.0", "1121.2.2", []string{"13.3"}},
			{"19.3.0", "1121.2.2", []string{"13.3.1"}},
			{"19.4.0", "1125.2", []string{"13.4"}},
			{"19.5.0", "1126", []string{"13.4.5", "13.4.6"}},
			{"19.6.0", "1128.0.1", []string{"13.4.8"}},
			{"20.0.0", "1197", []string{"14.0", "14.0.1", "14.0.2"}},
			{"20.1.0", "1206", []string{"14.2"}},
			{"20.2.0", "1209", []string{"14.3"}},
			{"20.3.0", "1220.1", []string{"14.4", "14.4.1", "14.4.2"}},
			{"20.4.0", "1237", []string{"14.5", "14.5.1"}},
			{"20.5.0", "1240.0.4", []string{"14.6"}},
			{"20.6.0", "1240.0.4", []string{"14.7", "14.7.1", "14.8", "14.8.1"}},
			{"21.0.0", "1312", []string{"15.0", "15.0.1", "15.0.2"}},
			{"21.1.0", "1325.0.1", []string{"15.1", "15.1.1"}},
			{"21.2.0", "1327.0.4", []string{"15.2", "15.2.1"}},
			{"21.3.0", "1329", []string{"15.3", "15.3.1"}},
			{"21.4.0", "1331.0.7", []string{"15.4", "15.4.1"}},
			{"21.5.0", "1333.0.4", []string{"15.5"}},
			{"21.6.0", "1335.0.3", []string{"15.6"}},
			{"22.0.0", "1388", []string{"16.0"}},
		}
	})

	return DarwinSystemManager{}
}

// GetCFNetworkAndDarwinVersion .
func (d DarwinSystemManager) GetCFNetworkAndDarwinVersion(iosVersion string) (string, string) {
	for i := range darwinSystemList {
		osVersionList := darwinSystemList[i].IOSVersions

		for _, v := range osVersionList {
			if v == iosVersion {
				return darwinSystemList[i].CFNetworkVersion, darwinSystemList[i].DarwinVersion
			}
		}
	}

	return darwinSystemList[0].CFNetworkVersion, darwinSystemList[0].DarwinVersion
}
