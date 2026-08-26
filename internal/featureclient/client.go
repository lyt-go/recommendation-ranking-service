package featureclient

import "context"

type Backend interface {
	Load(context.Context, string) error
}

type Client struct {
	backend Backend
}

func New(backend Backend) *Client { return &Client{backend: backend} }

func (c *Client) Load(ctx context.Context, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return c.backend.Load(context.Background(), userID)
}
