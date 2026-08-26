package featureclient

import "context"

type Backend interface {
	Load(context.Context, string) error
}

type Client struct {
	backend Backend
}

func New(backend Backend) *Client { return &Client{backend: backend} }

// Load 把调用方的 ctx 透传给特征源，使请求取消能即时传播到底层后端，
// 不再被 context.Background() 截断。详见 refreshjob.Coordinator.Refresh。
func (c *Client) Load(ctx context.Context, userID string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return c.backend.Load(ctx, userID)
}
