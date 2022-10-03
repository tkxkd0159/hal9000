package config

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"sync"
)

// IBCInfo -> field name : IBC Port, field value : IBC Channel
type IBCInfo struct {
	Transfer string
	ICA      string
}

type HostChainInfo struct {
	GrpcAddr    string
	Validator   string
	HostAccount string
	Denom       string
	Decimal     uint32
	mu          sync.RWMutex
}

func (hci *HostChainInfo) Set(host string) {
	hci.mu.Lock()
	defer hci.mu.Unlock()

	ip := viper.GetString(fmt.Sprintf("net.ip.%s", host))
	hci.GrpcAddr = ip + ":" + viper.GetString("net.port.grpc")
	hci.Validator = viper.GetString(fmt.Sprintf("%s.val_addr", host))
	hci.HostAccount = viper.GetString(fmt.Sprintf("%s.host_addr", host))
	hci.Denom = viper.GetString(fmt.Sprintf("%s.denom", host))
	hci.Decimal = viper.GetUint32(fmt.Sprintf("%s.decimal", host))
}

type NovaInfo struct {
	ChainID string
	IP      string
	TmRPC   *url.URL
	TmWsRPC *url.URL
	mu      sync.RWMutex
}

func NewNovaInfo() (ni NovaInfo) {
	ni.mu.Lock()
	defer ni.mu.Unlock()

	ni.ChainID = "nova"
	ni.IP = viper.GetString(fmt.Sprintf("net.ip.%s", ni.ChainID))
	ni.TmRPC = &url.URL{Scheme: "tcp", Host: ni.IP + ":" + viper.GetString("net.port.tmrpc")}
	ni.TmWsRPC = &url.URL{Scheme: "ws", Host: ni.IP + ":" + viper.GetString("net.port.tmrpc"), Path: "/websocket"}

	return
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
