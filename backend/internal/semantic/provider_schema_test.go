package semantic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/santhosh-tekuri/jsonschema/v6"
)

// Validate counterexamples against the actual schema sent over the wire, using
// a real JSON Schema engine. This catches a provider/local-contract mismatch
// without changing the local validation guard or making a network request.
func TestProviderSchemaRequiresDecisiveEvidence(t *testing.T) {
	model, err := NewGemini(SupportedModel, "fake-schema-key", func(CallLog) error { return nil })
	if err != nil {
		t.Fatal(err)
	}
	model.wait = func(context.Context) error { return nil }
	model.client.Transport = fakeTransport(func(req *http.Request) (*http.Response, error) {
		var sent struct {
			Config struct {
				Schema json.RawMessage `json:"responseJsonSchema"`
			} `json:"generationConfig"`
		}
		body, _ := io.ReadAll(req.Body)
		if err := json.Unmarshal(body, &sent); err != nil {
			t.Fatal(err)
		}
		doc, err := jsonschema.UnmarshalJSON(bytes.NewReader(sent.Config.Schema))
		if err != nil {
			t.Fatal(err)
		}
		compiler := jsonschema.NewCompiler()
		if err = compiler.AddResource("schema.json", doc); err != nil {
			t.Fatal(err)
		}
		schema, err := compiler.Compile("schema.json")
		if err != nil {
			t.Fatal(err)
		}
		for _, tc := range []struct {
			name, raw string
			valid     bool
		}{
			{"violation without evidence", `{"verdicts":[{"constraintId":"a","verdict":"violation","evidence":[]}]}`, false},
			{"pass without evidence", `{"verdicts":[{"constraintId":"a","verdict":"pass","evidence":[]}]}`, false},
			{"unsure without evidence", `{"verdicts":[{"constraintId":"a","verdict":"unsure","evidence":[]}]}`, true},
			{"evidenced violation", `{"verdicts":[{"constraintId":"a","verdict":"violation","evidence":[{"line":1,"text":"int item = 1;"}]}]}`, true},
			{"evidenced pass", `{"verdicts":[{"constraintId":"a","verdict":"pass","evidence":[{"line":1,"text":"int item = 1;"}]}]}`, true},
			{"unknown verdict", `{"verdicts":[{"constraintId":"a","verdict":"accepted","evidence":[]}]}`, false},
			{"injected deduction", `{"verdicts":[{"constraintId":"a","verdict":"unsure","evidence":[],"deduction":0}]}`, false},
			{"null evidence", `{"verdicts":[{"constraintId":"a","verdict":"unsure","evidence":null}]}`, false},
			{"line zero", `{"verdicts":[{"constraintId":"a","verdict":"violation","evidence":[{"line":0,"text":"code"}]}]}`, false},
		} {
			instance, err := jsonschema.UnmarshalJSON(strings.NewReader(tc.raw))
			if err != nil {
				t.Fatal(err)
			}
			if valid := schema.Validate(instance) == nil; valid != tc.valid {
				t.Errorf("%s: provider schema accepted=%v; want %v", tc.name, valid, tc.valid)
			}
		}
		return &http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{"responseId":"fake-schema-response","modelVersion":"fake-version","usageMetadata":{"promptTokenCount":1,"candidatesTokenCount":1},"candidates":[{"finishReason":"STOP","content":{"parts":[{"text":"{\"verdicts\":[]}"}]}}]}`))}, nil
	})
	if _, err = model.Generate(context.Background(), Request{Statement: "Review", Source: "int item = 1;"}); err != nil {
		t.Fatal(err)
	}
}
