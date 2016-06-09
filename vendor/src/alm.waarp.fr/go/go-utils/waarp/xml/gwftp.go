package xml

import (
	"encoding/xml"
)

type Gwftp struct {
	XMLName xml.Name `xml:"config"`
	Comment string   `xml:"config,omitempty"`

	Identity  *GwFtpIdentity    `xml:"identity"`
	Server    *GwFtpServerBlock `xml:"server"`
	Network   *GwFtpNetwork     `xml:"network"`
	Exec      *GwFtpExec        `xml:"exec"`
	Directory *GwFtpDirectory   `xml:"directory"`
	Limit     *GwFtpLimit       `xml:"limit"`
	Ssl       *GwFtpSsl         `xml:"ssl,omitempty"`
	Db        *Db               `xml:"db"`
}

type GwFtpIdentity struct {
	HostId      string `xml:"hostid"`
	CryptoKey   string `xml:"cryptokey"`
	AuthentFile string `xml:"authentfile"`
}

type GwFtpServerBlock struct {
	ServerAdmin     string `xml:"serveradmin"`
	ServerPasswd    string `xml:"serverpasswd"`
	UseHttpComp     bool   `xml:"usehttpcomp,omitempty"`
	UseLocalExec    bool   `xml:"uselocalexec,omitempty"`
	LExecAddr       string `xml:"lexecaddr,omitempty"`
	LExecPort       Port   `xml:"lexecport,omitempty"`
	HttpAdmin       string `xml:"httpadmin"`
	AdmKeyPath      string `xml:"admkeypath"`
	AdmKeystorePass string `xml:"admkeystorepass"`
	AdmKeyPass      string `xml:"admkeypass"`
}

type GwFtpNetwork struct {
	ServerPort      Port `xml:"serverport"`
	ServerAddress   Port `xml:"serveraddress,omitempty"`
	ServerHttpPort  Port `xml:"serverhttpport"`
	ServerHttpsPort Port `xml:"serverhttpsport"`
	PortMin         Port `xml:"portmin"`
	PortMax         Port `xml:"portmax"`
}

type GwFtpExec struct {
	RetrieveCmd   string `xml:"retrievecmd"`
	RetrieveDelay uint   `xml:"retrievedelay,omitempty"`
	StoreCmd      string `xml:"storecmd"`
	StoreDelay    uint   `xml:"storedelay,omitempty"`
}

type GwFtpDirectory struct {
	ServerHome string `xml:"serverhome"`
}

type GwFtpLimit struct {
	DeleteOnAbort  bool    `xml:"deleteonabort,omitempty"`
	ServerThread   uint    `xml:"serverthread,omitempty"`
	ClientThread   uint    `xml:"clientthread,omitempty"`
	MemoryLimit    uint    `xml:"memorylimit,omitempty"`
	SessionLimit   uint    `xml:"sessionlimit,omitempty"`
	GlobalLimit    uint    `xml:"globallimit,omitempty"`
	DelayLimit     uint    `xml:"delaylimit,omitempty"`
	TimeoutCon     uint    `xml:"timeoutcon,omitempty"`
	BlockSize      uint    `xml:"blocksize,omitempty"`
	UseNIO         bool    `xml:"usenio,omitempty"`
	UseCpuLimit    bool    `xml:"usecpulimit,omitempty"`
	UseJdkCpuLimit bool    `xml:"usejdkcpulimit,omitempty"`
	CpuLimit       float32 `xml:"cpulimit,omitempty"`
	ConnLimit      uint    `xml:"connlimit,omitempty"`
	UseFastMd5     bool    `xml:"usefastmd5,omitempty"`
	FastMd5        string  `xml:"fastmd5,omitempty"`
	GlobalDigest   bool    `xml:"globaldigest,omitempty"`
}

type GwFtpSsl struct {
	KeyPath                    string `xml:"keypath"`
	KeystorePass               string `xml:"keystorepass"`
	KeyPass                    string `xml:"keypass"`
	TrustKeyPath               string `xml:"trustkeypath"`
	TrustKeystorePass          string `xml:"trustkeystorepass"`
	TrustUseClientAuthenticate string `xml:"trustuseclientauthenticate"`
	Implicit                   bool   `xml:"useimplicitftps,omitempty"`
	Explicit                   bool   `xml:"useexplicitftps,omitempty"`
}

func NewGwftp() (s *Gwftp) {
	s = &Gwftp{}

	s.Identity = &GwFtpIdentity{}
	s.Server = &GwFtpServerBlock{}
	s.Network = &GwFtpNetwork{
		ServerPort:      6621,
		ServerHttpPort:  8068,
		ServerHttpsPort: 8069,
	}
	s.Directory = &GwFtpDirectory{}
	s.Limit = &GwFtpLimit{
		UseFastMd5:  false,
		TimeoutCon:  10000,
		MemoryLimit: 4000000000,
	}
	s.Exec = &GwFtpExec{
		RetrieveCmd: "REFUSED",
		StoreCmd:    "REFUSED",
	}
	s.Ssl = &GwFtpSsl{}
	s.Db = &Db{}
	return
}
