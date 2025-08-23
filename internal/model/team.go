package model

import (
	"time"

	"gorm.io/gorm"
)

// TeamStatus 团队状态枚举
type TeamStatus string

const (
	TeamStatusActive    TeamStatus = "active"    // 活跃
	TeamStatusInactive  TeamStatus = "inactive"  // 非活跃
	TeamStatusSuspended TeamStatus = "suspended" // 暂停
	TeamStatusDeleted   TeamStatus = "deleted"   // 删除
)

// teamMemberRole 团队成员角色枚举 (私有)
type teamMemberRole string

const (
	TeamMemberRoleOwner  teamMemberRole = "owner"  // 所有者
	TeamMemberRoleAdmin  teamMemberRole = "admin"  // 管理员
	TeamMemberRoleMember teamMemberRole = "member" // 成员
	TeamMemberRoleViewer teamMemberRole = "viewer" // 查看者
)

// teamMemberStatus 团队成员状态枚举 (私有)
type teamMemberStatus string

const (
	TeamMemberStatusActive   teamMemberStatus = "active"   // 活跃
	TeamMemberStatusInvited  teamMemberStatus = "invited"  // 已邀请
	TeamMemberStatusInactive teamMemberStatus = "inactive" // 非活跃
	TeamMemberStatusLeft     teamMemberStatus = "left"     // 已离开
)

// Team 团队模型
type Team struct {
	// time.Time 字段放在最前面 (8字节对齐)
	CreatedAt time.Time `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// 结构体字段 (8字节)
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 切片字段 (24字节)
	Members []teamMember `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
	Files   []teamFile   `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"files,omitempty"`

	// 字符串字段 (16字节)
	Name        string     `gorm:"type:varchar(100);not null;index" json:"name"`
	Description string     `gorm:"type:text" json:"description"`
	Avatar      string     `gorm:"type:varchar(500)" json:"avatar"`
	Status      TeamStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`

	// int64字段 (8字节)
	StorageUsed  int64 `gorm:"default:0" json:"storage_used"`
	StorageLimit int64 `gorm:"default:10737418240" json:"storage_limit"` // 10GB

	// uint字段 (4字节)
	ID        uint `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatorID uint `gorm:"not null;index" json:"creator_id"`

	// int字段 (4字节)
	MaxMembers  int `gorm:"default:100" json:"max_members"`
	MemberCount int `gorm:"default:0" json:"member_count"`
	FileCount   int `gorm:"default:0" json:"file_count"`

	// bool字段 (1字节)
	IsPublic bool `gorm:"default:false;index" json:"is_public"`

	// 关联关系
	Creator User `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT" json:"creator,omitempty"`
}

// TableName 指定表名
func (Team) TableName() string {
	return "teams"
}

// teamMember 团队成员模型 (私有)
type teamMember struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID uint `gorm:"not null;index" json:"team_id"`
	Team   Team `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"team,omitempty"`
	UserID uint `gorm:"not null;index" json:"user_id"`
	User   User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// 成员信息
	Role     teamMemberRole   `gorm:"type:varchar(20);default:'member';index" json:"role"`
	Status   teamMemberStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`
	Nickname string           `gorm:"type:varchar(100);comment:团队内昵称" json:"nickname"`

	// 邀请信息
	InvitedBy *uint      `gorm:"index;comment:邀请人ID" json:"invited_by"`
	Inviter   *User      `gorm:"foreignKey:InvitedBy;constraint:OnDelete:SET NULL" json:"inviter,omitempty"`
	InvitedAt *time.Time `gorm:"comment:邀请时间" json:"invited_at"`
	JoinedAt  *time.Time `gorm:"comment:加入时间" json:"joined_at"`
	LeftAt    *time.Time `gorm:"comment:离开时间" json:"left_at"`

	// 权限设置
	Permissions string `gorm:"type:text;comment:成员权限(JSON)" json:"permissions"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TeamMember 团队成员模型 (公共类型别名)
type TeamMember = teamMember

// TableName 指定表名
func (teamMember) TableName() string {
	return "team_members"
}

// teamFile 团队文件关联模型 (私有)
type teamFile struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID uint `gorm:"not null;index" json:"team_id"`
	Team   Team `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"team,omitempty"`
	FileID uint `gorm:"not null;index" json:"file_id"`
	File   File `gorm:"foreignKey:FileID;constraint:OnDelete:CASCADE" json:"file,omitempty"`

	// 共享信息
	SharedBy uint      `gorm:"not null;index;comment:分享人ID" json:"shared_by"`
	Sharer   User      `gorm:"foreignKey:SharedBy;constraint:OnDelete:RESTRICT" json:"sharer,omitempty"`
	SharedAt time.Time `gorm:"autoCreateTime" json:"shared_at"`

	// 权限设置
	Permissions string `gorm:"type:text;comment:文件权限(JSON)" json:"permissions"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TeamFile 团队文件关联模型 (公共类型别名)
type TeamFile = teamFile

// TableName 指定表名
func (teamFile) TableName() string {
	return "team_files"
}

// teamRole 团队角色关联 (私有)
type teamRole struct {
	ID     uint `gorm:"primaryKey;autoIncrement" json:"id"`
	TeamID uint `gorm:"not null;index" json:"team_id"`
	Team   Team `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"team,omitempty"`
	RoleID uint `gorm:"not null;index" json:"role_id"`
	Role   Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`

	// 权限来源
	GrantedBy *uint     `gorm:"index;comment:授权人ID" json:"granted_by"`
	Granter   *User     `gorm:"foreignKey:GrantedBy" json:"granter,omitempty"`
	GrantedAt time.Time `gorm:"autoCreateTime" json:"granted_at"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TeamRole 团队角色关联 (公共类型别名)
type TeamRole = teamRole

// TableName 指定表名
func (teamRole) TableName() string {
	return "team_roles"
}

// BeforeCreate GORM钩子：创建前
func (t *Team) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if t.Status == "" {
		t.Status = TeamStatusActive
	}
	if t.MaxMembers == 0 {
		t.MaxMembers = 50
	}
	if t.StorageLimit == 0 {
		t.StorageLimit = 53687091200 // 50GB
	}
	return nil
}

// BeforeCreate GORM钩子：创建前
func (tm *TeamMember) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if tm.Role == "" {
		tm.Role = TeamMemberRoleMember
	}
	if tm.Status == "" {
		tm.Status = TeamMemberStatusActive
	}
	return nil
}

// IsActive 检查团队是否活跃
func (t *Team) IsActive() bool {
	return t.Status == TeamStatusActive
}

// IsStorageExceeded 检查团队存储是否超限
func (t *Team) IsStorageExceeded() bool {
	return t.StorageUsed >= t.StorageLimit
}

// GetAvailableStorage 获取团队可用存储空间
func (t *Team) GetAvailableStorage() int64 {
	return t.StorageLimit - t.StorageUsed
}

// IsOwner 检查成员是否为团队所有者
func (tm *TeamMember) IsOwner() bool {
	return tm.Role == TeamMemberRoleOwner
}

// IsAdmin 检查成员是否为团队管理员
func (tm *TeamMember) IsAdmin() bool {
	return tm.Role == TeamMemberRoleAdmin || tm.Role == TeamMemberRoleOwner
}

// IsActive 检查成员是否活跃
func (tm *TeamMember) IsActive() bool {
	return tm.Status == TeamMemberStatusActive
}

// CanManageTeam 检查成员是否可以管理团队
func (tm *TeamMember) CanManageTeam() bool {
	return tm.IsActive() && (tm.Role == TeamMemberRoleOwner || tm.Role == TeamMemberRoleAdmin)
}
