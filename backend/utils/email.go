package utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

type EmailConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// 获取邮件配置
func GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		Host:     os.Getenv("EMAIL_HOST"),
		Port:     os.Getenv("EMAIL_PORT"),
		Username: os.Getenv("EMAIL_USERNAME"),
		Password: os.Getenv("EMAIL_PASSWORD"),
		From:     os.Getenv("EMAIL_FROM"),
	}
}

// 打印邮件配置（调试用）
func PrintEmailConfig() {
	config := GetEmailConfig()
	log.Printf("邮件配置信息：")
	log.Printf("Host: %s", config.Host)
	log.Printf("Port: %s", config.Port)
	log.Printf("Username: %s", config.Username)
	log.Printf("From: %s", config.From)
}

// SendEmail 发送邮件
func SendEmail(to, subject, body string) error {
	config := GetEmailConfig()

	// 构建邮件头
	headers := make(map[string]string)
	headers["From"] = config.From
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// 构建邮件内容
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// 服务器地址
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	// 如果是465端口，使用SSL连接
	if config.Port == "465" {
		return sendMailSSL(addr, config.Host, config.Username, config.Password, config.From, []string{to}, []byte(message))
	}

	// 创建认证
	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	// 发送邮件
	return smtp.SendMail(addr, auth, config.From, []string{to}, []byte(message))
}

// 通过SSL/TLS发送邮件
func sendMailSSL(addr, host, username, password, from string, to []string, msg []byte) error {
	// 创建TLS配置
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// 连接到服务器
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("连接到邮件服务器失败: %v", err)
	}
	defer conn.Close()

	// 创建SMTP客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("创建SMTP客户端失败: %v", err)
	}
	defer client.Close()

	// 验证身份
	auth := smtp.PlainAuth("", username, password, host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("身份验证失败: %v", err)
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	for _, recipient := range to {
		if err = client.Rcpt(recipient); err != nil {
			return fmt.Errorf("设置收件人 %s 失败: %v", recipient, err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取数据写入器失败: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	if err = w.Close(); err != nil {
		return fmt.Errorf("关闭写入器失败: %v", err)
	}

	return client.Quit()
}
