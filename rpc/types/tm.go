package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bytes"
	tmstrings "github.com/tendermint/tendermint/libs/strings"
	"net"
	"regexp"
	"strconv"
	"strings"
)

const (
	//maxNodeInfoSize = 10240 // 10KB
	maxNumChannels = 16 // plenty of room for upgrades, for now
)

// ProtocolVersion contains the protocol versions for the software.
type ProtocolVersion struct {
	P2P   uint64 `json:"p2p"`
	Block uint64 `json:"block"`
	App   uint64 `json:"app"`
}

//-------------------------------------------------------------

// NodeInfo is the basic node information exchanged
// between two peers during the Tendermint P2P handshake.
type NodeInfo struct {
	ProtocolVersion ProtocolVersion `json:"protocol_version"`

	// Authenticate
	NodeID     NodeID `json:"id"`          // authenticated identifier
	ListenAddr string `json:"listen_addr"` // accepting incoming

	// Check compatibility.
	// Channels are HexBytes so easier to read as JSON
	Network  string         `json:"network"`  // network/chain ID
	Version  string         `json:"version"`  // major.minor.revision
	Channels bytes.HexBytes `json:"channels"` // channels this node knows about

	// ASCIIText fields
	Moniker string        `json:"moniker"` // arbitrary moniker
	Other   NodeInfoOther `json:"other"`   // other application specific data
}

// NodeInfoOther is the misc. applcation specific data
type NodeInfoOther struct {
	TxIndex    string `json:"tx_index"`
	RPCAddress string `json:"rpc_address"`
}

// ID returns the node's peer ID.
func (info NodeInfo) ID() NodeID {
	return info.NodeID
}

// Validate checks the self-reported NodeInfo is safe.
// It returns an error if there
// are too many Channels, if there are any duplicate Channels,
func (info NodeInfo) Validate() error {

	// ID is already validated.

	// Validate ListenAddr.
	_, err := NewNetAddressString(info.ID().AddressString(info.ListenAddr))
	if err != nil {
		return err
	}

	// Network is validated in CompatibleWith.

	// Validate Version
	if len(info.Version) > 0 &&
		(!tmstrings.IsASCIIText(info.Version) || tmstrings.ASCIITrim(info.Version) == "") {

		return fmt.Errorf("info.Version must be valid ASCII text without tabs, but got %v", info.Version)
	}

	// Validate Channels - ensure max and check for duplicates.
	if len(info.Channels) > maxNumChannels {
		return fmt.Errorf("info.Channels is too long (%v). Max is %v", len(info.Channels), maxNumChannels)
	}
	channels := make(map[byte]struct{})
	for _, ch := range info.Channels {
		_, ok := channels[ch]
		if ok {
			return fmt.Errorf("info.Channels contains duplicate channel id %v", ch)
		}
		channels[ch] = struct{}{}
	}

	// Validate Moniker.
	if !tmstrings.IsASCIIText(info.Moniker) || tmstrings.ASCIITrim(info.Moniker) == "" {
		return fmt.Errorf("info.Moniker must be valid non-empty ASCII text without tabs, but got %v", info.Moniker)
	}

	// Validate Other.
	other := info.Other
	txIndex := other.TxIndex
	switch txIndex {
	case "", "on", "off":
	default:
		return fmt.Errorf("info.Other.TxIndex should be either 'on', 'off', or empty string, got '%v'", txIndex)
	}
	// XXX: Should we be more strict about address formats?
	rpcAddr := other.RPCAddress
	if len(rpcAddr) > 0 && (!tmstrings.IsASCIIText(rpcAddr) || tmstrings.ASCIITrim(rpcAddr) == "") {
		return fmt.Errorf("info.Other.RPCAddress=%v must be valid ASCII text without tabs", rpcAddr)
	}

	return nil
}

func ParseAddressString(addr string) (net.IP, uint16, error) {
	addrWithoutProtocol := removeProtocolIfDefined(addr)
	spl := strings.Split(addrWithoutProtocol, "@")
	if len(spl) != 2 {
		return nil, 0, errors.New("invalid address")
	}

	id, err := NewNodeID(spl[0])
	if err != nil {
		return nil, 0, err
	}

	if err := id.Validate(); err != nil {
		return nil, 0, err
	}

	addrWithoutProtocol = spl[1]

	// get host and port
	host, portStr, err := net.SplitHostPort(addrWithoutProtocol)
	if err != nil {
		return nil, 0, err
	}
	if len(host) == 0 {
		return nil, 0, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, 0, err
		}
		ip = ips[0]
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, 0, err
	}

	return ip, uint16(port), nil
}

