package flip

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"moul.io/http2curl"
	"net/http"
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
	Logger     Logger
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

	logOption := LogOption{
		Format:          "text",
		Level:           "info",
		TimestampFormat: "2006-01-02T15:04:05-0700",
		CallerToggle:    false,
	}

	// default HTTP Option
	httpOption := HTTPOption{
		Timeout:           10 * time.Second,
		BackoffInterval:   2 * time.Millisecond,
		MaxJitterInterval: 5 * time.Millisecond,
		RetryCount:        0,
	}

	logger := *NewLogger(logOption)

	return Client{
		// LogLevel is the logging level used by the LinkAja library
		// 0: No logging
		// 1: Errors only
		// 2: Errors + informational (default)
		// 3: Errors + informational + debug
		LogLevel:   2,
		Logger:     logger,
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
	req, err := http.NewRequest(method, fullPath, body)
	if err != nil {
		c.Logger.Error("Request creation failed: %v", err)
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
	command, _ := http2curl.GetCurlCommand(req)
	start := time.Now()
	c.Logger.Info("Start requesting: %v ", req.URL)

	res, err := getHTTPClient(c.HTTPOption).Do(req)
	if err != nil {
		c.Logger.Error("Request failed. Error : %v , Curl Request : %v", err, command)
		return
	}
	defer res.Body.Close()

	c.Logger.Info("Completed in %v", time.Since(start))
	c.Logger.Info("Curl Request: %v ", command)

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.Logger.Error("Cannot read response body: %v ", err)
		return
	}

	c.Logger.Info("Flip HTTP status response : %d", res.StatusCode)
	c.Logger.Info("Flip response body : %s", string(resBody))

	// General Error like auth error and data not found
	var generalErrorResponse GeneralErrorResponse

	if v != nil && res.StatusCode == http.StatusOK {
		if err = json.Unmarshal(resBody, v); err != nil {
			respErr.Code = res.StatusCode

			if err = json.Unmarshal(resBody, &generalErrorResponse); err != nil {
				respErr.Message = err.Error()
				return
			}

			respErr.Message = generalErrorResponse.Message
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
