package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/apernet/hysteria/core/server"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
	"time"
)

var SSPConfig SSPanelConfig

func getSSPanelConfig(config *ServerConfig) {
	SSPConfig.ConfigUrl = fmt.Sprintf(
		"%s?key=%s&muKey=%d",
		config.Panel.ApiHost+fmt.Sprintf(SSPanel_uri_conf, config.Panel.NodeID),
		config.Panel.ApiKey,
		config.Panel.ApiKey)
	SSPConfig.UserUrl = fmt.Sprintf(
		"%s?key=%s&muKey=%d",
		config.Panel.ApiHost+fmt.Sprintf(SSPanel_uri_user, config.Panel.NodeID),
		config.Panel.ApiKey,
		config.Panel.ApiKey)

	resp, err := http.Get(SSPConfig.ConfigUrl)
	if err != nil {
		// 处理错误
		fmt.Println("HTTP GET 请求出错:", err)
		Logger("fatal", "failed to client sspanel api to get nodeInfo", zap.Error(err))
	}
	defer resp.Body.Close()
	// 读取响应数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		Logger("fatal", "failed to read sspanel reaponse", zap.Error(err))
	}
	// 解析JSON数据
	var responseNodeInfo map[string]interface{}
	if err := json.Unmarshal(body, &responseNodeInfo); err != nil {
		Logger("fatal", "failed to unmarshal sspanel reaponse", zap.Error(err))
		return
	}
	fmt.Println("responseNodeInfo", responseNodeInfo)
	//var responseNodeInfo ResponseNodeInfo
	//err = json.Unmarshal(body, &responseNodeInfo)
	//if err != nil {
	//	Logger("fatal", "failed to unmarshal sspanel reaponse", zap.Error(err))
	//}
	// 给 hy的端口、obfs、上行下行进行赋值
	if responseNodeInfo["ServerPort"] != 0 {
		config.Listen = ":" + strconv.Itoa(responseNodeInfo["ServerPort"].(int))
	}
	if responseNodeInfo["DownMbps"] != 0 {
		config.Bandwidth.Down = strconv.Itoa(responseNodeInfo["DownMbps"].(int)) + "Mbps"
	}
	if responseNodeInfo["UpMbps"] != 0 {
		config.Bandwidth.Up = strconv.Itoa(responseNodeInfo["UpMbps"].(int)) + "Mbps"
	}
	if responseNodeInfo["Obfs"] != "" {
		config.Obfs.Type = "salamander"
		config.Obfs.Salamander.Password = responseNodeInfo["Obfs"].(string)
	}
}

func GetSSPanelApiProvider(config *ServerConfig) server.Authenticator {
	// 创建定时更新用户UUID协程
	return &V2boardApiProvider{URL: SSPConfig.UserUrl}
}

func UpdateSSPanelUsers(interval time.Duration, trafficlogger server.TrafficLogger) {

}
