// Package sat is a Go API client for both the SAT API v2 REST.
// Most methods should be implemented, and it's recommended to use
// To debug responses from the API, you can toggle IsDebug = true
package sat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/google/jsonapi"
	"github.com/tokopedia/digital-b2b-client-library/golang-sat/logger"
	"github.com/tokopedia/digital-b2b-client-library/golang-sat/signature"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

// Client contains dependencies need by the SDK
type Client struct {
	http           *http.Client
	logger         *log.Logger
	satBaseURL     string
	accessTokenURL string
	signature      *signature.Signature
	isDebug        bool
}

// Callback contains interface Handler the callback from the SAT
type Callback interface {
	Do(ctx context.Context, request *OrderDetail) error
}

// NewClient will return a new instance client
func NewClient(
	clientID,
	clientSecret,
	privateKey string,
	opts ...ClientOptionFunc,
) (*Client, error) {
	if clientID == "" {
		return nil, errors.New(EMPTY_CLIENT_ID)
	}

	if clientSecret == "" {
		return nil, errors.New(EMPTY_CLIENT_SECRET)
	}

	if privateKey == "" {
		return nil, errors.New(EMPTY_CLIENT_PRIVATE_KEY)
	}

	opt := defaultOption
	opt.clientID = clientID
	opt.clientSecret = clientSecret
	opt.clientPrivateKey = privateKey

	for _, option := range opts {
		option(&opt)
	}

	return &Client{
		http:           initHttpClient(&opt),
		logger:         opt.logger,
		satBaseURL:     opt.satBaseURL,
		accessTokenURL: opt.accessTokenURL,
		signature: signature.Init(signature.Options{
			PrivateKeyString: opt.clientPrivateKey,
			PublicKeyString:  opt.serverPublicKey,
			PaddingType:      opt.paddingType,
		}),
		isDebug: opt.isDebug,
	}, nil
}

func initHttpClient(
	cfg *Option,
) *http.Client {
	if cfg.http == nil {
		cfg.http = &http.Client{}
	}

	logr := logger.Config{
		Logger:  cfg.logger,
		IsDebug: cfg.isDebug,
	}

	cfg.http.Transport = &logger.Transport{
		Source: logr.GetLogger(),
		Base:   cfg.http.Transport,
	}

	cc := clientcredentials.Config{
		ClientID:     cfg.clientID,
		ClientSecret: cfg.clientSecret,
		TokenURL:     cfg.accessTokenURL,
	}

	cfg.http.Transport = &oauth2.Transport{
		Source: cc.TokenSource(context.Background()),
		Base:   cfg.http.Transport,
	}

	return cfg.http
}

// Ping is a method to check the SAT server health
func (c *Client) Ping(ctx context.Context) (*PingResponse, error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.satBaseURL+PING_PATH, nil)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	c.applyCustomHeader(hreq)

	resp, err := c.http.Do(hreq)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	response := new(PingResponse)
	err = json.Unmarshal(body, response)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}
	return response, nil
}

// Account is a method to check account balance
func (c *Client) Account(ctx context.Context) (*Account, error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.satBaseURL+ACCOUNT_PATH, nil)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	c.applyCustomHeader(hreq)

	resp, err := c.http.Do(hreq)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	response := new(Account)
	err = jsonapi.UnmarshalPayload(resp.Body, response)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	return response, nil
}

// Inquiry is a method to get user bills based on client number and product code
func (c *Client) Inquiry(ctx context.Context, req *InquiryRequest) (*InquiryResponse, error) {
	body := &bytes.Buffer{}
	err := jsonapi.MarshalPayload(body, req)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.satBaseURL+INQUIRY_PATH, body)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	c.applyCustomHeader(hreq)

	resp, err := c.http.Do(hreq)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	response := new(InquiryResponse)
	err = jsonapi.UnmarshalPayload(resp.Body, response)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	return response, nil
}

