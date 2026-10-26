package main

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type releaseRoundTripper func(*http.Request) (*http.Response, error)

func (f releaseRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestForkReleaseChannel(t *testing.T) {
	// Literal expectations deliberately catch accidentally restoring the upstream channel.
	if releaseLatestURL != "https://github.com/lyl7410/CFData-WEB/releases/latest" {
		t.Fatalf("unexpected download channel: %s", releaseLatestURL)
	}
	for _, tc := range []struct {
		name      string
		status    int
		body      string
		wantError bool
	}{
		{"release", 200, `{"tag_name":"v1.8.6","html_url":"https://github.com/lyl7410/CFData-WEB/releases/tag/v1.8.6"}`, false},
		{"missing-url", 200, `{"tag_name":"v1.8.6"}`, false},
		{"no-fork-release-yet", 404, `{"message":"Not Found"}`, true},
		{"rate-limit", 403, `{"message":"rate limit"}`, true},
		{"invalid-json", 200, `{`, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := &http.Client{Transport: releaseRoundTripper(func(r *http.Request) (*http.Response, error) {
				if r.URL.String() != "https://api.github.com/repos/lyl7410/CFData-WEB/releases/latest" {
					t.Fatalf("unexpected API channel: %s", r.URL)
				}
				if r.Method != http.MethodGet || r.Header.Get("Accept") != "application/vnd.github+json" || r.Header.Get("User-Agent") != "CFData-WEB/"+appVersion {
					t.Fatal("unexpected request method or headers")
				}
				if _, ok := r.Context().Deadline(); !ok {
					t.Fatal("request needs a timeout")
				}
				return &http.Response{StatusCode: tc.status, Body: io.NopCloser(strings.NewReader(tc.body)), Header: make(http.Header)}, nil
			})}
			info, err := getLatestReleaseWithClient(context.Background(), client)
			if (err != nil) != tc.wantError {
				t.Fatalf("error = %v; want error %v", err, tc.wantError)
			}
			if !tc.wantError && (info.TagName != "v1.8.6" || !strings.HasPrefix(info.HTMLURL, "https://github.com/lyl7410/CFData-WEB/releases/")) {
				t.Fatalf("unexpected release: %+v", info)
			}
		})
	}
}
