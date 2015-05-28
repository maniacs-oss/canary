package sampler

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSample(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
	}

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}
}

func TestSampleWithHeaders(t *testing.T) {
	headerName := "X-Request-Id"
	headerVal := "abcd-1234"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerName, headerVal)

		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
	}

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}

	if sample.ResponseHeaders.Get(headerName) != headerVal {
		t.Fatalf("Expected header %s to equal %s but got %s", headerName, headerVal, sample.ResponseHeaders.Get(headerName))
	}
}

func TestSampleWithCanonicalizedHeaderName(t *testing.T) {
	headerName := "x-request-id"
	headerVal := "abcd-1234"

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Request-Id", headerVal)

		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
	}

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}

	if sample.ResponseHeaders.Get(headerName) != headerVal {
		t.Fatalf("Expected header %s to equal %s but got %s", headerName, headerVal, sample.ResponseHeaders.Get(headerName))
	}
}

func TestSampleWithMissingHeader(t *testing.T) {
	headerName := "X-Request-Id"

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
	}

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}

	if val, ok := sample.ResponseHeaders[headerName]; ok {
		t.Fatalf("Expected header %s with missing value to be empty but was '%+v'", headerName, val)
	}
}

func TestSampleWithRequestHeaders(t *testing.T) {
	var header http.Header

	handler := func(w http.ResponseWriter, r *http.Request) {
		header = r.Header

		fmt.Fprintf(w, "ok")
	}
	ts := httptest.NewServer(http.HandlerFunc(handler))
	defer ts.Close()

	target := Target{
		URL: ts.URL,
		RequestHeaders: map[string]string{
			"X-Foo": "bar",
		},
	}

	sample, err := Ping(target, 1)
	if err != nil {
		t.Fatal(err)
	}

	if sample.StatusCode != 200 {
		t.Fatalf("Expected sampleStatus == 200, but got %d\n", sample.StatusCode)
	}

	h := header.Get("X-Foo")
	if h != "bar" {
		t.Fatalf("Expected request header X-Foo to be 'bar' but was '%s'", h)
	}
}

func TestGenRequest(t *testing.T) {
	expected := "GET / HTTP/1.1\r\nHost: canary.io\r\n\r\n"
	target := Target{
		URL: "http://canary.io",
	}

	req, err := genRequest(target)
	if err != nil {
		t.Fatalf("err while generating request: %v\n", err)
	}

	if req != expected {
		t.Fatalf("Expected request to look like:\n%s\n but got:\n%s\n", expected, req)
	}

}

func TestGenRequestWithCustomHost(t *testing.T) {
	expected := "GET / HTTP/1.1\r\nHost: canary.io\r\n\r\n"

	headers := make(map[string]string)
	headers["Host"] = "canary.io"

	target := Target{
		URL:            "http://192.168.1.1",
		RequestHeaders: headers,
	}

	req, err := genRequest(target)
	if err != nil {
		t.Fatalf("err while generating request: %v\n", err)
	}

	if req != expected {
		t.Fatalf("Expected request to look like:\n%s\n but got:\n%s\n", expected, req)
	}

}
