package model

import (
	"time"

	"gorm.io/gorm"
)

// messageType 消息类型枚举
type messageType string

const (
	MessageTypeText    messageType = "text"    // 文本消息
	MessageTypeImage   messageType = "image"   // 图片消息
	MessageTypeFile    messageType = "file"    // 文件消息
	MessageTypeAudio   messageType = "audio"   // 音频消息
	MessageTypeVideo   messageType = "video"   // 视频消息
	MessageTypeSystem  messageType = "system"  // 系统消息
	MessageTypeNotice  messageType = "notice"  // 通知消息
	MessageTypeReply   messageType = "reply"   // 回复消息
	MessageTypeForward messageType = "forward" // 转发消息
)

// messageStatus 消息状态枚举
type messageStatus string

const (
	MessageStatusSent      messageStatus = "sent"      // 已发送
	MessageStatusDelivered messageStatus = "delivered" // 已送达
	MessageStatusRead      messageStatus = "read"      // 已读
	MessageStatusDeleted   messageStatus = "deleted"   // 已删除
	MessageStatusRecalled  messageStatus = "recalled"  // 已撤回
)

// conversationType 会话类型枚举
type conversationType string

const (
	ConversationTypePrivate conversationType = "private" // 私聊
	ConversationTypeGroup   conversationType = "group"   // 群聊
	ConversationTypeTeam    conversationType = "team"    // 团队会话
	ConversationTypeSystem  conversationType = "system"  // 系统会话
)

// ConversationStatus 会话状态枚举
type ConversationStatus string

const (
	ConversationStatusActive   ConversationStatus = "active"   // 活跃
	ConversationStatusArchived ConversationStatus = "archived" // 已归档
	ConversationStatusMuted    ConversationStatus = "muted"    // 已静音
	ConversationStatusDeleted  ConversationStatus = "deleted"  // 已删除
)

// Conversation 会话模型
type Conversation struct {
	ID          uint               `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string             `gorm:"type:varchar(200);index" json:"title"`
	Description string             `gorm:"type:text" json:"description"`
	Avatar      string             `gorm:"type:varchar(500)" json:"avatar"`
	Type        conversationType   `gorm:"type:varchar(20);not null;index" json:"type"`
	Status      ConversationStatus `gorm:"type:varchar(20);default:'active';index" json:"status"`

	// 创建者信息
	CreatorID uint `gorm:"not null;index" json:"creator_id"`
	Creator   User `gorm:"foreignKey:CreatorID;constraint:OnDelete:RESTRICT" json:"creator,omitempty"`

	// 团队关联(可选)
	TeamID *uint `gorm:"index;comment:关联团队ID" json:"team_id"`
	Team   *Team `gorm:"foreignKey:TeamID;constraint:OnDelete:CASCADE" json:"team,omitempty"`

	// 会话设置
	IsPublic    bool `gorm:"default:false;index;comment:是否公开会话" json:"is_public"`
	MaxMembers  int  `gorm:"default:100;comment:最大成员数" json:"max_members"`
	AllowInvite bool `gorm:"default:true;comment:允许邀请新成员" json:"allow_invite"`

	// 消息设置
	MessageRetentionDays int  `gorm:"default:0;comment:消息保留天数(0表示永久)" json:"message_retention_days"`
	AllowFileShare       bool `gorm:"default:true;comment:允许文件分享" json:"allow_file_share"`

	// 最后消息信息
	LastMessageID *uint      `gorm:"index;comment:最后一条消息ID" json:"last_message_id"`
	LastMessage   *Message   `gorm:"foreignKey:LastMessageID;constraint:OnDelete:SET NULL" json:"last_message,omitempty"`
	LastMessageAt *time.Time `gorm:"index;comment:最后消息时间" json:"last_message_at"`

	// 时间戳
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 关联关系
	Messages []Message            `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"messages,omitempty"`
	Members  []conversationMember `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
}

// TableName 指定表名
func (Conversation) TableName() string {
	return "conversations"
}

