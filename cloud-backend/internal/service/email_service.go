// Package service provides email sending functionality
package service

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"math/big"
	"net"
	"net/smtp"
	"strings"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/config"
)

// Email purpose constants
const (
	EmailPurposeRegistration  = "registration"
	EmailPurposePasswordReset = "password_reset"
)

// EmailService handles email operations
type EmailService struct {
	config *config.Config
}

// NewEmailService creates a new email service
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		config: cfg,
	}
}

// EmailTemplate represents email template data
type EmailTemplate struct {
	Subject string
	Body    string
}

// VerificationCodeEmail generates verification code email template
func (s *EmailService) VerificationCodeEmail(code, _ /* email */, purpose string) *EmailTemplate {
	var subject, bodyTemplate string

	switch purpose {
	case EmailPurposeRegistration:
		subject = "YCG Cloud - 邮箱验证码"
		bodyTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>邮箱验证</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; background-color: #f8f9fa; padding: 20px; border-radius: 8px; }
        .content { padding: 20px; }
        .code { font-size: 32px; font-weight: bold; color: #007bff; text-align: center; 
                background-color: #f8f9fa; padding: 15px; border-radius: 8px; margin: 20px 0; }
        .footer { text-align: center; color: #666; font-size: 14px; margin-top: 30px; }
        .warning { color: #dc3545; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🔐 YCG Cloud 邮箱验证</h2>
        </div>
        <div class="content">
            <p>您好！</p>
            <p>感谢您注册 YCG Cloud 云盘服务。请使用以下验证码完成邮箱验证：</p>
            
            <div class="code">%s</div>
            
            <p><strong>重要提醒：</strong></p>
            <ul>
                <li>验证码有效期为 <span class="warning">5分钟</span></li>
                <li>请勿向他人泄露此验证码</li>
                <li>如果您没有申请注册，请忽略此邮件</li>
            </ul>
            
            <p>如有疑问，请联系我们的客服团队。</p>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复</p>
            <p>© 2024 YCG Cloud. 保留所有权利。</p>
        </div>
    </div>
</body>
</html>`
	case EmailPurposePasswordReset:
		subject = "YCG Cloud - 密码重置验证码"
		bodyTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>密码重置</title>
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; }
        .header { text-align: center; background-color: #f8f9fa; padding: 20px; border-radius: 8px; }
        .content { padding: 20px; }
        .code { font-size: 32px; font-weight: bold; color: #dc3545; text-align: center; 
                background-color: #f8f9fa; padding: 15px; border-radius: 8px; margin: 20px 0; }
        .footer { text-align: center; color: #666; font-size: 14px; margin-top: 30px; }
        .warning { color: #dc3545; font-weight: bold; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>🔒 YCG Cloud 密码重置</h2>
        </div>
        <div class="content">
            <p>您好！</p>
            <p>您正在重置 YCG Cloud 账户密码。请使用以下验证码：</p>
            
            <div class="code">%s</div>
            
            <p><span class="warning">安全提醒：</span></p>
            <ul>
                <li>验证码有效期为 <span class="warning">10分钟</span></li>
                <li>仅用于密码重置，请勿向他人提供</li>
                <li>如果您没有申请重置密码，请立即联系客服</li>
            </ul>
        </div>
        <div class="footer">
            <p>此邮件由系统自动发送，请勿回复</p>
            <p>© 2024 YCG Cloud. 保留所有权利。</p>
        </div>
    </div>
</body>
</html>`
	default:
		subject = "YCG Cloud - 验证码"
		bodyTemplate = `
<p>您的验证码是：<strong>%s</strong></p>
<p>有效期5分钟，请及时使用。</p>`
	}

	body := fmt.Sprintf(bodyTemplate, code)
	return &EmailTemplate{
		Subject: subject,
		Body:    body,
	}
}

// SendEmail sends an email with the given template
func (s *EmailService) SendEmail(to, subject, body string) error {
	// Check if email service is enabled
	if !s.config.Email.Enabled {
		fmt.Printf("\n=== 📧 邮件发送已禁用 ===\n")
		fmt.Printf("收件人: %s\n", to)
		fmt.Printf("主题: %s\n", subject)
		fmt.Printf("提示: 邮件服务已禁用，如需启用请设置 EMAIL_ENABLED=true\n")
		fmt.Printf("========================\n\n")
		return nil
	}

	// Check if SMTP configuration is available
	if s.config.Email.SMTPHost == "" || s.config.Email.SMTPUser == "" || s.config.Email.SMTPPassword == "" {
		// Fall back to console output for development
		fmt.Printf("\n=== 📧 邮件发送模拟 (缺少SMTP配置) ===\n")
		fmt.Printf("收件人: %s\n", to)
		fmt.Printf("主题: %s\n", subject)
		fmt.Printf("内容: %s\n", body)
		fmt.Printf("提示: 请配置环境变量: EMAIL_SMTP_HOST, EMAIL_SMTP_USER, EMAIL_SMTP_PASSWORD\n")
		fmt.Printf("========================\n\n")
		return nil
	}

	// Use real SMTP sending
	return s.sendSMTPEmail(to, subject, body)
}

// SendVerificationCode sends a verification code email
func (s *EmailService) SendVerificationCode(email, purpose string) (string, error) {
	// Generate 6-digit verification code
	code, err := s.generateVerificationCode()
	if err != nil {
		return "", fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Get email template
	template := s.VerificationCodeEmail(code, email, purpose)

	// Send email
	if err := s.SendEmail(email, template.Subject, template.Body); err != nil {
		return "", fmt.Errorf("failed to send verification email: %w", err)
	}

	return code, nil
}

// generateVerificationCode generates a 6-digit random verification code
func (s *EmailService) generateVerificationCode() (string, error) {
	// Generate random number between 100000 and 999999
	minVal := int64(100000)
	maxVal := int64(999999)

	n, err := rand.Int(rand.Reader, big.NewInt(maxVal-minVal+1))
	if err != nil {
		return "", err
	}

	code := n.Int64() + minVal
	return fmt.Sprintf("%06d", code), nil
}

// sendSMTPEmail sends email using SMTP with STARTTLS
func (s *EmailService) sendSMTPEmail(to, subject, body string) error {
	// Build SMTP server address
	smtpAddr := fmt.Sprintf("%s:%d", s.config.Email.SMTPHost, s.config.Email.SMTPPort)

	// Set up authentication
	auth := smtp.PlainAuth("", s.config.Email.SMTPUser, s.config.Email.SMTPPassword, s.config.Email.SMTPHost)

	// Build sender email
	from := s.config.Email.SMTPFrom
	if s.config.Email.SMTPFromName != "" {
		from = fmt.Sprintf("%s <%s>", s.config.Email.SMTPFromName, s.config.Email.SMTPFrom)
	}

	// Build email message
	message := s.buildEmailMessage(from, to, subject, body)

	// For QQ mail and other SMTP servers that require STARTTLS
	if s.config.Email.SMTPPort == 587 && !s.config.Email.SMTPSecure {
		return s.sendWithSTARTTLS(smtpAddr, auth, from, []string{to}, message)
	}

	// For direct SSL connection (port 465)
	if s.config.Email.SMTPSecure {
		return s.sendWithSSL(smtpAddr, auth, from, []string{to}, message)
	}

	// Standard SMTP (usually port 25, not recommended for external services)
	return smtp.SendMail(smtpAddr, auth, s.config.Email.SMTPFrom, []string{to}, message)
}

// sendWithSTARTTLS sends email using STARTTLS (commonly used for port 587)
func (s *EmailService) sendWithSTARTTLS(addr string, auth smtp.Auth, _ string, to []string, message []byte) error {
	// Connect to SMTP server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			// Log connection close error but don't fail the operation
			fmt.Printf("Failed to close SMTP connection: %v\n", err)
		}
	}()

	// Create SMTP client
	client, err := smtp.NewClient(conn, s.config.Email.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer func() {
		if err := client.Quit(); err != nil {
			// Log client quit error but don't fail the operation
			fmt.Printf("Failed to quit SMTP client: %v\n", err)
		}
	}()

	// Start TLS if supported
	if ok, _ := client.Extension("STARTTLS"); ok {
		tlsConfig := &tls.Config{
			ServerName: s.config.Email.SMTPHost,
			MinVersion: tls.VersionTLS12, // G402: Set minimum TLS version to 1.2
		}
		if err := client.StartTLS(tlsConfig); err != nil {
			return fmt.Errorf("failed to start TLS: %w", err)
		}
	}

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender
	if err := client.Mail(s.config.Email.SMTPFrom); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			// Log writer close error but don't fail the operation
			fmt.Printf("Failed to close SMTP writer: %v\n", err)
		}
	}()

	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("failed to write email data: %w", err)
	}

	fmt.Printf("\n=== ✅ 邮件发送成功 ===\n")
	fmt.Printf("收件人: %s\n", strings.Join(to, ", "))
	fmt.Printf("主题: %s\n", s.extractSubject(string(message)))
	fmt.Printf("SMTP服务器: %s\n", addr)
	fmt.Printf("===================\n\n")

	return nil
}

