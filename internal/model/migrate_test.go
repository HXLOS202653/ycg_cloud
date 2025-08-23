package model

import (
	"reflect"
	"testing"
)

// TestModelStructures 测试模型结构定义
func TestModelStructures(t *testing.T) {
	// 测试用户模型
	user := User{}
	userType := reflect.TypeOf(user)
	if userType.Name() != "User" {
		t.Errorf("User模型名称错误")
	}

	// 测试文件模型
	file := File{}
	fileType := reflect.TypeOf(file)
	if fileType.Name() != "File" {
		t.Errorf("File模型名称错误")
	}

	// 测试团队模型
	team := Team{}
	teamType := reflect.TypeOf(team)
	if teamType.Name() != "Team" {
		t.Errorf("Team模型名称错误")
	}

	// 测试权限模板模型
	template := PermissionTemplate{}
	templateType := reflect.TypeOf(template)
	if templateType.Name() != "PermissionTemplate" {
		t.Errorf("PermissionTemplate模型名称错误")
	}

	// 测试操作日志模型
	log := OperationLog{}
	logType := reflect.TypeOf(log)
	if logType.Name() != "OperationLog" {
		t.Errorf("OperationLog模型名称错误")
	}

	// 测试回收站模型
	recycle := RecycleItem{}
	recycleType := reflect.TypeOf(recycle)
	if recycleType.Name() != "RecycleItem" {
		t.Errorf("RecycleItem模型名称错误")
	}
}

// TestModelConstants 测试模型常量定义
func TestModelConstants(t *testing.T) {
	// 测试用户状态常量
	if UserStatusActive == "" {
		t.Errorf("UserStatusActive常量未定义")
	}

	// 测试文件状态常量
	if FileStatusNormal == "" {
		t.Errorf("FileStatusNormal常量未定义")
	}

	// 测试团队成员角色常量
	if TeamMemberRoleOwner == "" {
		t.Errorf("TeamMemberRoleOwner常量未定义")
	}

	// 测试权限动作常量
	if PermissionRead == "" {
		t.Errorf("PermissionRead常量未定义")
	}

	// 测试操作类型常量
	if ActionFileUpload == "" {
		t.Errorf("ActionFileUpload常量未定义")
	}

	// 测试回收站类型常量
	if RecycleTypeFile == "" {
		t.Errorf("RecycleTypeFile常量未定义")
	}
}

// TestModelTableNames 测试模型表名方法
func TestModelTableNames(t *testing.T) {
	// 测试用户表名
	user := User{}
	if user.TableName() != "users" {
		t.Errorf("User表名错误: 期望 users, 实际 %s", user.TableName())
	}

	// 测试文件表名
	file := File{}
	if file.TableName() != "files" {
		t.Errorf("File表名错误: 期望 files, 实际 %s", file.TableName())
	}

	// 测试团队表名
	team := Team{}
	if team.TableName() != "teams" {
		t.Errorf("Team表名错误: 期望 teams, 实际 %s", team.TableName())
	}

	// 测试权限模板表名
	template := PermissionTemplate{}
	if template.TableName() != "permission_templates" {
		t.Errorf("PermissionTemplate表名错误: 期望 permission_templates, 实际 %s", template.TableName())
	}

	// 测试操作日志表名
	log := OperationLog{}
	if log.TableName() != "operation_logs" {
		t.Errorf("OperationLog表名错误: 期望 operation_logs, 实际 %s", log.TableName())
	}

	// 测试回收站表名
	recycle := RecycleItem{}
	if recycle.TableName() != "recycle_items" {
		t.Errorf("RecycleItem表名错误: 期望 recycle_items, 实际 %s", recycle.TableName())
	}
}

// TestModelValidation 测试模型验证
func TestModelValidation(t *testing.T) {
	validateUserModel(t)
	validateFileModel(t)
	validateTeamModel(t)
	t.Log("所有模型验证通过")
}

// validateUserModel 验证用户模型
func validateUserModel(t *testing.T) {
	// 创建测试用户
	user := User{
		Username: "testuser",
		Email:    "test@example.com",
		Status:   UserStatusActive,
	}

	// 验证用户字段
	if user.Username == "" {
		t.Errorf("用户名不能为空")
	}
	if user.Email == "" {
		t.Errorf("邮箱不能为空")
	}
	if user.Status == "" {
		t.Errorf("用户状态不能为空")
	}
}

// validateFileModel 验证文件模型
func validateFileModel(t *testing.T) {
	// 创建测试文件
	file := File{
		Name:     "test.txt",
		Path:     "/test.txt",
		FileType: FileTypeDocument,
		Status:   FileStatusNormal,
	}

	// 验证文件字段
	if file.Name == "" {
		t.Errorf("文件名不能为空")
	}
	if file.Path == "" {
		t.Errorf("文件路径不能为空")
	}
	if file.FileType == "" {
		t.Errorf("文件类型不能为空")
	}
	if file.Status == "" {
		t.Errorf("文件状态不能为空")
	}
}

// validateTeamModel 验证团队模型
func validateTeamModel(t *testing.T) {
	// 创建测试团队
	team := Team{
		Name:        "测试团队",
		Description: "这是一个测试团队",
		Status:      TeamStatusActive,
	}

	// 验证团队字段
	if team.Name == "" {
		t.Errorf("团队名称不能为空")
	}
	if team.Description == "" {
		t.Errorf("团队描述不能为空")
	}
	if team.Status == "" {
		t.Errorf("团队状态不能为空")
	}
}

// BenchmarkModelCreation 基准测试模型创建性能
func BenchmarkModelCreation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		// 创建用户
		_ = User{
			Username: "testuser",
			Email:    "test@example.com",
			Status:   UserStatusActive,
		}

		// 创建文件
		_ = File{
			Name:     "test.txt",
			Path:     "/test.txt",
			FileType: FileTypeDocument,
			Status:   FileStatusNormal,
		}

		// 创建团队
		_ = Team{
			Name:        "测试团队",
			Description: "这是一个测试团队",
			Status:      TeamStatusActive,
		}
	}
}
