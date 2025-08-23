package model

import (
	"time"

	"gorm.io/gorm"
)

// configType 配置类型枚举
type configType string

const (
	ConfigTypeSystem   configType = "system"   // 系统配置
	ConfigTypeStorage  configType = "storage"  // 存储配置
	ConfigTypeSecurity configType = "security" // 安全配置
	ConfigTypeEmail    configType = "email"    // 邮件配置
	ConfigTypeIM       configType = "im"       // 即时通讯配置
	ConfigTypeUpload   configType = "upload"   // 上传配置
	ConfigTypePreview  configType = "preview"  // 预览配置
	ConfigTypeBackup   configType = "backup"   // 备份配置
	ConfigTypeMonitor  configType = "monitor"  // 监控配置
	ConfigTypeCustom   configType = "custom"   // 自定义配置
)

// ConfigStatus 配置状态枚举
type ConfigStatus string

const (
	ConfigStatusActive   ConfigStatus = "active"   // 激活
	ConfigStatusInactive ConfigStatus = "inactive" // 未激活
	ConfigStatusTesting  ConfigStatus = "testing"  // 测试中
	ConfigStatusError    ConfigStatus = "error"    // 错误
)

// SystemConfig 系统配置模型
type SystemConfig struct {
	// 时间戳 (24 bytes each)
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 指针字段 (8 bytes each)
	MinValue  *float64 `gorm:"comment:最小值" json:"min_value"`
	MaxValue  *float64 `gorm:"comment:最大值" json:"max_value"`
	UpdatedBy *uint    `gorm:"index;comment:更新人ID" json:"updated_by"`
	Updater   *User    `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`

	// 结构体字段 (size varies)
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`

	// 字符串字段 (24 bytes each)
	Key            string       `gorm:"type:varchar(100);not null;uniqueIndex;comment:配置键" json:"key"`
	Value          string       `gorm:"type:text;not null;comment:配置值" json:"value"`
	DefaultValue   string       `gorm:"type:text;comment:默认值" json:"default_value"`
	Name           string       `gorm:"type:varchar(200);not null;comment:配置名称" json:"name"`
	Description    string       `gorm:"type:varchar(500);comment:配置描述" json:"description"`
	Group          string       `gorm:"type:varchar(100);index;comment:配置分组" json:"group"`
	DataType       string       `gorm:"type:varchar(20);default:'string';comment:数据类型" json:"data_type"`
	ValidationRule string       `gorm:"type:varchar(500);comment:验证规则" json:"validation_rule"`
	Options        string       `gorm:"type:text;comment:可选值(JSON)" json:"options"`
	Type           configType   `gorm:"type:varchar(20);not null;index" json:"type"`
	Status         ConfigStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`

	// uint字段 (8 bytes each)
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy uint `gorm:"not null;index;comment:创建人ID" json:"created_by"`

	// int字段 (8 bytes)
	SortOrder int `gorm:"default:0;comment:排序" json:"sort_order"`

	// bool字段 (1 byte each)
	IsRequired   bool `gorm:"default:false;comment:是否必需" json:"is_required"`
	SecretFlag   bool `gorm:"default:false;comment:是否敏感信息" json:"is_secret"`
	ReadonlyFlag bool `gorm:"default:false;comment:是否只读" json:"is_readonly"`
	IsSystem     bool `gorm:"default:false;comment:是否系统配置" json:"is_system"`
	IsVisible    bool `gorm:"default:true;comment:是否可见" json:"is_visible"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}

// StorageProvider 存储提供商枚举
type StorageProvider string

const (
	StorageProviderLocal      StorageProvider = "local"      // 本地存储
	StorageProviderAliOSS     StorageProvider = "alioss"     // 阿里云OSS
	StorageProviderTencentCOS StorageProvider = "tencentcos" // 腾讯云COS
	StorageProviderQiniuKodo  StorageProvider = "qiniukodo"  // 七牛云Kodo
	StorageProviderAWSS3      StorageProvider = "awss3"      // AWS S3
	StorageProviderMinIO      StorageProvider = "minio"      // MinIO
)

// storageConfig 存储配置模型 (私有)
type storageConfig struct {
	// 时间戳 (24 bytes each)
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 指针字段 (8 bytes each)
	UpdatedBy *uint `gorm:"index;comment:更新人ID" json:"updated_by"`
	Updater   *User `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`

	// 结构体字段
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`

	// int64字段 (8 bytes each)
	MaxFileSize int64 `gorm:"default:104857600;comment:最大文件大小(字节,默认100MB)" json:"max_file_size"`
	TotalFiles  int64 `gorm:"default:0;comment:总文件数" json:"total_files"`
	TotalSize   int64 `gorm:"default:0;comment:总大小(字节)" json:"total_size"`
	UsedSize    int64 `gorm:"default:0;comment:已使用大小(字节)" json:"used_size"`
	QuotaSize   int64 `gorm:"default:0;comment:配额大小(字节,0表示无限制)" json:"quota_size"`

	// 字符串字段 (24 bytes each)
	Name           string          `gorm:"type:varchar(100);not null;uniqueIndex;comment:配置名称" json:"name"`
	Endpoint       string          `gorm:"type:varchar(200);comment:服务端点" json:"endpoint"`
	Region         string          `gorm:"type:varchar(50);comment:区域" json:"region"`
	Bucket         string          `gorm:"type:varchar(100);comment:存储桶名称" json:"bucket"`
	AccessKey      string          `gorm:"type:varchar(200);comment:访问密钥" json:"access_key"`
	SecretKey      string          `gorm:"type:varchar(500);comment:密钥" json:"secret_key"`
	BasePath       string          `gorm:"type:varchar(200);default:'/';comment:基础路径" json:"base_path"`
	Domain         string          `gorm:"type:varchar(200);comment:自定义域名" json:"domain"`
	AllowedTypes   string          `gorm:"type:text;comment:允许的文件类型(JSON)" json:"allowed_types"`
	EncryptionKey  string          `gorm:"type:varchar(500);comment:加密密钥" json:"encryption_key"`
	BackupProvider string          `gorm:"type:varchar(20);comment:备份提供商" json:"backup_provider"`
	BackupConfig   string          `gorm:"type:text;comment:备份配置(JSON)" json:"backup_config"`
	Provider       StorageProvider `gorm:"type:varchar(20);not null;index" json:"provider"`
	Status         ConfigStatus    `gorm:"type:varchar(20);default:'active';index" json:"status"`

	// uint字段 (8 bytes each)
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedBy uint `gorm:"not null;index;comment:创建人ID" json:"created_by"`

	// int字段 (8 bytes each)
	ChunkSize       int `gorm:"default:5242880;comment:分片大小(字节,默认5MB)" json:"chunk_size"`
	SignatureExpiry int `gorm:"default:3600;comment:签名过期时间(秒)" json:"signature_expiry"`
	CacheExpiry     int `gorm:"default:86400;comment:缓存过期时间(秒)" json:"cache_expiry"`
	MonitorInterval int `gorm:"default:300;comment:监控间隔(秒)" json:"monitor_interval"`

	// bool字段 (1 byte each)
	DefaultFlag      bool `gorm:"default:false;index;comment:是否默认存储" json:"is_default"`
	EnabledFlag      bool `gorm:"default:true;comment:是否启用" json:"is_enabled"`
	IsHTTPS          bool `gorm:"default:true;comment:是否使用HTTPS" json:"is_https"`
	EnableChunk      bool `gorm:"default:true;comment:是否启用分片上传" json:"enable_chunk"`
	EnableEncryption bool `gorm:"default:false;comment:是否启用加密" json:"enable_encryption"`
	EnableSignature  bool `gorm:"default:true;comment:是否启用签名" json:"enable_signature"`
	EnableCache      bool `gorm:"default:true;comment:是否启用缓存" json:"enable_cache"`
	EnableBackup     bool `gorm:"default:false;comment:是否启用备份" json:"enable_backup"`
	EnableMonitor    bool `gorm:"default:true;comment:是否启用监控" json:"enable_monitor"`
}

