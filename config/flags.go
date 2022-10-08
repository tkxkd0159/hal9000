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

func addBaseFlags(cmd *flag.FlagSet) {
	cmd.BoolVar(&isTest, "test", false, "Decide whether it's test with localnet")
	cmd.BoolVar(&isNew, "new", false, "Start client with making new account")
	cmd.BoolVar(&disp, "display", false, "Show context log through stdout")
	cmd.StringVar(&apiAddr, "api", "127.0.0.1:3334", "Set bot api address")
	cmd.StringVar(&keyname, "name", "nova_bot", "Set unique key name (uid)")
	cmd.StringVar(&hostchainName, "host", "gaia", "Name of the host chain from which to obtain oracle info")
	cmd.IntVar(&intv, "interval", 15*60, "Oracle update interval (sec)")
	cmd.StringVar(&logloc, "logloc", "logs", "Where All Logs Are Stored from project root")
}

func SetOracleFlags(cmd *flag.FlagSet) OracleFlags {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		panic("Something went wrong while parse flags")
	}

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

func SetRestakeFlags(cmd *flag.FlagSet) RestakeFlags {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		panic("Something went wrong while parse flags")
	}

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

func SetStakeFlags(cmd *flag.FlagSet) StakeFlags {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		panic("Something went wrong while parse flags")
	}

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

func SetWithdrawFlags(cmd *flag.FlagSet) WithdrawFlags {
	chanID := cmd.String("ch", "channel-225", "Host Transfer Channel ID")
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		panic("Something went wrong while parse flags")
	}

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
	actCmd := flag.NewFlagSet(fmt.Sprintf("%s bot", action), flag.ExitOnError)
	actCmd.Usage = func() {
		w := actCmd.Output()
		fmt.Fprintf(w, "  hal %s [flags]\n\nflags:\n", action)
		actCmd.PrintDefaults()
	}

	addBaseFlags(actCmd)

	switch action {
	case ActOracle:
		bf = SetOracleFlags(actCmd)
	case ActStake:
		bf = SetStakeFlags(actCmd)
	case ActRestake:
		bf = SetRestakeFlags(actCmd)
	case ActWithdraw:
		bf = SetWithdrawFlags(actCmd)
	}

	return
}
