package globals

type config struct {
	PostgresURI              string
	PostgresMaxConnections   int64
	AESKey                   string
	GRPCAPIKey               string `envconf:"optional"`
	GRPCServerPort           uint16
	AuthenticationServerPort uint16
	SecureServerHost         string
	SecureServerPort         uint16
	AccountGRPCHost          string
	AccountGRPCPort          uint16
	AccountGRPCAPIKey        string `envconf:"optional"`
	HealthCheckPort          uint16 `envconf:"optional"`
	EnableBella              bool   `envconf:"optional"`
	MiiDecryptKey            string
	PIDHmacKey               string
}

var Config *config = &config{}
