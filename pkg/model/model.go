package model

type EMContextKey struct {
	ContextKey string
}

type EMConfigHolder interface {
	GetEMConfig() (*EMConfig, error)
}

type URLProvider interface {
	DocWithPath(path string) string
	AuthWithPath(path string) string
	SyncWithPath(path string) string
}

type URLProviderHolder interface {
	GetURLProvider() (URLProvider, error)
}

type ConcurrencyProvider interface {
	GetConcurrency() int
}
