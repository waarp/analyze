package xml

import "encoding/xml"

type FtpAuthent struct {
	XMLName   xml.Name     `xml:"authent"`
	Comment   string       `xml:"comment,omitempty"`
	FtpAccess []*FtpAccess `xml:"entry"`
}

func (f *FtpAuthent) AddAccess(a *FtpAccess) {
	f.FtpAccess = append(f.FtpAccess, a)
}

type FtpAccess struct {
	User          string   `xml:"user"`
	PasswordFile  string   `xml:"passwdfile,omitempty"`
	Password      string   `xml:"passwd,omitempty"`
	Accounts      []string `xml:"account,omitempty"`
	Admin         bool     `xml:"admin"`
	RetrieveCmd   string   `xml:"retrievecmd,omitempty"`
	RetrieveDelay int      `xml:"retrievedelay,omitempty"`
	StoreCmd      string   `xml:"storecmd,omitempty"`
	StoreDelay    int      `xml:"storedelay,omitempty"`
}

func (f *FtpAccess) AddAccount(name string) {
	f.Accounts = append(f.Accounts, name)
}
