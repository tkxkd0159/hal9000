package config

import "flag"

var (
	isTest        bool
	isNew         bool
	disp          bool
	apiAddr       string
	keyname       string
	hostchainName string
	intv          int
	logloc        string
)

var (
	_ BotCommon = OracleFlags{}
	_ BotCommon = RestakeFlags{}
	_ BotCommon = StakeFlags{}
	_ BotCommon = WithdrawFlags{}
)

type BotCommon interface {
	GetBase() BaseFlags
}

type MonitorFlag interface {
	GetExtIP() string
}

type BaseFlags struct {
	IsTest      bool
	New         bool
	Disp        bool
	ExtIP       string
	Kn          string
	HostChain   string
	Period      int
	LogLocation string
}

func (bf BaseFlags) GetExtIP() string {
	return bf.ExtIP
}

type OracleFlags struct {
	BaseFlags
}

func (of OracleFlags) GetBase() BaseFlags {
	return BaseFlags{
		of.IsTest,
		of.New,
		of.Disp,
		of.ExtIP,
		of.Kn,
		of.HostChain,
		of.Period,
		of.LogLocation,
	}
}

type RestakeFlags struct {
	BaseFlags
}

func (rf RestakeFlags) GetBase() BaseFlags {
	return BaseFlags{
		rf.IsTest,
		rf.New,
		rf.Disp,
		rf.ExtIP,
		rf.Kn,
		rf.HostChain,
		rf.Period,
		rf.LogLocation,
	}
}

type StakeFlags struct {
	BaseFlags
}

func (sf StakeFlags) GetBase() BaseFlags {
	return BaseFlags{
		sf.IsTest,
		sf.New,
		sf.Disp,
		sf.ExtIP,
		sf.Kn,
		sf.HostChain,
		sf.Period,
		sf.LogLocation,
	}
}

type WithdrawFlags struct {
	BaseFlags
	HostIBC IBCInfo
}

func (wf WithdrawFlags) GetBase() BaseFlags {
	return BaseFlags{
		wf.IsTest,
		wf.New,
		wf.Disp,
		wf.ExtIP,
		wf.Kn,
		wf.HostChain,
		wf.Period,
		wf.LogLocation,
	}
}

func setBaseFlags() {
	flag.BoolVar(&isTest, "test", false, "Decide whether it's test with localnet")
	flag.BoolVar(&isNew, "new", false, "Start client with making new account")
	flag.BoolVar(&disp, "display", false, "Show context log through stdout")
	flag.StringVar(&apiAddr, "api", "127.0.0.1:3334", "Set bot api address")
	flag.StringVar(&keyname, "name", "nova_bot", "Set unique key name (uid)")
	flag.StringVar(&hostchainName, "host", "gaia", "Name of the host chain from which to obtain oracle info")
	flag.IntVar(&intv, "interval", 15*60, "Oracle update interval (sec)")
	flag.StringVar(&logloc, "logloc", "logs", "Where All Logs Are Stored from project root")
}

func SetOracleFlags() OracleFlags {
	setBaseFlags()
	flag.Parse()
	return OracleFlags{
		BaseFlags{
			IsTest:      isTest,
			New:         isNew,
			Disp:        disp,
			ExtIP:       apiAddr,
			Kn:          keyname,
			HostChain:   hostchainName,
			Period:      intv,
			LogLocation: logloc,
		},
	}
}

func SetRestakeFlags() RestakeFlags {
	setBaseFlags()
	flag.Parse()
	return RestakeFlags{
		BaseFlags{
			IsTest:      isTest,
			New:         isNew,
			Disp:        disp,
			ExtIP:       apiAddr,
			Kn:          keyname,
			HostChain:   hostchainName,
			Period:      intv,
			LogLocation: logloc,
		},
	}
}

func SetStakeFlags() StakeFlags {
	setBaseFlags()
	flag.Parse()
	return StakeFlags{
		BaseFlags{
			IsTest:      isTest,
			New:         isNew,
			Disp:        disp,
			ExtIP:       apiAddr,
			Kn:          keyname,
			HostChain:   hostchainName,
			Period:      intv,
			LogLocation: logloc,
		},
	}
}

func SetWithdrawFlags() WithdrawFlags {
	setBaseFlags()
	chanID := flag.String("ch", "channel-2", "Host Transfer Channel ID")
	flag.Parse()
	return WithdrawFlags{
		BaseFlags{
			IsTest:      isTest,
			New:         isNew,
			Disp:        disp,
			ExtIP:       apiAddr,
			Kn:          keyname,
			HostChain:   hostchainName,
			Period:      intv,
			LogLocation: logloc,
		},
		IBCInfo{Transfer: *chanID},
	}
}
