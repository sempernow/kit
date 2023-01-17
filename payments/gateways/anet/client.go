package anet

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"
)

/************************************************************************
	TODO : Migrate to more robust client, e.g., github.com/imroc/req/v3,
	but mind the BOM issue; see Do(..).
*************************************************************************/

const (
	EndpointSandbox = "https://apitest.authorize.net/xml/v1/request.api"
	EndpointLive    = "https://api.authorize.net/xml/v1/request.api"

	DefaultTimeOut               = 5 * time.Second
	DefaultKeepAlive             = 5 * time.Second
	DefaultMaxIdleConns          = 100
	DefaultMaxIdleConnsPerHost   = 100
	DefaultIdleConnTimeout       = 90 * time.Second
	DefaultTLSHandshakeTimeout   = 5 * time.Second
	DefaultExpectContinueTimeout = 1 * time.Second
	DefaultMaxConnsPerHost       = 50

	ContentType     = "Content-Type"
	Accept          = "Accept"
	ApplicationJSON = "application/json; charset=utf-8"
	UserAgent       = "User-Agent"
	Authorization   = "Authorization"
	Bearer          = "Bearer "
)

type Logger interface {
	Printf(format string, v ...interface{})
}

type HttpClientCfg struct {
	BasePath      string            `json:"basePath,omitempty"`
	Host          string            `json:"host,omitempty"`
	Scheme        string            `json:"scheme,omitempty"`
	DefaultHeader map[string]string `json:"defaultHeader,omitempty"`
	UserAgent     string            `json:"userAgent,omitempty"`
	HTTPClient    *http.Client
}

type Client struct {
	cfg *HttpClientCfg
	log Logger
}

type API struct {
	MerchantName   string
	MerchantSecret string
	TestMode       bool
	URL            string
}

func SetAPI(name string, key string, mode string) API {
	api := API{}
	api.MerchantName = name
	api.MerchantSecret = key
	if mode == "live" {
		api.TestMode = true
		api.URL = EndpointLive
	} else {
		api.TestMode = false
		api.URL = EndpointSandbox
	}
	return api
}

//	func Endpoint() string {
//		if testMode {
//			return EndpointSandbox
//		}
//		return EndpointLive
//	}
func NewAPIClient(cfg *HttpClientCfg, logger Logger) *Client {
	if cfg == nil {
		cfg = &HttpClientCfg{}
	}
	if cfg.HTTPClient == nil {
		// build a "default" http client
		cfg.HTTPClient = &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: DefaultMaxIdleConnsPerHost,
				Proxy:               http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   DefaultTimeOut,
					KeepAlive: DefaultKeepAlive,
				}).DialContext,
				MaxIdleConns:          DefaultMaxIdleConns,
				IdleConnTimeout:       DefaultIdleConnTimeout,
				TLSHandshakeTimeout:   DefaultTLSHandshakeTimeout,
				ExpectContinueTimeout: DefaultExpectContinueTimeout,
				MaxConnsPerHost:       DefaultMaxConnsPerHost,
			},
		}
	}

	result := Client{
		cfg: cfg,
		log: logger,
	}

	return &result
}

// c.Send the request to Authorize.net API, returning the response by reference.
func (c Client) Send(ctx context.Context, url string, req, rtn interface{}) error {

	request, err := c.prepareRequest(ctx, url, req)
	if err != nil {
		return errors.Wrap(err, "client send : prepare request")
	}
	return c.Do(request, rtn, false)

}

// c.prepareRequest builds the http.Request object for our client.
func (c Client) prepareRequest(
	ctx context.Context,
	path string,
	postBody interface{},
) (*http.Request, error) {

	var (
		body   *bytes.Buffer
		err    error
		method = http.MethodPost
	)
	// Detect postBody type and post.
	if postBody != nil {
		body = &bytes.Buffer{}
		if err = json.NewEncoder(body).Encode(postBody); err != nil {
			return nil, err
		}
	}

	// Setup path and query parameters
	uri, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	// Generate a new request
	var request *http.Request
	if body != nil {
		if request, err = http.NewRequest(method, uri.String(), body); err != nil {
			return nil, err
		}
	} else {
		if request, err = http.NewRequest(method, uri.String(), nil); err != nil {
			return nil, err
		}
	}

	// set Content-Type header
	request.Header.Set(ContentType, ApplicationJSON)
	// set Accept header
	request.Header.Set(Accept, ApplicationJSON)

	// Override request host, if applicable
	if c.cfg.Host != "" {
		request.Host = c.cfg.Host
	}

	// Add the user agent to the request.
	request.Header.Add(UserAgent, c.cfg.UserAgent)

	// yeah, context should never be nil, but we never know
	if ctx != nil {
		// add context to the request, for cancellation
		request = request.WithContext(ctx)
	}
	// add rest of the "default" headers (if any)
	for header, value := range c.cfg.DefaultHeader {
		request.Header.Add(header, value)
	}
	return request, nil
}

// Do the prepared HTTP request, returning the result by reference.
func (c Client) Do(req *http.Request, result interface{}, printRaw bool) error {

	response, err := c.cfg.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "client do : resp")
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return errors.Wrap(err, "client do : reading resp body")
	}

	// if response.StatusCode != http.StatusOK {
	// 	fmt.Printf("HTTP Response: %v (%s)\n",
	// 		response.StatusCode, response.Status,
	// 	)
	// }

	// if printRaw {
	// 	fmt.Println("---")
	// 	fmt.Println(string(body))
	// 	fmt.Println("---")
	// }

	// Must remove the BOM (inserted by Anet API) prior to decoding:
	// https://en.wikipedia.org/wiki/Byte_order_mark
	body = bytes.TrimPrefix(body, []byte{239, 187, 191}) // or []byte("\xef\xbb\xbf")

	if err := json.Unmarshal(body, &result); err != nil {
		// var x interface{} // @ response is invalid result (struct).
		// if err := json.Unmarshal(body, &x); err != nil {
		//return errors.Wrapf(err, "unmarshal body: %s", string(body))
		return errors.Wrap(err, "client do : unmarshal body")
		// }
		// //return errors.New(PrettyPrint(x))
		// //... our result struct expects AccountTypeEnum uint,
		// // but @ AcceptPayment() returns field: `"accountType": ""`.
		// // So, we get: Error : invalid AccountTypeEnum : "\"\""
		// //... UPDATE : Added AccountTypeEmpty  ("")
	}

	return nil
}
