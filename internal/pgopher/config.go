package pgopher

import (
	"fmt"
	"os"
	"time"

	"github.com/mycreepy/box"
	"gopkg.in/yaml.v3"
)

type Config struct {
	box.Config `yaml:",inline"`

	PprofListenAddress string            `yaml:"pprofListenAddress"`
	ProfilingTargets   []ProfilingTarget `yaml:"profilingTargets"`
	Sink               Sink              `yaml:"sink"`
}

type ProfilingTarget struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Duration time.Duration `yaml:"duration"`
	Schedule string        `yaml:"schedule"`
}

type Sink struct {
	Type                  string                `yaml:"type"`
	FileSinkOptions       FileSinkOptions       `yaml:"fileSinkOptions"`
	S3SinkOptions         S3SinkOptions         `yaml:"s3SinkOptions"`
	KubernetesSinkOptions KubernetesSinkOptions `yaml:"kubernetesSinkOptions"`
}

type FileSinkOptions struct {
	Folder string `yaml:"folder"`
}

type S3SinkOptions struct {
	Bucket string `yaml:"bucket"`
}

type KubernetesSinkOptions struct {
	APIServerURL string `yaml:"apiServerURL"`
	Namespace    string `yaml:"namespace"`
}

func DefaultConfig() Config {
	return Config{
		PprofListenAddress: "localhost:8010",
		ProfilingTargets:   make([]ProfilingTarget, 0),
		Sink: Sink{
			Type: "file",
			FileSinkOptions: FileSinkOptions{
				Folder: "profiles",
			},
			KubernetesSinkOptions: KubernetesSinkOptions{
				Namespace: "pgopher",
			},
		},
	}
}

func LoadConfig(path string) (Config, error) {
	cfg := DefaultConfig()

	file, err := os.Open(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to open config file: %w", err)
	}

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	err = decoder.Decode(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}
