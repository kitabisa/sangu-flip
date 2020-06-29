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

func NewClient() Client {
	// default HTTP Option
	httpOption := HTTPOption{
		Timeout:           10 * time.Second,
		BackoffInterval:   2 * time.Millisecond,
		MaxJitterInterval: 5 * time.Millisecond,
		RetryCount:        0,
	}

	return Client{
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
func (c *Client) ExecuteRequest(req *http.Request, v interface{}) (err error, respErr ErrorResponse) {
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
		return
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
		return
	}

	if logLevel > 2 {
		logger.Println("Flip HTTP status response: ", res.StatusCode)
		logger.Println("Flip body response: ", string(resBody))
	}

	if v != nil && res.StatusCode == http.StatusOK {
		if err = json.Unmarshal(resBody, v); err != nil {
			respErr.Code = res.StatusCode
			respErr.Message = err.Error()
			return
		}
	}

	if res.StatusCode != http.StatusOK {
		// Disbursement validation error
		if res.StatusCode == http.StatusUnprocessableEntity {
			var disbursementError DisbursementError
			if err = json.Unmarshal(resBody, &disbursementError); err != nil {
				respErr.Message = err.Error()
				return
			}
			respErr.Code = res.StatusCode
			respErr.Message = disbursementError.Code
			respErr.Error = &disbursementError
			return
		}

		// General Error like auth error and data not found
		var generalErrorResponse GeneralErrorResponse
		if err = json.Unmarshal(resBody, &generalErrorResponse); err != nil {
			respErr.Message = err.Error()
			return
		}
		respErr.Code = res.StatusCode
		respErr.Message = generalErrorResponse.Message
		respErr.Error = &generalErrorResponse
	}

	return
}

// Call the BigFlip API at specific `path` using the specified HTTP `method`. The result will be
// given to `v` if there is no error. If any error occurred and Bigflip send error response, the result will be error code, error message and data, otherwise only error code and error message
func (c *Client) Call(method, path string, header map[string]string, body io.Reader, v interface{}) (err error, respErr ErrorResponse) {
	req, err := c.NewRequest(method, path, header, body)
	if err != nil {
		return
	}

	err, respErr = c.ExecuteRequest(req, v)
	return
}

// ===================== END HTTP CLIENT ================================================
