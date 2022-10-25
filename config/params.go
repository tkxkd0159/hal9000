package config

const (
	ActOracle     = "oracle"
	ActStake      = "stake"
	ActAutoStake  = "restake"
	ActWithdraw   = "withdraw"
	ActUndelegate = "undelegate"

	ActPusher = "push"
	ActMinter = "mint"
	ActTaker  = "take"
)

const (
	ControlChain = "nova"
	Gas          = "auto"
	NovaGasPrice = "0unova"
	GasWeight    = 1.1
)

const (
	StdLogFile            = "ctxlog.txt"
	LocalErrlogFile       = "nova_err.txt"
	ExtRedirectErrlogFile = "other_err.txt"
)

// related to Viper
const (
	NovaBotAddrKey = "bot_addr"
)
