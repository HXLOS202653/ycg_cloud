# 网络云盘系统 WebSocket协议规范

## 📋 文档信息

| 项目 | 信息 |
|------|------|
| 文档标题 | 网络云盘系统 WebSocket协议规范 |
| 文档版本 | v1.0 |
| 创建日期 | 2024年12月 |
| 最后更新 | 2024年12月 |
| 作者 | 开发团队 |

## 📖 目录

1. [概述](#概述)
2. [连接管理](#连接管理)
3. [消息格式](#消息格式)
4. [协议类型](#协议类型)
5. [实时文件同步](#实时文件同步)
6. [即时通讯](#即时通讯)
7. [协作编辑](#协作编辑)
8. [语音通话](#语音通话)
9. [系统通知](#系统通知)
10. [错误处理](#错误处理)
11. [安全机制](#安全机制)
12. [性能优化](#性能优化)
13. [客户端实现](#客户端实现)
14. [服务端实现](#服务端实现)

## 概述

### 设计目标
本WebSocket协议规范旨在为网络云盘系统提供高效、稳定、安全的实时通信能力，支持以下核心功能：

- **实时文件同步**：文件操作的实时通知和状态同步
- **即时通讯**：用户间的实时消息传递
- **协作编辑**：多用户实时协作编辑文档
- **语音通话**：WebRTC信令传输
- **系统通知**：系统事件和状态推送

### 技术规范
- **WebSocket版本**：RFC 6455
- **传输层协议**：TCP
- **应用层协议**：自定义JSON格式
- **编码格式**：UTF-8
- **压缩支持**：Per-message-deflate

### 服务端点
```
开发环境：ws://localhost:8080/ws
测试环境：wss://test-api.yunpan.com/ws
生产环境：wss://api.yunpan.com/ws
```

## 连接管理

### 连接建立

#### 请求格式
```http
GET /ws HTTP/1.1
Host: api.yunpan.com
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
Sec-WebSocket-Version: 13
Sec-WebSocket-Protocol: yunpan-v1
Authorization: Bearer <JWT_TOKEN>
User-Agent: YunPan-Client/1.0
```

#### 响应格式
```http
HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
Sec-WebSocket-Protocol: yunpan-v1
```

#### 认证机制
```javascript
// 连接建立后立即发送认证消息
{
  "type": "auth",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "client_id": "web_client_123",
    "client_type": "web", // web, mobile, desktop
    "version": "1.0.0"
  },
  "timestamp": 1703123456789
}
```

#### 认证响应
```javascript
// 成功响应
{
  "type": "auth_response",
  "data": {
    "status": "success",
    "user_id": "12345",
    "session_id": "sess_67890",
    "permissions": ["file_read", "file_write", "chat", "collaborate"]
  },
  "timestamp": 1703123456790
}

// 失败响应
{
  "type": "auth_response",
  "data": {
    "status": "error",
    "error_code": "INVALID_TOKEN",
    "message": "Token已过期或无效"
  },
  "timestamp": 1703123456790
}
```

### 连接维护

#### 心跳机制
```javascript
// 客户端发送心跳（每30秒）
{
  "type": "ping",
  "timestamp": 1703123456789
}

// 服务端响应心跳
{
  "type": "pong",
  "timestamp": 1703123456790
}
```

#### 重连机制
- **初始重连间隔**：1秒
- **最大重连间隔**：30秒
- **退避策略**：指数退避（1s → 2s → 4s → 8s → 16s → 30s）
- **最大重连次数**：无限制（用户主动断开除外）

#### 连接状态
```javascript
// 连接状态枚举
const ConnectionState = {
  CONNECTING: "connecting",
  CONNECTED: "connected", 
  AUTHENTICATED: "authenticated",
  DISCONNECTING: "disconnecting",
  DISCONNECTED: "disconnected",
  RECONNECTING: "reconnecting",
  ERROR: "error"
}
```

### 连接断开

#### 正常断开
```javascript
// 客户端主动断开
{
  "type": "disconnect",
  "data": {
    "reason": "user_logout"
  },
  "timestamp": 1703123456789
}
```

#### 异常断开
- **1000**：Normal Closure（正常关闭）
- **1001**：Going Away（端点离开）
- **1002**：Protocol Error（协议错误）
- **1003**：Unsupported Data（不支持的数据类型）
- **1006**：Abnormal Closure（异常关闭）
- **1011**：Internal Server Error（服务器内部错误）
- **4000**：Authentication Failed（认证失败）
- **4001**：Permission Denied（权限不足）
- **4002**：Rate Limited（频率限制）

## 消息格式

### 基础消息结构
```javascript
{
  "type": "message_type",        // 消息类型（必需）
  "id": "msg_12345",            // 消息ID（可选，用于响应关联）
  "data": {                     // 消息数据（必需）
    // 具体数据内容
  },
  "timestamp": 1703123456789,   // 时间戳（必需）
  "version": "1.0"              // 协议版本（可选）
}
```

### 消息类型分类

#### 系统消息
- `auth` - 认证请求
- `auth_response` - 认证响应
- `ping` - 心跳检测
- `pong` - 心跳响应
- `disconnect` - 断开连接
- `error` - 错误消息

#### 文件同步消息
- `file_created` - 文件创建
- `file_updated` - 文件更新
- `file_deleted` - 文件删除
- `file_moved` - 文件移动
- `folder_created` - 文件夹创建
- `folder_updated` - 文件夹更新
- `folder_deleted` - 文件夹删除

#### 聊天消息
- `chat_message` - 聊天消息
- `chat_typing` - 正在输入
- `chat_read` - 消息已读
- `chat_online` - 用户上线
- `chat_offline` - 用户下线

#### 协作消息
- `collab_join` - 加入协作
- `collab_leave` - 离开协作
- `collab_operation` - 协作操作
- `collab_cursor` - 光标位置
- `collab_selection` - 选中内容

#### 通话消息
- `call_offer` - 通话邀请
- `call_answer` - 接受通话
- `call_ice_candidate` - ICE候选
- `call_hangup` - 挂断通话

### 消息编码
- **字符编码**：UTF-8
- **序列化格式**：JSON
- **压缩算法**：Per-message-deflate（可选）

### 消息大小限制
- **单条消息最大大小**：1MB
- **文本消息建议大小**：64KB
- **二进制消息建议大小**：256KB

## 协议类型

### 文本协议
用于传输结构化数据，采用JSON格式：

```javascript
{
  "type": "chat_message",
  "data": {
    "room_id": "room_123",
    "sender_id": "user_456",
    "content": "Hello, World!",
    "message_type": "text"
  },
  "timestamp": 1703123456789
}
```

### 二进制协议
用于传输文件块、音频数据等：

```javascript
// 二进制消息头部（JSON）
{
  "type": "file_chunk",
  "data": {
    "file_id": "file_789",
    "chunk_index": 0,
    "chunk_size": 1024,
    "total_chunks": 100,
    "chunk_hash": "sha256_hash"
  },
  "timestamp": 1703123456789
}
// 后跟二进制数据块
```

## 实时文件同步

### 文件操作通知

#### 文件创建
```javascript
{
  "type": "file_created",
  "data": {
    "file_id": "file_12345",
    "filename": "document.pdf",
    "file_size": 1048576,
    "mime_type": "application/pdf",
    "parent_folder_id": "folder_67890",
    "created_by": "user_123",
    "created_at": "2024-12-21T10:30:00Z",
    "file_hash": "sha256_hash_value"
  },
  "timestamp": 1703123456789
}
```

#### 文件更新
```javascript
{
  "type": "file_updated",
  "data": {
    "file_id": "file_12345",
    "changes": {
      "filename": "new_document.pdf",
      "file_size": 1048576,
      "updated_at": "2024-12-21T10:35:00Z",
      "version": 2
    },
    "updated_by": "user_123",
    "change_type": "content" // content, metadata, permissions
  },
  "timestamp": 1703123456789
}
```

#### 文件删除
```javascript
{
  "type": "file_deleted",
  "data": {
    "file_id": "file_12345",
    "filename": "document.pdf",
    "deleted_by": "user_123",
    "deleted_at": "2024-12-21T10:40:00Z",
    "is_permanent": false // true为永久删除，false为移入回收站
  },
  "timestamp": 1703123456789
}
```

#### 文件移动
```javascript
{
  "type": "file_moved",
  "data": {
    "file_id": "file_12345",
    "old_parent_id": "folder_111",
    "new_parent_id": "folder_222",
    "old_path": "/Documents/file.pdf",
    "new_path": "/Projects/file.pdf",
    "moved_by": "user_123",
    "moved_at": "2024-12-21T10:45:00Z"
  },
  "timestamp": 1703123456789
}
```

### 文件夹操作通知

#### 文件夹创建
```javascript
{
  "type": "folder_created",
  "data": {
    "folder_id": "folder_12345",
    "folder_name": "新建文件夹",
    "parent_folder_id": "folder_67890",
    "created_by": "user_123",
    "created_at": "2024-12-21T10:30:00Z",
    "permissions": {
      "read": true,
      "write": true,
      "delete": true
    }
  },
  "timestamp": 1703123456789
}
```

### 同步状态管理

#### 同步状态枚举
```javascript
const SyncStatus = {
  SYNCED: "synced",           // 已同步
  SYNCING: "syncing",         // 同步中
  CONFLICT: "conflict",       // 冲突
  ERROR: "error",             // 错误
  OFFLINE: "offline"          // 离线
}
```

#### 同步状态通知
```javascript
{
  "type": "sync_status",
  "data": {
    "entity_id": "file_12345",
    "entity_type": "file", // file, folder
    "status": "syncing",
    "progress": 75, // 同步进度百分比
    "message": "正在上传文件...",
    "error": null
  },
  "timestamp": 1703123456789
}
```

## 即时通讯

### 聊天消息

#### 文本消息
```javascript
{
  "type": "chat_message",
  "data": {
    "message_id": "msg_12345",
    "room_id": "room_67890",
    "sender_id": "user_123",
    "sender_name": "张三",
    "sender_avatar": "https://cdn.yunpan.com/avatars/user_123.jpg",
    "content": "大家好！",
    "message_type": "text",
    "reply_to": null, // 回复消息ID
    "mentions": [], // @提及的用户ID列表
    "created_at": "2024-12-21T10:30:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 图片消息
```javascript
{
  "type": "chat_message",
  "data": {
    "message_id": "msg_12346",
    "room_id": "room_67890",
    "sender_id": "user_123",
    "sender_name": "张三",
    "content": "",
    "message_type": "image",
    "attachments": [{
      "file_id": "file_78901",
      "filename": "screenshot.png",
      "file_size": 524288,
      "mime_type": "image/png",
      "thumbnail_url": "https://cdn.yunpan.com/thumbnails/file_78901_thumb.jpg",
      "download_url": "https://cdn.yunpan.com/files/file_78901.png",
      "width": 1920,
      "height": 1080
    }],
    "created_at": "2024-12-21T10:32:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 语音消息
```javascript
{
  "type": "chat_message",
  "data": {
    "message_id": "msg_12347",
    "room_id": "room_67890",
    "sender_id": "user_123",
    "sender_name": "张三",
    "content": "",
    "message_type": "audio",
    "attachments": [{
      "file_id": "file_78902",
      "filename": "voice_message.m4a",
      "file_size": 65536,
      "mime_type": "audio/m4a",
      "duration": 15, // 语音时长（秒）
      "download_url": "https://cdn.yunpan.com/audio/file_78902.m4a",
      "transcript": "这是语音转文字的内容" // 可选
    }],
    "created_at": "2024-12-21T10:33:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 文件消息
```javascript
{
  "type": "chat_message",
  "data": {
    "message_id": "msg_12348",
    "room_id": "room_67890",
    "sender_id": "user_123",
    "sender_name": "张三",
    "content": "分享一个文档给大家",
    "message_type": "file",
    "attachments": [{
      "file_id": "file_78903",
      "filename": "项目方案.docx",
      "file_size": 2097152,
      "mime_type": "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
      "download_url": "https://cdn.yunpan.com/files/file_78903.docx",
      "preview_url": "https://api.yunpan.com/preview/file_78903"
    }],
    "created_at": "2024-12-21T10:35:00Z"
  },
  "timestamp": 1703123456789
}
```

### 聊天状态

#### 正在输入
```javascript
{
  "type": "chat_typing",
  "data": {
    "room_id": "room_67890",
    "user_id": "user_123",
    "user_name": "张三",
    "is_typing": true // true开始输入，false停止输入
  },
  "timestamp": 1703123456789
}
```

#### 消息已读
```javascript
{
  "type": "chat_read",
  "data": {
    "room_id": "room_67890",
    "user_id": "user_123",
    "message_id": "msg_12345", // 最后已读消息ID
    "read_at": "2024-12-21T10:36:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 用户在线状态
```javascript
{
  "type": "chat_online",
  "data": {
    "user_id": "user_123",
    "user_name": "张三",
    "status": "online", // online, away, busy, offline
    "last_seen": "2024-12-21T10:36:00Z",
    "device_type": "web" // web, mobile, desktop
  },
  "timestamp": 1703123456789
}

{
  "type": "chat_offline",
  "data": {
    "user_id": "user_123",
    "user_name": "张三",
    "last_seen": "2024-12-21T10:40:00Z"
  },
  "timestamp": 1703123456789
}
```

### 群聊管理

#### 加入群聊
```javascript
{
  "type": "chat_room_joined",
  "data": {
    "room_id": "room_67890",
    "user_id": "user_456",
    "user_name": "李四",
    "joined_by": "user_123", // 邀请者ID
    "joined_at": "2024-12-21T10:40:00Z",
    "role": "member" // admin, member, readonly
  },
  "timestamp": 1703123456789
}
```

#### 离开群聊
```javascript
{
  "type": "chat_room_left",
  "data": {
    "room_id": "room_67890",
    "user_id": "user_456",
    "user_name": "李四",
    "left_at": "2024-12-21T10:45:00Z",
    "reason": "voluntary" // voluntary, kicked, banned
  },
  "timestamp": 1703123456789
}
```

## 协作编辑

### 协作会话管理

#### 加入协作
```javascript
{
  "type": "collab_join",
  "data": {
    "document_id": "doc_12345",
    "session_id": "collab_67890",
    "user_id": "user_123",
    "user_name": "张三",
    "user_avatar": "https://cdn.yunpan.com/avatars/user_123.jpg",
    "user_color": "#FF5733", // 用户在协作中的颜色
    "permissions": ["read", "write", "comment"],
    "joined_at": "2024-12-21T10:30:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 离开协作
```javascript
{
  "type": "collab_leave",
  "data": {
    "document_id": "doc_12345",
    "session_id": "collab_67890",
    "user_id": "user_123",
    "user_name": "张三",
    "left_at": "2024-12-21T10:45:00Z",
    "reason": "voluntary" // voluntary, timeout, forced
  },
  "timestamp": 1703123456789
}
```

### 协作操作

#### 文档操作（OT算法）
```javascript
{
  "type": "collab_operation",
  "data": {
    "document_id": "doc_12345",
    "session_id": "collab_67890",
    "operation_id": "op_78901",
    "user_id": "user_123",
    "operation": {
      "type": "text_insert", // text_insert, text_delete, text_replace, format_apply
      "position": 150,
      "content": "插入的文本内容",
      "length": 7,
      "attributes": {
        "bold": true,
        "italic": false,
        "color": "#000000"
      }
    },
    "document_version": 42,
    "timestamp": "2024-12-21T10:30:15Z"
  },
  "timestamp": 1703123456789
}
```

#### 光标位置同步
```javascript
{
  "type": "collab_cursor",
  "data": {
    "document_id": "doc_12345",
    "session_id": "collab_67890",
    "user_id": "user_123",
    "user_name": "张三",
    "user_color": "#FF5733",
    "cursor": {
      "position": 150,
      "selection_start": 140,
      "selection_end": 160,
      "is_selecting": true
    },
    "viewport": {
      "top": 100,
      "height": 800
    }
  },
  "timestamp": 1703123456789
}
```

#### 选中内容同步
```javascript
{
  "type": "collab_selection",
  "data": {
    "document_id": "doc_12345",
    "session_id": "collab_67890",
    "user_id": "user_123",
    "user_name": "张三",
    "user_color": "#FF5733",
    "selection": {
      "start": 140,
      "end": 160,
      "content": "选中的文本内容"
    }
  },
  "timestamp": 1703123456789
}
```

### 版本控制

#### 文档快照
```javascript
{
  "type": "collab_snapshot",
  "data": {
    "document_id": "doc_12345",
    "version": 50,
    "content": "完整的文档内容...",
    "created_by": "system",
    "created_at": "2024-12-21T10:40:00Z",
    "reason": "auto_save" // auto_save, manual_save, conflict_resolution
  },
  "timestamp": 1703123456789
}
```

#### 冲突解决
```javascript
{
  "type": "collab_conflict",
  "data": {
    "document_id": "doc_12345",
    "conflict_id": "conflict_12345",
    "operations": [
      {
        "operation_id": "op_78901",
        "user_id": "user_123",
        "operation": { /* 操作详情 */ }
      },
      {
        "operation_id": "op_78902", 
        "user_id": "user_456",
        "operation": { /* 冲突操作详情 */ }
      }
    ],
    "resolution_strategy": "last_write_wins", // last_write_wins, manual, auto_merge
    "resolved_at": "2024-12-21T10:35:00Z"
  },
  "timestamp": 1703123456789
}
```

## 语音通话

### WebRTC信令

#### 通话邀请
```javascript
{
  "type": "call_offer",
  "data": {
    "call_id": "call_12345",
    "caller_id": "user_123",
    "caller_name": "张三",
    "callee_id": "user_456",
    "call_type": "audio", // audio, video
    "sdp": {
      "type": "offer",
      "sdp": "v=0\r\no=- 1234567890 1234567890 IN IP4 192.168.1.100\r\n..."
    },
    "created_at": "2024-12-21T10:30:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 应答通话
```javascript
{
  "type": "call_answer",
  "data": {
    "call_id": "call_12345",
    "caller_id": "user_123",
    "callee_id": "user_456",
    "callee_name": "李四",
    "action": "accept", // accept, reject, busy
    "sdp": {
      "type": "answer",
      "sdp": "v=0\r\no=- 0987654321 0987654321 IN IP4 192.168.1.101\r\n..."
    },
    "answered_at": "2024-12-21T10:30:15Z"
  },
  "timestamp": 1703123456789
}
```

#### ICE候选交换
```javascript
{
  "type": "call_ice_candidate",
  "data": {
    "call_id": "call_12345",
    "sender_id": "user_123",
    "receiver_id": "user_456",
    "candidate": {
      "candidate": "candidate:1 1 UDP 2113667326 192.168.1.100 54400 typ host",
      "sdpMid": "audio",
      "sdpMLineIndex": 0
    }
  },
  "timestamp": 1703123456789
}
```

#### 挂断通话
```javascript
{
  "type": "call_hangup",
  "data": {
    "call_id": "call_12345",
    "initiator_id": "user_123",
    "reason": "normal", // normal, busy, network_error, timeout
    "duration": 180, // 通话时长（秒）
    "ended_at": "2024-12-21T10:33:00Z"
  },
  "timestamp": 1703123456789
}
```

### 群组通话

#### 加入会议室
```javascript
{
  "type": "conference_join",
  "data": {
    "room_id": "conf_12345",
    "user_id": "user_123",
    "user_name": "张三",
    "user_avatar": "https://cdn.yunpan.com/avatars/user_123.jpg",
    "media_constraints": {
      "audio": true,
      "video": false
    },
    "joined_at": "2024-12-21T10:30:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 离开会议室
```javascript
{
  "type": "conference_leave",
  "data": {
    "room_id": "conf_12345",
    "user_id": "user_123",
    "user_name": "张三",
    "reason": "voluntary", // voluntary, kicked, network_error
    "left_at": "2024-12-21T10:45:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 媒体控制
```javascript
{
  "type": "conference_media_control",
  "data": {
    "room_id": "conf_12345",
    "user_id": "user_123",
    "controls": {
      "audio_muted": false,
      "video_enabled": true,
      "screen_sharing": false
    },
    "updated_at": "2024-12-21T10:35:00Z"
  },
  "timestamp": 1703123456789
}
```

## 系统通知

### 通知类型

#### 系统公告
```javascript
{
  "type": "system_announcement",
  "data": {
    "announcement_id": "announce_12345",
    "title": "系统维护通知",
    "content": "系统将于今晚22:00-24:00进行维护，期间服务可能中断。",
    "priority": "high", // low, medium, high, urgent
    "target_users": "all", // all, specific_users, user_groups
    "valid_until": "2024-12-22T00:00:00Z",
    "created_at": "2024-12-21T10:00:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 文件分享通知
```javascript
{
  "type": "file_share_notification",
  "data": {
    "share_id": "share_12345",
    "file_id": "file_67890",
    "filename": "重要文档.pdf",
    "sharer_id": "user_123",
    "sharer_name": "张三",
    "recipient_id": "user_456",
    "permissions": ["read", "download"],
    "message": "这个文档很重要，请查阅。",
    "expires_at": "2024-12-28T23:59:59Z",
    "shared_at": "2024-12-21T10:30:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 协作邀请通知
```javascript
{
  "type": "collaboration_invitation",
  "data": {
    "invitation_id": "invite_12345",
    "document_id": "doc_67890",
    "document_title": "项目计划书",
    "inviter_id": "user_123",
    "inviter_name": "张三",
    "invitee_id": "user_456",
    "permissions": ["read", "write", "comment"],
    "message": "请协助完成这个项目计划书。",
    "expires_at": "2024-12-28T23:59:59Z",
    "invited_at": "2024-12-21T10:30:00Z"
  },
  "timestamp": 1703123456789
}
```

#### 存储空间警告
```javascript
{
  "type": "storage_warning",
  "data": {
    "user_id": "user_123",
    "current_usage": 9663676416, // 当前使用量（字节）
    "total_quota": 10737418240, // 总配额（字节）
    "usage_percentage": 90,
    "warning_level": "high", // low(80%), medium(90%), high(95%), critical(98%)
    "suggested_actions": [
      "删除不必要的文件",
      "清空回收站",
      "升级存储计划"
    ]
  },
  "timestamp": 1703123456789
}
```

### 通知配置

#### 用户通知偏好
```javascript
{
  "type": "notification_preferences",
  "data": {
    "user_id": "user_123",
    "preferences": {
      "file_operations": {
        "enabled": true,
        "types": ["shared_with_me", "collaboration_invite"]
      },
      "chat_messages": {
        "enabled": true,
        "types": ["direct_message", "group_message", "mentions"]
      },
      "system_announcements": {
        "enabled": true,
        "priority_filter": "medium" // all, high, medium, low
      },
      "do_not_disturb": {
        "enabled": false,
        "start_time": "22:00",
        "end_time": "08:00",
        "timezone": "Asia/Shanghai"
      }
    }
  },
  "timestamp": 1703123456789
}
```

## 错误处理

### 错误消息格式
```javascript
{
  "type": "error",
  "data": {
    "error_code": "FILE_NOT_FOUND",
    "error_message": "指定的文件不存在",
    "error_details": {
      "file_id": "file_12345",
      "requested_at": "2024-12-21T10:30:00Z"
    },
    "retry_after": 5000, // 建议重试间隔（毫秒）
    "is_recoverable": true
  },
  "timestamp": 1703123456789
}
```

### 错误代码

#### 认证错误 (4000-4099)
- **4000**: `AUTH_REQUIRED` - 需要认证
- **4001**: `AUTH_FAILED` - 认证失败
- **4002**: `TOKEN_EXPIRED` - Token已过期
- **4003**: `TOKEN_INVALID` - Token无效
- **4004**: `AUTH_PERMISSION_DENIED` - 权限不足

#### 请求错误 (4100-4199)
- **4100**: `BAD_REQUEST` - 请求格式错误
- **4101**: `INVALID_MESSAGE_TYPE` - 无效的消息类型
- **4102**: `MISSING_REQUIRED_FIELD` - 缺少必需字段
- **4103**: `INVALID_FIELD_VALUE` - 字段值无效
- **4104**: `MESSAGE_TOO_LARGE` - 消息过大

#### 资源错误 (4200-4299)
- **4200**: `RESOURCE_NOT_FOUND` - 资源不存在
- **4201**: `FILE_NOT_FOUND` - 文件不存在
- **4202**: `ROOM_NOT_FOUND` - 聊天室不存在
- **4203**: `USER_NOT_FOUND` - 用户不存在
- **4204**: `DOCUMENT_NOT_FOUND` - 文档不存在

#### 业务逻辑错误 (4300-4399)
- **4300**: `OPERATION_NOT_ALLOWED` - 操作不被允许
- **4301**: `FILE_LOCKED` - 文件被锁定
- **4302**: `STORAGE_QUOTA_EXCEEDED` - 存储配额超限
- **4303**: `CONCURRENT_EDIT_CONFLICT` - 并发编辑冲突
- **4304**: `CALL_ALREADY_EXISTS` - 通话已存在

#### 系统错误 (5000-5099)
- **5000**: `INTERNAL_SERVER_ERROR` - 服务器内部错误
- **5001**: `SERVICE_UNAVAILABLE` - 服务不可用
- **5002**: `DATABASE_ERROR` - 数据库错误
- **5003**: `STORAGE_ERROR` - 存储服务错误
- **5004**: `NETWORK_ERROR` - 网络错误

### WebSocket与HTTP错误码映射表

| WebSocket错误码 | WebSocket错误名 | API错误码 | HTTP状态码 | HTTP状态描述 | 说明 |
|----------------|-----------------|-----------|-------------|-------------|------|
| 4000 | AUTH_REQUIRED | AUTH_VALIDATION_TOKEN_REQUIRED | 401 | Unauthorized | 需要身份认证 |
| 4001 | AUTH_FAILED | AUTH_BUSINESS_LOGIN_FAILED | 401 | Unauthorized | 认证失败 |
| 4002 | TOKEN_EXPIRED | AUTH_VALIDATION_TOKEN_EXPIRED | 401 | Unauthorized | Token过期 |
| 4003 | TOKEN_INVALID | AUTH_VALIDATION_TOKEN_INVALID | 401 | Unauthorized | Token无效 |
| 4004 | AUTH_PERMISSION_DENIED | AUTH_PERMISSION_ACCESS_DENIED | 403 | Forbidden | 权限不足 |
| 4100 | BAD_REQUEST | SYSTEM_VALIDATION_REQUEST_INVALID | 400 | Bad Request | 请求格式错误 |
| 4101 | INVALID_MESSAGE_TYPE | SYSTEM_VALIDATION_MESSAGE_TYPE_INVALID | 400 | Bad Request | 消息类型无效 |
| 4102 | MISSING_REQUIRED_FIELD | SYSTEM_VALIDATION_FIELD_REQUIRED | 400 | Bad Request | 缺少必需字段 |
| 4103 | INVALID_FIELD_VALUE | SYSTEM_VALIDATION_FIELD_VALUE_INVALID | 400 | Bad Request | 字段值无效 |
| 4104 | MESSAGE_TOO_LARGE | SYSTEM_VALIDATION_MESSAGE_TOO_LARGE | 413 | Payload Too Large | 消息过大 |
| 4200 | RESOURCE_NOT_FOUND | SYSTEM_RESOURCE_NOT_FOUND | 404 | Not Found | 资源不存在 |
| 4201 | FILE_NOT_FOUND | FILE_RESOURCE_NOT_FOUND | 404 | Not Found | 文件不存在 |
| 4202 | ROOM_NOT_FOUND | TEAM_RESOURCE_ROOM_NOT_FOUND | 404 | Not Found | 聊天室不存在 |
| 4203 | USER_NOT_FOUND | USER_RESOURCE_NOT_FOUND | 404 | Not Found | 用户不存在 |
| 4204 | DOCUMENT_NOT_FOUND | FILE_RESOURCE_DOCUMENT_NOT_FOUND | 404 | Not Found | 文档不存在 |
| 4300 | OPERATION_NOT_ALLOWED | SYSTEM_BUSINESS_OPERATION_NOT_ALLOWED | 405 | Method Not Allowed | 操作不被允许 |
| 4301 | FILE_LOCKED | FILE_BUSINESS_LOCKED | 423 | Locked | 文件被锁定 |
| 4302 | STORAGE_QUOTA_EXCEEDED | STORAGE_BUSINESS_QUOTA_EXCEEDED | 422 | Unprocessable Entity | 存储配额超限 |
| 4303 | CONCURRENT_EDIT_CONFLICT | FILE_BUSINESS_CONCURRENT_EDIT_CONFLICT | 409 | Conflict | 并发编辑冲突 |
| 4304 | CALL_ALREADY_EXISTS | TEAM_BUSINESS_CALL_ALREADY_EXISTS | 409 | Conflict | 通话已存在 |
| 5000 | INTERNAL_SERVER_ERROR | SYSTEM_SYSTEM_INTERNAL_ERROR | 500 | Internal Server Error | 服务器内部错误 |
| 5001 | SERVICE_UNAVAILABLE | SYSTEM_SYSTEM_SERVICE_UNAVAILABLE | 503 | Service Unavailable | 服务不可用 |
| 5002 | DATABASE_ERROR | SYSTEM_SYSTEM_DATABASE_ERROR | 500 | Internal Server Error | 数据库错误 |
| 5003 | STORAGE_ERROR | STORAGE_SYSTEM_ERROR | 500 | Internal Server Error | 存储服务错误 |
| 5004 | NETWORK_ERROR | SYSTEM_SYSTEM_NETWORK_ERROR | 502 | Bad Gateway | 网络错误 |

### 错误恢复策略

#### 自动重试
```javascript
const retryStrategy = {
  maxRetries: 3,
  baseDelay: 1000, // 基础延迟（毫秒）
  maxDelay: 30000, // 最大延迟（毫秒）
  backoffFactor: 2, // 退避因子
  retryableErrors: [
    "NETWORK_ERROR",
    "SERVICE_UNAVAILABLE", 
    "INTERNAL_SERVER_ERROR"
  ]
}
```

#### 降级策略
```javascript
const fallbackStrategy = {
  file_sync: "cache_operations", // 缓存操作，稍后同步
  chat: "queue_messages", // 队列消息，重连后发送
  collaboration: "readonly_mode", // 只读模式
  voice_call: "disable_feature" // 禁用功能
}
```

## 安全机制

### 认证安全

#### JWT Token验证
```javascript
// Token格式
{
  "alg": "HS256",
  "typ": "JWT"
}
{
  "sub": "user_123",
  "iat": 1703123456,
  "exp": 1703130656,
  "scope": ["file_read", "file_write", "chat", "collaborate"],
  "client_type": "web"
}
```

#### Token刷新机制
```javascript
{
  "type": "token_refresh",
  "data": {
    "access_token": "new_access_token",
    "expires_in": 7200, // 2小时
    "refresh_token": "new_refresh_token"
  },
  "timestamp": 1703123456789
}
```

### 传输安全

#### WSS加密
- **协议**：WSS (WebSocket Secure)
- **TLS版本**：TLS 1.2+
- **加密套件**：ECDHE-RSA-AES256-GCM-SHA384

#### 消息完整性
```javascript
// 消息签名（可选，用于重要消息）
{
  "type": "chat_message",
  "data": { /* 消息内容 */ },
  "signature": "HMAC-SHA256签名",
  "timestamp": 1703123456789
}
```

### 权限控制

#### 操作权限验证
```javascript
const permissions = {
  file_operations: {
    read: ["file_read"],
    write: ["file_write"],
    delete: ["file_delete"],
    share: ["file_share"]
  },
  chat: {
    send_message: ["chat_send"],
    create_room: ["chat_admin"],
    invite_user: ["chat_invite"]
  },
  collaboration: {
    join_session: ["collab_read"],
    edit_document: ["collab_write"],
    manage_session: ["collab_admin"]
  }
}
```

### 频率限制

#### 消息频率限制
```javascript
const rateLimits = {
  chat_message: {
    rate: 10, // 每分钟10条
    window: 60000, // 1分钟窗口
    burst: 3 // 突发允许3条
  },
  file_operation: {
    rate: 100, // 每分钟100次操作
    window: 60000,
    burst: 10
  },
  collaboration_operation: {
    rate: 1000, // 每分钟1000次操作（打字频率高）
    window: 60000,
    burst: 50
  }
}
```

## 性能优化

### 消息压缩

#### Per-Message-Deflate
```javascript
// 启用压缩的消息头
{
  "type": "chat_message",
  "compressed": true,
  "original_size": 1024,
  "compressed_size": 256,
  "data": { /* 压缩后的数据 */ }
}
```

### 批量处理

#### 消息批量发送
```javascript
{
  "type": "batch_messages",
  "data": {
    "batch_id": "batch_12345",
    "messages": [
      { /* 消息1 */ },
      { /* 消息2 */ },
      { /* 消息3 */ }
    ],
    "batch_size": 3
  },
  "timestamp": 1703123456789
}
```

### 缓存策略

#### 客户端缓存
```javascript
const cacheStrategy = {
  chat_messages: {
    maxSize: 1000, // 最大缓存消息数
    ttl: 3600000 // 1小时TTL
  },
  user_profiles: {
    maxSize: 500,
    ttl: 1800000 // 30分钟TTL
  },
  file_metadata: {
    maxSize: 2000,
    ttl: 600000 // 10分钟TTL
  }
}
```

### 连接池管理

#### 服务端连接池
```javascript
const connectionPool = {
  maxConnections: 10000, // 最大并发连接数
  keepAliveInterval: 30000, // 30秒心跳间隔
  connectionTimeout: 60000, // 60秒连接超时
  idleTimeout: 300000, // 5分钟空闲超时
  grouping: {
    byUser: true, // 按用户分组
    byRoom: true, // 按房间分组
    byDocument: true // 按文档分组
  }
}
```

## 客户端实现

### JavaScript客户端示例

#### 连接管理
```javascript
class YunPanWebSocket {
  constructor(url, options = {}) {
    this.url = url;
    this.options = {
      reconnectInterval: 1000,
      maxReconnectInterval: 30000,
      reconnectBackoffRate: 1.5,
      maxReconnectAttempts: Infinity,
      ...options
    };
    this.ws = null;
    this.reconnectAttempts = 0;
    this.listeners = new Map();
    this.messageQueue = [];
    this.isAuthenticated = false;
  }

  connect(token) {
    return new Promise((resolve, reject) => {
      this.ws = new WebSocket(this.url, 'yunpan-v1');
      
      this.ws.onopen = () => {
        console.log('WebSocket连接已建立');
        this.authenticate(token)
          .then(resolve)
          .catch(reject);
      };

      this.ws.onmessage = (event) => {
        this.handleMessage(JSON.parse(event.data));
      };

      this.ws.onclose = (event) => {
        console.log('WebSocket连接已关闭', event.code, event.reason);
        this.isAuthenticated = false;
        this.scheduleReconnect();
      };

      this.ws.onerror = (error) => {
        console.error('WebSocket错误', error);
        reject(error);
      };
    });
  }

  authenticate(token) {
    return new Promise((resolve, reject) => {
      const authMessage = {
        type: 'auth',
        data: {
          token: token,
          client_id: this.generateClientId(),
          client_type: 'web',
          version: '1.0.0'
        },
        timestamp: Date.now()
      };

      this.send(authMessage);

      const authTimeout = setTimeout(() => {
        reject(new Error('认证超时'));
      }, 10000);

      this.once('auth_response', (response) => {
        clearTimeout(authTimeout);
        if (response.data.status === 'success') {
          this.isAuthenticated = true;
          this.flushMessageQueue();
          resolve(response.data);
        } else {
          reject(new Error(response.data.message));
        }
      });
    });
  }

  send(message) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      if (this.isAuthenticated || message.type === 'auth') {
        this.ws.send(JSON.stringify(message));
      } else {
        this.messageQueue.push(message);
      }
    } else {
      this.messageQueue.push(message);
    }
  }

  on(type, callback) {
    if (!this.listeners.has(type)) {
      this.listeners.set(type, []);
    }
    this.listeners.get(type).push(callback);
  }

  once(type, callback) {
    const wrapper = (data) => {
      this.off(type, wrapper);
      callback(data);
    };
    this.on(type, wrapper);
  }

  off(type, callback) {
    if (this.listeners.has(type)) {
      const callbacks = this.listeners.get(type);
      const index = callbacks.indexOf(callback);
      if (index > -1) {
        callbacks.splice(index, 1);
      }
    }
  }

  handleMessage(message) {
    console.log('收到消息', message);
    
    // 处理心跳
    if (message.type === 'ping') {
      this.send({ type: 'pong', timestamp: Date.now() });
      return;
    }

    // 触发监听器
    if (this.listeners.has(message.type)) {
      this.listeners.get(message.type).forEach(callback => {
        try {
          callback(message);
        } catch (error) {
          console.error('消息处理错误', error);
        }
      });
    }
  }

  scheduleReconnect() {
    if (this.reconnectAttempts < this.options.maxReconnectAttempts) {
      const interval = Math.min(
        this.options.reconnectInterval * Math.pow(this.options.reconnectBackoffRate, this.reconnectAttempts),
        this.options.maxReconnectInterval
      );

      setTimeout(() => {
        this.reconnectAttempts++;
        console.log(`尝试重连 (${this.reconnectAttempts}/${this.options.maxReconnectAttempts})`);
        this.connect(this.lastToken);
      }, interval);
    }
  }

  flushMessageQueue() {
    while (this.messageQueue.length > 0) {
      const message = this.messageQueue.shift();
      this.send(message);
    }
  }

  generateClientId() {
    return 'web_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9);
  }

  close() {
    if (this.ws) {
      this.ws.close(1000, 'Client initiated close');
    }
  }
}
```

#### 文件同步客户端
```javascript
class FileSyncManager {
  constructor(websocket) {
    this.ws = websocket;
    this.syncedFiles = new Map();
    this.setupEventListeners();
  }

  setupEventListeners() {
    this.ws.on('file_created', (message) => {
      this.handleFileCreated(message.data);
    });

    this.ws.on('file_updated', (message) => {
      this.handleFileUpdated(message.data);
    });

    this.ws.on('file_deleted', (message) => {
      this.handleFileDeleted(message.data);
    });

    this.ws.on('sync_status', (message) => {
      this.handleSyncStatus(message.data);
    });
  }

  handleFileCreated(data) {
    console.log('文件已创建', data);
    this.syncedFiles.set(data.file_id, data);
    this.notifyUI('file_created', data);
  }

  handleFileUpdated(data) {
    console.log('文件已更新', data);
    const existingFile = this.syncedFiles.get(data.file_id);
    if (existingFile) {
      Object.assign(existingFile, data.changes);
      this.notifyUI('file_updated', existingFile);
    }
  }

  handleFileDeleted(data) {
    console.log('文件已删除', data);
    this.syncedFiles.delete(data.file_id);
    this.notifyUI('file_deleted', data);
  }

  handleSyncStatus(data) {
    console.log('同步状态更新', data);
    this.notifyUI('sync_status', data);
  }

  notifyUI(event, data) {
    // 通知UI更新
    window.dispatchEvent(new CustomEvent(`filesync:${event}`, {
      detail: data
    }));
  }
}
```

#### 聊天客户端
```javascript
class ChatManager {
  constructor(websocket) {
    this.ws = websocket;
    this.rooms = new Map();
    this.typingUsers = new Map();
    this.setupEventListeners();
  }

  setupEventListeners() {
    this.ws.on('chat_message', (message) => {
      this.handleChatMessage(message.data);
    });

    this.ws.on('chat_typing', (message) => {
      this.handleTyping(message.data);
    });

    this.ws.on('chat_read', (message) => {
      this.handleMessageRead(message.data);
    });
  }

  sendMessage(roomId, content, messageType = 'text', attachments = null) {
    const message = {
      type: 'chat_message',
      data: {
        room_id: roomId,
        content: content,
        message_type: messageType,
        attachments: attachments,
        created_at: new Date().toISOString()
      },
      timestamp: Date.now()
    };

    this.ws.send(message);
  }

  sendTyping(roomId, isTyping) {
    const message = {
      type: 'chat_typing',
      data: {
        room_id: roomId,
        is_typing: isTyping
      },
      timestamp: Date.now()
    };

    this.ws.send(message);
  }

  markAsRead(roomId, messageId) {
    const message = {
      type: 'chat_read',
      data: {
        room_id: roomId,
        message_id: messageId,
        read_at: new Date().toISOString()
      },
      timestamp: Date.now()
    };

    this.ws.send(message);
  }

  handleChatMessage(data) {
    if (!this.rooms.has(data.room_id)) {
      this.rooms.set(data.room_id, { messages: [] });
    }

    this.rooms.get(data.room_id).messages.push(data);
    this.notifyUI('new_message', data);
  }

  handleTyping(data) {
    const key = `${data.room_id}_${data.user_id}`;
    
    if (data.is_typing) {
      this.typingUsers.set(key, {
        ...data,
        timestamp: Date.now()
      });
    } else {
      this.typingUsers.delete(key);
    }

    this.notifyUI('typing_status', {
      room_id: data.room_id,
      typing_users: Array.from(this.typingUsers.values())
        .filter(user => user.room_id === data.room_id)
    });
  }

  handleMessageRead(data) {
    this.notifyUI('message_read', data);
  }

  notifyUI(event, data) {
    window.dispatchEvent(new CustomEvent(`chat:${event}`, {
      detail: data
    }));
  }
}
```

## 服务端实现

### Go服务端示例

#### WebSocket处理器
```go
package websocket

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"

    "github.com/gorilla/websocket"
)

type Hub struct {
    clients    map[*Client]bool
    broadcast  chan []byte
    register   chan *Client
    unregister chan *Client
    mutex      sync.RWMutex
}

type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    send     chan []byte
    userID   string
    clientID string
    authenticated bool
}

type Message struct {
    Type      string      `json:"type"`
    ID        string      `json:"id,omitempty"`
    Data      interface{} `json:"data"`
    Timestamp int64       `json:"timestamp"`
    Version   string      `json:"version,omitempty"`
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        // 在生产环境中应该检查Origin
        return true
    },
    EnableCompression: true,
}

func NewHub() *Hub {
    return &Hub{
        clients:    make(map[*Client]bool),
        broadcast:  make(chan []byte),
        register:   make(chan *Client),
        unregister: make(chan *Client),
    }
}

func (h *Hub) Run() {
    for {
        select {
        case client := <-h.register:
            h.mutex.Lock()
            h.clients[client] = true
            h.mutex.Unlock()
            log.Printf("客户端已连接: %s", client.clientID)

        case client := <-h.unregister:
            h.mutex.Lock()
            if _, ok := h.clients[client]; ok {
                delete(h.clients, client)
                close(client.send)
            }
            h.mutex.Unlock()
            log.Printf("客户端已断开: %s", client.clientID)

        case message := <-h.broadcast:
            h.mutex.RLock()
            for client := range h.clients {
                select {
                case client.send <- message:
                default:
                    close(client.send)
                    delete(h.clients, client)
                }
            }
            h.mutex.RUnlock()
        }
    }
}

func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Printf("WebSocket升级失败: %v", err)
        return
    }

    client := &Client{
        hub:     h,
        conn:    conn,
        send:    make(chan []byte, 256),
        clientID: generateClientID(),
    }

    client.hub.register <- client

    go client.writePump()
    go client.readPump()
}

func (c *Client) readPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(1024 * 1024) // 1MB
    c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
        return nil
    })

    for {
        _, messageData, err := c.conn.ReadMessage()
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket错误: %v", err)
            }
            break
        }

        var msg Message
        if err := json.Unmarshal(messageData, &msg); err != nil {
            log.Printf("消息解析失败: %v", err)
            continue
        }

        c.handleMessage(&msg)
    }
}

func (c *Client) writePump() {
    ticker := time.NewTicker(30 * time.Second)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case message, ok := <-c.send:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }

            w, err := c.conn.NextWriter(websocket.TextMessage)
            if err != nil {
                return
            }
            w.Write(message)

            // 批量发送队列中的消息
            n := len(c.send)
            for i := 0; i < n; i++ {
                w.Write([]byte{'\n'})
                w.Write(<-c.send)
            }

            if err := w.Close(); err != nil {
                return
            }

        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
            if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                return
            }
        }
    }
}

func (c *Client) handleMessage(msg *Message) {
    switch msg.Type {
    case "auth":
        c.handleAuth(msg)
    case "ping":
        c.handlePing(msg)
    case "chat_message":
        if c.authenticated {
            c.handleChatMessage(msg)
        }
    case "file_operation":
        if c.authenticated {
            c.handleFileOperation(msg)
        }
    case "collab_operation":
        if c.authenticated {
            c.handleCollabOperation(msg)
        }
    default:
        log.Printf("未知消息类型: %s", msg.Type)
    }
}

func (c *Client) handleAuth(msg *Message) {
    authData := msg.Data.(map[string]interface{})
    token := authData["token"].(string)

    // 验证JWT Token
    userID, err := validateJWTToken(token)
    if err != nil {
        c.sendError("AUTH_FAILED", "认证失败: "+err.Error())
        return
    }

    c.userID = userID
    c.authenticated = true

    response := &Message{
        Type: "auth_response",
        Data: map[string]interface{}{
            "status":  "success",
            "user_id": userID,
            "session_id": generateSessionID(),
            "permissions": getUserPermissions(userID),
        },
        Timestamp: time.Now().UnixMilli(),
    }

    c.sendMessage(response)
}

func (c *Client) handlePing(msg *Message) {
    response := &Message{
        Type:      "pong",
        Timestamp: time.Now().UnixMilli(),
    }
    c.sendMessage(response)
}

func (c *Client) sendMessage(msg *Message) {
    data, err := json.Marshal(msg)
    if err != nil {
        log.Printf("消息序列化失败: %v", err)
        return
    }

    select {
    case c.send <- data:
    default:
        close(c.send)
        delete(c.hub.clients, c)
    }
}

func (c *Client) sendError(code, message string) {
    errorMsg := &Message{
        Type: "error",
        Data: map[string]interface{}{
            "error_code":    code,
            "error_message": message,
        },
        Timestamp: time.Now().UnixMilli(),
    }
    c.sendMessage(errorMsg)
}

// 辅助函数
func generateClientID() string {
    return fmt.Sprintf("client_%d_%s", time.Now().UnixNano(), randomString(8))
}

func generateSessionID() string {
    return fmt.Sprintf("sess_%d_%s", time.Now().UnixNano(), randomString(12))
}

func validateJWTToken(token string) (string, error) {
    // 实现JWT Token验证逻辑
    // 返回用户ID或错误
    return "", nil
}

func getUserPermissions(userID string) []string {
    // 获取用户权限
    return []string{"file_read", "file_write", "chat", "collaborate"}
}

func randomString(length int) string {
    // 生成随机字符串
    return ""
}
```

#### 文件同步服务
```go
package fileSync

import (
    "encoding/json"
    "log"
    "time"
)

type FileSyncService struct {
    hub *websocket.Hub
}

type FileEvent struct {
    Type     string      `json:"type"`
    FileID   string      `json:"file_id"`
    Filename string      `json:"filename"`
    UserID   string      `json:"user_id"`
    Data     interface{} `json:"data"`
}

func NewFileSyncService(hub *websocket.Hub) *FileSyncService {
    return &FileSyncService{
        hub: hub,
    }
}

func (fs *FileSyncService) NotifyFileCreated(fileID, filename, userID string, fileData interface{}) {
    event := &FileEvent{
        Type:     "file_created",
        FileID:   fileID,
        Filename: filename,
        UserID:   userID,
        Data:     fileData,
    }

    fs.broadcastToAuthorizedUsers(event, fileID)
}

func (fs *FileSyncService) NotifyFileUpdated(fileID, filename, userID string, changes interface{}) {
    event := &FileEvent{
        Type:     "file_updated",
        FileID:   fileID,
        Filename: filename,
        UserID:   userID,
        Data:     changes,
    }

    fs.broadcastToAuthorizedUsers(event, fileID)
}

func (fs *FileSyncService) NotifyFileDeleted(fileID, filename, userID string) {
    event := &FileEvent{
        Type:     "file_deleted",
        FileID:   fileID,
        Filename: filename,
        UserID:   userID,
    }

    fs.broadcastToAuthorizedUsers(event, fileID)
}

func (fs *FileSyncService) broadcastToAuthorizedUsers(event *FileEvent, fileID string) {
    message := &websocket.Message{
        Type:      event.Type,
        Data:      event.Data,
        Timestamp: time.Now().UnixMilli(),
    }

    data, err := json.Marshal(message)
    if err != nil {
        log.Printf("文件事件序列化失败: %v", err)
        return
    }

    // 获取有权限访问该文件的用户列表
    authorizedUsers := getAuthorizedUsers(fileID)

    // 向授权用户广播消息
    fs.hub.BroadcastToUsers(data, authorizedUsers)
}

func getAuthorizedUsers(fileID string) []string {
    // 实现获取有权限访问文件的用户列表
    return []string{}
}
```

## 结语

本WebSocket协议规范为网络云盘系统提供了完整的实时通信解决方案，涵盖了文件同步、即时通讯、协作编辑、语音通话等核心功能。通过标准化的消息格式、完善的错误处理机制和强大的安全保障，确保系统的稳定性、可扩展性和用户体验。

### 关键特性总结

1. **统一的消息格式**：JSON格式，结构清晰，易于解析
2. **完整的认证机制**：JWT Token认证，权限细分控制
3. **可靠的连接管理**：心跳检测，自动重连，优雅降级
4. **高效的性能优化**：消息压缩，批量处理，连接池管理
5. **全面的错误处理**：错误码规范，重试策略，降级方案
6. **强大的安全机制**：WSS加密，权限验证，频率限制

通过遵循本规范，开发团队可以构建出高质量、高性能的实时通信功能，为用户提供流畅的云盘使用体验。

---

**文档版本**: v1.0  
**最后更新**: 2024年12月  
**维护团队**: 网络云盘开发团队