// sendWithSSL sends email using direct SSL connection (port 465)
func (s *EmailService) sendWithSSL(addr string, auth smtp.Auth, _ string, to []string, message []byte) error {
	tlsConfig := &tls.Config{
		ServerName: s.config.Email.SMTPHost,
		MinVersion: tls.VersionTLS12, // G402: Set minimum TLS version to 1.2
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server with SSL: %w", err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			// Log connection close error but don't fail the operation
		}
	}()

	client, err := smtp.NewClient(conn, s.config.Email.SMTPHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client with SSL: %w", err)
	}
	defer func() {
		if err := client.Quit(); err != nil {
			// Log client quit error but don't fail the operation
		}
	}()

	// Authenticate
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// Set sender
	if err := client.Mail(s.config.Email.SMTPFrom); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Set recipients
	for _, recipient := range to {
		if err := client.Rcpt(recipient); err != nil {
			return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
		}
	}

	// Send message
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data writer: %w", err)
	}
	defer func() {
		if err := writer.Close(); err != nil {
			// Log writer close error but don't fail the operation
			fmt.Printf("Failed to close SMTP writer: %v\n", err)
		}
	}()

	if _, err := writer.Write(message); err != nil {
		return fmt.Errorf("failed to write email data: %w", err)
	}

	fmt.Printf("\n=== ✅ 邮件发送成功 (SSL) ===\n")
	fmt.Printf("收件人: %s\n", strings.Join(to, ", "))
	fmt.Printf("主题: %s\n", s.extractSubject(string(message)))
	fmt.Printf("SMTP服务器: %s\n", addr)
	fmt.Printf("========================\n\n")

	return nil
}

