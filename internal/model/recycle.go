package model

import (
	"time"

	"gorm.io/gorm"
)

// RecycleStatus 回收站状态枚举
type RecycleStatus string

const (
	RecycleStatusDeleted   RecycleStatus = "deleted"   // 已删除
	RecycleStatusRestored  RecycleStatus = "restored"  // 已恢复
	RecycleStatusPermanent RecycleStatus = "permanent" // 永久删除
)

// RecycleType 回收类型枚举
type RecycleType string

const (
	RecycleTypeFile   RecycleType = "file"   // 文件
	RecycleTypeFolder RecycleType = "folder" // 文件夹
)

// RecycleItem 回收站项目模型
type RecycleItem struct {
	// 时间戳字段 (24 bytes each)
	DeletedAt time.Time `gorm:"not null;index;comment:删除时间" json:"deleted_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 指针字段 (8 bytes each)
	OriginalParentID   *uint      `gorm:"index;comment:原父目录ID" json:"original_parent_id"`
	OriginalParent     *File      `gorm:"foreignKey:OriginalParentID" json:"original_parent,omitempty"`
	RestoredAt         *time.Time `gorm:"index;comment:恢复时间" json:"restored_at"`
	RestoredBy         *uint      `gorm:"index;comment:恢复操作人ID" json:"restored_by"`
	Restorer           *User      `gorm:"foreignKey:RestoredBy" json:"restorer,omitempty"`
	PermanentDeletedAt *time.Time `gorm:"index;comment:永久删除时间" json:"permanent_deleted_at"`
	PermanentDeletedBy *uint      `gorm:"index;comment:永久删除操作人ID" json:"permanent_deleted_by"`
	PermanentDeleter   *User      `gorm:"foreignKey:PermanentDeletedBy" json:"permanent_deleter,omitempty"`
	ExpiresAt          *time.Time `gorm:"index;comment:过期时间" json:"expires_at"`

	// 结构体字段
	User         User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	OriginalFile File `gorm:"foreignKey:OriginalFileID;constraint:OnDelete:CASCADE" json:"original_file,omitempty"`
	Deleter      User `gorm:"foreignKey:DeletedBy;constraint:OnDelete:RESTRICT" json:"deleter,omitempty"`

	// 字符串字段 (24 bytes each)
	OriginalPath    string        `gorm:"type:varchar(1000);not null;comment:原文件路径" json:"original_path"`
	FileName        string        `gorm:"type:varchar(255);not null;index" json:"file_name"`
	FileType        string        `gorm:"type:varchar(100);index" json:"file_type"`
	MimeType        string        `gorm:"type:varchar(200)" json:"mime_type"`
	DeletedReason   string        `gorm:"type:varchar(500);comment:删除原因" json:"deleted_reason"`
	RestoredPath    string        `gorm:"type:varchar(1000);comment:恢复后路径" json:"restored_path"`
	StoragePath     string        `gorm:"type:varchar(1000);comment:存储路径" json:"storage_path"`
	StorageProvider string        `gorm:"type:varchar(50);comment:存储提供商" json:"storage_provider"`
	Metadata        string        `gorm:"type:text;comment:文件元数据(JSON)" json:"metadata"`
	Tags            string        `gorm:"type:text;comment:文件标签(JSON)" json:"tags"`
	Type            RecycleType   `gorm:"type:varchar(20);not null;index" json:"type"`
	Status          RecycleStatus `gorm:"type:varchar(20);default:'deleted';index" json:"status"`

	// int64字段 (8 bytes each)
	FileSize int64 `gorm:"not null;comment:文件大小(字节)" json:"file_size"`

	// uint字段 (8 bytes each)
	ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID         uint `gorm:"not null;index" json:"user_id"`
	OriginalFileID uint `gorm:"not null;index;comment:原文件ID" json:"original_file_id"`
	DeletedBy      uint `gorm:"not null;index;comment:删除操作人ID" json:"deleted_by"`

	// int字段 (8 bytes each)
	AutoDeleteDays int `gorm:"default:30;comment:自动删除天数" json:"auto_delete_days"`

	// bool字段 (1 byte each)
	IsEncrypted bool `gorm:"default:false;comment:是否加密" json:"is_encrypted"`
}

// TableName 指定表名
func (RecycleItem) TableName() string {
	return "recycle_items"
}

// RecycleBin 回收站配置模型
type RecycleBin struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// 用户信息
	UserID uint `gorm:"not null;uniqueIndex" json:"user_id"`
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// 配置信息
	IsEnabled          bool  `gorm:"default:true;comment:是否启用回收站" json:"is_enabled"`
	AutoDeleteDays     int   `gorm:"default:30;comment:自动删除天数" json:"auto_delete_days"`
	MaxStorageSize     int64 `gorm:"default:1073741824;comment:最大存储大小(字节,默认1GB)" json:"max_storage_size"`
	CurrentStorageSize int64 `gorm:"default:0;comment:当前存储大小(字节)" json:"current_storage_size"`
	MaxItemCount       int   `gorm:"default:1000;comment:最大项目数量" json:"max_item_count"`
	CurrentItemCount   int   `gorm:"default:0;comment:当前项目数量" json:"current_item_count"`

	// 通知设置
	NotifyBeforeDelete bool `gorm:"default:true;comment:删除前通知" json:"notify_before_delete"`
	NotifyDays         int  `gorm:"default:7;comment:提前通知天数" json:"notify_days"`

	// 统计信息
	TotalDeletedFiles   int64 `gorm:"default:0;comment:总删除文件数" json:"total_deleted_files"`
	TotalRestoredFiles  int64 `gorm:"default:0;comment:总恢复文件数" json:"total_restored_files"`
	TotalPermanentFiles int64 `gorm:"default:0;comment:总永久删除文件数" json:"total_permanent_files"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (RecycleBin) TableName() string {
	return "recycle_bins"
}

// RecycleLog 回收站操作日志模型
type RecycleLog struct {
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// 关联信息
	RecycleItemID uint        `gorm:"not null;index" json:"recycle_item_id"`
	RecycleItem   RecycleItem `gorm:"foreignKey:RecycleItemID;constraint:OnDelete:CASCADE" json:"recycle_item,omitempty"`
	UserID        uint        `gorm:"not null;index" json:"user_id"`
	User          User        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// 操作信息
	Action      string        `gorm:"type:varchar(50);not null;index;comment:操作类型" json:"action"`
	Description string        `gorm:"type:varchar(500);comment:操作描述" json:"description"`
	OldStatus   RecycleStatus `gorm:"type:varchar(20);comment:原状态" json:"old_status"`
	NewStatus   RecycleStatus `gorm:"type:varchar(20);comment:新状态" json:"new_status"`

	// 客户端信息
	IPAddress string `gorm:"type:varchar(45);index" json:"ip_address"`
	UserAgent string `gorm:"type:varchar(500)" json:"user_agent"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
}

// TableName 指定表名
func (RecycleLog) TableName() string {
	return "recycle_logs"
}

// BeforeCreate GORM钩子：创建前
func (ri *RecycleItem) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if ri.Status == "" {
		ri.Status = RecycleStatusDeleted
	}
	if ri.AutoDeleteDays == 0 {
		ri.AutoDeleteDays = 30
	}
	// 设置过期时间
	if ri.ExpiresAt == nil {
		expiresAt := time.Now().AddDate(0, 0, ri.AutoDeleteDays)
		ri.ExpiresAt = &expiresAt
	}
	return nil
}

// BeforeCreate GORM钩子：创建前
func (rb *RecycleBin) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if rb.AutoDeleteDays == 0 {
		rb.AutoDeleteDays = 30
	}
	if rb.MaxStorageSize == 0 {
		rb.MaxStorageSize = 1073741824 // 1GB
	}
	if rb.MaxItemCount == 0 {
		rb.MaxItemCount = 1000
	}
	if rb.NotifyDays == 0 {
		rb.NotifyDays = 7
	}
	return nil
}

// IsDeleted 检查项目是否已删除
func (ri *RecycleItem) IsDeleted() bool {
	return ri.Status == RecycleStatusDeleted
}

// IsRestored 检查项目是否已恢复
func (ri *RecycleItem) IsRestored() bool {
	return ri.Status == RecycleStatusRestored
}

// IsPermanentDeleted 检查项目是否已永久删除
func (ri *RecycleItem) IsPermanentDeleted() bool {
	return ri.Status == RecycleStatusPermanent
}

// IsExpired 检查项目是否已过期
func (ri *RecycleItem) IsExpired() bool {
	return ri.ExpiresAt != nil && time.Now().After(*ri.ExpiresAt)
}

// CanRestore 检查项目是否可以恢复
func (ri *RecycleItem) CanRestore() bool {
	return ri.Status == RecycleStatusDeleted && !ri.IsExpired()
}

// IsFile 检查是否为文件
func (ri *RecycleItem) IsFile() bool {
	return ri.Type == RecycleTypeFile
}

// IsFolder 检查是否为文件夹
func (ri *RecycleItem) IsFolder() bool {
	return ri.Type == RecycleTypeFolder
}

// IsStorageFull 检查存储是否已满
func (rb *RecycleBin) IsStorageFull() bool {
	return rb.CurrentStorageSize >= rb.MaxStorageSize
}

// IsItemCountFull 检查项目数量是否已满
func (rb *RecycleBin) IsItemCountFull() bool {
	return rb.CurrentItemCount >= rb.MaxItemCount
}

// GetStorageUsagePercent 获取存储使用百分比
func (rb *RecycleBin) GetStorageUsagePercent() float64 {
	if rb.MaxStorageSize == 0 {
		return 0
	}
	return float64(rb.CurrentStorageSize) / float64(rb.MaxStorageSize) * 100
}

// GetItemUsagePercent 获取项目数量使用百分比
func (rb *RecycleBin) GetItemUsagePercent() float64 {
	if rb.MaxItemCount == 0 {
		return 0
	}
	return float64(rb.CurrentItemCount) / float64(rb.MaxItemCount) * 100
}

// ShouldNotifyBeforeDelete 检查是否应该在删除前通知
func (rb *RecycleBin) ShouldNotifyBeforeDelete() bool {
	return rb.NotifyBeforeDelete
}

// GetNotifyDate 获取通知日期
func (ri *RecycleItem) GetNotifyDate(notifyDays int) *time.Time {
	if ri.ExpiresAt == nil {
		return nil
	}
	notifyDate := ri.ExpiresAt.AddDate(0, 0, -notifyDays)
	return &notifyDate
}
