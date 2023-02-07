package worker

import (
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/hibiken/asynq"
	"golang.org/x/net/context"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerify(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	store  db.Store
}

func (r *RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, r.ProcessTaskSendVerify)

	return nil
}

func NewRedisTaskProcessor(redisOtps asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(
		redisOtps,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 70, // 70%
				QueueDefault:  30,
			},
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
