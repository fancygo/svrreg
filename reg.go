package svrreg

type Register interface {
	SvrRegInit(cfg *RegCfg) bool
	RegSvr() bool
	UnregSvr() bool
}

type RegCfg struct {
	LocalSvrID   string
	LocalSvrName string
	LocalSvrDNS  string
	LocalSvrPort int

	CoreSvrDNS  string
	CoreSvrPort int

	SvrCheckTimeout  int
	SvrCheckInterval int
}

func Reginit(reg Register, cfg *RegCfg) bool {
	return reg.SvrRegInit(cfg)
}

func Reg(reg Register) bool {
	return reg.RegSvr()
}

func Unreg(reg Register) bool {
	return reg.UnregSvr()
}