// conversationMember 会话成员模型
type conversationMember struct {
	// 结构体字段 (最大的字段放在前面)
	Conversation Conversation `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"conversation,omitempty"`
	User         User         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// 时间戳字段 (24 bytes each)
	JoinedAt  time.Time      `gorm:"autoCreateTime;comment:加入时间" json:"joined_at"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 字符串字段 (24 bytes each)
	Nickname string `gorm:"type:varchar(100);comment:群昵称" json:"nickname"`
	Role     string `gorm:"type:varchar(20);default:'member';comment:成员角色" json:"role"`

	// 指针字段 (8 bytes each)
	LastReadAt        *time.Time `gorm:"comment:最后阅读时间" json:"last_read_at"`
	Inviter           *User      `gorm:"foreignKey:InvitedBy" json:"inviter,omitempty"`
	LastReadMessage   *Message   `gorm:"foreignKey:LastReadMessageID" json:"last_read_message,omitempty"`
	InvitedBy         *uint      `gorm:"index;comment:邀请人ID" json:"invited_by"`
	LastReadMessageID *uint      `gorm:"index;comment:最后已读消息ID" json:"last_read_message_id"`

	// uint字段 (8 bytes each)
	ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID uint `gorm:"not null;index;comment:会话ID" json:"conversation_id"`
	UserID         uint `gorm:"not null;index;comment:用户ID" json:"user_id"`

	// int字段 (8 bytes each)
	UnreadCount int `gorm:"default:0;comment:未读消息数" json:"unread_count"`

	// bool字段 (1 byte each) - 放在最后
	IsMuted   bool `gorm:"default:false;comment:是否静音" json:"is_muted"`
	AdminFlag bool `gorm:"default:false;index;comment:是否为管理员" json:"is_admin"`
}

// TableName 指定表名
func (conversationMember) TableName() string {
	return "conversation_members"
}

// ConversationMember 公共类型别名
type ConversationMember = conversationMember

// Message 消息模型
type Message struct {
	// 时间戳字段 (24 bytes each)
	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// 指针字段 (8 bytes each)
	ReplyToID     *uint      `gorm:"index;comment:回复的消息ID" json:"reply_to_id"`
	ReplyTo       *Message   `gorm:"foreignKey:ReplyToID;constraint:OnDelete:SET NULL" json:"reply_to,omitempty"`
	ForwardFromID *uint      `gorm:"index;comment:转发来源消息ID" json:"forward_from_id"`
	ForwardFrom   *Message   `gorm:"foreignKey:ForwardFromID;constraint:OnDelete:SET NULL" json:"forward_from,omitempty"`
	FileID        *uint      `gorm:"index;comment:关联文件ID" json:"file_id"`
	File          *File      `gorm:"foreignKey:FileID" json:"file,omitempty"`
	EditedAt      *time.Time `gorm:"comment:编辑时间" json:"edited_at"`
	RecalledAt    *time.Time `gorm:"comment:撤回时间" json:"recalled_at"`
	RecalledBy    *uint      `gorm:"index;comment:撤回人ID" json:"recalled_by"`
	Recaller      *User      `gorm:"foreignKey:RecalledBy" json:"recalled,omitempty"`

	// 结构体字段
	Conversation Conversation `gorm:"foreignKey:ConversationID;constraint:OnDelete:CASCADE" json:"conversation,omitempty"`
	Sender       User         `gorm:"foreignKey:SenderID;constraint:OnDelete:RESTRICT" json:"sender,omitempty"`

	// 字符串字段 (24 bytes each)
	Content    string        `gorm:"type:text;not null" json:"content"`
	RawContent string        `gorm:"type:text;comment:原始内容(用于编辑历史)" json:"raw_content"`
	Metadata   string        `gorm:"type:text;comment:消息元数据(JSON)" json:"metadata"`
	Mentions   string        `gorm:"type:text;comment:提及的用户(JSON)" json:"mentions"`
	Type       messageType   `gorm:"type:varchar(20);not null;index" json:"type"`
	Status     messageStatus `gorm:"type:varchar(20);default:'sent';index" json:"status"`

	// uint字段 (8 bytes each)
	ID             uint `gorm:"primaryKey;autoIncrement" json:"id"`
	ConversationID uint `gorm:"not null;index" json:"conversation_id"`
	SenderID       uint `gorm:"not null;index" json:"sender_id"`

	// bool字段 (1 byte each)
	IsEdited     bool `gorm:"default:false;index;comment:是否已编辑" json:"is_edited"`
	RecalledFlag bool `gorm:"default:false;index;comment:是否已撤回" json:"is_recalled"`

	// 关联关系
	Replies      []Message            `gorm:"foreignKey:ReplyToID;constraint:OnDelete:SET NULL" json:"replies,omitempty"`
	Forwards     []Message            `gorm:"foreignKey:ForwardFromID;constraint:OnDelete:SET NULL" json:"forwards,omitempty"`
	ReadReceipts []messageReadReceipt `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE" json:"read_receipts,omitempty"`
}

