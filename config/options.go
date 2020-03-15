// Copyright 2019 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package config // import "miniflux.app/config"

import (
	"fmt"
	"strings"
)

const (
	defaultHTTPS                     = false
	defaultLogDateTime               = false
	defaultHSTS                      = true
	defaultHTTPService               = true
	defaultSchedulerService          = true
	defaultDebug                     = false
	defaultBaseURL                   = "http://localhost"
	defaultRootURL                   = "http://localhost"
	defaultBasePath                  = ""
	defaultWorkerPoolSize            = 5
	defaultPollingFrequency          = 60
	defaultBatchSize                 = 10
	defaultRunMigrations             = false
	defaultDatabaseURL               = "user=postgres password=postgres dbname=miniflux2 sslmode=disable"
	defaultDatabaseMaxConns          = 20
	defaultDatabaseMinConns          = 1
	defaultListenAddr                = "127.0.0.1:8080"
	defaultCertFile                  = ""
	defaultKeyFile                   = ""
	defaultCertDomain                = ""
	defaultCertCache                 = "/tmp/cert_cache"
	defaultCleanupFrequencyHours     = 24
	defaultCleanupArchiveReadDays    = 60
	defaultCleanupRemoveSessionsDays = 30
	defaultProxyImages               = "http-only"
	defaultCreateAdmin               = false
	defaultOAuth2UserCreation        = false
	defaultOAuth2ClientID            = ""
	defaultOAuth2ClientSecret        = ""
	defaultOAuth2RedirectURL         = ""
	defaultOAuth2Provider            = ""
	defaultPocketConsumerKey         = ""
	defaultHTTPClientTimeout         = 20
	defaultHTTPClientMaxBodySize     = 15
)

// Options contains configuration options.
type Options struct {
	HTTPS                     bool
	logDateTime               bool
	hsts                      bool
	httpService               bool
	schedulerService          bool
	debug                     bool
	baseURL                   string
	rootURL                   string
	basePath                  string
	databaseURL               string
	databaseMaxConns          int
	databaseMinConns          int
	runMigrations             bool
	listenAddr                string
	certFile                  string
	certDomain                string
	certCache                 string
	certKeyFile               string
	cleanupFrequencyHours     int
	cleanupArchiveReadDays    int
	cleanupRemoveSessionsDays int
	pollingFrequency          int
	batchSize                 int
	workerPoolSize            int
	createAdmin               bool
	proxyImages               string
	oauth2UserCreationAllowed bool
	oauth2ClientID            string
	oauth2ClientSecret        string
	oauth2RedirectURL         string
	oauth2Provider            string
	pocketConsumerKey         string
	httpClientTimeout         int
	httpClientMaxBodySize     int64
}

// NewOptions returns Options with default values.
func NewOptions() *Options {
	return &Options{
		HTTPS:                     defaultHTTPS,
		logDateTime:               defaultLogDateTime,
		hsts:                      defaultHSTS,
		httpService:               defaultHTTPService,
		schedulerService:          defaultSchedulerService,
		debug:                     defaultDebug,
		baseURL:                   defaultBaseURL,
		rootURL:                   defaultRootURL,
		basePath:                  defaultBasePath,
		databaseURL:               defaultDatabaseURL,
		databaseMaxConns:          defaultDatabaseMaxConns,
		databaseMinConns:          defaultDatabaseMinConns,
		runMigrations:             defaultRunMigrations,
		listenAddr:                defaultListenAddr,
		certFile:                  defaultCertFile,
		certDomain:                defaultCertDomain,
		certCache:                 defaultCertCache,
		certKeyFile:               defaultKeyFile,
		cleanupFrequencyHours:     defaultCleanupFrequencyHours,
		cleanupArchiveReadDays:    defaultCleanupArchiveReadDays,
		cleanupRemoveSessionsDays: defaultCleanupRemoveSessionsDays,
		pollingFrequency:          defaultPollingFrequency,
		batchSize:                 defaultBatchSize,
		workerPoolSize:            defaultWorkerPoolSize,
		createAdmin:               defaultCreateAdmin,
		proxyImages:               defaultProxyImages,
		oauth2UserCreationAllowed: defaultOAuth2UserCreation,
		oauth2ClientID:            defaultOAuth2ClientID,
		oauth2ClientSecret:        defaultOAuth2ClientSecret,
		oauth2RedirectURL:         defaultOAuth2RedirectURL,
		oauth2Provider:            defaultOAuth2Provider,
		pocketConsumerKey:         defaultPocketConsumerKey,
		httpClientTimeout:         defaultHTTPClientTimeout,
		httpClientMaxBodySize:     defaultHTTPClientMaxBodySize * 1024 * 1024,
	}
}

