package resty

import (
	"context"
	"fmt"
	"net/http"

	"github.com/akemoon/crowdfunding-app-auth/cluster/user"
	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client *resty.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		client: resty.New().SetBaseURL(baseURL),
	}
}

func (c *Client) CreateUser(ctx context.Context, in user.CreateUserReq) error {
	resp, err := c.client.R().
		SetContext(ctx).
		SetBody(in).
		SetResult(&user.CreateUserResp{}).
		Post("/user")

	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	if resp.StatusCode() != http.StatusCreated {
		if resp.StatusCode() == http.StatusConflict {
			return user.ErrUsernameExists
		}
		return fmt.Errorf("%w: unexpected status code: %d", domain.ErrInternal, resp.StatusCode())
	}

	return nil
}
