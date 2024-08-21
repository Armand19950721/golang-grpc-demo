package EmailUtils

import (
	// "service/protos/ArContent"
	"service/utils"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	// SES Init
	sess, _ = session.NewSession(&aws.Config{
		Region:      aws.String("ap-northeast-1"),
		Credentials: credentials.NewStaticCredentials(utils.GetEnv("AWS_SES_IP"), utils.GetEnv("AWS_SES_SECRET"), ""),
	})

	svc    = ses.New(sess)
	sender = utils.GetEnv("AWS_SES_SENDER")

	// Mail Template
	head = ` 
	<div style="width: 100%;font-family:Microsoft JhengHei;display: block;">
		<div style="width: 800px;margin: auto;">
			<div style="text-align: left;">
				<img src="https://dev-api-v2.metarcommerce.com/api/static/static/MetARlogo.png" width="300"
					style="left: 20;padding-bottom: 40px;padding-top: 70px;">
				<div style="font-size: 16px;font-weight: bold;line-height: 40px;">
					嗨 您好，
				</div>
	`

	footer = `
				<div style="font-weight: 1000;padding: 30px 0 20px 0;color: #d55d28;">
				此為系統自動通知信，請勿直接回信！
				</div>
				<div style="border-bottom: 1.6px solid #cacccc;">

				</div>
				<div style="font-size: 16px;padding: 30px 0 0px 0;font-weight: bold;">
					METAR commerce
				</div>
				<a style="padding: 3px 0 20px 0;color: #83858e;font-size: 11px;">
               		www.metarcommerce.com
            	</a>
		</div>
	</div>
	`
	textBody = "METAR commerce"
	charSet  = "UTF-8"

	// backup
	// 				<img src="` + utils.GetImagePath("MetARlogo.png", ArContent.ArContentImageType_STATIC_IMAGE) + `" width="300"
)

type EmailTemplate struct {
	Subject  string
	HtmlBody string
	TextBody string
	CharSet  string
}

func ValidEmail(recipient string) bool {
	utils.PrintObj(recipient, "ValidEmail")

	// Valid the email.
	_, err := svc.VerifyEmailAddress(&ses.VerifyEmailAddressInput{EmailAddress: aws.String(recipient)})

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			utils.PrintObj(aerr.Code()+","+aerr.Error(), "")
		} else {
			utils.PrintObj(err.Error(), "")
		}

		return false
	}
	utils.PrintObj("success")

	return true
}

func GetTemplateForgotPasswod(newPassword string) EmailTemplate {
	utils.PrintObj(newPassword, "GetTemplateForgotPasswod")
	return EmailTemplate{
		Subject: "Meta Commerce - 重設密碼",
		HtmlBody: head + `
		<div style="font-size: 16px;font-weight: bold;line-height: 40px;">
			我們己收到您申請重新設定密碼的請求！
		</div>
		<div style="font-size: 16px;font-weight: bold;line-height: 40px;">
			新密碼：` + newPassword + `
		</div>
		<div style="font-size: 16px;font-weight: bold;line-height: 40px;">
			請您使用這組臨時密碼登入忝統並重新設定您的密碼。
		</div>
		` + footer,
		TextBody: textBody,
		CharSet:  charSet,
	}
}

func GetTemplateRegister(url string) EmailTemplate {
	utils.PrintObj(url, "GetTemplateAddUserChild")
	return EmailTemplate{
		Subject: "驗證 METAR commerce 使用的電子郵件地址",
		HtmlBody: head + `
		<div style="font-size: 16px;font-weight: bold;line-height: 40px;">
			感謝您在
			<span style="font-size: 16px;font-weight:unset;"> METAR commerce </span>
			註冊了一個新帳號，成為元宇宙商城建築師的一員！
		</div>
		<div style="font-size: 16px;font-weight: bold;line-height: 40px;">
			在開始前，請點選底下的連結以驗證您的電子郵件地址：
		</div>
		<div style="line-height: 40px;">
		` + url + `
		</div>
		` + footer,
		TextBody: textBody,
		CharSet:  charSet,
	}
}

func GetTemplateAddUserChild(url string) EmailTemplate {
	utils.PrintObj(url, "GetTemplateAddUserChild")
	return EmailTemplate{
		Subject: "Meta Commerce - 會員邀請信",
		HtmlBody: `<h1>這是您的邀請網址</h1>
				  <p>` + url + `</p>`,
		TextBody: textBody,
		CharSet:  charSet,
	}
}

func SendEmail(recipient string, emailTemplate EmailTemplate) error {
	utils.PrintObj(recipient, "SendEmail")
	// utils.PrintObj(emailTemplate.HtmlBody, "emailTemplate")
	utils.SaveLog("SendEmail", recipient)

	// for test use
	if strings.Contains(recipient, "@metacommerce.test.com") {
		recipient = "armand@spe3d.co"
	}

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(emailTemplate.CharSet),
					Data:    aws.String(emailTemplate.HtmlBody),
				},
				Text: &ses.Content{
					Charset: aws.String(emailTemplate.CharSet),
					Data:    aws.String(emailTemplate.TextBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(emailTemplate.CharSet),
				Data:    aws.String(emailTemplate.Subject),
			},
		},
		Source: aws.String(sender),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			utils.PrintObj(aerr.Code()+","+aerr.Error(), "")
		} else {
			utils.PrintObj(err.Error(), "")
		}

		return err
	}

	utils.PrintObj(result, "")

	return nil
}
