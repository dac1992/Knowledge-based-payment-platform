// Package queue 封装基于 Redis 的异步任务队列（asynq）。
//
// Phase 1 只提供「生产者」侧的最小能力：未配置 REDIS_ADDR 时不阻断启动，
// 队列降级为禁用状态，调用 Enqueue 会返回明确的业务错误。
// 消费者（worker）在引入具体异步任务后（发券 / 对账 / 推送）再接入。
package queue

import (
	"github.com/hibiken/asynq"
	"go.uber.org/zap"

	"github.com/yanmai/api/pkg/errs"
)

// Options 队列连接参数。
type Options struct {
	Addr     string
	Password string
	DB       int
}

// Client 队列生产者。
type Client struct {
	client *asynq.Client
	log    *zap.Logger
}

// New 创建队列生产者；addr 为空时返回禁用态的 Client，不报错。
func New(opts Options, log *zap.Logger) *Client {
	c := &Client{log: log}
	if opts.Addr == "" {
		log.Warn("未配置 REDIS_ADDR，异步任务队列处于禁用状态（Phase 1 允许）")
		return c
	}
	c.client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})
	log.Info("异步任务队列已就绪", zap.String("redis", opts.Addr), zap.Int("db", opts.DB))
	return c
}

// Enabled 队列是否可用。
func (c *Client) Enabled() bool { return c != nil && c.client != nil }

// Enqueue 投递任务。
func (c *Client) Enqueue(task *asynq.Task, opts ...asynq.Option) error {
	if !c.Enabled() {
		return errs.New(errs.CodeNotImplemented).WithMessage("异步任务队列未启用：缺少 REDIS_ADDR 配置")
	}
	if _, err := c.client.Enqueue(task, opts...); err != nil {
		return errs.Wrap(errs.CodeInternal, err).WithMessage("投递异步任务失败")
	}
	return nil
}

// Close 释放连接。
func (c *Client) Close() error {
	if !c.Enabled() {
		return nil
	}
	return c.client.Close()
}
