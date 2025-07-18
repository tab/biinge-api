package server

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"biinge-api/internal/app/controllers"
	"biinge-api/internal/config"
	"biinge-api/internal/config/middlewares"
	"biinge-api/internal/config/router"
)

func Test_NewServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		AppEnv:  "test",
		AppAddr: "localhost:8080",
	}

	mockTraceMiddleware := middlewares.NewMockTraceMiddleware(ctrl)
	mockLoggerMiddleware := middlewares.NewMockLoggerMiddleware(ctrl)
	mockAuthenticationMiddleware := middlewares.NewMockAuthenticationMiddleware(ctrl)
	mockHealthController := controllers.NewMockHealthController(ctrl)
	mockSessionsController := controllers.NewMockAuthenticationController(ctrl)
	mockAccountsController := controllers.NewMockAccountsController(ctrl)
	mockMoviesController := controllers.NewMockMoviesController(ctrl)
	mockPeopleController := controllers.NewMockPeopleController(ctrl)

	mockAuthenticationMiddleware.EXPECT().
		Authenticate(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(next http.Handler) http.Handler {
			return next
		})
	mockTraceMiddleware.EXPECT().
		Trace(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(next http.Handler) http.Handler {
			return next
		})
	mockLoggerMiddleware.EXPECT().
		Log(gomock.Any()).
		AnyTimes().
		DoAndReturn(func(next http.Handler) http.Handler {
			return next
		})

	appRouter := router.NewRouter(
		cfg,
		mockAuthenticationMiddleware,
		mockTraceMiddleware,
		mockLoggerMiddleware,
		mockHealthController,
		mockSessionsController,
		mockAccountsController,
		mockMoviesController,
		mockPeopleController,
	)

	srv := NewServer(cfg, appRouter)
	assert.NotNil(t, srv)

	s, ok := srv.(*server)
	assert.True(t, ok)

	assert.Equal(t, cfg.AppAddr, s.httpServer.Addr)
	assert.Equal(t, appRouter, s.httpServer.Handler)
	assert.Equal(t, 5*time.Second, s.httpServer.ReadTimeout)
	assert.Equal(t, 10*time.Second, s.httpServer.WriteTimeout)
	assert.Equal(t, 120*time.Second, s.httpServer.IdleTimeout)
}

func Test_Server_RunAndShutdown(t *testing.T) {
	cfg := &config.Config{
		AppEnv:  "test",
		AppAddr: "localhost:5000",
	}
	handler := http.NewServeMux()
	srv := NewServer(cfg, handler)

	runErrCh := make(chan error, 1)
	go func() {
		err := srv.Run()
		if err != nil && err != http.ErrServerClosed {
			runErrCh <- err
		}
		close(runErrCh)
	}()

	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	assert.NoError(t, err)

	err = <-runErrCh
	assert.NoError(t, err)
}