// buildEmailMessage builds the email message with proper headers
func (s *EmailService) buildEmailMessage(from, to, subject, body string) []byte {
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"
	headers["Content-Transfer-Encoding"] = "base64"

	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n"

	// Base64 encode the body for proper UTF-8 support
	message += s.base64Encode(body)

	return []byte(message)
}

// base64Encode encodes text to base64
func (s *EmailService) base64Encode(text string) string {
	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	data := []byte(text)
	result := ""

	for i := 0; i < len(data); i += 3 {
		chunk := make([]byte, 3)
		chunkLen := 0

		for j := 0; j < 3 && i+j < len(data); j++ {
			chunk[j] = data[i+j]
			chunkLen++
		}

		encoded := make([]byte, 4)
		encoded[0] = base64Table[chunk[0]>>2]
		encoded[1] = base64Table[((chunk[0]&0x03)<<4)|((chunk[1]&0xf0)>>4)]

		if chunkLen > 1 {
			encoded[2] = base64Table[((chunk[1]&0x0f)<<2)|((chunk[2]&0xc0)>>6)]
		} else {
			encoded[2] = '='
		}

		if chunkLen > 2 {
			encoded[3] = base64Table[chunk[2]&0x3f]
		} else {
			encoded[3] = '='
		}

		result += string(encoded)

		// Add line breaks every 76 characters as per RFC
		if len(result)%76 == 0 {
			result += "\r\n"
		}
	}

	return result
}

// extractSubject extracts subject from email message for logging
func (s *EmailService) extractSubject(message string) string {
	lines := strings.Split(message, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Subject: ") {
			return strings.TrimPrefix(line, "Subject: ")
		}
	}
	return "Unknown Subject"
}

// ValidateEmail performs basic email format validation
func (s *EmailService) ValidateEmail(email string) bool {
	// Basic validation
	if email == "" {
		return false
	}

	// Check for @ symbol
	if !strings.Contains(email, "@") {
		return false
	}

	// Split by @
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	// Check local part (before @)
	local := parts[0]
	if local == "" || len(local) > 64 {
		return false
	}

	// Check domain part (after @)
	domain := parts[1]
	if domain == "" || len(domain) > 255 {
		return false
	}

	// Check for dot in domain
	if !strings.Contains(domain, ".") {
		return false
	}

	return true
}

// GetEmailDomain extracts domain from email address
func (s *EmailService) GetEmailDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}
