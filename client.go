//go:generate easyjson -all $GOFILE

// Package twitch provides a client and a set of functions to retrieve data
// from the Twitch API.
package twitch

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const (
	// ClientID is the string used to specify the Client-ID header for authorization.
	ClientID = "Client-ID"
)

const (
	// EnvTwitchClientID is a string specifying the Client-ID utilized to make Twitch API calls.
	EnvTwitchClientID = "twitch_client_id"
)

// Client specifies a set of process dependencies for retrieving data from the
// Twitch API.
type Client struct {
	viper      *viper.Viper
	httpClient *http.Client
}

// NewClient intializes a new twitch Client based on the ClientOption(s).
func NewClient(options ...ClientOption) *Client {
	var c = new(Client)

	for _, option := range options {
		option(c)
	}

	// required fields
	if c.httpClient == nil {
		log.Fatal("failed to initialize httpClient field")
	}
	if c.viper == nil {
		log.Fatal("failed to initialize viper field")
	}
	return c
}

// ClientOption modifies a Client object.
type ClientOption func(*Client)

// WithHttpClient initializes the Client object with a http.Client.
func WithHttpClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithViper initializes the Client object with a viper.Viper object.
func WithViper(v *viper.Viper) ClientOption {
	return func(c *Client) {
		v.SetDefault(EnvTwitchClientID, "kkjo0gmafpxng9jlmod6g8x0z7rjjj")
		c.viper = v
	}
}

// Close closes resources hanging off of the Client object.
func (c *Client) Close() {
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
	}
}

// StreamsURL specifies the twitch API Streams URL
const StreamsURL = "https://api.twitch.tv/helix/streams"

// Streams retrieves streams as specified by the by function.
func (c *Client) Streams(by StreamsBy) (*bytes.Buffer, error) {
	return by()
}

// StreamsBy is a function that retrieves a set of Stream objects.
type StreamsBy func() (*bytes.Buffer, error)

// ByUserLogin returns a method by which to retrieve a set of Stream objects by
// the specified userLogin.
func (c *Client) ByUserLogin(userLogin string) StreamsBy {
	return func() (*bytes.Buffer, error) {
		var url = StreamsURL + "?user_login=" + userLogin

		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to create streams http.Request")
		}
		req.Header.Set(ClientID, c.viper.GetString(EnvTwitchClientID))

		resp, err := c.httpClient.Do(req)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to retrieve Streams by userLogin %s", userLogin)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.Errorf("failed to retrieve Streams by userLogin %s, StatusCode = %v", userLogin, resp.StatusCode)
		}

		var streams = new(bytes.Buffer)
		if _, err := io.Copy(streams, resp.Body); err != nil {
			return nil, errors.Wrap(err, "failed to Copy streams into buffer")
		}
		return streams, nil
	}
}
