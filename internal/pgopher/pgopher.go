package pgopher

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Server struct {
	cfg        Config
	mux        *echo.Echo
	s3Client   *s3.Client
	kubeClient *kubernetes.Clientset
}

func NewServer(cfg Config) (*Server, error) {
	s := &Server{
		cfg: cfg,
		mux: echo.New(),
	}

	s.mux.HideBanner = true
	s.mux.HidePort = true

	s.mux.GET("/_ready", readinessProbe)
	s.mux.GET("/_live", livenessProbe)
	s.mux.GET("/api/v1/profile/:profile", s.handleProfile)

	if cfg.Sink.Type == "s3" {
		sdkConfig, err := awsConfig.LoadDefaultConfig(context.Background())
		if err != nil {
			return nil, fmt.Errorf("failed to load aws config: %w", err)
		}

		s.s3Client = s3.NewFromConfig(sdkConfig)
	}

	if cfg.Sink.Type == "kubernetes" {
		kubeCfg, err := clientcmd.BuildConfigFromFlags(cfg.Sink.KubernetesSinkOptions.APIServerURL, os.Getenv("KUBECONFIG"))
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes config: %w", err)
		}

		s.kubeClient, err = kubernetes.NewForConfig(kubeCfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
		}
	}

	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go s.startScheduler(ctx, wg)

	go func() {
		<-ctx.Done()

		err := s.mux.Shutdown(context.Background())
		if err != nil {
			slog.Error("failed to shutdown http server", slog.String("err", err.Error()))
		}
	}()

	slog.Info("starting http server", slog.String("listenAddr", s.cfg.ListenAddress))

	err := s.mux.Start(s.cfg.ListenAddress)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server failed: %w", err)
	}

	wg.Wait()

	slog.Info("graceful shutdown completed, see you next time!")

	return nil
}
