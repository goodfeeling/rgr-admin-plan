package email

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
)

// EmailService 邮件服务接口
type EmailService interface {
	SendEmail(to, subject, body string) error
	SendHTMLEmail(to, subject, htmlBody string) error
}

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	IsSSL    bool
}

// SMTPEmailService SMTP邮件服务实现
type SMTPEmailService struct {
	config SMTPConfig
}

// NewSMTPEmailService 创建SMTP邮件服务
func NewSMTPEmailService() EmailService {
	config := SMTPConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
		IsSSL:    os.Getenv("EMAIL_SSL") == "true",
	}
	return &SMTPEmailService{config: config}
}

// SendEmail 发送纯文本邮件
func (s *SMTPEmailService) SendEmail(to, subject, body string) error {
	return s.sendEmail(to, subject, body, "text/plain")
}

// SendHTMLEmail 发送HTML邮件
func (s *SMTPEmailService) SendHTMLEmail(to, subject, htmlBody string) error {
	return s.sendEmail(to, subject, htmlBody, "text/html")
}

// sendEmail 发送邮件的内部实现
func (s *SMTPEmailService) sendEmail(to, subject, body, contentType string) error {
	// 构建邮件内容
	message := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: %s; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		s.config.From,
		to,
		subject,
		contentType,
		body,
	)

	// 设置认证信息
	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	// 发送邮件
	if s.config.IsSSL {
		return s.sendEmailSSL(auth, to, message)
	}

	return s.sendEmailTLS(auth, to, message)
}

// 手动SSL连接方式
func (s *SMTPEmailService) sendEmailSSL(auth smtp.Auth, to, message string) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// 手动建立TLS连接
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		ServerName: s.config.Host,
	})
	if err != nil {
		return fmt.Errorf("failed to establish TLS connection: %w", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// 认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// 发送邮件
	if err = client.Mail(s.config.From); err != nil {
		return err
	}

	if err = client.Rcpt(to); err != nil {
		return err
	}

	writer, err := client.Data()
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		writer.Close()
		return err
	}

	return writer.Close()
}

// sendEmailTLS 通过TLS发送邮件
func (s *SMTPEmailService) sendEmailTLS(auth smtp.Auth, to, message string) error {
	addr := fmt.Sprintf("%s:%s", s.config.Host, s.config.Port)

	// 连接到服务器
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Quit()

	// 启用认证
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP authentication failed: %w", err)
	}

	// 设置发件人
	if err = client.Mail(s.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// 设置收件人
	if err = client.Rcpt(to); err != nil { // 使用正确的收件人地址
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// 发送邮件内容
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to create data writer: %w", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("failed to close writer: %w", err)
	}

	return nil
}