// LogDateTime returns true if the date/time should be displayed in log messages.
func (o *Options) LogDateTime() bool {
	return o.logDateTime
}

// HasDebugMode returns true if debug mode is enabled.
func (o *Options) HasDebugMode() bool {
	return o.debug
}

// BaseURL returns the application base URL with path.
func (o *Options) BaseURL() string {
	return o.baseURL
}

// RootURL returns the base URL without path.
func (o *Options) RootURL() string {
	return o.rootURL
}

// BasePath returns the application base path according to the base URL.
func (o *Options) BasePath() string {
	return o.basePath
}

// IsDefaultDatabaseURL returns true if the default database URL is used.
func (o *Options) IsDefaultDatabaseURL() bool {
	return o.databaseURL == defaultDatabaseURL
}

// DatabaseURL returns the database URL.
func (o *Options) DatabaseURL() string {
	return o.databaseURL
}

// DatabaseMaxConns returns the maximum number of database connections.
func (o *Options) DatabaseMaxConns() int {
	return o.databaseMaxConns
}

// DatabaseMinConns returns the minimum number of database connections.
func (o *Options) DatabaseMinConns() int {
	return o.databaseMinConns
}

// ListenAddr returns the listen address for the HTTP server.
func (o *Options) ListenAddr() string {
	return o.listenAddr
}

// CertFile returns the SSL certificate filename if any.
func (o *Options) CertFile() string {
	return o.certFile
}

// CertKeyFile returns the private key filename for custom SSL certificate.
func (o *Options) CertKeyFile() string {
	return o.certKeyFile
}

// CertDomain returns the domain to use for Let's Encrypt certificate.
func (o *Options) CertDomain() string {
	return o.certDomain
}

// CertCache returns the directory to use for Let's Encrypt session cache.
func (o *Options) CertCache() string {
	return o.certCache
}

// CleanupFrequencyHours returns the interval in hours for cleanup jobs.
func (o *Options) CleanupFrequencyHours() int {
	return o.cleanupFrequencyHours
}

// CleanupArchiveReadDays returns the number of days after which marking read items as removed.
func (o *Options) CleanupArchiveReadDays() int {
	return o.cleanupArchiveReadDays
}

// CleanupRemoveSessionsDays returns the number of days after which to remove sessions.
func (o *Options) CleanupRemoveSessionsDays() int {
	return o.cleanupRemoveSessionsDays
}

// WorkerPoolSize returns the number of background worker.
func (o *Options) WorkerPoolSize() int {
	return o.workerPoolSize
}

// PollingFrequency returns the interval to refresh feeds in the background.
func (o *Options) PollingFrequency() int {
	return o.pollingFrequency
}

// BatchSize returns the number of feeds to send for background processing.
func (o *Options) BatchSize() int {
	return o.batchSize
}

// IsOAuth2UserCreationAllowed returns true if user creation is allowed for OAuth2 users.
func (o *Options) IsOAuth2UserCreationAllowed() bool {
	return o.oauth2UserCreationAllowed
}

// OAuth2ClientID returns the OAuth2 Client ID.
func (o *Options) OAuth2ClientID() string {
	return o.oauth2ClientID
}

// OAuth2ClientSecret returns the OAuth2 client secret.
func (o *Options) OAuth2ClientSecret() string {
	return o.oauth2ClientSecret
}

// OAuth2RedirectURL returns the OAuth2 redirect URL.
func (o *Options) OAuth2RedirectURL() string {
	return o.oauth2RedirectURL
}

// OAuth2Provider returns the name of the OAuth2 provider configured.
func (o *Options) OAuth2Provider() string {
	return o.oauth2Provider
}

// HasHSTS returns true if HTTP Strict Transport Security is enabled.
func (o *Options) HasHSTS() bool {
	return o.hsts
}

// RunMigrations returns true if the environment variable RUN_MIGRATIONS is not empty.
func (o *Options) RunMigrations() bool {
	return o.runMigrations
}

// CreateAdmin returns true if the environment variable CREATE_ADMIN is not empty.
func (o *Options) CreateAdmin() bool {
	return o.createAdmin
}

// ProxyImages returns "none" to never proxy, "http-only" to proxy non-HTTPS, "all" to always proxy.
func (o *Options) ProxyImages() string {
	return o.proxyImages
}

// HasHTTPService returns true if the HTTP service is enabled.
func (o *Options) HasHTTPService() bool {
	return o.httpService
}

// HasSchedulerService returns true if the scheduler service is enabled.
func (o *Options) HasSchedulerService() bool {
	return o.schedulerService
}

