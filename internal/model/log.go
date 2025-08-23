package model

import (
	"time"

	"gorm.io/gorm"
)

// LogLevel 日志级别枚举
type LogLevel string

const (
	LogLevelDebug LogLevel = "debug" // 调试
	LogLevelInfo  LogLevel = "info"  // 信息
	LogLevelWarn  LogLevel = "warn"  // 警告
	LogLevelError LogLevel = "error" // 错误
	LogLevelFatal LogLevel = "fatal" // 致命错误
)

// LogType 日志类型枚举
type LogType string

const (
	LogTypeUser     LogType = "user"     // 用户操作日志
	LogTypeSystem   LogType = "system"   // 系统日志
	LogTypeSecurity LogType = "security" // 安全日志
	LogTypeAPI      LogType = "api"      // API调用日志
	LogTypeFile     LogType = "file"     // 文件操作日志
	LogTypeAuth     LogType = "auth"     // 认证日志
	LogTypeAdmin    LogType = "admin"    // 管理员操作日志
	LogTypeError    LogType = "error"    // 错误日志
	LogTypeAudit    LogType = "audit"    // 审计日志
)

// actionType 操作类型
type actionType string

const (
	// 用户操作
	ActionLogin         actionType = "login"          // 登录
	ActionLogout        actionType = "logout"         // 登出
	ActionRegister      actionType = "register"       // 注册
	ActionPasswordReset actionType = "password_reset" // 密码重置

	// 文件操作
	ActionFileUpload   actionType = "file_upload"   // 文件上传
	ActionFileDownload actionType = "file_download" // 文件下载
	ActionFileDelete   actionType = "file_delete"   // 文件删除
	ActionFileMove     actionType = "file_move"     // 文件移动
	ActionFileRename   actionType = "file_rename"   // 文件重命名
	ActionFileCopy     actionType = "file_copy"     // 文件复制
	ActionFileShare    actionType = "file_share"    // 文件分享
	ActionFilePreview  actionType = "file_preview"  // 文件预览

	// 文件夹操作
	ActionFolderCreate actionType = "folder_create" // 创建文件夹
	ActionFolderDelete actionType = "folder_delete" // 删除文件夹
	ActionFolderMove   actionType = "folder_move"   // 移动文件夹
	ActionFolderRename actionType = "folder_rename" // 重命名文件夹

	// 权限操作
	ActionPermissionGrant  actionType = "permission_grant"  // 授权
	ActionPermissionRevoke actionType = "permission_revoke" // 撤销权限
	ActionPermissionUpdate actionType = "permission_update" // 更新权限

	// 团队操作
	ActionTeamCreate actionType = "team_create" // 创建团队
	ActionTeamJoin   actionType = "team_join"   // 加入团队
	ActionTeamLeave  actionType = "team_leave"  // 离开团队
	ActionTeamDelete actionType = "team_delete" // 删除团队

	// 系统操作
	ActionSystemStart   actionType = "system_start"   // 系统启动
	ActionSystemStop    actionType = "system_stop"    // 系统停止
	ActionSystemRestart actionType = "system_restart" // 系统重启
	ActionConfigUpdate  actionType = "config_update"  // 配置更新

	// 管理员操作
	ActionAdminUserCreate actionType = "admin_user_create" // 管理员创建用户
	ActionAdminUserUpdate actionType = "admin_user_update" // 管理员更新用户
	ActionAdminUserDelete actionType = "admin_user_delete" // 管理员删除用户
	ActionAdminUserBlock  actionType = "admin_user_block"  // 管理员封禁用户
)