// Checkout is a method to do payment an order based on client number, product code and request id.
// Request ID should use unique identifier for each transaction
func (c *Client) Checkout(ctx context.Context, req *OrderRequest) (*OrderDetail, error) {
	body := &bytes.Buffer{}
	err := jsonapi.MarshalPayload(body, req)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	hreq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.satBaseURL+CHECKOUT_PATH, body)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	sign, err := c.signature.Sign(body.Bytes())
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	hreq.Header.Add(SIGNATURE_HEADER_KEY, sign)
	c.applyCustomHeader(hreq)

	resp, err := c.http.Do(hreq)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	response := new(OrderDetail)
	err = jsonapi.UnmarshalPayload(resp.Body, response)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	return response, nil
}

// CheckStatus is a method to check the final status of an order.
// request id is must be filled
func (c *Client) CheckStatus(ctx context.Context, requestID string) (*OrderDetail, error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.satBaseURL+fmt.Sprintf(CHECK_STATUS_PATH, requestID), nil)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	c.applyCustomHeader(hreq)

	resp, err := c.http.Do(hreq)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	resp.Body = io.NopCloser(bytes.NewBuffer(body))
	err = c.signature.Verify(string(body), resp.Header.Get(SIGNATURE_HEADER_KEY))
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	response := new(OrderDetail)
	err = jsonapi.UnmarshalPayload(resp.Body, response)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	return response, nil
}

// ListProduct is a method to get all the product list enabled on your credentials.
// you can also specify the product code, to get only one product detail.
// specify product code will be very beneficial to sync product status on your engine
// it will come with low bandwidth and fast response
func (c *Client) ListProduct(ctx context.Context, code string) ([]*Product, error) {
	hreq, err := http.NewRequestWithContext(ctx, http.MethodGet, c.satBaseURL+PRODUCT_LIST_PATH, nil)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	q := hreq.URL.Query()
	q.Add("product_code", code)
	hreq.URL.RawQuery = q.Encode()
	c.applyCustomHeader(hreq)

	resp, err := c.http.Do(hreq)
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, c.handleErrorResponse(resp)
	}

	items, err := jsonapi.UnmarshalManyPayload(resp.Body, reflect.TypeOf(new(Product)))
	if err != nil {
		c.logger.Println(err)
		return nil, err
	}

	var response []*Product
	for _, item := range items {
		response = append(response, item.(*Product))
	}

	return response, nil
}

// GetHTTTPClient will return http client which already wrapped to support oauth2
// for custom integration to SAT Service
func (c *Client) GetHTTTPClient() *http.Client {
	return c.http
}

// HandleCallback is method http.HandlerFunc to handle callback request from SAT
// you can customize the implementation based on this interface Callback
func (c *Client) HandleCallback(impl Callback) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			c.logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = c.signature.Verify(string(body), req.Header.Get(SIGNATURE_HEADER_KEY))
		if err != nil {
			c.logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(INVALID_SIGNATURE))
			return
		}

		request := new(OrderDetail)
		err = jsonapi.UnmarshalPayload(bytes.NewReader(body), request)
		if err != nil {
			c.logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(INVALID_PAYLOAD))
			return
		}

		err = impl.Do(req.Context(), request)
		if err != nil {
			c.logger.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(SUCCESS_OK))
		return
	}
}

func (c *Client) applyCustomHeader(hreq *http.Request) {
	hreq.Header.Add("Date", time.Now().Format(http.TimeFormat))
	hreq.Header.Add("X-Sat-Sdk-Version", SAT_SDK_VERSION)
}

func (c *Client) handleErrorResponse(resp *http.Response) error {
	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "application/json") {
		return &InternalError{resp: resp}
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Println(err)
		return err
	}

	var errorResponse *ErrorResponse
	err = json.Unmarshal(b, &errorResponse)
	if err != nil {
		c.logger.Println(err)
		return err
	}

	return errorResponse
}
