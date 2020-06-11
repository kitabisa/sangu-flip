package flip

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gojek/heimdall"
	"github.com/gojek/heimdall/httpclient"
)

// Client - Bigflip Client data
// UserKey is base64 BasicAuth
type Client struct {
	BaseURL    string
	UserKey    string
	LogLevel   int
	Logger     *log.Logger
	HTTPOption HTTPOption
}

// HTTPOption for heimdall properties
type HTTPOption struct {
	Timeout           time.Duration
	BackoffInterval   time.Duration
	MaxJitterInterval time.Duration
	RetryCount        int
}

func NewClient(baseUrl string, userKey string) Client {
	// default HTTP Option
	httpOption := HTTPOption{
		Timeout:           10 * time.Second,
		BackoffInterval:   2 * time.Millisecond,
		MaxJitterInterval: 5 * time.Millisecond,
		RetryCount:        3,
	}

	return Client{
		BaseURL: baseUrl,
		UserKey: userKey,
		// LogLevel is the logging level used by the LinkAja library
		// 0: No logging
		// 1: Errors only
		// 2: Errors + informational (default)
		// 3: Errors + informational + debug
		LogLevel:   2,
		Logger:     log.New(os.Stderr, "", log.LstdFlags),
		HTTPOption: httpOption,
	}
}

// getHTTPClient will get heimdall http client
func getHTTPClient(opt HTTPOption) *httpclient.Client {
	backoff := heimdall.NewConstantBackoff(opt.BackoffInterval, opt.MaxJitterInterval)
	retrier := heimdall.NewRetrier(backoff)

	return httpclient.NewClient(
		httpclient.WithHTTPTimeout(opt.Timeout),
		httpclient.WithRetrier(retrier),
		httpclient.WithRetryCount(opt.RetryCount),
	)
}

// NewRequest : send new request
func (c *Client) NewRequest(method string, fullPath string, headers map[string]string, body io.Reader) (*http.Request, error) {
	logLevel := c.LogLevel
	logger := c.Logger

	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Request creation failed: ", err)
		}
		return nil, err
	}

	if headers != nil {
		for k, vv := range headers {
			req.Header.Set(k, vv)
		}
	}

	req.SetBasicAuth(c.UserKey, "")

	return req, nil
}

// ExecuteRequest : execute request
func (c *Client) ExecuteRequest(req *http.Request, v interface{}, x interface{}) error {
	logLevel := c.LogLevel
	logger := c.Logger

	if logLevel > 1 {
		logger.Println("Request ", req.Method, ": ", req.URL.Host, req.URL.Path)
	}

	start := time.Now()
	res, err := getHTTPClient(c.HTTPOption).Do(req)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Request failed: ", err)
		}
		return err
	}
	defer res.Body.Close()

	if logLevel > 2 {
		logger.Println("Completed in ", time.Since(start))
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if logLevel > 0 {
			logger.Println("Cannot read response body: ", err)
		}
		return err
	}

	if logLevel > 2 {
		logger.Println("Flip HTTP status response: ", res.StatusCode)
		logger.Println("Flip body response: ", string(resBody))
	}

	if v != nil && res.StatusCode == 200 {
		if err = json.Unmarshal(resBody, v); err != nil {
			return err
		}
	}

	if x != nil && res.StatusCode != 200 {
		if err = json.Unmarshal(resBody, x); err != nil {
			return err
		}
	}

	return nil
}

// Call the LinkAja API at specific `path` using the specified HTTP `method`. The result will be
// given to `v` if there is no error. If any error occurred, the return of this function is the error
// itself, otherwise nil.
func (c *Client) Call(method, path string, header map[string]string, body io.Reader, v interface{}, x interface{}) error {
	req, err := c.NewRequest(method, path, header, body)
	if err != nil {
		return err
	}

	return c.ExecuteRequest(req, v, x)
}

// ===================== END HTTP CLIENT ================================================
