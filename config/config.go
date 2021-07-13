package config

import (
	"log"

	"github.com/go-ini/ini"
)

type WebServer struct {
	Host string
}

var WebServerSetting = &WebServer{}

type TcpServer struct {
	Host string
}

var TcpServerSetting = &TcpServer{}

// Setup 启动配置
func Setup() {
	cfg, err := ini.Load("./my.ini")
	if err != nil {
		log.Fatalf("Fail to parse '../my.ini': %v", err)
	}

	mapTo(cfg, "webserver", WebServerSetting)
	mapTo(cfg, "tcpserver", TcpServerSetting)
}

func mapTo(cfg *ini.File, section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo RedisSetting err: %v", err)
	}
}
