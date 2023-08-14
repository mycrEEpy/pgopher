package pgopher

import "time"

type Config struct {
	ListenAddress    string            `yaml:"listenAddress"`
	ProfilingTargets []ProfilingTarget `yaml:"profilingTargets"`
	Sink             Sink              `yaml:"sink"`
}

type ProfilingTarget struct {
	Name     string        `yaml:"name"`
	URL      string        `yaml:"url"`
	Duration time.Duration `yaml:"duration"`
	Schedule string        `yaml:"schedule"`
}

type Sink struct {
	Type            string          `yaml:"type"`
	FileSinkOptions FileSinkOptions `yaml:"fileSinkOptions"`
}

type FileSinkOptions struct {
	Folder string `yaml:"folder"`
}

func DefaultConfig() Config {
	return Config{
		ListenAddress:    ":8000",
		ProfilingTargets: make([]ProfilingTarget, 0),
		Sink: Sink{
			Type: "file",
			FileSinkOptions: FileSinkOptions{
				Folder: "profiles",
			},
		},
	}
}