// TableName 指定表名
func (storageConfig) TableName() string {
	return "storage_configs"
}

// StorageConfig 存储配置模型 (公共类型别名)
type StorageConfig = storageConfig

// configHistory 配置历史模型 (私有)
type configHistory struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	ConfigType configType `gorm:"type:varchar(20);not null;index;comment:配置类型" json:"config_type"`
	ConfigID   uint       `gorm:"not null;index;comment:配置ID" json:"config_id"`
	ConfigKey  string     `gorm:"type:varchar(100);not null;index;comment:配置键" json:"config_key"`

	// 变更信息
	Action       string `gorm:"type:varchar(20);not null;index;comment:操作类型" json:"action"`
	OldValue     string `gorm:"type:text;comment:旧值" json:"old_value"`
	NewValue     string `gorm:"type:text;comment:新值" json:"new_value"`
	ChangeReason string `gorm:"type:varchar(500);comment:变更原因" json:"change_reason"`

	// 操作信息
	OperatorID uint   `gorm:"not null;index;comment:操作人ID" json:"operator_id"`
	Operator   User   `gorm:"foreignKey:OperatorID;constraint:OnDelete:RESTRICT" json:"operator,omitempty"`
	IPAddress  string `gorm:"type:varchar(45);index" json:"ip_address"`
	UserAgent  string `gorm:"type:varchar(500)" json:"user_agent"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 指定表名
func (configHistory) TableName() string {
	return "config_histories"
}

// configHistory 配置历史模型现在是私有的，通过方法访问

// NewConfigHistory 创建新的配置历史记录
func NewConfigHistory() *configHistory {
	return &configHistory{}
}

// ConfigHistoryQuery 获取配置历史（用于查询）
type ConfigHistoryQuery = configHistory

// BeforeCreate GORM钩子：创建前
func (sc *SystemConfig) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if sc.Status == "" {
		sc.Status = ConfigStatusActive
	}
	if sc.DataType == "" {
		sc.DataType = "string"
	}
	return nil
}

// BeforeCreate GORM钩子：创建前
func (stc *storageConfig) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if stc.Status == "" {
		stc.Status = ConfigStatusActive
	}
	if stc.BasePath == "" {
		stc.BasePath = "/"
	}
	if stc.MaxFileSize == 0 {
		stc.MaxFileSize = 104857600 // 100MB
	}
	if stc.ChunkSize == 0 {
		stc.ChunkSize = 5242880 // 5MB
	}
	if stc.SignatureExpiry == 0 {
		stc.SignatureExpiry = 3600 // 1小时
	}
	if stc.CacheExpiry == 0 {
		stc.CacheExpiry = 86400 // 24小时
	}
	if stc.MonitorInterval == 0 {
		stc.MonitorInterval = 300 // 5分钟
	}
	return nil
}

// IsActive 检查配置是否激活
func (sc *SystemConfig) IsActive() bool {
	return sc.Status == ConfigStatusActive
}

// IsSecret 检查是否为敏感配置
func (sc *SystemConfig) IsSecret() bool {
	return sc.SecretFlag
}

// IsReadonly 检查是否为只读配置
func (sc *SystemConfig) IsReadonly() bool {
	return sc.ReadonlyFlag
}

// IsSystemConfig 检查是否为系统配置
func (sc *SystemConfig) IsSystemConfig() bool {
	return sc.IsSystem
}

// IsActive 检查存储配置是否激活
func (stc *storageConfig) IsActive() bool {
	return stc.Status == ConfigStatusActive
}

// IsEnabled 检查存储配置是否启用
func (stc *storageConfig) IsEnabled() bool {
	return stc.EnabledFlag
}

// IsDefault 检查是否为默认存储
func (stc *storageConfig) IsDefault() bool {
	return stc.DefaultFlag
}

// IsLocal 检查是否为本地存储
func (stc *storageConfig) IsLocal() bool {
	return stc.Provider == StorageProviderLocal
}

// IsCloudStorage 检查是否为云存储
func (stc *storageConfig) IsCloudStorage() bool {
	return stc.Provider != StorageProviderLocal
}

// GetUsagePercent 获取存储使用百分比
func (stc *storageConfig) GetUsagePercent() float64 {
	if stc.QuotaSize == 0 {
		return 0
	}
	return float64(stc.UsedSize) / float64(stc.QuotaSize) * 100
}

// IsQuotaExceeded 检查是否超出配额
func (stc *storageConfig) IsQuotaExceeded() bool {
	return stc.QuotaSize > 0 && stc.UsedSize >= stc.QuotaSize
}

// CanUpload 检查是否可以上传文件
func (stc *storageConfig) CanUpload(fileSize int64) bool {
	if !stc.IsEnabled() || !stc.IsActive() {
		return false
	}
	if fileSize > stc.MaxFileSize {
		return false
	}
	if stc.QuotaSize > 0 && (stc.UsedSize+fileSize) > stc.QuotaSize {
		return false
	}
	return true
}

// Config 应用配置结构体 (公共类型别名)
type Config = config

// config 应用配置结构体 (私有)
type config struct {
	App       appConfig       `json:"app" yaml:"app"`
	Server    serverConfig    `json:"server" yaml:"server"`
	Database  databaseConfig  `json:"database" yaml:"database"`
	Redis     redisConfig     `json:"redis" yaml:"redis"`
	JWT       jwtConfig       `json:"jwt" yaml:"jwt"`
	Log       logConfig       `json:"log" yaml:"log"`
	Upload    uploadConfig    `json:"upload" yaml:"upload"`
	CORS      corsConfig      `json:"cors" yaml:"cors"`
	RateLimit rateLimitConfig `json:"rate_limit" yaml:"rate_limit"`
	Cache     cacheConfig     `json:"cache" yaml:"cache"`
}

// appConfig 应用配置 (私有)
type appConfig struct {
	Name     string `json:"name" yaml:"name"`
	Version  string `json:"version" yaml:"version"`
	Env      string `json:"env" yaml:"env"`
	Debug    bool   `json:"debug" yaml:"debug"`
	Timezone string `json:"timezone" yaml:"timezone"`
}

// serverConfig 服务器配置 (私有)
type serverConfig struct {
	Host           string        `json:"host" yaml:"host"`
	Port           int           `json:"port" yaml:"port"`
	Mode           string        `json:"mode" yaml:"mode"`
	ReadTimeout    time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout   time.Duration `json:"write_timeout" yaml:"write_timeout"`
	MaxHeaderBytes int           `json:"max_header_bytes" yaml:"max_header_bytes"`
}

// databaseConfig 数据库配置 (私有)
type databaseConfig struct {
	Driver          string        `json:"driver" yaml:"driver"`
	Host            string        `json:"host" yaml:"host"`
	Port            int           `json:"port" yaml:"port"`
	Username        string        `json:"username" yaml:"username"`
	Password        string        `json:"password" yaml:"password"`
	DBName          string        `json:"dbname" yaml:"dbname"`
	Charset         string        `json:"charset" yaml:"charset"`
	ParseTime       bool          `json:"parse_time" yaml:"parse_time"`
	Loc             string        `json:"loc" yaml:"loc"`
	MaxIdleConns    int           `json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int           `json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	LogLevel        string        `json:"log_level" yaml:"log_level"`
}

// redisConfig Redis配置 (私有)
type redisConfig struct {
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	Password     string        `json:"password" yaml:"password"`
	DB           int           `json:"db" yaml:"db"`
	PoolSize     int           `json:"pool_size" yaml:"pool_size"`
	MinIdleConns int           `json:"min_idle_conns" yaml:"min_idle_conns"`
	DialTimeout  time.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout" yaml:"write_timeout"`
	PoolTimeout  time.Duration `json:"pool_timeout" yaml:"pool_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout" yaml:"idle_timeout"`
}

