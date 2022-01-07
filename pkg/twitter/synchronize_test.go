package twitter

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/stretchr/testify/assert"
)

type RewriteTransport struct {
	Transport http.RoundTripper
}

// RoundTrip rewrites the request scheme to http and calls through to the
// composed RoundTripper or if it is nil, to the http.DefaultTransport.
func (t *RewriteTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.URL.Scheme = "http"
	if t.Transport == nil {
		return http.DefaultTransport.RoundTrip(req)
	}
	return t.Transport.RoundTrip(req)
}

func testServer() (*http.Client, *http.ServeMux, *httptest.Server) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	transport := &RewriteTransport{&http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}}
	client := &http.Client{Transport: transport}
	return client, mux, server
}

func makeTestClient() *twitter.Client {

	httpClient, mux, _ := testServer()
	// httpClient, mux, server := testServer()
	// defer server.Close()

	mux.HandleFunc("/1.1/friends/ids.json", func(w http.ResponseWriter, r *http.Request) {
		// assertMethod(t, "GET", r)
		// assertQuery(t, map[string]string{"user_id": "623265148", "count": "5", "cursor": "1516933260114270762"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ids":[178082406,3318241001,1318020818,191714329,376703838],"next_cursor":1516837838944119498,"next_cursor_str":"1516837838944119498","previous_cursor":-1516924983503961435,"previous_cursor_str":"-1516924983503961435"}`)
	})
	mux.HandleFunc("/1.1/users/show.json", func(w http.ResponseWriter, r *http.Request) {
		// assertMethod(t, "GET", r)
		// assertQuery(t, map[string]string{"user_id": "623265148", "count": "5", "cursor": "1516933260114270762"}, r)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"ids":[178082406,3318241001,1318020818,191714329,376703838],"next_cursor":1516837838944119498,"next_cursor_str":"1516837838944119498","previous_cursor":-1516924983503961435,"previous_cursor_str":"-1516924983503961435"}`)
	})

	client := twitter.NewClient(httpClient)
	return client
}
func Test_synchronize(t *testing.T) {

	t.Run("Synchronize", func(t *testing.T) {
		expected := &twitter.FriendIDs{
			IDs:               []int64{178082406, 3318241001, 1318020818, 191714329, 376703838},
			NextCursor:        1516837838944119498,
			NextCursorStr:     "1516837838944119498",
			PreviousCursor:    -1516924983503961435,
			PreviousCursorStr: "-1516924983503961435",
		}

		err = Synchronize(makeTestClient())
		assert.Nil(t, err)

		// count persisted friends
		var c int64
		err = db.Model(&User{}).Count(&c).Error
		assert.Nil(t, err)
		assert.Equal(t, int64(len(expected.IDs)), c)
	})
}
