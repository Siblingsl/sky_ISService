package utils

import (
	"crypto/tls"
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/mail"
	"net/smtp"
	"os"
	"path/filepath"
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
	smtpHost := "smtp.163.com"
	smtpPort := "465" // 465 端口使用 SSL 加密
	senderEmail := "shilei07070707@163.com"
	senderPassword := "RLQaJRA5GzRPPzFR" // 163 邮箱 SMTP 授权码

	// 邮件主题
	subject := "【安全验证】您的邮箱验证码"

	// HTML 邮件内容
	htmlBody := fmt.Sprintf(`
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
		        <div class="code">%s</div>
		        <p>该验证码有效期为 1 分钟，请尽快使用。</p>
		        <p>如果您没有请求此验证码，请忽略此邮件。</p>
		        <div class="footer">© 2025 芳科技商贸有限公司</div>
		    </div>
		</body>
		</html>`, code)

	// 组装邮件消息
	message := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n\r\n%s",
		senderEmail, email, subject, htmlBody)

	// 连接 SMTP 服务器
	serverAddress := smtpHost + ":" + smtpPort
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false, // 生产环境建议验证 CA 证书
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", serverAddress, tlsConfig)
	if err != nil {
		return fmt.Errorf("SMTP 连接失败: %v", err)
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("创建 SMTP 客户端失败: %v", err)
	}
	defer client.Close()

	// 进行身份验证
	auth := smtp.PlainAuth("", senderEmail, senderPassword, smtpHost)
	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("SMTP 身份验证失败: %v", err)
	}

	// 设置发件人和收件人
	if err := client.Mail(senderEmail); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}
	if err := client.Rcpt(email); err != nil {
		return fmt.Errorf("设置收件人失败: %v", err)
	}

	// 发送邮件数据
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("获取写入器失败: %v", err)
	}

	_, err = writer.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("邮件写入失败: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return fmt.Errorf("关闭邮件写入失败: %v", err)
	}

	// 关闭 SMTP 连接
	client.Quit()
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

// ExtractConditions 从请求 URL 中提取分页和查询条件
func ExtractConditions(ctx *gin.Context) (int, int, map[string]interface{}) {
	// 获取请求中的分页参数
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))

	// 创建一个 map 来存储动态查询条件
	conditions := make(map[string]interface{})

	// 获取请求中的所有查询参数，并将分页参数从中分离出去
	for key, value := range ctx.Request.URL.Query() {
		// 分离分页参数，不作为查询条件
		if key == "page" || key == "size" {
			continue
		}
		// 如果值不为空，则将其作为查询条件
		if len(value) > 0 && value[0] != "" {
			conditions[key] = value[0]
		}
	}

	return page, size, conditions
}

// GetAbsolutePath 获取当前工作目录并返回与相对路径组合后的绝对路径
func GetAbsolutePath(relativePath string) string {
	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return "获取当前目录错误"
	}
	// 构造绝对路径
	absolutePath := filepath.Join(wd, relativePath)
	return absolutePath
}

// Contains 判断数组中是否包含某个元素
func Contains(slice []int, item int) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
