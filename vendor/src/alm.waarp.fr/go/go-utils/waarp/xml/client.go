package xml

import (
	"encoding/xml"
)

type Client struct {
	XMLName xml.Name `xml:"config"`
	Comment string   `xml:"config,omitempty"`

	Identity  *Identity    `xml:"identity"`
	Client    *ClientBlock `xml:"client,omitempty"`
	Ssl       *Ssl         `xml:"ssl,omitempty"`
	Directory *Directory   `xml:"directory"`
	Limit     *Limit       `xml:"limit,omitempty"`
	Db        *Db          `xml:"db,omitempty"`
	Business  *Business    `xml:"business,omitempty"`
	Aliases   []*Alias     `xml:"aliases,omitempty"`
}

func NewClient() (s *Client) {
	s = &Client{}

	s.Identity = &Identity{}
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
	return
}
