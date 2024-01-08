package config

const (
	AppName    = "Fun Banking"
	AppVersion = "1.0.0"
	AdminRole  = 10
)

// Banking constants
const (
	MAX_BANKING_TRANSFER_AMOUNT = 250_000_000
)

var JwtKey []byte
var AppBaseURL string
