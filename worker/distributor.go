package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type TaskDistributor interface {
	DistributeTaskSendEmail(context.Context,
		*PayloadSendVerifyEmail, ...asynq.Option) error
}

type RedisTaskDistributor struct {
	client *asynq.Client
}

func NewRedisTaskDistributor(options asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(options)
	return &RedisTaskDistributor{client}
}
