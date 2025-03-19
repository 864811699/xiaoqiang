package server

import (
	"github.com/jordan-wright/email"
	"log"
	"net/smtp"
	"testing"
)

func Test_SendEmail(t *testing.T)  {
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "smj <864811699@qq.com>"
	// 设置接收方的邮箱
	e.To = []string{"shumingjun@minghsiim.com"}
	//设置主题
	e.Subject = "这是主题"
	//设置文件发送的内容
	e.Text = []byte("www.topgoer.com是个不错的go语言中文文档")
	//设置服务器相关的配置
	err := e.Send("smtp.qq.com:463", smtp.PlainAuth("", "864811699@qq.com", "gaymswketnhobecf", "smtp.qq.com"))
	if err != nil {
		log.Fatal(err)
	}
}