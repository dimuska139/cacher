package main

import (
	"fmt"
	grpc2 "github.com/dimuska139/cacher/internal/api/grpc"
	v1 "github.com/dimuska139/cacher/internal/api/grpc/gen/cacher/cache/v1"
	"github.com/dimuska139/cacher/internal/cache/embedded"
	memcache2 "github.com/dimuska139/cacher/internal/cache/memcache"
	"github.com/dimuska139/cacher/pkg/config"
	"github.com/dimuska139/cacher/pkg/logging"
	"github.com/dimuska139/cacher/pkg/memcache"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const applicationName = "Cacher"

func main() {
	app := &cli.App{
		Name: applicationName,
		Authors: []*cli.Author{
			{
				Name:  "Sviridov Dmitriy",
				Email: "dimuska139@gmail.com",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config",
				Value: "./config.yml",
				Usage: "path to the config file",
			},
		},
		Action: func(c *cli.Context) error {
			cfg, err := config.NewConfig(c.String("config"))
			logger := logging.NewLogger(cfg)
			logger.Info(applicationName + " starting...")
			if err != nil {
				logger.Error("can't initialize config", err)
				return err
			}

			grpcServer := grpc.NewServer()

			if cfg.Storage == "memcache" {
				memcacheClient, err := memcache.NewClient(cfg)
				if err != nil {
					return fmt.Errorf("can't initialize memcache client: %w", err)
				}
				v1.RegisterCacheAPIServer(grpcServer,
					grpc2.NewCacheServer(logger, memcache2.NewMemcacheStorage(memcacheClient)))
			} else {
				v1.RegisterCacheAPIServer(grpcServer,
					grpc2.NewCacheServer(logger, embedded.NewEmbeddedStorage(time.Millisecond*50)))
			}
			reflection.Register(grpcServer)

			go func(conf *config.Config) {
				lis, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.GrpcPort))
				if err != nil {
					logger.Fatalf("failed to listen: %v", err)
				}

				if err := grpcServer.Serve(lis); err != nil {
					logger.Fatalf("failed to serve: %v", err)
				}
			}(cfg)

			stopSignal := make(chan os.Signal)
			signal.Notify(stopSignal, syscall.SIGTERM)
			signal.Notify(stopSignal, syscall.SIGINT)
			signal.Notify(stopSignal, syscall.SIGKILL)

			reloadSignal := make(chan os.Signal)
			signal.Notify(reloadSignal, syscall.SIGUSR1)
			logger.Info(fmt.Sprintf("%s started at 127.0.0.1:%d (grpc)", applicationName, cfg.GrpcPort))
			for {
				select {
				case <-stopSignal:
					logger.Info(fmt.Sprintf("%s shutdown started...", applicationName))
					grpcServer.GracefulStop()
					logger.Info(fmt.Sprintf("%s shutdown finished", applicationName))
					os.Exit(0)

				case <-reloadSignal:
					break
				}
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger := logging.NewLogger(nil)
		logger.Fatal(fmt.Errorf("can't run application: %w", err))
	}
}
