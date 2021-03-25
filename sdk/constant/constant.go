package constant

const (
	//k8s模板的标签
	Tag = "{{%s}}"
	//私钥后缀
	PriKeySuf = "_sk"
	//msp后缀
	MspSuf = "MSP"
)
const (
	//cryptogen工具生成证书目录
	CryptoConfigDir = "crypto-config"
	//configtxgen工具生成创世区块,channel交易保存的目录
	ChannelArtifactsDir = "channel-artifacts"
	//cryptogen配置
	CryptoConfigYaml = CryptoConfigDir + ".yaml"
	//configtxgen配置
	ConfigtxYaml = "configtx.yaml"
)
const (
	OrdererSuffix      = "orderer"
	OrdererMsp         = "OrdererMSP"
	OrdererSolo        = "solo"
	OrdererKafka       = "kafka"
	OrdererEtcdraft    = "etcdraft"
	KafkaSuffix        = "kafka"
	TypeImplicitMeta   = "ImplicitMeta"
	TypeSignature      = "Signature"
	RuleAnyReaders     = "ANY Readers"
	RuleAnyWriters     = "ANY Writers"
	RuleMajorityAdmins = "MAJORITY Admins"
	Country            = "CN"
	Province           = "GuangDong"
	Locality           = "GuangZhou"
)
