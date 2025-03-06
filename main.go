package main

import (
	"context"
	"fmt"
	"log"
	"my-project/infrastructure/cache"
	tulushost "my-project/infrastructure/clients/tulustech"
	"my-project/infrastructure/configuration"
	"my-project/infrastructure/logger"
	"my-project/infrastructure/persistence"
	"my-project/infrastructure/pubsub"
	"my-project/infrastructure/servicebus"
	youtube_client "my-project/infrastructure/youtubeclient"
	httpHandler "my-project/interfaces/http"
	"my-project/usecase"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	httpServer *http.Server
)

func recoverPanic() {
	if err := recover(); err != nil {
		fmt.Printf("RECOVERED: %v\n", err)
	}
}

func main() {
	//InitiateGoroutine()
	defer recoverPanic()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	g, ctx := errgroup.WithContext(ctx)

	// configuration.LoadConfig()

	app := configuration.C.App

	youtubeClient, err := youtube_client.NewYoutubeClient()
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while instantiate YoutubeClient")
		panic(err)
	}
	testYoutubeClient := youtube_client.NewTestYoutubeClient(youtubeClient)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while instantiate YoutubeClient")
		panic(err)
	}

	mysqlDb, psqlDb, err := InitiateDatabase()
	if err != nil {
		fmt.Println(err)
	}

	mongoDb, err := persistence.NewMongoDb(configuration.C.Database.Mongo.Host, configuration.C.Database.Mongo.Port, configuration.C.Database.Mongo.User, configuration.C.Database.Mongo.Password, configuration.C.Database.Mongo.Name)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while instantiate MongoDB")
		panic(err)
	}
	err = mongoDb.Ping(ctx, nil)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while Ping MongoDB")
		panic(err)
	}
	fmt.Println("MongoDB connected")

	logger.GetLogger().WithField("MySQLDb", mysqlDb.Ping()).WithField("PSQLDb", psqlDb.Ping()).Info("Database connected.")

	pubSubClient, err := pubsub.NewPubSub(ctx, configuration.C.Pubsub.ProjectID)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while instantiate PubSub")
		// panic(err)
	}

	azServiceBusClient, err := servicebus.NewServiceBus(ctx, configuration.C.ServiceBus.Namespace)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while instantiate ServiceBus")
		panic(err)
	}
	redisClient, _ := cache.NewCache(ctx, fmt.Sprintf("%s:%s", configuration.C.RedisClient.Host, configuration.C.RedisClient.Port), configuration.C.RedisClient.Username, configuration.C.RedisClient.Password)

	testRepository := persistence.NewTestRepository(mongoDb)
	project, err := testRepository.Test(ctx)
	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while fetching data")
	}
	fmt.Printf("Project %v\n", project)
	testCache := cache.NewTestCache(redisClient)

	logger.GetLogger().Info("Redis client initialized successfully.")

	tulusTechHost := tulushost.NewTulusHost(configuration.C.TulusTech.Host)

	testPubSub := pubsub.NewTestPubSub(pubSubClient)
	testServiceBus := servicebus.NewTestServiceBus(azServiceBusClient)

	userRepository := persistence.NewUserRepository(psqlDb)
	userUsecase := usecase.NewUserUsecase(userRepository)
	testUsecase := usecase.NewTestUsecase(tulusTechHost, testPubSub, testServiceBus, testCache, testYoutubeClient)
	testRes := testUsecase.Test(ctx)
	fmt.Println("Test response", testRes)

	userHandler := httpHandler.NewUserHandler(userUsecase)
	testHandler := httpHandler.NewTestHandler(testUsecase)

	router := InitiateRouter(userHandler, testHandler, userRepository)

	if err != nil {
		logger.GetLogger().WithField("error", err).Error("Error while StartSubscription")
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()

	port := app.Port
	logger.GetLogger().WithField("port", port).Info("Starting application")
	g.Go(func() error {
		httpServer := &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      router,
			ReadTimeout:  0,
			WriteTimeout: 0,
		}
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return err
		}
		logger.GetLogger().WithField("port", port).Error("Application start")
		return nil
	})

	select {
	case <-interrupt:
		fmt.Println("Exit")
		os.Exit(1)
	case <-ctx.Done():
		break
	}

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if httpServer != nil {
		_ = httpServer.Shutdown(shutdownCtx)
	}

	err = g.Wait()
	if err != nil {
		log.Printf("server returning an error %v", err)
		os.Exit(2)
	}
}
