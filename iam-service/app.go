package main

import (
	"context"
	"errors"
	"fmt"
	"iam-service/configs"
	"iam-service/model"
	"iam-service/proto1"
	"iam-service/model/db"
	"iam-service/service"
	"iam-service/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

type app struct {
	config                    configs.Config
	grpcServer                *grpc.Server
	authServiceServer         proto1.AuthServiceServer
	authService 			  *service.AuthService
	authRepo                  model.UserRepo		// da li mi treba ovaj ili db/repo
	shutdownProcesses         []func()
	gracefulShutdownProcesses []func(wg *sync.WaitGroup)
}

func NewAppWithConfig(config configs.Config) (*app, error) {
	if config == nil {
		return nil, errors.New("config is nil")
	}
	return &app{
		config:                    config,
		shutdownProcesses:         make([]func(), 0),
		gracefulShutdownProcesses: make([]func(wg *sync.WaitGroup), 0),
	}, nil
}

func (a *app) Start() error {
	a.init()

	return a.startGrpcServer()
}

func (a *app) GracefulStop(ctx context.Context) {
	// call all shutdown processes after a timeout or graceful shutdown processes completion
	defer a.shutdown()

	// wait for all graceful shutdown processes to complete
	wg := &sync.WaitGroup{}
	wg.Add(len(a.gracefulShutdownProcesses))

	for _, gracefulShutdownProcess := range a.gracefulShutdownProcesses {
		go gracefulShutdownProcess(wg)
	}

	// notify when graceful shutdown processes are done
	gracefulShutdownDone := make(chan struct{})
	go func() {
		wg.Wait()
		gracefulShutdownDone <- struct{}{}
	}()

	// wait for graceful shutdown processes to complete or for ctx timeout
	select {
	case <-ctx.Done():
		log.Println("ctx timeout ... shutting down")
	case <-gracefulShutdownDone:
		log.Println("app gracefully stopped")
	}
}

func (a *app) init() {
	manager, err := db.NewTransactionManager(
		a.config.Neo4j().Uri(),
		a.config.Neo4j().DbName())
	if err != nil {
		log.Fatalln(err)
	}
	a.shutdownProcesses = append(a.shutdownProcesses, func() {
		log.Println("closing neo4j conn")
		manager.Stop()
	})

	a.initUserRepo(manager)

	a.initAuthService()

	a.initAuthServiceServer()
	a.initGrpcServer()
}

func (a *app) initGrpcServer() {
	
	if a.authServiceServer == nil {
		log.Fatalln("eval grpc server is nil")
	}
	s := grpc.NewServer()
	proto1.RegisterAuthServiceServer(s, a.authServiceServer)
	reflection.Register(s)
	a.grpcServer = s
}

func (a *app) initAuthServiceServer() {
	if a.authService == nil {
		log.Fatalln("eval service is nil")
	}
	server, err := server.NewAuthServiceServer(*a.authService)
	if err != nil {
		log.Fatalln(err)
	}
	a.authServiceServer = server
}

func (a *app) initAuthService() {
	if a.authRepo == nil {
		log.Fatalln("rhabac repo is nil")
	}
	authService, err := service.NewAuthService(a.authRepo)
	if err != nil {
		log.Fatalln(err)
	}
	a.authService = authService
}


func (a *app) initUserRepo(manager *db.TransactionManager) {
	a.authRepo = db.NewUserRepo(manager, db.NewSimpleCypherFactory())
}


func (a *app) startGrpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", a.config.Server().Port()))
	if err != nil {
		return err
	}
	go func() {
		log.Printf("server listening at PROMENAAAAA %v", lis.Addr())
		if err := a.grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	a.gracefulShutdownProcesses = append(a.gracefulShutdownProcesses, func(wg *sync.WaitGroup) {
		a.grpcServer.GracefulStop()
		log.Println("iam-service server gracefully stopped")
		wg.Done()
	})
	return nil
}

func (a *app) shutdown() {
	for _, shutdownProcess := range a.shutdownProcesses {
		shutdownProcess()
	}
}
