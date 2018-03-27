package svrreg

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"net"
	"strconv"
)

type RegConsul struct {
	consulCfg   *consul.Config
	consulCli   *consul.Client
	localSvrReg *consul.AgentServiceRegistration
}

func NewRegConsul() *RegConsul {
	return &RegConsul{}
}

func (this *RegConsul) SvrRegInit(cfg *RegCfg) bool {
	//连接consul
	//获取本机ip
	localIp, err := net.LookupHost(cfg.LocalSvrDNS)
	if err != nil {
		fmt.Printf("LookupIP local svr err = %v\n", err)
		return false
	}
	localAdd := string(localIp[0])

	//获取consul中心服务器ip
	coresvrIp, err := net.LookupHost(cfg.CoreSvrDNS)
	if err != nil {
		fmt.Printf("LookupIP core svr err = %v\n", err)
		return false
	}
	coresvrAdd := string(coresvrIp[0])
	coresvrPort := ":" + strconv.Itoa(cfg.CoreSvrPort)
	fmt.Printf("coresvrAdd = %s\n", coresvrAdd)

	//生成consul中心服务器配置
	this.consulCfg = consul.DefaultConfig()
	this.consulCfg.Address = coresvrAdd + coresvrPort
	//建立consul cli
	this.consulCli, err = consul.NewClient(this.consulCfg)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//设置本机localsvr配置信息
	this.localSvrReg = &consul.AgentServiceRegistration{
		ID:      cfg.LocalSvrID,
		Name:    cfg.LocalSvrName,
		Address: localAdd,
		Port:    cfg.LocalSvrPort,
		Tags:    []string{cfg.LocalSvrName},
		Check: &consul.AgentServiceCheck{
			HTTP:     "http://" + cfg.LocalSvrDNS + ":" + strconv.Itoa(cfg.LocalSvrPort) + "/health",
			Interval: strconv.Itoa(cfg.SvrCheckInterval) + "s",
			Timeout:  strconv.Itoa(cfg.SvrCheckTimeout) + "s",
		},
	}
	return true
}

func (this *RegConsul) RegSvr() bool {
	//向consul中心服务器注册
	if err := this.consulCli.Agent().ServiceRegister(this.localSvrReg); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (this *RegConsul) UnregSvr() bool {
	//向consul中心服务器注销
	if err := this.consulCli.Agent().ServiceDeregister(this.localSvrReg.ID); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
