package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/yousoon/apps/services/notification-service/internal/domain"
)

// AWSClient client pour AWS SES et SNS
type AWSClient struct {
	sesClient   *ses.Client
	snsClient   *sns.Client
	fromEmail   string
	fromName    string
	smsOriginal string
}

// NewAWSClient crée un nouveau client AWS
func NewAWSClient(ctx context.Context, region, fromEmail, fromName, smsOriginator string) (*AWSClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %w", err)
	}

	return &AWSClient{
		sesClient:   ses.NewFromConfig(cfg),
		snsClient:   sns.NewFromConfig(cfg),
		fromEmail:   fromEmail,
		fromName:    fromName,
		smsOriginal: smsOriginator,
	}, nil
}

// SendEmail envoie un email via SES
func (c *AWSClient) SendEmail(ctx context.Context, notification *domain.Notification, toEmail string) error {
	input := &ses.SendEmailInput{
		Source: aws.String(fmt.Sprintf("%s <%s>", c.fromName, c.fromEmail)),
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(notification.Title),
			},
			Body: &types.Body{
				Html: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(c.buildHTMLBody(notification)),
				},
				Text: &types.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(notification.Body),
				},
			},
		},
	}

	_, err := c.sesClient.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

// SendSMS envoie un SMS via SNS
func (c *AWSClient) SendSMS(ctx context.Context, notification *domain.Notification, phoneNumber string) error {
	input := &sns.PublishInput{
		PhoneNumber: aws.String(phoneNumber),
		Message:     aws.String(notification.Body),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
			"AWS.SNS.SMS.SenderID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(c.smsOriginal),
			},
			"AWS.SNS.SMS.SMSType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Transactional"),
			},
		},
	}

	_, err := c.snsClient.Publish(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}

func (c *AWSClient) buildHTMLBody(notification *domain.Notification) string {
	html := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f5f5f5;
        }
        .container {
            background-color: #000;
            border-radius: 12px;
            padding: 30px;
            color: #fff;
        }
        .logo {
            text-align: center;
            margin-bottom: 30px;
        }
        .logo img {
            height: 40px;
        }
        h1 {
            color: #E99B27;
            font-size: 24px;
            margin-bottom: 20px;
        }
        .content {
            font-size: 16px;
            margin-bottom: 30px;
        }
        .button {
            display: inline-block;
            background-color: #E99B27;
            color: #000;
            text-decoration: none;
            padding: 12px 30px;
            border-radius: 8px;
            font-weight: bold;
        }
        .footer {
            text-align: center;
            margin-top: 30px;
            font-size: 12px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="logo">
            <img src="https://cdn.yousoon.com/logo.png" alt="Yousoon">
        </div>
        <h1>%s</h1>
        <div class="content">
            %s
        </div>
    </div>
    <div class="footer">
        © 2024 Yousoon. Tous droits réservés.
    </div>
</body>
</html>
`, notification.Title, notification.Body)

	return html
}

// Type pour les attributs SMS SNS
type snstypes = sns.types
