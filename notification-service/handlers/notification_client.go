package handlers

import (
	"log"

	"github.com/go-resty/resty/v2"
)

type NotifyPayload struct {
	AuthorID      uint   `json:"author_id"`
	AuthorEmail   string `json:"author_email"`
	AuthorName    string `json:"author_name"`
	PostTitle     string `json:"post_title"`
	CommenterName string `json:"commenter_name"`
}

var notifyClient = resty.New().SetBaseURL("http://localhost:8081")

func init() {
	notifyClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		log.Printf("[Resty] → %s %s", req.Method, req.URL)
		return nil
	})

	notifyClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		log.Printf("[Resty] ← %d %s", resp.StatusCode(), resp.String())
		return nil
	})
}

func SendNotification(payload NotifyPayload) {
	resp, err := notifyClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(payload).
		Post("/notify")

	if err != nil {
		log.Printf("[Resty] Ошибка отправки уведомления: %v", err)
		return
	}

	if resp.IsError() {
		log.Printf("[Resty] NotificationService вернул ошибку: %s", resp.String())
	}
}
