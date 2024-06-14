package cmd

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/apernet/hysteria/core/server"
)

var V2bConfig V2boardConfig

type ResponseV2board struct {
	Host       string `json:"host"`
	ServerPort uint   `json:"server_port"`
	ServerName string `json:"server_name"`
	UpMbps     uint   `json:"down_mbps"`
	DownMbps   uint   `json:"up_mbps"`
	Obfs       string `json:"obfs"`
	BaseConfig struct {
		PushInterval int `json:"push_interval"`
		PullInterval int `json:"pull_interval"`
	} `json:"base_config"`
}

func getV2boardConfig(config *ServerConfig) {
	V2bConfig.UserUrl = fmt.Sprintf(
		"%s?token=%s&node_id=%d&node_type=hysteria",
		config.Panel.ApiHost+V2board_uri_user,
		config.Panel.ApiKey,
		config.Panel.NodeID)

	V2bConfig.PushUrl = fmt.Sprintf(
		"%s?token=%s&node_id=%d&node_type=hysteria",
		config.Panel.ApiHost+V2board_uri_push,
		config.Panel.ApiKey,
		config.Panel.NodeID)
	V2bConfig.ConfigUrl = fmt.Sprintf(
		"%s?token=%s&node_id=%d&node_type=hysteria",
		config.Panel.ApiHost+V2board_uri_conf,
		config.Panel.ApiKey,
		config.Panel.NodeID)

	// 发起 HTTP GET 请求
	//proxy, err := url.Parse("http://127.0.0.1:10809")
	//client := http.Client{
	//	Transport: &http.Transport{
	//		// 设置代理
	//		Proxy: http.ProxyURL(proxy),
	//	},
	//}

	resp, err := http.Get(V2bConfig.ConfigUrl)
	if err != nil {
		// 处理错误
		fmt.Println("HTTP GET 请求出错:", err)
		Logger("fatal", "failed to client v2board api to get nodeInfo", zap.Error(err))
	}
	defer resp.Body.Close()
	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger("fatal", "failed to read v2board reaponse", zap.Error(err))
	}
	// 解析JSON数据
	var responseNodeInfo ResponseV2board
	err = json.Unmarshal(body, &responseNodeInfo)
	if err != nil {
		Logger("fatal", "failed to unmarshal v2board reaponse", zap.Error(err))
	}
	// 给 hy的端口、obfs、上行下行进行赋值
	if responseNodeInfo.ServerPort != 0 {
		config.Listen = ":" + strconv.Itoa(int(responseNodeInfo.ServerPort))
	}
	if responseNodeInfo.DownMbps != 0 {
		config.Bandwidth.Down = strconv.Itoa(int(responseNodeInfo.DownMbps)) + "Mbps"
	}
	if responseNodeInfo.UpMbps != 0 {
		config.Bandwidth.Up = strconv.Itoa(int(responseNodeInfo.UpMbps)) + "Mbps"
	}
	if responseNodeInfo.Obfs != "" {
		config.Obfs.Type = "salamander"
		config.Obfs.Salamander.Password = responseNodeInfo.Obfs
	}
}

func GetV2boardApiProvider(config *ServerConfig) server.Authenticator {
	// 创建定时更新用户UUID协程
	return &V2boardApiProvider{URL: V2bConfig.UserUrl}
}

var _ server.Authenticator = &V2boardApiProvider{}

type V2boardApiProvider struct {
	Client *http.Client
	URL    string
}

// 用户列表
var (
	usersMap map[string]User
	lock     sync.Mutex
)

type User struct {
	ID         int     `json:"id"`
	UUID       string  `json:"uuid"`
	SpeedLimit *uint32 `json:"speed_limit"`
}
type ResponseData struct {
	Users []User `json:"users"`
}

func getUserList(url string) ([]User, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var responseData ResponseData
	err = json.NewDecoder(resp.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return responseData.Users, nil
}

func UpdateV2boardUsers(interval time.Duration, trafficlogger server.TrafficLogger) {
	fmt.Println("用户列表自动更新服务已激活")
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		userList, err := getUserList(V2bConfig.UserUrl)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		lock.Lock()
		newUsersMap := make(map[string]User)
		for _, user := range userList {
			newUsersMap[user.UUID] = user
		}
		if trafficlogger != nil {
			for uuid := range usersMap {
				if _, exists := newUsersMap[uuid]; !exists {
					trafficlogger.NewKick(strconv.Itoa(usersMap[uuid].ID))
				}
			}
		}

		usersMap = newUsersMap
		lock.Unlock()
	}

}

// 验证代码
func (v *V2boardApiProvider) Authenticate(addr net.Addr, auth string, tx uint64) (ok bool, id string) {

	// 获取判断连接用户是否在用户列表内
	lock.Lock()
	defer lock.Unlock()

	if user, exists := usersMap[auth]; exists {
		return true, strconv.Itoa(user.ID)
	}
	return false, ""
}
