package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/inasknh/simple-poke-app/internal/api"
	cache2 "github.com/inasknh/simple-poke-app/internal/cache"
	"github.com/inasknh/simple-poke-app/internal/config"
	db2 "github.com/inasknh/simple-poke-app/internal/db"
	handler2 "github.com/inasknh/simple-poke-app/internal/handler"
	repository2 "github.com/inasknh/simple-poke-app/internal/repository"
	service2 "github.com/inasknh/simple-poke-app/internal/service"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read config")
	}

	var configuration config.Configurations
	if err := viper.Unmarshal(&configuration); err != nil {
		log.Fatalf("Couldn't unmarshal configuration")
	}

	db := db2.NewMySql(configuration)
	dbRepository := repository2.NewRepository(db)
	cache := cache2.NewRedis(configuration.Cache)
	redisRepository := repository2.NewRedisRepository(cache, configuration)

	restyClient := resty.New().
		SetTimeout(5 * time.Second).
		SetRetryCount(3).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			return err != nil || response.StatusCode() >= 500 || response.StatusCode() == http.StatusTooManyRequests
		})

	client := api.NewClient(configuration.Api, restyClient)
	service := service2.NewService(dbRepository, redisRepository, client)
	handler := handler2.NewHandler(service)

	http.HandleFunc("/sync", handler.SyncData)
	http.HandleFunc("/items", handler.GetItems)

	port := configuration.App.Port
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("%s is running on port %d", configuration.App.Name, port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	<-done
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxTimeout); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("All server stopped!")
}
