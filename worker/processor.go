package worker

import (
	db "github.com/TranQuocToan1996/backendMaster/db/sqlc"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
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
	var (
		errHandle = func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().
				Err(err).
				Str("type", task.Type()).
				Bytes("payload", task.Payload()).
				Msg("process start fail")
		}

		ratioPriorityQueues = map[string]int{
			QueueCritical: 8,
			QueueDefault:  2,
		}
	)

	server := asynq.NewServer(
		redisOtps,
		asynq.Config{
			Queues:       ratioPriorityQueues,
			ErrorHandler: asynq.ErrorHandlerFunc(errHandle),
			Logger:       NewLoggerAsynq(),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		store:  store,
	}
}
