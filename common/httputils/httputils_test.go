// Copyright 2022-2023 The Cloudprober Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httputils

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsHandled(t *testing.T) {
	srvMux := http.NewServeMux()

	srvMux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {})
	srvMux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {})
	srvMux.Handle("/", http.RedirectHandler("/status", http.StatusFound))

	tests := map[string]bool{
		"/":            true,
		"/probestatus": false,
		"/status":      true,
		"/config":      true,
		"/config2":     false,
	}

	for url, wantResult := range tests {
		assert.Equal(t, wantResult, IsHandled(srvMux, url))
	}
}

func TestNewRequest(t *testing.T) {
	tests := []struct {
		name        string
		data        []string
		wantReqBody string
		wantCT      string
		wantErr     bool
	}{
		{
			name:        "json_body",
			data:        []string{`{"clientId":"testID", "clientSecret":"testSecret"}`},
			wantReqBody: `{"clientId":"testID", "clientSecret":"testSecret"}`,
			wantCT:      "application/json",
		},
		{
			name:        "query_body",
			data:        []string{"clientId=testID", "clientSecret=testSecret"},
			wantReqBody: "clientId=testID&clientSecret=testSecret",
			wantCT:      "application/x-www-form-urlencoded",
		},
		{
			name:        "no_data",
			data:        []string{},
			wantReqBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := NewRequest("method", "test-url", NewRequestBody(tt.data...))
			if (err != nil) != tt.wantErr {
				t.Errorf("newHTTPTokenSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(tt.data) == 0 {
				assert.Equal(t, nil, req.Body, "request Body not nil")
				assert.Nil(t, req.GetBody, "request GetBody not nil")
				return
			}

			assert.NotEqual(t, nil, req.Body, "request Body nil")
			assert.NotNil(t, req.GetBody, "request GetBody nil")

			got, _ := io.ReadAll(req.Body)
			assert.Equal(t, tt.wantReqBody, string(got))
			assert.Equal(t, tt.wantCT, req.Header.Get("Content-Type"), "Content-Type Header")

			// We want an empty read next time.
			got, _ = io.ReadAll(req.Body)
			assert.Equal(t, "", string(got))
		})
	}
}

func TestRequestBody(t *testing.T) {
	data := make([]string, 100)
	for i := 0; i < len(data); i++ {
		data[i] = fmt.Sprintf("var-%d=value-%d", i, i)
	}
	tests := []struct {
		name           string
		data           []string
		wantData       string
		wantLen        int64
		wantCT         string
		wantNilReader  bool
		nilRequestBody bool
	}{
		{
			name:     "large_data",
			data:     data,
			wantLen:  1579,
			wantData: strings.Join(data, "&"),
			wantCT:   "",
		},
		{
			name:     "small_data",
			data:     data[:10],
			wantData: strings.Join(data[:10], "&"),
			wantLen:  139,
			wantCT:   "application/x-www-form-urlencoded",
		},
		{
			name:     "single_data_string",
			data:     []string{"clientId=testID&clientSecret=testSecret"},
			wantData: "clientId=testID&clientSecret=testSecret",
			wantLen:  39,
			wantCT:   "application/x-www-form-urlencoded",
		},
		{
			name:     "json_data",
			data:     []string{`{"clientId":"testID", "clientSecret":"testSecret"}`},
			wantData: `{"clientId":"testID", "clientSecret":"testSecret"}`,
			wantLen:  50,
			wantCT:   "application/json",
		},
		{
			name:          "no_data",
			data:          []string{},
			wantNilReader: true,
		},
		{
			name:           "nil_request_body", // Verify nil RequestBody works
			nilRequestBody: true,
			wantNilReader:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var body *RequestBody
			if !tt.nilRequestBody {
				body = NewRequestBody(tt.data...)
			}

			assert.Equal(t, tt.wantLen, body.Len(), "length mismatch")
			assert.Equal(t, tt.wantCT, body.ContentType(), "content-type mismatch")

			// Verify reader is good
			reader := body.Reader()
			if tt.wantNilReader {
				assert.Equal(t, nil, reader, "reader not nil")
				return
			}

			gotData, _ := io.ReadAll(reader)
			assert.Equal(t, tt.wantData, string(gotData), "body data mismatch 1st read")
			reader = body.Reader()
			gotData, _ = io.ReadAll(reader)
			assert.Equal(t, tt.wantData, string(gotData), "body data mismatch 2nd read")
		})
	}
}
