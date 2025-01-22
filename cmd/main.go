package main

import (
	"context"
	"fmt"
	"mail/internal/config"
	"mail/internal/mail"
	"mail/internal/repository/redis"
)

var queue_names = []string{"otp", "email"}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg := config.NewConfigService()
	mailService := mail.NewMailService(
		cfg.SMTP.Username,
		cfg.SMTP.Password,
		cfg.SMTP.Host,
		cfg.ServerURL,
	)
	consumer := redis.NewRepositoryConsumer(
		cfg.Redis.Host,
		cfg.Redis.RedisPassword,
		cfg.Redis.Port,
		cfg.Redis.DB,
	)
	go func() {
		consumer.ConsumerMessages(ctx, queue_names, mailService.MailerFactory)
	}()
	select {
	case <-ctx.Done():
		fmt.Println("shutting down...")
	}
}