// OperationLog 操作日志模型
type OperationLog struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// 用户信息
	UserID   *uint  `gorm:"index;comment:操作用户ID" json:"user_id"`
	User     *User  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Username string `gorm:"type:varchar(100);index;comment:用户名" json:"username"`

	// 日志基本信息
	Type   LogType    `gorm:"type:varchar(20);not null;index" json:"type"`
	Level  LogLevel   `gorm:"type:varchar(10);not null;index" json:"level"`
	Action actionType `gorm:"type:varchar(50);not null;index" json:"action"`
	Module string     `gorm:"type:varchar(50);index;comment:模块名称" json:"module"`

	// 操作内容
	Title       string `gorm:"type:varchar(200);not null;comment:操作标题" json:"title"`
	Description string `gorm:"type:text;comment:操作描述" json:"description"`
	Content     string `gorm:"type:text;comment:详细内容" json:"content"`

	// 资源信息
	ResourceType string `gorm:"type:varchar(50);index;comment:资源类型" json:"resource_type"`
	ResourceID   *uint  `gorm:"index;comment:资源ID" json:"resource_id"`
	ResourceName string `gorm:"type:varchar(255);comment:资源名称" json:"resource_name"`

	// 操作结果
	Status       string `gorm:"type:varchar(20);not null;index;comment:操作状态" json:"status"`
	Result       string `gorm:"type:text;comment:操作结果" json:"result"`
	ErrorCode    string `gorm:"type:varchar(50);index;comment:错误代码" json:"error_code"`
	ErrorMessage string `gorm:"type:text;comment:错误信息" json:"error_message"`

	// 请求信息
	Method    string `gorm:"type:varchar(10);index;comment:请求方法" json:"method"`
	URL       string `gorm:"type:varchar(500);comment:请求URL" json:"url"`
	IPAddress string `gorm:"type:varchar(45);index;comment:IP地址" json:"ip_address"`
	UserAgent string `gorm:"type:varchar(500);comment:用户代理" json:"user_agent"`
	Referer   string `gorm:"type:varchar(500);comment:来源页面" json:"referer"`

	// 性能信息
	Duration     int64 `gorm:"comment:执行时长(毫秒)" json:"duration"`
	RequestSize  int64 `gorm:"comment:请求大小(字节)" json:"request_size"`
	ResponseSize int64 `gorm:"comment:响应大小(字节)" json:"response_size"`

	// 地理位置信息
	Country string `gorm:"type:varchar(100);comment:国家" json:"country"`
	Region  string `gorm:"type:varchar(100);comment:地区" json:"region"`
	City    string `gorm:"type:varchar(100);comment:城市" json:"city"`

	// 设备信息
	Device  string `gorm:"type:varchar(100);comment:设备类型" json:"device"`
	OS      string `gorm:"type:varchar(100);comment:操作系统" json:"os"`
	Browser string `gorm:"type:varchar(100);comment:浏览器" json:"browser"`

	// 标记信息
	ImportantFlag bool `gorm:"default:false;index;comment:是否重要操作" json:"is_important"`
	AuditedFlag   bool `gorm:"default:false;comment:是否已审计" json:"is_audited"`

	// 扩展信息
	Metadata string `gorm:"type:text;comment:元数据(JSON)" json:"metadata"`
	Tags     string `gorm:"type:varchar(500);comment:标签" json:"tags"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}

// systemLog 系统日志模型 (私有)
type systemLog struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// 日志基本信息
	Level     LogLevel `gorm:"type:varchar(10);not null;index" json:"level"`
	Type      LogType  `gorm:"type:varchar(20);not null;index" json:"type"`
	Module    string   `gorm:"type:varchar(50);not null;index;comment:模块名称" json:"module"`
	Component string   `gorm:"type:varchar(50);index;comment:组件名称" json:"component"`

	// 日志内容
	Title      string `gorm:"type:varchar(200);not null;comment:日志标题" json:"title"`
	Message    string `gorm:"type:text;not null;comment:日志消息" json:"message"`
	StackTrace string `gorm:"type:text;comment:堆栈跟踪" json:"stack_trace"`

	// 错误信息
	ErrorCode string `gorm:"type:varchar(50);index;comment:错误代码" json:"error_code"`
	ErrorType string `gorm:"type:varchar(100);index;comment:错误类型" json:"error_type"`

	// 上下文信息
	RequestID string `gorm:"type:varchar(100);index;comment:请求ID" json:"request_id"`
	SessionID string `gorm:"type:varchar(100);index;comment:会话ID" json:"session_id"`
	TraceID   string `gorm:"type:varchar(100);index;comment:追踪ID" json:"trace_id"`

	// 服务器信息
	Hostname string `gorm:"type:varchar(100);index;comment:主机名" json:"hostname"`
	PID      int    `gorm:"comment:进程ID" json:"pid"`
	ThreadID string `gorm:"type:varchar(50);comment:线程ID" json:"thread_id"`

	// 性能信息
	MemoryUsage int64   `gorm:"comment:内存使用(字节)" json:"memory_usage"`
	CPUUsage    float64 `gorm:"comment:CPU使用率" json:"cpu_usage"`
	DiskUsage   int64   `gorm:"comment:磁盘使用(字节)" json:"disk_usage"`

	// 扩展信息
	Metadata string `gorm:"type:text;comment:元数据(JSON)" json:"metadata"`
	Tags     string `gorm:"type:varchar(500);comment:标签" json:"tags"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 指定表名
func (systemLog) TableName() string {
	return "system_logs"
}

// SystemLog 系统日志模型 (公共类型别名)
type SystemLog = systemLog

