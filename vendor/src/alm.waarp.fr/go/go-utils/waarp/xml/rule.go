package xml

import (
	"encoding/xml"
)

type Rules struct {
	XMLName xml.Name `xml:"rules"`
	Rules   []*Rule
}

func (r *Rules) AddRule(rule *Rule) {
	r.Rules = append(r.Rules, rule)
}

type Rule struct {
	XMLName xml.Name `xml:"rule"`
	Id      string   `xml:"idrule"`
	Comment string   `xml:"comment,omitempty"`
	HostIds []string `xml:"hostids>hostid,omitempty"`
	Mode    int      `xml:"mode"`

	// Path
	RecvPath    string `xml:"recvpath,omitempty"`
	SendPath    string `xml:"sendpath,omitempty"`
	ArchivePath string `xml:"archivepath,omitempty"`
	WorkPath    string `xml:"workpath,omitempty"`

	// Tasks
	RPreTasks   []Task `xml:"rpretasks>tasks>task"`
	RPostTasks  []Task `xml:"rposttasks>tasks>task"`
	RErrorTasks []Task `xml:"rerrortasks>tasks>task"`
	SPreTasks   []Task `xml:"spretasks>tasks>task"`
	SPostTasks  []Task `xml:"sposttasks>tasks>task"`
	SErrorTasks []Task `xml:"serrortasks>tasks>task"`
}

type Task struct {
	Type  string `xml:"type"`
	Path  string `xml:"path"`
	Delay int    `xml:"delay"`
}
