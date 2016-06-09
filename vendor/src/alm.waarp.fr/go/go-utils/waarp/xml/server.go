package xml

import (
	"encoding/xml"
)

type Server struct {
	XMLName xml.Name `xml:"config"`
	Comment string   `xml:"config,omitempty"`

	Identity  *Identity    `xml:"identity"`
	Server    *ServerBlock `xml:"server"`
	Network   *Network     `xml:"network"`
	Ssl       *Ssl         `xml:"ssl,omitempty"`
	Directory *Directory   `xml:"directory"`
	Rest      *Rest        `xml:"rest,omitempty"`
	Limit     *Limit       `xml:"limit"`
	Db        *Db          `xml:"db"`
	Business  *Business    `xml:"business,omitempty"`
	Roles     []*Role      `xml:"roles,omitempty"`
	Aliases   []*Alias     `xml:"aliases,omitempty"`
}

type Rest struct {
	XMLName       xml.Name     `xml:"rest"`
	Address       string       `xml:"restaddress"`
	Port          Port         `xml:"restport"`
	Ssl           bool         `xml:"restssl"`
	Authenticated bool         `xml:"restauthenticated"`
	TimeLimit     int          `xml:"resttimelimit"`
	Signature     bool         `xml:"restsignature"`
	SignatureKey  string       `xml:"restsigkey"`
	RestMethod    []RestMethod `xml:"restmethod"`
}

type RestMethod struct {
	XMLName xml.Name `xml:"restmethod"`
	Name    string   `xml:"restname"`
	Crud    string   `xml:"restcrud"`
}

func NewServer() (s *Server) {
	s = &Server{}

	s.Identity = &Identity{}
	s.Server = &ServerBlock{
		UseNoSsl: true,
	}
	s.Network = &Network{
		ServerPort:      6666,
		ServerSslPort:   6667,
		ServerHttpPort:  8066,
		ServerHttpsPort: 8067,
	}
	s.Directory = &Directory{
		In:   "in",
		Out:  "out",
		Arch: "arch",
		Work: "work",
		Conf: "conf",
	}
	s.Limit = &Limit{
		UseFastMd5:   false,
		DelayRetry:   10000,
		TimeoutCon:   10000,
		MemoryLimit:  4000000000,
		RunLimit:     10000,
		DelayCommand: 5000,
	}
	s.Db = &Db{}
	return
}
