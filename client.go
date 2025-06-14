package detectlanguage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"
)

const defaultUserAgent = "detectlanguage-go/" + Version

const defaultTimeout = 10 * time.Second

var apiBaseURL = &url.URL{
	Scheme: "https", Host: "ws.detectlanguage.com", Path: "/0.2/",
}

// A Client provides an HTTP client for DetectLanguage API operations.
type Client struct {
	// BaseURL specifies the location of the API. It is used with
	// ResolveReference to create request URLs. (If 'Path' is specified, it
	// should end with a trailing slash.) If nil, the default will be used.
	BaseURL *url.URL
	// Client is an HTTP client used to make API requests. If nil,
	// default will be used.
	Client *http.Client
	// APIKey is the user's API key. It is required.
	// Note: Treat your API Keys as passwords—keep them secret. API Keys give
	// full read/write access to your account, so they should not be included in
	// public repositories, emails, client side code, etc.
	APIKey string
	// UserAgent is a User-Agent to be sent with API HTTP requests. If empty,
	// a default will be used.
	UserAgent string
}

// New returns a new Client with the given API key.
func New(apiKey string) *Client {
	return &Client{APIKey: apiKey}
}

func (c *Client) baseURL() *url.URL {
	if c.BaseURL != nil {
		return c.BaseURL
	}
	return apiBaseURL
}

func (c *Client) userAgent() string {
	if c.UserAgent != "" {
		return c.UserAgent
	}
	return defaultUserAgent
}

func (c *Client) client() *http.Client {
	if c.Client == nil {
		c.Client = &http.Client{Timeout: defaultTimeout}
	}
	return c.Client
}

func (c *Client) setBody(req *http.Request, in interface{}) error {
	if in != nil {
		buf, err := json.Marshal(in)
		if err != nil {
			return err
		}
		req.Body = io.NopCloser(bytes.NewReader(buf))
		req.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(buf)), nil
		}
		req.Header.Set("Content-Type", "application/json")
		req.ContentLength = int64(len(buf))
	}
	return nil
}

func (c *Client) do(ctx context.Context, method, path string, in, out interface{}) error {
	req := &http.Request{
		Method: method,
		URL:    c.baseURL().ResolveReference(&url.URL{Path: path}),
		Header: make(http.Header, 2),
	}
	req.Header.Set("User-Agent", c.userAgent())
	if err := c.setBody(req, in); err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+c.APIKey)
	if ctx != nil {
		req = req.WithContext(ctx)
	}
	res, err := c.client().Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		if out != nil {
			return json.NewDecoder(res.Body).Decode(out)
		}
		return nil
	}

	buf, _ := io.ReadAll(res.Body)
	apiErr := &APIError{Status: res.Status, StatusCode: res.StatusCode}
	if json.Unmarshal(buf, &apiErrorResponse{Error: apiErr}) != nil {
		apiErr.Message = string(buf)
	}
	return apiErr
}

func (c *Client) get(ctx context.Context, path string, out interface{}) error {
	return c.do(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) post(ctx context.Context, path string, in, out interface{}) error {
	return c.do(ctx, http.MethodPost, path, in, out)
}
