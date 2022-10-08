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
	mu sync.RWMutex
}

func NewHostChainInfo(zone string) *HostChainInfo {
	return &HostChainInfo{Name: zone}
}

func (hci *HostChainInfo) Set() {
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
	switch botTypes {
	case ActWithdraw:
		hci.IBCInfo = bc.(WithdrawFlags).HostIBC
	}
}

type ChainNetInfo struct {
	ChainID string
	IP      string
	TmRPC   *url.URL
	TmWsRPC *url.URL
}

func NewChainNetInfo(zone string) (ni *ChainNetInfo) {
	ip := viper.GetString(fmt.Sprintf("net.ip.%s", zone))
	return &ChainNetInfo{
		ChainID: zone,
		IP:      ip,
		TmRPC:   &url.URL{Scheme: "tcp", Host: ip + ":" + viper.GetString("net.port.tmrpc")},
		TmWsRPC: &url.URL{Scheme: "ws", Host: ip + ":" + viper.GetString("net.port.tmrpc"), Path: "/websocket"}}
}

type BotScrt struct {
	addr       string
	passphrase string
}

func NewBotScrt(zone string, addrTarget string, keyname ...string) (bi BotScrt) {
	if len(keyname) == 1 {
		bi.passphrase = GetPassphrase(Sviper)
	}
	bi.addr = viper.GetString(fmt.Sprintf("%s.%s", zone, addrTarget))
	return
}

func (b BotScrt) Address() string {
	return b.addr
}

func (b BotScrt) Passphrase() string {
	return b.passphrase
}
