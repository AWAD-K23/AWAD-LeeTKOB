package semantic

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"leetkob/backend/internal/hardrules"
)

type fakeTransport func(*http.Request) (*http.Response, error)

func (f fakeTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestGeminiRequestAndAudit(t *testing.T) {
	var saved CallLog
	model, err := NewGemini(SupportedModel, "fake-secret", func(entry CallLog) error {
		saved = entry
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	model.client.Transport = fakeTransport(func(req *http.Request) (*http.Response, error) {
		if req.URL.Scheme != "https" || req.URL.Host != "generativelanguage.googleapis.com" ||
			req.Header.Get("x-goog-api-key") != "fake-secret" || req.URL.RawQuery != "" {
			t.Fatal("incorrect authentication or endpoint")
		}

		body, _ := io.ReadAll(req.Body)
		var sent map[string]json.RawMessage
		if json.Unmarshal(body, &sent) != nil || len(sent) != 3 ||
			!strings.Contains(string(body), "responseJsonSchema") ||
			strings.Contains(string(body), "fake-secret") ||
			!strings.Contains(string(body), "untrusted assignment data") {
			t.Fatalf("unexpected request: %s", body)
		}

		return &http.Response{
			StatusCode: 200, Header: http.Header{},
			Body: io.NopCloser(strings.NewReader(`{
			  "responseId":"provider-123", "modelVersion":"version-1",
			  "usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":20,"thoughtsTokenCount":10},
			  "candidates":[{"finishReason":"STOP","content":{"parts":[
			    {"thought":true,"text":"private reasoning"},
			    {"text":"{\"verdicts\":[{\"constraintId\":\"a\",\"verdict\":\"unsure\",\"evidence\":[]}]}"}
			  ]}}]
			}`)),
		}, nil
	})

	result, err := model.Generate(context.Background(), Request{
		Statement: "Use recursion", Source: "int main() {}",
		Constraints: []hardrules.Constraint{{ID: "a", Text: "Use recursion"}},
	})
	if err != nil || !json.Valid(result) || strings.Contains(string(result), "reasoning") ||
		saved.ProviderID != "provider-123" || saved.CostMicroUSD != 105 ||
		!saved.UsageKnown || saved.Outcome != "response_received" {
		t.Fatalf("result=%s log=%+v err=%v", result, saved, err)
	}

	audit, _ := json.Marshal(saved)
	if strings.Contains(string(audit), "fake-secret") {
		t.Fatal("key leaked into call log")
	}
}

func TestGeminiProviderFailures(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
	}{
		{"rate limit", 429, "fake-secret must not be logged"},
		{"redirect", 302, ""},
		{"bad envelope", 200, "invalid"},
		{"missing usage", 200, `{"responseId":"r","candidates":[]}`},
		{"oversized", 200, strings.Repeat("x", 1024*1024+1)},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			calls, records := 0, 0
			model, err := NewGemini(SupportedModel, "fake-secret", func(log CallLog) error {
				records++
				data, _ := json.Marshal(log)
				if strings.Contains(string(data), "fake-secret") {
					t.Fatal("provider error body leaked into audit")
				}
				return nil
			})
			if err != nil {
				t.Fatal(err)
			}

			model.client.Transport = fakeTransport(func(*http.Request) (*http.Response, error) {
				calls++
				return &http.Response{
					StatusCode: test.status, Header: http.Header{},
					Body: io.NopCloser(strings.NewReader(test.body)),
				}, nil
			})
			if _, err := model.Generate(context.Background(), Request{}); err == nil ||
				strings.Contains(err.Error(), "fake-secret") || calls != 1 || records != 1 {
				t.Fatalf("failure not bounded/audited: calls=%d records=%d err=%v", calls, records, err)
			}
		})
	}
}

func TestGeminiCancellationAndAuditFailure(t *testing.T) {
	model, err := NewGemini(SupportedModel, "fake", func(CallLog) error {
		return errors.New("audit unavailable")
	})
	if err != nil {
		t.Fatal(err)
	}

	model.client.Transport = fakeTransport(func(req *http.Request) (*http.Response, error) {
		<-req.Context().Done()
		return nil, req.Context().Err()
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := model.Generate(ctx, Request{}); err == nil {
		t.Fatal("cancelled request succeeded")
	}

	if _, err := NewGemini("unpriced-model", "fake", func(CallLog) error { return nil }); err == nil {
		t.Fatal("unknown model accepted without reviewed rates")
	}
}
