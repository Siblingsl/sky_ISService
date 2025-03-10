package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/mail"
	"net/smtp"
	"regexp"
	"strconv"
	"time"
)

// GenerateRandomCode 生成指定位数的随机验证码
func GenerateRandomCode(length int) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	for i := 0; i < length; i++ {
		code += strconv.Itoa(rand.Intn(10)) // 生成0-9的随机数
	}
	return code
}

// SendEmail 发送验证码的邮件
func SendEmail(email string, code string) error {
	// 设置邮箱 SMTP 配置
	smtpHost := "smtp.163.com"              // 163 邮箱 SMTP 服务器
	smtpPort := "587"                       // SMTP 端口
	senderEmail := "shilei07070707@163.com" // 发件人邮箱地址
	senderPassword := "Sl1035515807"        // 发件人邮箱授权码

	// 邮件主题
	subject := "【安全验证】您的邮箱验证码"

	// 读取 HTML 模板并替换验证码
	htmlBody := `
		<!DOCTYPE html>
		<html lang="zh-CN">
		<head>
		    <meta charset="UTF-8">
		    <meta name="viewport" content="width=device-width, initial-scale=1.0">
		    <title>邮箱验证码</title>
		    <style>
		        body { font-family: Arial, sans-serif; background-color: #f4f4f4; margin: 0; padding: 0; }
		        .container { max-width: 500px; margin: 50px auto; background: #ffffff; padding: 20px;
		            border-radius: 8px; box-shadow: 0 0 10px rgba(0, 0, 0, 0.1); text-align: center; }
		        h2 { color: #333; }
		        .code { font-size: 24px; font-weight: bold; color: #ff5722; padding: 10px; 
		            background: #f8f8f8; display: inline-block; border-radius: 5px; margin: 20px 0; }
		        p { color: #666; font-size: 14px; }
		        .footer { margin-top: 20px; font-size: 12px; color: #999; }
		    </style>
		</head>
		<body>
		    <div class="container">
		        <h2>您的邮箱验证码</h2>
		        <p>您好，您的验证码是：</p>
		        <div class="code">` + code + `</div>
		        <p>该验证码有效期为 5 分钟，请尽快使用。</p>
		        <p>如果您没有请求此验证码，请忽略此邮件。</p>
		        <div class="footer">© 2024 您的公司名称.芳科技商贸有限公司</div>
		    </div>
		</body>
		</html>`

	// 组装邮件消息
	message := "From: " + senderEmail + "\r\n" +
		"To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" + htmlBody

	// 认证信息
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)

	// 发送邮件
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, senderEmail, []string{email}, []byte(message))
	if err != nil {
		return fmt.Errorf("发送邮件失败: %v", err)
	}

	return nil
}

// IsValidEmail 校验邮箱格式
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email) // 使用标准库解析邮箱格式
	if err != nil {
		return false
	}

	// 进一步使用正则校验，防止邮箱伪造
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// GetClientIP 获取客户端的真实 IP 地址
func GetClientIP(c *gin.Context) string {
	ip := c.ClientIP() // Gin提供的内置方法，直接获取客户端 IP
	return ip
}
