package endpoint

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPGet(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `response from the mock server`)
	}))
	defer ts.Close()

	resp, _ := http.Get(ts.URL)

	// assert over resp and err here
	if resp.StatusCode != 200 {
		t.Errorf("Incorrect response")
	}
}