func removeProtocolIfDefined(addr string) string {
	if strings.Contains(addr, "://") {
		return strings.Split(addr, "://")[1]
	}
	return addr

}

// CompatibleWith checks if two NodeInfo are compatible with each other.
// CONTRACT: two nodes are compatible if the Block version and network match
// and they have at least one channel in common.
func (info NodeInfo) CompatibleWith(other NodeInfo) error {
	if info.ProtocolVersion.Block != other.ProtocolVersion.Block {
		return fmt.Errorf("peer is on a different Block version. Got %v, expected %v",
			other.ProtocolVersion.Block, info.ProtocolVersion.Block)
	}

	// nodes must be on the same network
	if info.Network != other.Network {
		return fmt.Errorf("peer is on a different network. Got %v, expected %v", other.Network, info.Network)
	}

	// if we have no channels, we're just testing
	if len(info.Channels) == 0 {
		return nil
	}

	// for each of our channels, check if they have it
	found := false
OuterLoop:
	for _, ch1 := range info.Channels {
		for _, ch2 := range other.Channels {
			if ch1 == ch2 {
				found = true
				break OuterLoop // only need one
			}
		}
	}
	if !found {
		return fmt.Errorf("peer has no common channels. Our channels: %v ; Peer channels: %v", info.Channels, other.Channels)
	}
	return nil
}

// NetAddress returns a NetAddress derived from the NodeInfo -
// it includes the authenticated peer ID and the self-reported
// ListenAddr. Note that the ListenAddr is not authenticated and
// may not match that address actually dialed if its an outbound peer.
func (info NodeInfo) NetAddress() (*NetAddress, error) {
	idAddr := info.ID().AddressString(info.ListenAddr)
	return NewNetAddressString(idAddr)
}

// AddChannel is used by the router when a channel is opened to add it to the node info
func (info *NodeInfo) AddChannel(channel uint16) {
	// check that the channel doesn't already exist
	for _, ch := range info.Channels {
		if ch == byte(channel) {
			return
		}
	}

	info.Channels = append(info.Channels, byte(channel))
}

func (info NodeInfo) Copy() NodeInfo {
	return NodeInfo{
		ProtocolVersion: info.ProtocolVersion,
		NodeID:          info.NodeID,
		ListenAddr:      info.ListenAddr,
		Network:         info.Network,
		Version:         info.Version,
		Channels:        info.Channels,
		Moniker:         info.Moniker,
		Other:           info.Other,
	}
}

// NodeIDByteLength is the length of a crypto.Address. Currently only 20.
const NodeIDByteLength = crypto.AddressSize

// reNodeID is a regexp for valid node IDs.
var reNodeID = regexp.MustCompile(`^[\da-f]{40}$`)

// NodeID is a hex-encoded crypto.Address. It must be lowercased
// (for uniqueness) and of length 2*NodeIDByteLength.
type NodeID string

// NewNodeID returns a lowercased (normalized) NodeID, or errors if the
// node ID is invalid.
func NewNodeID(nodeID string) (NodeID, error) {
	n := NodeID(strings.ToLower(nodeID))
	return n, n.Validate()
}

// AddressString returns id@hostPort. It strips the leading
// protocol from protocolHostPort if it exists.
func (id NodeID) AddressString(protocolHostPort string) string {
	hostPort := removeProtocolIfDefined(protocolHostPort)
	return fmt.Sprintf("%s@%s", id, hostPort)
}

// NodeIDFromPubKey creates a node ID from a given PubKey address.
func NodeIDFromPubKey(pubKey crypto.PubKey) NodeID {
	return NodeID(hex.EncodeToString(pubKey.Address()))
}

// Bytes converts the node ID to its binary byte representation.
func (id NodeID) Bytes() ([]byte, error) {
	bz, err := hex.DecodeString(string(id))
	if err != nil {
		return nil, fmt.Errorf("invalid node ID encoding: %w", err)
	}
	return bz, nil
}

// Validate validates the NodeID.
func (id NodeID) Validate() error {
	switch {
	case len(id) == 0:
		return errors.New("empty node ID")

	case len(id) != 2*NodeIDByteLength:
		return fmt.Errorf("invalid node ID length %d, expected %d", len(id), 2*NodeIDByteLength)

	case !reNodeID.MatchString(string(id)):
		return fmt.Errorf("node ID can only contain lowercased hex digits")

	default:
		return nil
	}
}
