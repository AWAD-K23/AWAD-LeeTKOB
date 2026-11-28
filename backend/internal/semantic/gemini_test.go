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

// The model must copy line numbers, not count them: on 2026-11-24 the deployed
// model cited every line one too low, and strict evidence validation turned
// three correct verdicts into unsure.
func TestGeminiSendsNumberedLines(t *testing.T) {
	var saved CallLog
	model, err := NewGemini(SupportedModel, "fake-secret", func(entry CallLog) error {
		saved = entry
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	source := "int main() {\n\n    return 0;\n}"
	model.client.Transport = fakeTransport(func(req *http.Request) (*http.Response, error) {
		var sent struct {
			SystemInstruction struct {
				Parts []struct{ Text string } `json:"parts"`
			} `json:"systemInstruction"`
			Contents []struct {
				Parts []struct{ Text string } `json:"parts"`
			} `json:"contents"`
		}
		body, _ := io.ReadAll(req.Body)
		if json.Unmarshal(body, &sent) != nil || len(sent.Contents) != 1 ||
			len(sent.Contents[0].Parts) != 1 || len(sent.SystemInstruction.Parts) != 1 {
			t.Fatalf("unexpected request: %s", body)
		}
		if !strings.Contains(sent.SystemInstruction.Parts[0].Text, "numbered lines") {
			t.Fatalf("instruction does not explain numbered lines: %q", sent.SystemInstruction.Parts[0].Text)
		}

		var prompt struct {
			Code []struct {
				Line int    `json:"line"`
				Text string `json:"text"`
			} `json:"code"`
		}
		if err := json.Unmarshal([]byte(sent.Contents[0].Parts[0].Text), &prompt); err != nil {
			t.Fatalf("code is not a list of numbered lines: %v", err)
		}
		want := []string{"int main() {", "", "    return 0;", "}"}
		if len(prompt.Code) != len(want) {
			t.Fatalf("got %d lines, want %d", len(prompt.Code), len(want))
		}
		for i, line := range prompt.Code {
			if line.Line != i+1 || line.Text != want[i] {
				t.Fatalf("line %d = %+v, want {%d %q}", i, line, i+1, want[i])
			}
		}

		return &http.Response{
			StatusCode: 200, Header: http.Header{},
			Body: io.NopCloser(strings.NewReader(`{
			  "responseId":"provider-1",
			  "usageMetadata":{"promptTokenCount":100,"candidatesTokenCount":20},
			  "candidates":[{"finishReason":"STOP","content":{"parts":[
			    {"text":"{\"verdicts\":[{\"constraintId\":\"a\",\"verdict\":\"unsure\",\"evidence\":[]}]}"}
			  ]}}]
			}`)),
		}, nil
	})

	if _, err := model.Generate(context.Background(), Request{
		Statement: "Use recursion", Source: source,
		Constraints: []hardrules.Constraint{{ID: "a", Text: "Use recursion"}},
	}); err != nil {
		t.Fatal(err)
	}
	// The audit keeps the plain source: the same data, without the wire framing.
	if saved.Prompt.Source != source {
		t.Fatalf("audit prompt changed: %q", saved.Prompt.Source)
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

func TestGeminiConfirmationDoesNotRevealPriorFinding(t *testing.T) {
	var saved CallLog
	model, err := NewGemini(SupportedModel, "fake-test-key", func(entry CallLog) error { saved = entry; return nil })
	if err != nil {
		t.Fatal(err)
	}
	model.wait = func(context.Context) error { return nil }
	model.client.Transport = fakeTransport(func(req *http.Request) (*http.Response, error) {
		var sent struct {
			Contents []struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"contents"`
		}
		data, _ := io.ReadAll(req.Body)
		if err := json.Unmarshal(data, &sent); err != nil {
			t.Fatal(err)
		}
		var prompt wireRequest
		if len(sent.Contents) != 1 || len(sent.Contents[0].Parts) != 1 {
			t.Fatalf("missing prompt: %s", data)
		}
		if err := json.Unmarshal([]byte(sent.Contents[0].Parts[0].Text), &prompt); err != nil {
			t.Fatal(err)
		}
		if len(prompt.Constraints) != 1 || prompt.Constraints[0].ID != "a" || len(prompt.Code) != 1 || prompt.Code[0].Text != "int quantity = 1;" {
			t.Fatalf("confirmation data lost: %+v", prompt)
		}
		if strings.Contains(sent.Contents[0].Parts[0].Text, "proposedViolations") {
			t.Fatal("confirmation is anchored by the primary finding")
		}
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{"responseId":"fake-confirmation-id","usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1},"candidates":[{"finishReason":"STOP","content":{"parts":[{"text":"{\"verdicts\":[{\"constraintId\":\"a\",\"verdict\":\"unsure\",\"evidence\":[]}]}"}]}}]}`))}, nil
	})
	_, err = model.Generate(context.Background(), Request{Statement: "Review", Source: "int quantity = 1;", Constraints: []hardrules.Constraint{{ID: "a", Text: "Use floating point"}},
		ProposedViolations: []ProposedViolation{{ConstraintID: "a", Evidence: []hardrules.Evidence{{Line: 1, Text: "int quantity = 1;"}}}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(saved.Prompt.ProposedViolations) != 0 {
		t.Fatal("audit claims prior findings were sent during independent confirmation")
	}
}
