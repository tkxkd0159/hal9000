package config

import (
	"flag"
	"fmt"
	"os"
)

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
	Observable
}

type Observable interface {
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
	return of.BaseFlags
}

type RestakeFlags struct {
	BaseFlags
}

func (rf RestakeFlags) GetBase() BaseFlags {
	return rf.BaseFlags
}

type StakeFlags struct {
	BaseFlags
}

func (sf StakeFlags) GetBase() BaseFlags {
	return sf.BaseFlags
}

type WithdrawFlags struct {
	BaseFlags
	HostIBC IBCInfo
}

func (wf WithdrawFlags) GetBase() BaseFlags {
	return wf.BaseFlags
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
	chanID := flag.String("ch", "channel-225", "Host Transfer Channel ID")
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

func SetFlags(action string) (bf BotCommon) {
	//actCmd := flag.NewFlagSet("action", flag.ExitOnError)
	switch action {
	case ActOracle:
		bf = SetOracleFlags()
	case ActStake:
		bf = SetStakeFlags()
	case ActRestake:
		bf = SetRestakeFlags()
	case ActWithdraw:
		bf = SetWithdrawFlags()
	}

	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "  hal %s [flags]\n\nflags:\n", action)
		flag.PrintDefaults()
	}

	if len(os.Args) >= 3 {
		switch os.Args[2] {
		case "-h", "-help", "--help":
			flag.Usage()
			os.Exit(0)
		}
	}

	return
}
