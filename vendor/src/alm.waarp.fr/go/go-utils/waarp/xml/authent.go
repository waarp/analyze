package xml

import (
	"encoding/xml"
)

type Authent struct {
	XMLName   xml.Name    `xml:"authent"`
	Comment   string      `xml:"comment,omitempty"`
	Instances []*Instance `xml:"entry"`
}

func (a *Authent) AddInstance(i *Instance) {
	a.Instances = append(a.Instances, i)
}

type Instance struct {
	HostId   string `xml:"hostid"`
	Address  string `xml:"address"`
	Port     int    `xml:"port"`
	IsSsl    bool   `xml:"isssl"`
	Admin    bool   `xml:"admin,omitempty"`
	IsClient bool   `xml:"isclient,omitempty"`
	KeyFile  string `xml:"keyfile,omitempty"`
	Key      string `xml:"key,omitempty"`
}
