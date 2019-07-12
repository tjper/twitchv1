package twitch_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/spf13/viper"
	twitch "github.com/tjper/twitchv1"
)

func TestStreamsByUserLogin(t *testing.T) {
	var tests = []struct {
		name      string
		userLogin string
	}{
		{name: "summit1g test", userLogin: "summit1g"},
		{name: "penutty test", userLogin: "penutty"},
	}

	for _, test := range tests {
		var client = twitch.NewClient(
			twitch.WithHttpClient(&http.Client{}),
			twitch.WithViper(viper.New()),
		)
		t.Run(test.name, func(t *testing.T) {
			streams, err := client.Streams(client.ByUserLogin(test.userLogin))
			if err != nil {
				t.Fatalf("unexpected err \"%s\"", err)
			}
			defer streams.Close()
			b, err := ioutil.ReadAll(streams)
			if err != nil {
				t.Fatalf("unexpected err \"%s\"", err)
			}
			if len(b) == 0 {
				t.Fatalf("expected non-empty streams to be returned\nstreams = %v", streams)
			}
		})
	}
}
