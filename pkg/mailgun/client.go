package mailgun

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	cfg *Config
	mg  *mailgun.MailgunImpl
}

func New(cfg *Config) *Client {
	mg := mailgun.NewMailgun(cfg.
		Domain, cfg.APIKey)
	mg.SetAPIBase(mailgun.APIBaseEU)

	return &Client{
		mg:  mg,
		cfg: cfg,
	}
}

func (c *Client) Send(subject, to, from, text string, template *string, attachments []interface{}) {
	message := c.mg.NewMessage(from, subject, text, to)

	for _, attachment := range attachments {
		if _, ok := attachment.(string); ok {
			message.AddAttachment(attachment.(string))
		}

		if _, ok := attachment.([]byte); ok {
			message.AddBufferAttachment(subject, attachment.([]byte))
		}
	}

	if template != nil {
		message.SetHtml(*template)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := c.mg.Send(ctx, message)

	if err != nil {
		zap.S().Error(resp, id, err)

		return
	}
}
