package pgopher

import "time"

type Config struct {
	ListenAddress    string                     `yaml:"listenAddress"`
	ProfilingTargets map[string]ProfilingTarget `yaml:"profilingTargets"`
	Sink             string                     `yaml:"sink"`
}

type ProfilingTarget struct {
	URL      string        `yaml:"url"`
	Duration time.Duration `yaml:"duration"`
	Schedule string        `yaml:"schedule"`
}

func DefaultConfig() Config {
	return Config{
		ListenAddress:    ":8000",
		ProfilingTargets: make(map[string]ProfilingTarget),
		Sink:             "file",
	}
}
