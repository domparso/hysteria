package cmd

import (
	"fmt"
	"github.com/apernet/hysteria/core/server"
	"time"
)

var sspanelConfig SSPanelConfig

func getSSPanelConfig(param map[string][]string, config ServerConfig) {
	sspanelConfig.UserUrl = fmt.Sprintf("%s?token=%s&node_id=%d&node_type=hysteria", config.Panel.ApiHost+V2board_uri_user, config.Panel.ApiKey, config.Panel.NodeID)
}

func GetSSPanelApiProvider(config *ServerConfig) server.Authenticator {
	// 创建定时更新用户UUID协程
	return &V2boardApiProvider{URL: sspanelConfig.UserUrl}
}

func UpdateSSPanelUsers(interval time.Duration, trafficlogger server.TrafficLogger) {

}
