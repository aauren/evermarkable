package model

const (
	// ******************* GENERAL CONSTANTS *******************
	// Constants related to CLI and Config
	EnvPrefix                 = "EM_"
	EMAppName                 = "evermarkable"
	EMDefaultCacheDirModeSet  = 0700
	EMDTCacheFile             = "dt.cache"
	EMDefaultCacheFileModeSet = 0600
	EMCacheVersion            = 3

	// ******************* REMARKABLE SPECIFIC CONSTANTS *******************
	// URI's for remarkable cloud endpoints which can be overridden
	DocHost  = "https://document-storage-production-dot-remarkable-production.appspot.com"
	AuthHost = "https://webapp-prod.cloud.remarkable.engineering"
	SyncHost = "https://internal.cloud.remarkable.com"

	// Paths for requesting tokens
	//nolint:gosec // the following are not hardcoded credentials
	DeviceTokenPath = "/token/json/2/device/new"
	UserDevicePath  = "/token/json/2/user/new"

	// Paths for document storage
	ListDocsPath      = "/document-storage/json/2/docs"
	UpdateStatusPath  = "/document-storage/json/2/upload/update-status"
	UploadRequestPath = "/document-storage/json/2/upload/request"
	DeleteEntryPath   = "/document-storage/json/2/delete"

	// Paths for sync endpoints
	UploadBlobPath   = "/sync/v2/signed-urls/uploads"
	DownloadBlobPath = "/sync/v2/signed-urls/downloads"
	SyncCompletePath = "/sync/v2/sync-complete"

	// User agent used when making web requests
	EMUserAgent = "evermarkable"

	// Constants related to authentication to remarkable API
	DeviceTokenSecName = "device-token"
	UserTokenSecName   = "user-token"
	DefaultDeviceDesc  = "desktop-linux"

	// Constants related to remarkable nodes and documents
	DirectoryType = "CollectionType"
	DocumentType  = "DocumentType"

	// The root path name for the remarkable cloud
	RemRootPathName       = "root"
	RemDefaultConcurrency = 99

	// ******************* EVERNOTE CONSTANTS *******************
)

var (
	ContextConfigSet = EMContextKey{ContextKey: "config"}
)
