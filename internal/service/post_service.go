package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/entity"
	"github.com/ariefzainuri96/go-logstream/cmd/api/dto/request"
	"github.com/ariefzainuri96/go-logstream/cmd/api/middleware"
	"github.com/ariefzainuri96/go-logstream/internal/store"
	"github.com/ariefzainuri96/go-logstream/internal/utils"
	"go.uber.org/zap"
)

type WebhookPayload map[string]interface{}

type PostService struct {
	logger *zap.Logger
	store  store.Storage
}

func NewPostService(store store.Storage, logger *zap.Logger) *PostService {
	return &PostService{
		logger: logger,
		store:  store,
	}
}

func (s *PostService) CreatePost(ctx context.Context, req request.AddPostRequest) (entity.Post, error) {
	reqID, ok := ctx.Value(middleware.CtxRequestID).(string)

	// 2. Safety check: Context values are optional!
	if !ok {
		reqID = "unknown-request" // Fallback if missing
	}

	post, err := s.store.IPost.CreatePost(ctx, req)

	if err != nil {
		return entity.Post{}, err
	}

	if post.Project.WebhookUrl != "" {
		go func(p entity.Post) {
			err := callWebhook(p)

			if err != nil {
				s.logger.Error("⚠️ Failed to trigger webhook", zap.String("RequestId", reqID), zap.Uint("PostId", p.ID), zap.Error(err))
			} else {
				s.logger.Info("✅ Webhook triggered successfully", zap.String("RequestId", reqID), zap.Uint("PostId", p.ID))
			}
		}(post)
	}

	return post, nil
}

func callWebhook(post entity.Post) error {
	// 1. Create a client with a timeout (CRITICAL for stability)
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 2. Prepare payload (if doing POST)
	jsonPayload, _ := json.Marshal(webhookPayload(post))

	// 3. Create the Custom Request (POST example)
	req, err := http.NewRequest("POST", post.Project.WebhookUrl, bytes.NewBuffer(jsonPayload))

	if err != nil {
		return err
	}

	// 4. Set Headers (like Auth or Content-Type)
	req.Header.Set("Content-Type", "application/json")

	// 5. Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// 2. IMPORTANT: Close the body when done to prevent memory leaks
	defer resp.Body.Close()

	// 3. Read the body
	_, err = io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("Return code not 200")
	}

	return nil
}

func webhookPayload(post entity.Post) map[string]interface{} {
	switch post.Project.WebhookProvider {
	case "discord":
		// Discord expects "content" or "embeds"
		return WebhookPayload{
			"username": "LogStream", // Custom bot name
			"embeds": []map[string]interface{}{
				{
					"title":       post.Title,
					"description": post.Content, // Or a snippet if it's too long
					"color":       5814783,      // A nice blue color (Decimal value)
					"fields": []map[string]interface{}{
						{
							"name":   "Category",
							"value":  post.Category,
							"inline": true,
						},
						{
							"name":   "Status",
							"value":  post.Status,
							"inline": true,
						},
						{
							"name":   "Post ID",
							"value":  fmt.Sprintf("%d", post.ID),
							"inline": true,
						},
					},
					"footer": map[string]string{
						"text": "Sent via LogStream",
					},
					"timestamp": time.Now().Format(time.RFC3339),
				},
			},
		}

	case "slack":
		// Slack expects "text"
		return WebhookPayload{
			"text": fmt.Sprintf("*New Update: %s*\nCategory: %s\n<%s>", post.Title, post.Category, post.Content),
		}

	default: // "generic"
		// Send the raw data for custom integrations (Zapier, n8n, custom backends)
		return WebhookPayload{
			"event": "post_created",
			"data": map[string]interface{}{
				"id":         post.ID,
				"title":      post.Title,
				"content":    post.Content,
				"category":   post.Category,
				"created_at": post.CreatedAt,
			},
		}
	}
}

func (s *PostService) GetPost(ctx context.Context, req request.GetPostRequest) (utils.PaginateResult[entity.Post], error) {
	post, err := s.store.IPost.GetPost(ctx, req)

	if err != nil {
		return utils.PaginateResult[entity.Post]{}, err
	}

	return post, nil
}
