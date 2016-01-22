package connman

import (
	"fmt"

	"github.com/godbus/dbus"
)

type IPv4Config struct {
	Method  string
	Address string
	Netmask string
	Gateway string
}

type IPv6Config struct {
	Method       string
	Address      string
	PrefixLength uint8
	Gateway      string
	Privacy      string
}

type EthConfig struct {
	Method    string
	Interface string
	Address   string
	MTU       uint16
}

type ProxyConfig struct {
	Method   string
	URL      string
	Servers  []string
	Excludes []string
}

type Provider struct {
	Host   string
	Domain string
	Name   string
	Type   string
}

type Service struct {
	Path        dbus.ObjectPath
	Name        string
	Type        string
	State       string
	Error       string
	Security    []string
	Strength    uint8
	Favorite    bool
	AutoConnect bool
	Immutable   bool
	Roaming     bool

	Ethernet           EthConfig
	IPv4               IPv4Config
	IPv4Configuration  IPv4Config
	IPv6               IPv6Config
	IPv6Configuration  IPv6Config
	Proxy              ProxyConfig
	ProxyConfiguration ProxyConfig
	Provider           Provider

	Domains                  []string
	DomainsConfiguration     []string
	Nameservers              []string
	NameserversConfiguration []string
	Timeservers              []string
	TimeserversConfiguration []string
}

func (s *Service) Connect(psk string) error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	secure := false
	for _, s := range s.Security {
		if s == "psk" || s == "wep" {
			secure = true
			break
		}
	}

	if !secure {
		_, err = db.Call("Connect")
		return err
	}

	ag := NewAgent(psk)
	if ag == nil {
		return fmt.Errorf("Could not spawn a new agent\n")
	}

	if err := RegisterAgent(ag); err != nil {
		return err
	}
	defer func() {
		UnregisterAgent(ag)
		ag.Destroy()
	}()

	_, err = db.Call("Connect")
	return err
}

func (s *Service) Disconnect() error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	_, err = db.Call("Disconnect")
	return err
}

func (s *Service) ApplyIP() error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	arg, err := structToDict(s.IPv4Configuration)
	if err != nil {
		return err
	}

	return db.Set("IPv4.Configuration", arg)
}

func (s *Service) ApplyDNS() error {
	db, err := DBusService(s.Path)
	if err != nil {
		return err
	}

	return db.Set("Nameservers.Configuration", s.NameserversConfiguration)
}
