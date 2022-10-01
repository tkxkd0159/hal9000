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
