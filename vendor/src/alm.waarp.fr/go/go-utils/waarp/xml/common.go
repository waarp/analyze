package xml

type Port uint16

type Identity struct {
	HostId      string `xml:"hostid"`
	SslHostId   string `xml:"sslhostid"`
	CryptoKey   string `xml:"cryptokey"`
	AuthentFile string `xml:"authentfile"`
}

type ServerBlock struct {
	ServerAdmin        string `xml:"serveradmin"`
	ServerPasswd       string `xml:"serverpasswd"`
	UseNoSsl           bool   `xml:"usenossl"`
	UseSsl             bool   `xml:"usessl"`
	UseHttpComp        bool   `xml:"usehttpcomp,omitempty"`
	UseLocalExec       bool   `xml:"uselocalexec,omitempty"`
	LExecAddr          string `xml:"lexecaddr,omitempty"`
	LExecPort          Port   `xml:"lexecport,omitempty"`
	HttpAdmin          string `xml:"httpadmin"`
	AdmKeyPath         string `xml:"admkeypath"`
	AdmKeystorePass    string `xml:"admkeystorepass"`
	AdmKeyPass         string `xml:"admkeypass"`
	CheckAddress       string `xml:"checkaddress,omitempty"`
	CheckClientAddress string `xml:"checkclientaddress,omitempty"`
	MultipleMonitors   uint   `xml:"multiplemonitors,omitempty"`
}

type ClientBlock struct {
	TaskRunnerNoDb bool `xml:"taskrunnernodb,omitempty"`
}

type Network struct {
	ServerPort      Port `xml:"serverport"`
	ServerSslPort   Port `xml:"serversslport"`
	ServerHttpPort  Port `xml:"serverhttpport"`
	ServerHttpsPort Port `xml:"serverhttpsport"`
}

type Ssl struct {
	KeyPath                    string `xml:"keypath"`
	KeystorePass               string `xml:"keystorepass"`
	KeyPass                    string `xml:"keypass"`
	TrustKeyPath               string `xml:"trustkeypath"`
	TrustKeystorePass          string `xml:"trustkeystorepass"`
	TrustUseClientAuthenticate string `xml:"trustuseclientauthenticate"`
}

type Directory struct {
	ServerHome string `xml:"serverhome"`
	In         string `xml:"in"`
	Out        string `xml:"out"`
	Arch       string `xml:"arch"`
	Work       string `xml:"work"`
	Conf       string `xml:"conf"`
}

type Limit struct {
	ServerThread   uint    `xml:"serverthread,omitempty"`
	ClientThread   uint    `xml:"clientthread,omitempty"`
	MemoryLimit    uint    `xml:"memorylimit,omitempty"`
	SessionLimit   uint    `xml:"sessionlimit,omitempty"`
	GlobalLimit    uint    `xml:"globallimit,omitempty"`
	DelayLimit     uint    `xml:"delaylimit,omitempty"`
	RunLimit       uint    `xml:"runlimit,omitempty"`
	DelayCommand   uint    `xml:"delaycommand,omitempty"`
	DelayRetry     uint    `xml:"delayretry,omitempty"`
	TimeoutCon     uint    `xml:"timeoutcon,omitempty"`
	BlockSize      uint    `xml:"blocksize,omitempty"`
	GapRestart     uint    `xml:"gaprestart,omitempty"`
	UseNIO         bool    `xml:"usenio,omitempty"`
	UseCpuLimit    bool    `xml:"usecpulimit,omitempty"`
	UseJdkCpuLimit bool    `xml:"usejdkcpulimit,omitempty"`
	CpuLimit       float32 `xml:"cpulimit,omitempty"`
	ConnLimit      uint    `xml:"connlimit,omitempty"`
	Digest         uint    `xml:"digest,omitempty"`
	UseFastMd5     bool    `xml:"usefastmd5,omitempty"`
	FastMd5        string  `xml:"fastmd5,omitempty"`
	CheckVersion   bool    `xml:"checkversion,omitempty"`
	GlobalDigest   bool    `xml:"globaldigest,omitempty"`
}

type Db struct {
	Driver         string `xml:"dbdriver"`
	Server         string `xml:"dbserver"`
	User           string `xml:"dbuser"`
	Passwd         string `xml:"dbpasswd"`
	TaskRunnerNoDb bool   `xml:"taskrunnernodb,omitempty"`
}

type Business struct {
	BusinessId string `xml:"businessid,omitempty"`
}

type Role struct {
	RoleId  string `xml:"roleid"`
	RoleSet string `xml:"roleset"`
}

type Alias struct {
	RealId  string `xml:"realid"`
	AliasId string `xml:"aliasid"`
}

type BandwidthIdentity struct {
	HostId string `xml:"hostid"`
}

type BandwidthLimit struct {
	SessionLimit uint `xml:"sessionlimit"`
	GlobalLimit  uint `xml:"globallimit"`
	DelayLimit   uint `xml:"delaylimit"`
	RunLimit     uint `xml:"runlimit"`
	DelayCommand uint `xml:"delaycommand"`
	DelayRetry   uint `xml:"delayretry"`
}