// PocketConsumerKey returns the Pocket Consumer Key if configured.
func (o *Options) PocketConsumerKey(defaultValue string) string {
	if o.pocketConsumerKey != "" {
		return o.pocketConsumerKey
	}
	return defaultValue
}

// HTTPClientTimeout returns the time limit in seconds before the HTTP client cancel the request.
func (o *Options) HTTPClientTimeout() int {
	return o.httpClientTimeout
}

// HTTPClientMaxBodySize returns the number of bytes allowed for the HTTP client to transfer.
func (o *Options) HTTPClientMaxBodySize() int64 {
	return o.httpClientMaxBodySize
}

func (o *Options) String() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("LOG_DATE_TIME: %v\n", o.logDateTime))
	builder.WriteString(fmt.Sprintf("DEBUG: %v\n", o.debug))
	builder.WriteString(fmt.Sprintf("HTTP_SERVICE: %v\n", o.httpService))
	builder.WriteString(fmt.Sprintf("SCHEDULER_SERVICE: %v\n", o.schedulerService))
	builder.WriteString(fmt.Sprintf("HTTPS: %v\n", o.HTTPS))
	builder.WriteString(fmt.Sprintf("HSTS: %v\n", o.hsts))
	builder.WriteString(fmt.Sprintf("BASE_URL: %v\n", o.baseURL))
	builder.WriteString(fmt.Sprintf("ROOT_URL: %v\n", o.rootURL))
	builder.WriteString(fmt.Sprintf("BASE_PATH: %v\n", o.basePath))
	builder.WriteString(fmt.Sprintf("LISTEN_ADDR: %v\n", o.listenAddr))
	builder.WriteString(fmt.Sprintf("DATABASE_URL: %v\n", o.databaseURL))
	builder.WriteString(fmt.Sprintf("DATABASE_MAX_CONNS: %v\n", o.databaseMaxConns))
	builder.WriteString(fmt.Sprintf("DATABASE_MIN_CONNS: %v\n", o.databaseMinConns))
	builder.WriteString(fmt.Sprintf("RUN_MIGRATIONS: %v\n", o.runMigrations))
	builder.WriteString(fmt.Sprintf("CERT_FILE: %v\n", o.certFile))
	builder.WriteString(fmt.Sprintf("KEY_FILE: %v\n", o.certKeyFile))
	builder.WriteString(fmt.Sprintf("CERT_DOMAIN: %v\n", o.certDomain))
	builder.WriteString(fmt.Sprintf("CERT_CACHE: %v\n", o.certCache))
	builder.WriteString(fmt.Sprintf("CLEANUP_FREQUENCY_HOURS: %v\n", o.cleanupFrequencyHours))
	builder.WriteString(fmt.Sprintf("CLEANUP_ARCHIVE_READ_DAYS: %v\n", o.cleanupArchiveReadDays))
	builder.WriteString(fmt.Sprintf("CLEANUP_REMOVE_SESSIONS_DAYS: %v\n", o.cleanupRemoveSessionsDays))
	builder.WriteString(fmt.Sprintf("WORKER_POOL_SIZE: %v\n", o.workerPoolSize))
	builder.WriteString(fmt.Sprintf("POLLING_FREQUENCY: %v\n", o.pollingFrequency))
	builder.WriteString(fmt.Sprintf("BATCH_SIZE: %v\n", o.batchSize))
	builder.WriteString(fmt.Sprintf("PROXY_IMAGES: %v\n", o.proxyImages))
	builder.WriteString(fmt.Sprintf("CREATE_ADMIN: %v\n", o.createAdmin))
	builder.WriteString(fmt.Sprintf("POCKET_CONSUMER_KEY: %v\n", o.pocketConsumerKey))
	builder.WriteString(fmt.Sprintf("OAUTH2_USER_CREATION: %v\n", o.oauth2UserCreationAllowed))
	builder.WriteString(fmt.Sprintf("OAUTH2_CLIENT_ID: %v\n", o.oauth2ClientID))
	builder.WriteString(fmt.Sprintf("OAUTH2_CLIENT_SECRET: %v\n", o.oauth2ClientSecret))
	builder.WriteString(fmt.Sprintf("OAUTH2_REDIRECT_URL: %v\n", o.oauth2RedirectURL))
	builder.WriteString(fmt.Sprintf("OAUTH2_PROVIDER: %v\n", o.oauth2Provider))
	builder.WriteString(fmt.Sprintf("HTTP_CLIENT_TIMEOUT: %v\n", o.httpClientTimeout))
	builder.WriteString(fmt.Sprintf("HTTP_CLIENT_MAX_BODY_SIZE: %v\n", o.httpClientMaxBodySize))
	return builder.String()
}
