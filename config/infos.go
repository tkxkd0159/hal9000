package config

import (
	"fmt"
	"net/url"
	"sync"

	"github.com/spf13/viper"
)

// IBCInfo -> field name : IBC Port, field value : IBC Channel
type IBCInfo struct {
	Transfer string
	ICA      string
}

type HostChainInfo struct {
	Name        string
	GrpcAddr    string
	Validator   string
	HostAccount string
	Denom       string
	Decimal     uint32
	IBCInfo
	IBCTimeout uint64
	mu         sync.RWMutex
}

func NewHostChainInfo(zone string) *HostChainInfo {
	return &HostChainInfo{Name: zone}
}

func (hci *HostChainInfo) SetDetailInfos() {
	hci.mu.Lock()
	host := hci.Name
	defer hci.mu.Unlock()

	ip := viper.GetString(fmt.Sprintf("net.ip.%s", host))
	hci.GrpcAddr = ip + ":" + viper.GetString("net.port.grpc")
	hci.Validator = viper.GetString(fmt.Sprintf("%s.val_addr", host))
	hci.HostAccount = viper.GetString(fmt.Sprintf("%s.host_addr", host))
	hci.Denom = viper.GetString(fmt.Sprintf("%s.denom", host))
	hci.Decimal = viper.GetUint32(fmt.Sprintf("%s.decimal", host))
}

func (hci *HostChainInfo) WithIBCInfo(bc BotCommon, botTypes string) {
	hci.mu.Lock()
	defer hci.mu.Unlock()

	switch botTypes {
	case ActStake:
		sf := bc.(StakeFlags)
		hci.IBCTimeout = sf.IBCTimeout()
	case ActAutoStake:
		rf := bc.(RestakeFlags)
		hci.IBCTimeout = rf.IBCTimeout()
	case ActWithdraw:
		wf := bc.(WithdrawFlags)
		hci.IBCTimeout = wf.IBCTimeout()
		hci.IBCInfo = wf.HostIBC
	}
}

func SetHostChainInfo(flags BotCommon, btype string) *HostChainInfo {
	info := NewHostChainInfo(flags.GetBase().HostChain)
	switch btype {
	case ActOracle, ActStake, ActAutoStake, ActWithdraw:
		info.SetDetailInfos()
	default:
		return info
	}
	info.WithIBCInfo(flags, btype)

	return info
}

type ChainNetInfo struct {
	ChainID string
	IP      string
	Secure  bool
	GRPC    *url.URL
	TmRPC   *url.URL
	TmWsRPC *url.URL
}

func NewChainNetInfo(zone string) (ni *ChainNetInfo) {
	ip := viper.GetString(fmt.Sprintf("net.ip.%s", zone))
	chainId := viper.GetString(fmt.Sprintf("%s.chain_id", zone))
	return &ChainNetInfo{
		ChainID: chainId,
		IP:      ip,
		Secure:  viper.GetBool("net.connection.secure"),
		GRPC:    &url.URL{Scheme: "tcp", Host: ip + ":" + viper.GetString("net.port.grpc")},
		TmRPC:   &url.URL{Scheme: "tcp", Host: ip + ":" + viper.GetString("net.port.tmrpc")},
		TmWsRPC: &url.URL{Scheme: "ws", Host: ip + ":" + viper.GetString("net.port.tmrpc"), Path: "/websocket"},
	}
}
