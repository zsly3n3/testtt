package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

func main() {
	Cfg := loadStartConfig(config)
	port := Cfg.Section(`common`).Key(`server_port`).String()
	r := gin.Default()
	r.GET("/p3", func(c *gin.Context) {
		ipStr,err:=getLocalIP()
        if err!=nil{
			c.JSON(200, gin.H{
				"message": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"message": ipStr,
		})
	})
	r.GET("/p6", func(c *gin.Context) {
		bytes,_:=httpClientGet(`http://192.168.100.66:30080/p3`)
		c.JSON(200, gin.H{
			"message": string(bytes),
		})
	})
	addr := fmt.Sprintf(`:%s`, port)
	err:=r.Run(addr)
	if err!=nil{
		log.Fatal(`run_err:`,err.Error())
	}
}

func httpClientGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}
	err = errors.New(`not found`)
	return
}

//初始化配置
func loadStartConfig(env string) *ini.File {
	var err error
	pwd, _ := os.Getwd()
	cfg, err := ini.Load(fmt.Sprintf(`%s/config/%s.ini`, pwd, env))
	if err != nil {
		log.Fatalln(err)
	}
	return cfg
}

var (
	config      string
	serviceName string
	help        bool
)

func init() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	flag.StringVar(&config, "c", "local", "配置文件名")
	flag.StringVar(&serviceName, "s", "api", "服务功能")
	flag.BoolVar(&help, "h", false, "使用说明")
	flag.Parse()
	if help {
		flag.PrintDefaults()
		os.Exit(2)
	}
}
