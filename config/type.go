package config

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type redis struct {
	Addr     string
	Password string
	DB       int
}

type oss struct {
	Endpoint        string
	AccessKeyID     string `mapstructure:"accessKey-id"`
	AccessKeySecret string `mapstructure:"accessKey-secret"`
	BucketName      string
	Zone            string `mapstructure:"zone"`
	UseCdnDomains   bool   `mapstructure:"use_cdn_domains"`
}

// SMTP 配置
type smtp struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	FromName string `mapstructure:"from_name"`
}

// Verify（验证码）配置
type verify struct {
	CodeLength int `mapstructure:"code_length"`
	TTLSeconds int `mapstructure:"ttl_seconds"`
}

type server struct {
	Addr string
	Port int
}

type turnstile struct {
	SecretKey string `mapstructure:"secret_key"`
	Enabled   bool
}

type logger struct {
	Level string
	Dir   string
	Env   string // development, testing, production
}

type cors struct {
	AllowOrigins  []string `mapstructure:"allow_origins"`
	AllowMethods  []string `mapstructure:"allow_methods"`
	AllowHeaders  []string `mapstructure:"allow_headers"`
	ExposeHeaders []string `mapstructure:"expose_headers"`
	MaxAge        int      `mapstructure:"max_age"`
	AllowWildcard bool     `mapstructure:"allow_wildcard"`
}

type config struct {
	MySQL     mySQL
	Redis     redis
	OSS       oss
	Smtp      smtp   `mapstructure:"smtp"`
	Verify    verify `mapstructure:"verify"`
	Server    server
	Turnstile turnstile `mapstructure:"turnstile"`
	Logger    logger    `mapstructure:"logger"`
	Cors      cors      `mapstructure:"cors"`
}