// jwtConfig JWT配置 (私有)
type jwtConfig struct {
	Secret            string        `json:"secret" yaml:"secret"`
	ExpireTime        time.Duration `json:"expire_time" yaml:"expire_time"`
	RefreshExpireTime time.Duration `json:"refresh_expire_time" yaml:"refresh_expire_time"`
	Issuer            string        `json:"issuer" yaml:"issuer"`
}

// logConfig 日志配置 (私有)
type logConfig struct {
	Level      string `json:"level" yaml:"level"`
	Format     string `json:"format" yaml:"format"`
	Output     string `json:"output" yaml:"output"`
	FilePath   string `json:"file_path" yaml:"file_path"`
	MaxSize    int    `json:"max_size" yaml:"max_size"`
	MaxAge     int    `json:"max_age" yaml:"max_age"`
	MaxBackups int    `json:"max_backups" yaml:"max_backups"`
	Compress   bool   `json:"compress" yaml:"compress"`
}

// uploadConfig 上传配置 (私有)
type uploadConfig struct {
	MaxSize      int64    `json:"max_size" yaml:"max_size"`
	AllowedTypes []string `json:"allowed_types" yaml:"allowed_types"`
	UploadPath   string   `json:"upload_path" yaml:"upload_path"`
	URLPrefix    string   `json:"url_prefix" yaml:"url_prefix"`
}

