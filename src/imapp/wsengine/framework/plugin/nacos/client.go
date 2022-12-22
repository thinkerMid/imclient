package nacos

import (
	"context"
	"fmt"
	hertzClient "github.com/cloudwego/hertz/pkg/app/client"
	hertzProtocol "github.com/cloudwego/hertz/pkg/protocol"
	hertzConst "github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/goccy/go-yaml"
	"os"
	"strconv"
	"ws/framework/plugin/json"
	functionTools "ws/framework/utils/function_tools"
)

type nacosEnv struct {
	DataId      string
	NamespaceId string
	GroupId     string
	Ip          string
	Port        uint64
	Username    string
	Password    string
	Token       string
}

func getNacosEnv() nacosEnv {
	ip := os.Getenv("NACOS_IP")
	id := os.Getenv("NACOS_DATA_ID")
	port := os.Getenv("NACOS_PORT")
	group := os.Getenv("NACOS_GROUP")
	spaceId := os.Getenv("NACOS_NAMESPACE_ID")
	username := os.Getenv("NACOS_USERNAME")
	password := os.Getenv("NACOS_PASSWORD")

	if ip == "" || port == "" || group == "" || id == "" || spaceId == "" {
		panic("NACOS INVALID PARAMS")
	}

	ports, _ := strconv.Atoi(port)

	return nacosEnv{
		id,
		spaceId,
		group,
		ip,
		uint64(ports),
		username,
		password,
		"",
	}
}

// New 依赖os环境变量，interface{} 是具体的配置
func New(i interface{}) {
	nacosEnv := getNacosEnv()

	client, _ := hertzClient.NewClient()
	defer client.CloseIdleConnections()

	login(client, &nacosEnv)

	run(client, nacosEnv, i)
}

func login(client *hertzClient.Client, env *nacosEnv) {
	if env.Username == "" {
		return
	}

	url := fmt.Sprintf("http://%s:%v/nacos/v1/auth/users/login", env.Ip, env.Port)

	req := hertzProtocol.AcquireRequest()
	resp := hertzProtocol.AcquireResponse()

	defer hertzProtocol.ReleaseRequest(req)
	defer hertzProtocol.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(hertzConst.MethodPost)

	args := req.URI().QueryArgs()
	args.Set("username", env.Username)
	args.Set("password", env.Password)

	err := client.Do(context.Background(), req, resp)

	if err != nil {
		panic("can't connect to nacos server error. " + err.Error())
	}

	buffer := resp.Body()

	var result map[string]interface{}
	err = json.Unmarshal(buffer, &result)
	if err != nil {
		panic("can't unmarshal nacos server reply login response. " + err.Error())
	}

	token, ok := result["accessToken"]

	if !ok {
		panic("can't find nacos access token")
	}

	env.Token = token.(string)
}

func run(client *hertzClient.Client, env nacosEnv, dstConfig interface{}) {
	url := fmt.Sprintf("http://%s:%v/nacos/v1/cs/configs", env.Ip, env.Port)

	req := hertzProtocol.AcquireRequest()
	resp := hertzProtocol.AcquireResponse()

	defer hertzProtocol.ReleaseRequest(req)
	defer hertzProtocol.ReleaseResponse(resp)

	req.SetRequestURI(url)
	req.Header.SetMethod(hertzConst.MethodGet)

	args := req.URI().QueryArgs()
	args.Set("dataId", env.DataId)
	args.Set("group", env.GroupId)
	args.Set("tenant", env.NamespaceId)
	args.Set("accessToken", env.Token)

	err := client.Do(context.Background(), req, resp)

	if err != nil {
		panic("can't connect to nacos server error. " + err.Error())
	}

	buffer := resp.Body()

	err = yaml.Unmarshal(buffer, dstConfig)
	if err != nil {
		panic(err.Error())
	}

	envBody, _ := json.Marshal(env)
	configBody, _ := json.Marshal(dstConfig)
	fmt.Println("nacos env:", functionTools.B2S(envBody))
	fmt.Println("nacos config:", functionTools.B2S(configBody))
}
