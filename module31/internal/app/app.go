package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"os"
	"os/signal"
	"skillbox/internal/handler"
	"skillbox/internal/repository"
	"skillbox/internal/usecase"
	"syscall"
	"time"
)

func Run(port string) error {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	repositoryA, err := repository.NewMongodb()

	if err != nil {
		return err
	}
	useCase := usecase.NewUsecase(repositoryA)
	router := chi.NewRouter()
	handler.Build(router, useCase)

	go func() {
		log.Fatal(http.ListenAndServe(port, router))
	}()
	<-done
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer func() {
		repository.DisconnectDB(repositoryA)
		fmt.Println("App close")
		cancel()
	}()

	err = http.ListenAndServe(port, router)
	if err != nil {
		return err
	}
	return nil
}
