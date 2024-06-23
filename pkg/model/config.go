package model

type EMRootConfig struct {
	ConfigPath string
	Config     EMConfig
}

type EMConfig struct {
	Remarkable EMRemarkableConfig
}

type EMRemarkableConfig struct {
	URLs        EMURLConfig `yaml:"urls"`
	Concurrency int
}

func (e EMRemarkableConfig) GetConcurrency() int {
	if e.Concurrency == 0 {
		return RemDefaultConcurrency
	}
	return e.Concurrency
}

func (e EMRemarkableConfig) GetURLProvider() URLProvider {
	return e.URLs
}

type EMURLConfig struct {
	DocumentHost string `yaml:"documentHost"`
	AuthHost     string `yaml:"authHost"`
	SyncHost     string `yaml:"syncHost"`
}

type EMURLProvider interface {
	DocWithPath(path string) string
	AuthWithPath(path string) string
	SyncWithPath(path string) string
}

func (e EMURLConfig) DocWithPath(path string) string {
	if e.DocumentHost == "" {
		return DocHost + path
	}
	return e.DocumentHost + path
}

func (e EMURLConfig) AuthWithPath(path string) string {
	if e.AuthHost == "" {
		return AuthHost + path
	}
	return e.AuthHost + path
}

func (e EMURLConfig) SyncWithPath(path string) string {
	if e.SyncHost == "" {
		return SyncHost + path
	}
	return e.SyncHost + path
}
