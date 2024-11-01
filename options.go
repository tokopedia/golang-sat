package sat

import (
	"log"
	"net/http"

	"github.com/tokopedia/golang-sat/signature"
)

// Option contains field you can configure based on your SAT credentials
type Option struct {
	http             *http.Client
	logger           *log.Logger
	clientID         string
	clientSecret     string
	clientPrivateKey string
	serverPublicKey  string
	paddingType      signature.PaddingType
	isDebug          bool
	accessTokenURL   string
	satBaseURL       string
}

var defaultOption = Option{
	http:           &http.Client{},
	logger:         log.New(log.Writer(), "[sat] ", 0),
	paddingType:    signature.PaddingTypePSS,
	isDebug:        false,
	accessTokenURL: ACCESS_TOKEN_URL,
	satBaseURL:     PLAYGROUND_SAT_BASE_URL,
}

type ClientOptionFunc func(*Option)

// WithHTTPClient customizes http client
func WithHTTPClient(httpClient *http.Client) ClientOptionFunc {
	return func(o *Option) {
		o.http = httpClient
	}
}

// WithLogger override existing logger
func WithLogger(logger *log.Logger) ClientOptionFunc {
	return func(o *Option) {
		o.logger = logger
	}
}

// WithServerPublicKeyString load server public key
func WithServerPublicKeyString(serverPublicKey string) ClientOptionFunc {
	return func(o *Option) {
		o.serverPublicKey = serverPublicKey
	}
}

// WithPaddingType set specific padding type for sign & verify signature
func WithPaddingType(paddingType signature.PaddingType) ClientOptionFunc {
	return func(o *Option) {
		o.paddingType = paddingType
	}
}

// WithIsDebug is toggle debug log
func WithIsDebug(isDebug bool) ClientOptionFunc {
	return func(o *Option) {
		o.isDebug = isDebug
	}
}

// WithSatBaseURL override SAT Base URL
func WithSatBaseURL(satBaseURL string) ClientOptionFunc {
	return func(o *Option) {
		o.satBaseURL = satBaseURL
	}
}

// WithAccessTokenURL WithSatBaseURL override SAT Base URL
func WithAccessTokenURL(accessTokenURL string) ClientOptionFunc {
	return func(o *Option) {
		o.accessTokenURL = accessTokenURL
	}
}