// securityLog 安全日志模型 (私有)
type securityLog struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// 用户信息
	UserID   *uint  `gorm:"index;comment:用户ID" json:"user_id"`
	User     *User  `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
	Username string `gorm:"type:varchar(100);index;comment:用户名" json:"username"`

	// 安全事件信息
	EventType string   `gorm:"type:varchar(50);not null;index;comment:事件类型" json:"event_type"`
	Severity  LogLevel `gorm:"type:varchar(10);not null;index;comment:严重程度" json:"severity"`
	Status    string   `gorm:"type:varchar(20);not null;index;comment:事件状态" json:"status"`

	// 事件内容
	Title       string `gorm:"type:varchar(200);not null;comment:事件标题" json:"title"`
	Description string `gorm:"type:text;comment:事件描述" json:"description"`
	Details     string `gorm:"type:text;comment:详细信息" json:"details"`

	// 威胁信息
	ThreatLevel string `gorm:"type:varchar(20);index;comment:威胁级别" json:"threat_level"`
	ThreatType  string `gorm:"type:varchar(50);index;comment:威胁类型" json:"threat_type"`
	AttackType  string `gorm:"type:varchar(50);index;comment:攻击类型" json:"attack_type"`

	// 网络信息
	SourceIP string `gorm:"type:varchar(45);index;comment:源IP地址" json:"source_ip"`
	TargetIP string `gorm:"type:varchar(45);index;comment:目标IP地址" json:"target_ip"`
	Port     int    `gorm:"comment:端口号" json:"port"`
	Protocol string `gorm:"type:varchar(20);comment:协议" json:"protocol"`

	// 地理位置
	Country string `gorm:"type:varchar(100);comment:国家" json:"country"`
	Region  string `gorm:"type:varchar(100);comment:地区" json:"region"`
	City    string `gorm:"type:varchar(100);comment:城市" json:"city"`

	// 处理信息
	BlockedFlag  bool       `gorm:"default:false;index;comment:是否已阻止" json:"is_blocked"`
	ResolvedFlag bool       `gorm:"default:false;index;comment:是否已解决" json:"is_resolved"`
	ResolvedBy   *uint      `gorm:"index;comment:解决人ID" json:"resolved_by"`
	Resolver     *User      `gorm:"foreignKey:ResolvedBy;constraint:OnDelete:SET NULL" json:"resolver,omitempty"`
	ResolvedAt   *time.Time `gorm:"comment:解决时间" json:"resolved_at"`

	// 扩展信息
	Metadata string `gorm:"type:text;comment:元数据(JSON)" json:"metadata"`
	Tags     string `gorm:"type:varchar(500);comment:标签" json:"tags"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (securityLog) TableName() string {
	return "security_logs"
}

// SecurityLog 安全日志模型 (公共类型别名)
type SecurityLog = securityLog

// BeforeCreate GORM钩子：创建前
func (ol *OperationLog) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if ol.Level == "" {
		ol.Level = LogLevelInfo
	}
	if ol.Status == "" {
		ol.Status = "success"
	}
	return nil
}

// BeforeCreate GORM钩子：创建前
func (sl *systemLog) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if sl.Level == "" {
		sl.Level = LogLevelInfo
	}
	return nil
}

// BeforeCreate GORM钩子：创建前
func (sl *securityLog) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if sl.Severity == "" {
		sl.Severity = LogLevelWarn
	}
	if sl.Status == "" {
		sl.Status = "detected"
	}
	return nil
}

// IsSuccess 检查操作是否成功
func (ol *OperationLog) IsSuccess() bool {
	return ol.Status == "success"
}

// IsError 检查是否为错误日志
func (ol *OperationLog) IsError() bool {
	return ol.Level == LogLevelError || ol.Level == LogLevelFatal
}

// IsUserOperation 检查是否为用户操作
func (ol *OperationLog) IsUserOperation() bool {
	return ol.Type == LogTypeUser
}

// IsSystemOperation 检查是否为系统操作
func (ol *OperationLog) IsSystemOperation() bool {
	return ol.Type == LogTypeSystem
}

// IsSecurityOperation 检查是否为安全操作
func (ol *OperationLog) IsSecurityOperation() bool {
	return ol.Type == LogTypeSecurity
}

// IsError 检查是否为错误日志
func (sl *systemLog) IsError() bool {
	return sl.Level == LogLevelError || sl.Level == LogLevelFatal
}

// IsWarning 检查是否为警告日志
func (sl *systemLog) IsWarning() bool {
	return sl.Level == LogLevelWarn
}

// IsBlocked 检查是否为阻止操作
func (sl *securityLog) IsBlocked() bool {
	return sl.BlockedFlag
}

// IsResolved 检查安全事件是否已解决
func (sl *SecurityLog) IsResolved() bool {
	return sl.ResolvedFlag
}

// IsHighThreat 检查是否为高威胁
func (sl *securityLog) IsHighThreat() bool {
	return sl.ThreatLevel == "high" || sl.ThreatLevel == "critical"
}

// IsCritical 检查是否为严重安全事件
func (sl *SecurityLog) IsCritical() bool {
	return sl.Severity == LogLevelError || sl.Severity == LogLevelFatal
}

// IsImportant 检查是否为重要操作
func (ol *OperationLog) IsImportant() bool {
	return ol.ImportantFlag
}

// IsAudited 检查是否已审计
func (ol *OperationLog) IsAudited() bool {
	return ol.AuditedFlag
}