// TableName 指定表名
func (Message) TableName() string {
	return "messages"
}

// messageReadReceipt 消息已读回执模型
type messageReadReceipt struct {
	ID        uint    `gorm:"primaryKey;autoIncrement" json:"id"`
	MessageID uint    `gorm:"not null;index" json:"message_id"`
	Message   Message `gorm:"foreignKey:MessageID;constraint:OnDelete:CASCADE" json:"message,omitempty"`
	UserID    uint    `gorm:"not null;index" json:"user_id"`
	User      User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`

	// 已读信息
	ReadAt time.Time `gorm:"autoCreateTime" json:"read_at"`

	// 时间戳
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName 指定表名
func (messageReadReceipt) TableName() string {
	return "message_read_receipts"
}

// MessageReadReceipt 公共类型别名
type MessageReadReceipt = messageReadReceipt

// BeforeCreate GORM钩子：创建前
func (c *Conversation) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if c.Status == "" {
		c.Status = ConversationStatusActive
	}
	if c.MaxMembers == 0 {
		c.MaxMembers = 100
	}
	return nil
}

// BeforeCreate GORM钩子：创建前
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	// 设置默认值
	if m.Status == "" {
		m.Status = MessageStatusSent
	}
	return nil
}

// IsActive 检查会话是否活跃
func (c *Conversation) IsActive() bool {
	return c.Status == ConversationStatusActive
}

// IsPrivate 检查是否为私聊
func (c *Conversation) IsPrivate() bool {
	return c.Type == ConversationTypePrivate
}

// IsGroup 检查是否为群聊
func (c *Conversation) IsGroup() bool {
	return c.Type == ConversationTypeGroup
}

// IsTeamConversation 检查是否为团队会话
func (c *Conversation) IsTeamConversation() bool {
	return c.Type == ConversationTypeTeam
}

// IsAdmin 检查成员是否为管理员
func (cm *conversationMember) IsAdmin() bool {
	return cm.AdminFlag
}

// CanManageConversation 检查成员是否可以管理会话
func (cm *conversationMember) CanManageConversation() bool {
	return cm.AdminFlag
}

// IsTextMessage 检查是否为文本消息
func (m *Message) IsTextMessage() bool {
	return m.Type == MessageTypeText
}

// IsFileMessage 检查是否为文件消息
func (m *Message) IsFileMessage() bool {
	return m.Type == MessageTypeFile
}

// IsSystemMessage 检查是否为系统消息
func (m *Message) IsSystemMessage() bool {
	return m.Type == MessageTypeSystem
}

// IsRecalled 检查消息是否已撤回
func (m *Message) IsRecalled() bool {
	return m.RecalledFlag
}

// CanRecall 检查消息是否可以撤回(2分钟内)
func (m *Message) CanRecall() bool {
	return !m.RecalledFlag && time.Since(m.CreatedAt) <= 2*time.Minute
}
