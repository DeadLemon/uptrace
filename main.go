package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/otel"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	defer cancel()

	uptrace.ConfigureOpentelemetry(
		uptrace.WithServiceName("sandbox"),
		uptrace.WithServiceVersion("v0.0.1-dev"),
		uptrace.WithDeploymentEnvironment("local"),
		uptrace.WithDSN("http://XaweuoAphEdCeihgPXa9DpTUnzUuaNtUWgfyNJ95qxwW@localhost:14318?grpc=14317"),
	)
	defer uptrace.Shutdown(ctx)

	tracer := otel.Tracer("sandbox")

	ctx, main := tracer.Start(ctx, "main")
	defer main.End()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	if err := redisotel.InstrumentTracing(rdb); err != nil {
		log.Fatal(err)
	}

	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			_ctx, span := tracer.Start(ctx, "ticker")
			if _, err := rdb.Ping(_ctx).Result(); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Millisecond * 100)
			span.End()
		}
	}

}
