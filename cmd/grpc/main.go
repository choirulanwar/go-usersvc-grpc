package main

import (
	"os"
	"time"

	"go-micro.dev/v4"
	"go-micro.dev/v4/registry/cache"
	"go-micro.dev/v4/server"
	"go-micro.dev/v4/util/log"

	"github.com/asim/go-micro/plugins/registry/consul/v4"
	"github.com/asim/go-micro/plugins/server/grpc/v4"
	"github.com/joho/godotenv"
	"go-micro.dev/v4/registry"

	pb "choirulanwar/user-svc/internal/core/domain"
	userhdlgrpc "choirulanwar/user-svc/internal/core/handler/userhdl/grpc"
	"choirulanwar/user-svc/internal/core/service/usersrv"
	"choirulanwar/user-svc/internal/repositories/mongodb"
	"choirulanwar/user-svc/pkg/uidgen"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	reg := consul.NewRegistry(func(opts *registry.Options) {
		opts.Addrs = []string{
			os.Getenv("CONSUL_HOST"),
		}
	})

	service := micro.NewService(
		micro.Server(grpc.NewServer()),
		micro.Registry(cache.New(reg)),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
		micro.Name(os.Getenv("SVC_NAME")),
		micro.Address(os.Getenv("SVC_ADDRESS")),
	)

	service.Server().Init(
		server.Wait(nil),
	)

	userRepository, err := mongodb.NewMongoRepository(os.Getenv("DB_URL"), os.Getenv("DB_NAME"), 30)

	if err != nil {
		log.Fatal(err)
	}

	userService := usersrv.NewUserService(userRepository, uidgen.New())
	userHandler := userhdlgrpc.NewGrpcHandler(userService)

	pb.RegisterUserServiceHandler(service.Server(), userHandler)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
