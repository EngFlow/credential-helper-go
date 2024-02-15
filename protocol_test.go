// Copyright 2023 EngFlow, Inc. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package credentialhelper_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/EngFlow/credential-helper-go"
)

func asPointer[T any](value T) *T {
	return &value
}

func TestParseGetCredentialsRequestFromString(t *testing.T) {
	var request1 credentialhelper.GetCredentialsRequest
	if err := json.Unmarshal(
		[]byte(`{"uri": "grpcs://example.com"}`),
		&request1); err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		credentialhelper.GetCredentialsRequest{
			URI: "grpcs://example.com",
		},
		request1); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}

	var request2 credentialhelper.GetCredentialsRequest
	if err := json.Unmarshal(
		[]byte(`{"uri": "grpcs://example.org"}`),
		&request2); err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		credentialhelper.GetCredentialsRequest{
			URI: "grpcs://example.org",
		},
		request2); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestParseGetCredentialsRequestFromStringWithExtraFields(t *testing.T) {
	var request credentialhelper.GetCredentialsRequest
	if err := json.Unmarshal(
		[]byte(`{"foo": 1, "uri": "grpcs://example.com", "bar": 2}`),
		&request); err != nil {
		t.Error(err)
	}
	if diff := cmp.Diff(
		credentialhelper.GetCredentialsRequest{
			URI: "grpcs://example.com",
		},
		request); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestParseGetCredentialsResponseFromString(t *testing.T) {
	var response1 credentialhelper.GetCredentialsResponse
	if err := json.Unmarshal(
		[]byte(`
			{
				"headers": {
					"header1": ["value1"],
					"header2": ["value1", "value2"],
					"header3": ["value1", "value2", "value3"]
				}
			}`),
		&response1); err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		credentialhelper.GetCredentialsResponse{
			Headers: map[string][]string{
				"header1": {"value1"},
				"header2": {"value1", "value2"},
				"header3": {"value1", "value2", "value3"},
			},
		},
		response1); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestParseGetCredentialsResponseFromStringWithExtraFields(t *testing.T) {
	var response1 credentialhelper.GetCredentialsResponse
	if err := json.Unmarshal(
		[]byte(`{"foo": 1, "headers": {"foo": ["1"], "bar": ["2"]}, "bar": 2}`),
		&response1); err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		credentialhelper.GetCredentialsResponse{
			Headers: map[string][]string{
				"foo": {"1"},
				"bar": {"2"},
			},
		},
		response1); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestParseGetCredentialsResponseFromStringWithExpires(t *testing.T) {
	var response1 credentialhelper.GetCredentialsResponse
	if err := json.Unmarshal(
		[]byte(`
			{
				"headers": {
					"header1": ["value1"]
				},
				"expires": "1970-09-28T23:46:29-12:00"
			}`),
		&response1); err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(
		credentialhelper.GetCredentialsResponse{
			Headers: map[string][]string{
				"header1": {"value1"},
			},
			Expires: asPointer(time.UnixMilli(23456789 * 1000)),
		},
		response1); diff != "" {
		t.Errorf("(-want +got):\n%s", diff)
	}
}

func TestParseGetCredentialsResponseFromStringWithInvalidExpires(t *testing.T) {
	var response1 credentialhelper.GetCredentialsResponse
	if err := json.Unmarshal(
		[]byte(`
			{
				"expires": "foo"
			}`),
		&response1); err == nil {
		t.Error("Expected error, got nil")
	}
}