// corsConfig CORS配置 (私有)
type corsConfig struct {
	AllowOrigins     []string `json:"allow_origins" yaml:"allow_origins"`
	AllowMethods     []string `json:"allow_methods" yaml:"allow_methods"`
	AllowHeaders     []string `json:"allow_headers" yaml:"allow_headers"`
	ExposeHeaders    []string `json:"expose_headers" yaml:"expose_headers"`
	AllowCredentials bool     `json:"allow_credentials" yaml:"allow_credentials"`
	MaxAge           int      `json:"max_age" yaml:"max_age"`
}

// rateLimitConfig 限流配置 (私有)
type rateLimitConfig struct {
	Enabled           bool `json:"enabled" yaml:"enabled"`
	RequestsPerMinute int  `json:"requests_per_minute" yaml:"requests_per_minute"`
	Burst             int  `json:"burst" yaml:"burst"`
}

// cacheConfig 缓存配置 (私有)
type cacheConfig struct {
	DefaultExpiration time.Duration `json:"default_expiration"`
	CleanupInterval   time.Duration `json:"cleanup_interval"`
}

// cacheConfig 缓存配置现在是私有的，通过Config结构体访问

// GetCacheConfig 从Config中获取缓存配置
func (c *config) GetCacheConfig() cacheConfig {
	return c.Cache
}
