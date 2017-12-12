package alipay

// Alipay alipay struct
type Alipay struct {
	Config *Config
}

// Config alipay config
type Config struct {
	APPID          string `required:"true"`
	PrivateKey     string `required:"true"`
	PublicKey      string `required:"true"`
	ProductionMode bool
}
