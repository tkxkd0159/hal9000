package types

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type NetAddress struct {
	ID   NodeID `json:"id"`
	IP   net.IP `json:"ip"`
	Port uint16 `json:"port"`
}

const EmptyNetAddress = "<nil-NetAddress>"

// NewNetAddress returns a new NetAddress using the provided TCP
// address. When testing, other net.Addr (except TCP) will result in
// using 0.0.0.0:0. When normal run, other net.Addr (except TCP) will
// panic. Panics if ID is invalid.
func NewNetAddress(id NodeID, addr net.Addr) *NetAddress {
	tcpAddr, ok := addr.(*net.TCPAddr)
	if !ok {
		if flag.Lookup("test.v") == nil { // normal run
			panic(fmt.Sprintf("Only TCPAddrs are supported. Got: %v", addr))
		} else { // in testing
			netAddr := NewNetAddressIPPort(net.IP("127.0.0.1"), 0)
			netAddr.ID = id
			return netAddr
		}
	}

	if err := id.Validate(); err != nil {
		panic(fmt.Sprintf("Invalid ID %v: %v (addr: %v)", id, err, addr))
	}

	ip := tcpAddr.IP
	port := uint16(tcpAddr.Port)
	na := NewNetAddressIPPort(ip, port)
	na.ID = id
	return na
}

// NewNetAddressIPPort returns a new NetAddress using the provided IP
// and port number.
func NewNetAddressIPPort(ip net.IP, port uint16) *NetAddress {
	return &NetAddress{
		IP:   ip,
		Port: port,
	}
}

// NewNetAddressString returns a new NetAddress using the provided address in
// the form of "ID@IP:Port".
// Also resolves the host if host is not an IP.
// Errors are of type ErrNetAddressXxx where Xxx is in (NoID, Invalid, Lookup)
func NewNetAddressString(addr string) (*NetAddress, error) {
	addrWithoutProtocol := removeProtocolIfDefined(addr)
	spl := strings.Split(addrWithoutProtocol, "@")
	if len(spl) != 2 {
		return nil, ErrNetAddressNoID{addr}
	}

	id, err := NewNodeID(spl[0])
	if err != nil {
		return nil, ErrNetAddressInvalid{addrWithoutProtocol, err}
	}

	if err := id.Validate(); err != nil {
		return nil, ErrNetAddressInvalid{addrWithoutProtocol, err}
	}

	addrWithoutProtocol = spl[1]

	// get host and port
	host, portStr, err := net.SplitHostPort(addrWithoutProtocol)
	if err != nil {
		return nil, ErrNetAddressInvalid{addrWithoutProtocol, err}
	}
	if len(host) == 0 {
		return nil, ErrNetAddressInvalid{
			addrWithoutProtocol,
			errors.New("host is empty")}
	}

	ip := net.ParseIP(host)
	if ip == nil {
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, ErrNetAddressLookup{host, err}
		}
		ip = ips[0]
	}

	port, err := strconv.ParseUint(portStr, 10, 16)
	if err != nil {
		return nil, ErrNetAddressInvalid{portStr, err}
	}

	na := NewNetAddressIPPort(ip, uint16(port))
	na.ID = id
	return na, nil
}

// Equals reports whether na and other are the same addresses,
// including their ID, IP, and Port.
func (na *NetAddress) Equals(other interface{}) bool {
	if o, ok := other.(*NetAddress); ok {
		return na.String() == o.String()
	}
	return false
}

// Same returns true is na has the same non-empty ID or DialString as other.
func (na *NetAddress) Same(other interface{}) bool {
	if o, ok := other.(*NetAddress); ok {
		if na.DialString() == o.DialString() {
			return true
		}
		if na.ID != "" && na.ID == o.ID {
			return true
		}
	}
	return false
}

// String representation: <ID>@<IP>:<PORT>
func (na *NetAddress) String() string {
	if na == nil {
		return EmptyNetAddress
	}

	addrStr := na.DialString()
	if na.ID != "" {
		addrStr = na.ID.AddressString(addrStr)
	}

	return addrStr
}

func (na *NetAddress) DialString() string {
	if na == nil {
		return "<nil-NetAddress>"
	}
	return net.JoinHostPort(
		na.IP.String(),
		strconv.FormatUint(uint64(na.Port), 10),
	)
}

// Dial calls net.Dial on the address.
func (na *NetAddress) Dial() (net.Conn, error) {
	conn, err := net.Dial("tcp", na.DialString())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// DialTimeout calls net.DialTimeout on the address.
func (na *NetAddress) DialTimeout(timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout("tcp", na.DialString(), timeout)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

const (
	protoHTTP  = "http"
	protoHTTPS = "https"
	protoWSS   = "wss"
	protoWS    = "ws"
	protoTCP   = "tcp"
	protoUNIX  = "unix"
)

//-------------------------------------------------------------

// Parsed URL structure
type parsedURL struct {
	url.URL

	isUnixSocket bool
}

// Parse URL and set defaults
func newParsedURL(remoteAddr string) (*parsedURL, error) {
	u, err := url.Parse(remoteAddr)
	if err != nil {
		return nil, err
	}

	// default to tcp if nothing specified
	if u.Scheme == "" {
		u.Scheme = protoTCP
	}

	pu := &parsedURL{
		URL:          *u,
		isUnixSocket: false,
	}

	if u.Scheme == protoUNIX {
		pu.isUnixSocket = true
	}

	return pu, nil
}

// Change protocol to HTTP for unknown protocols and TCP protocol - useful for RPC connections
func (u *parsedURL) SetDefaultSchemeHTTP() {
	// protocol to use for http operations, to support both http and https
	switch u.Scheme {
	case protoHTTP, protoHTTPS, protoWS, protoWSS:
		// known protocols not changed
	default:
		// default to http for unknown protocols (ex. tcp)
		u.Scheme = protoHTTP
	}
}

func (u parsedURL) GetHostWithPath() string {
	// Remove protocol, userinfo and # fragment, assume opaque is empty
	return u.Host + u.EscapedPath()
}

func (u parsedURL) GetTrimmedHostWithPath() string {
	// if it's not an unix socket we return the normal URL
	if !u.isUnixSocket {
		return u.GetHostWithPath()
	}
	// if it's a unix socket we replace the host slashes with a period
	// this is because otherwise the http.Client would think that the
	// domain is invalid.
	return strings.ReplaceAll(u.GetHostWithPath(), "/", ".")
}

// GetDialAddress returns the endpoint to dial for the parsed URL
func (u parsedURL) GetDialAddress() string {
	// if it's not a unix socket we return the host, example: localhost:443
	if !u.isUnixSocket {
		return u.Host
	}
	// otherwise we return the path of the unix socket, ex /tmp/socket
	return u.GetHostWithPath()
}

func (u parsedURL) GetTrimmedURL() string {
	return u.Scheme + "://" + u.GetTrimmedHostWithPath()
}
