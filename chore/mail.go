package main

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

func sendMail(code uint, dest string) {
	Welcomemsg := `
	Hi<br />
	<div style="margin-left:2em">Welcome aboard, your verification code is %06d</div>
	Enjoy<br />
	Photo Tutor<br />
	`

	m := gomail.NewMessage()
	m.SetHeader("From", "phototutor@126.com") //发件人
	m.SetHeader("To", dest)                   //收件人
	// m.SetAddressHeader("Cc", "test@126.com", "test")     //抄送人
	m.SetHeader("Subject", fmt.Sprintf("%06d Verification Code - Photo Tutor", code)) //邮件标题
	m.SetBody("text/html", fmt.Sprintf(Welcomemsg, code))                             //邮件内容
	// m.Attach("E:\\IMGP0814.JPG")       //邮件附件

	d := gomail.NewDialer("smtp.126.com", 465, "phototutor@126.com", "YOWXVZRRBCEJQECJ")
	//邮件发送服务器信息,使用授权码而非密码
	if err := d.DialAndSend(m); err != nil {
		fmt.Print(err.Error())
	}
}

func main() {
	sendMail(54123, "phototutor@126.com")

}
