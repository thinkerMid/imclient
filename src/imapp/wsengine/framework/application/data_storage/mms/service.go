package mmsService

import (
	"fmt"
	"github.com/chenzhuoyu/base64x"
	hertzConst "github.com/cloudwego/hertz/pkg/protocol/consts"
	"net/url"
	containerInterface "ws/framework/application/container/abstract_interface"
	mmsConstant "ws/framework/application/data_storage/mms/constant"
	"ws/framework/lib/media_crypto"
	httpApi "ws/framework/plugin/network/http_api"
	"ws/framework/utils"
)

type uploadMediaResp struct {
	URL        string `json:"url"`
	DirectPath string `json:"direct_path"`
}

var _ containerInterface.IMultimediaMessagingService = &MultimediaMessaging{}

// MultimediaMessaging 多媒体管理
type MultimediaMessaging struct {
	containerInterface.BaseService

	mmsEndpoints mmsConstant.MMSEndpoints
}

// Init .
func (m *MultimediaMessaging) Init() {
	m.mmsEndpoints = mmsConstant.MMSEndpoints{Buckets: []mmsConstant.Bucket{}}
}

// UpdateMMSEndpoints .
func (m *MultimediaMessaging) UpdateMMSEndpoints(v mmsConstant.MMSEndpoints) {
	m.mmsEndpoints = v
}

// UploadMediaFile .
func (m *MultimediaMessaging) UploadMediaFile(file mediaCrypto.File) (fileUrl string, filePath string, err error) {
	auth := m.mmsEndpoints.Auth
	buckets := m.mmsEndpoints.Buckets
	if len(buckets) == 0 {
		err = fmt.Errorf("object bucket not initialization")
		return
	}

	mediaId := utils.RandInt64(0, 0xfffffffe)

	token := base64x.URLEncoding.EncodeToString(file.FileEncSHA256)
	queryParams := url.Values{
		"direct_ip": []string{"0"},
		"auth":      []string{auth},
		"token":     []string{token},
		//"resume":    []string{"1"}, // 会话内重复发的图片，用resume(清空会话重新计算)
		"media_id": []string{fmt.Sprintf("%v", mediaId)}, // 会话内第一次发某张图片,重复发时给resume字段
	}

	uploadURL := url.URL{
		Scheme:   "https",
		Path:     fmt.Sprintf("/%s/%s", file.UploadPath, token),
		RawQuery: queryParams.Encode(),
	}

	client := m.AppIocContainer.ResolveHttpClient()

	var resp uploadMediaResp

	// 采用多域名访问策略
	for _, bucket := range buckets {
		uploadURL.Host = bucket.Host

		err = httpApi.DoAndBind(
			client, &resp,
			httpApi.Url(uploadURL.String()),
			httpApi.Method(hertzConst.MethodPost),
			httpApi.UserAgent(m.AppIocContainer.ResolveDeviceService().DeviceAgent()),
			httpApi.Body(file.UploadBuff),
		)

		if err != nil && len(bucket.FallbackHost) > 0 {
			uploadURL.Host = bucket.FallbackHost

			err = httpApi.DoAndBind(
				client, &resp,
				httpApi.Url(uploadURL.String()),
				httpApi.Method(hertzConst.MethodPost),
				httpApi.UserAgent(m.AppIocContainer.ResolveDeviceService().DeviceAgent()),
				httpApi.Body(file.UploadBuff),
			)
		}

		if err == nil {
			fileUrl = resp.URL
			filePath = resp.DirectPath
			break
		}
	}

	return
}
