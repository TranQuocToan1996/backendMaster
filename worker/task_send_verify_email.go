package worker

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type PayloadSendVerifyEmail struct {
	Username string `json:"username"`
}

func (r *RedisTaskDistributor) DistributeTaskSendEmail(ctx context.Context,
	payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {

	buf, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("fail to marshal payload when send email %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, buf, opts...)

	info, err := r.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("fail to enqueue task when send email %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Str("queue", info.Queue).Int("max_retry", info.MaxRetry).
		Msgf("enqueued task send email for %v", payload.Username)
	return nil
}

func (r *RedisTaskProcessor) ProcessTaskSendVerify(ctx context.Context, task *asynq.Task) error {
	var payload *PayloadSendVerifyEmail
	err := json.Unmarshal(task.Payload(), payload)
	if err != nil {
		log.Error().Msgf("fail to unmarshal payload - skip retry for the task -  %s", err)
		return fmt.Errorf("fail to unmarshal payload -  %w", asynq.SkipRetry)
	}

	user, err := r.store.GetUser(ctx, payload.Username)
	if err != nil {
		// if err == sql.ErrNoRows {
		// 	return fmt.Errorf("user doesnt exist: %w", asynq.SkipRetry)
		// }

		return fmt.Errorf("fail to get user: %w", err)
	}

	log.Info().Str("type", task.Type()).Bytes("payload", task.Payload()).
		Msgf("sent email for user %v with email %v", payload.Username, user.Email)

	return nil
}
