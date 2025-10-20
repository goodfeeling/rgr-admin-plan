package config

type ConfigResponse struct {
	Data Config `json:"data"`
}
type Config struct {
	Site          SiteConfig          `yaml:"site" json:"site"`
	Server        ServerConfig        `yaml:"server" json:"server"`
	Postgres      PostgresConfig      `yaml:"postgres" json:"postgres"`
	JWT           JWTConfig           `yaml:"jwt" json:"jwt"`
	Redis         RedisConfig         `yaml:"redis" json:"redis"`
	CORS          CORSConfig          `yaml:"cors" json:"cors"`
	AliyunOSS     AliyunOSSConfig     `yaml:"aliyun_oss" json:"aliyun_oss"`
	RabbitMQ      RabbitMQConfig      `yaml:"rabbitmq" json:"rabbitmq"`
	Zap           ZapConfig           `yaml:"zap" json:"zap"`
	NativeStorage NativeStorageConfig `yaml:"native_storage" json:"native_storage"`
	Email         EmailConfig         `yaml:"email" json:"email"`
	Captcha       CaptchaConfig       `yaml:"captcha" json:"captcha"`
}

type SiteConfig struct {
	Name    string `json:"name" yaml:"name"`
	Logo    string `json:"logo" yaml:"logo"`
	Favicon string `json:"favicon" yaml:"favicon"`
	Login   string `json:"login_img" yaml:"login"`
}

type ServerConfig struct {
	FrontendUrl   string `yaml:"frontend_url" json:"frontend_url"`
	Port          int    `yaml:"port" json:"port"`
	StorageEngine string `yaml:"storage" json:"storage_engine"`
	EventBus      string `yaml:"event_bus" json:"event_bus"`
	Database      string `yaml:"database" json:"database"`
	LimitRate     int    `yaml:"limit_rate" json:"limit_rate"`
	LimitTime     int    `yaml:"limit_time" json:"limit_time"`
	SingleSignOn  bool   `yaml:"single_sign_on" json:"single_sign_on"`
}

type PostgresConfig struct {
	Host            string `yaml:"host" json:"host"`
	Port            int    `yaml:"port" json:"port"`
	User            string `yaml:"user" json:"user"`
	Password        string `yaml:"password" json:"password"`
	Name            string `yaml:"name" json:"name"`
	SSLMode         string `yaml:"sslmode" json:"sslmode"`
	MaxIdleConns    int    `yaml:"max_idle_conns" json:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns" json:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime" json:"conn_max_lifetime"`
}

type JWTConfig struct {
	AccessSecret     string `yaml:"access_secret" json:"access_secret"`
	AccessTimeMinute int    `yaml:"access_time_minute" json:"access_time_minute"`
	RefreshSecret    string `yaml:"refresh_secret" json:"refresh_secret"`
	RefreshTimeHour  int    `yaml:"refresh_time_hour" json:"refresh_time_hour"`
	ResetSecret      string `yaml:"reset_secret" json:"reset_secret"`
}

type RedisConfig struct {
	Host     string `yaml:"host" json:"host"`
	Port     int    `yaml:"port" json:"port"`
	Password string `yaml:"password" json:"password"`
	DB       int    `yaml:"db" json:"db"`
	PoolSize int    `yaml:"pool_size" json:"pool_size"`
}

type CORSConfig struct {
	AllowedOrigins string `yaml:"allowed_origins" json:"allowed_origins"`
	AllowedMethods string `yaml:"allowed_methods" json:"allowed_methods"`
	AllowedHeaders string `yaml:"allowed_headers" json:"allowed_headers"`
}

type AliyunOSSConfig struct {
	AccessKeyID            string `yaml:"access_key_id" json:"access_key_id"`
	AccessKeySecret        string `yaml:"access_key_secret" json:"access_key_secret"`
	BucketName             string `yaml:"bucket_name" json:"bucket_name"`
	RAMRoleARN             string `yaml:"ram_role_arn" json:"ram_role_arn"`
	SecurityRegionID       string `yaml:"security_region_id" json:"security_region_id"`
	SecurityServiceAddress string `yaml:"security_service_address" json:"security_service_address"`
	BaseURL                string `yaml:"base_url" json:"base_url"`
}

type RabbitMQConfig struct {
	URL      string `yaml:"url" json:"url"`
	Exchange string `yaml:"exchange" json:"exchange"`
	Queue    string `yaml:"queue" json:"queue"`
}

type CaptchaConfig struct {
	Enable bool `yaml:"enable" json:"enable"`

	Complexity    int `yaml:"complexity" json:"complexity"`
	Width         int `yaml:"width" json:"width"`
	Height        int `yaml:"height" json:"height"`
	Length        int `yaml:"length" json:"length"`
	TimeoutMinute int `yaml:"timeout_minute" json:"timeout_minute"`
}
type NativeStorageConfig struct {
	BaseURL    string `yaml:"base_url" json:"base_url"`
	AccessPath string `yaml:"access_dir" json:"access_path"`
	UploadDir  string `yaml:"upload_dir" json:"upload_dir"`
}

type ZapConfig struct {
	EncodeLevel   string `yaml:"encode_level" json:"encode_level"`
	Dirname       string `yaml:"dirname" json:"dirname"`
	MaxAge        int    `yaml:"max_age" json:"max_age"`
	StacktraceKey string `yaml:"stacktrace_key" json:"stacktrace_key"`
	Level         string `yaml:"level" json:"level"`
	LogInConsole  bool   `yaml:"log_in_console" json:"log_in_console"`
	Encoding      string `yaml:"encoding" json:"encoding"`
}
type EmailConfig struct {
	From     string `yaml:"from"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	SSL      bool   `yaml:"ssl"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
type IConfigService interface {
	GetConfig() (*ConfigResponse, error)
	Update(module string, dataMap map[string]interface{}) error
	GetConfigByModule(module string) (*map[string]string, error)
}
