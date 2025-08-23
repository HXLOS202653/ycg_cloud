// Package model 定义了应用程序的数据模型和数据库结构
package model

import (
	"time"

	"gorm.io/gorm"
)

// FileStatus 文件状态枚举
type FileStatus string

const (
	FileStatusNormal    FileStatus = "normal"    // 正常
	FileStatusDeleted   FileStatus = "deleted"   // 已删除
	FileStatusUploading FileStatus = "uploading" // 上传中
	FileStatusCorrupted FileStatus = "corrupted" // 损坏
)

// FileType 文件类型枚举
type FileType string

const (
	FileTypeFolder   FileType = "folder"   // 文件夹
	FileTypeDocument FileType = "document" // 文档
	FileTypeImage    FileType = "image"    // 图片
	FileTypeVideo    FileType = "video"    // 视频
	FileTypeAudio    FileType = "audio"    // 音频
	FileTypeArchive  FileType = "archive"  // 压缩包
	FileTypeOther    FileType = "other"    // 其他
)

// StorageType 存储类型枚举
type StorageType string

const (
	StorageTypeLocal StorageType = "local" // 本地存储
	StorageTypeOSS   StorageType = "oss"   // 对象存储
)

// File 文件模型
type File struct {
	// 时间戳字段 (8 bytes each)
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 指针字段 (8 bytes each)
	ShareExpiry    *time.Time `gorm:"comment:分享过期时间" json:"share_expiry"`
	ParentID       *uint      `gorm:"index;comment:父级目录ID" json:"parent_id"`
	Parent         *File      `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	OriginalFileID *uint      `gorm:"index;comment:原始文件ID" json:"original_file_id"`
	OriginalFile   *File      `gorm:"foreignKey:OriginalFileID;constraint:OnDelete:SET NULL" json:"original_file,omitempty"`

	// 切片字段 (24 bytes each - pointer + len + cap)
	Children []File `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"children,omitempty"`
	Versions []File `gorm:"foreignKey:OriginalFileID;constraint:OnDelete:SET NULL" json:"versions,omitempty"`

	// int64字段 (8 bytes)
	Size int64 `gorm:"default:0;comment:文件大小(字节)" json:"size"`

	// uint字段 (4 bytes)
	ID uint `gorm:"primaryKey;autoIncrement" json:"id"`

	// int字段 (4 bytes each)
	Version       int `gorm:"default:1;comment:文件版本号" json:"version"`
	DownloadCount int `gorm:"default:0;comment:下载次数" json:"download_count"`
	ViewCount     int `gorm:"default:0;comment:查看次数" json:"view_count"`

	// 字符串字段 (16 bytes each - pointer + len)
	Name          string `gorm:"type:varchar(255);not null;index" json:"name" validate:"required"`
	Path          string `gorm:"type:varchar(1000);not null;index" json:"path"`
	MimeType      string `gorm:"type:varchar(100);index" json:"mime_type"`
	MD5Hash       string `gorm:"type:varchar(32);index;comment:文件MD5哈希" json:"md5_hash"`
	SHA256Hash    string `gorm:"type:varchar(64);index;comment:文件SHA256哈希" json:"sha256_hash"`
	StoragePath   string `gorm:"type:varchar(1000);comment:实际存储路径" json:"storage_path"`
	BucketName    string `gorm:"type:varchar(100);comment:OSS桶名" json:"bucket_name"`
	ShareToken    string `gorm:"type:varchar(100);uniqueIndex;comment:分享令牌" json:"share_token"`
	SharePassword string `gorm:"type:varchar(255);comment:分享密码" json:"-"`
	Tags          string `gorm:"type:text;comment:文件标签(JSON)" json:"tags"`
	Category      string `gorm:"type:varchar(100);index;comment:文件分类" json:"category"`
	Description   string `gorm:"type:text;comment:文件描述" json:"description"`
	ThumbnailPath string `gorm:"type:varchar(1000);comment:缩略图路径" json:"thumbnail_path"`
	PreviewPath   string `gorm:"type:varchar(1000);comment:预览文件路径" json:"preview_path"`
	EncryptionKey string `gorm:"type:varchar(255);comment:加密密钥" json:"-"`

	// 枚举字段 (按字符串处理，16 bytes each)
	FileType    FileType    `gorm:"type:varchar(20);index" json:"file_type"`
	Status      FileStatus  `gorm:"type:varchar(20);default:'normal';index" json:"status"`
	StorageType StorageType `gorm:"type:varchar(20);default:'local';index" json:"storage_type"`

	// bool字段 (1 byte each)
	IsPublic     bool `gorm:"default:false;index;comment:是否公开" json:"is_public"`
	IsLatest     bool `gorm:"default:true;index;comment:是否最新版本" json:"is_latest"`
	CanPreview   bool `gorm:"default:false;index;comment:是否可预览" json:"can_preview"`
	IsEncrypted  bool `gorm:"default:false;index;comment:是否加密" json:"is_encrypted"`
	IsCompressed bool `gorm:"default:false;comment:是否压缩" json:"is_compressed"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}

// BeforeCreate GORM钩子：创建前
func (f *File) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if f.Status == "" {
		f.Status = FileStatusNormal
	}
	if f.StorageType == "" {
		f.StorageType = StorageTypeLocal
	}
	if f.Version == 0 {
		f.Version = 1
	}
	return nil
}

// IsFolder 检查是否为文件夹
func (f *File) IsFolder() bool {
	return f.FileType == FileTypeFolder
}

// IsDeleted 检查文件是否已删除
func (f *File) IsDeleted() bool {
	return f.Status == FileStatusDeleted
}

// IsShared 检查文件是否已分享
func (f *File) IsShared() bool {
	return f.ShareToken != "" && (f.ShareExpiry == nil || f.ShareExpiry.After(time.Now()))
}

// IsShareExpired 检查分享是否已过期
func (f *File) IsShareExpired() bool {
	return f.ShareExpiry != nil && f.ShareExpiry.Before(time.Now())
}

// GetFullPath 获取完整路径
func (f *File) GetFullPath() string {
	if f.Path == "" {
		return f.Name
	}
	return f.Path + "/" + f.Name
}

// CanPreviewFile 检查文件是否可以预览
func (f *File) CanPreviewFile() bool {
	return f.CanPreview && !f.IsDeleted() && f.Status == FileStatusNormal
}
