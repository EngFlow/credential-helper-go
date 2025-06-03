// Copyright 2025 EngFlow, Inc. All rights reserved.
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

//go:build darwin || linux

package credentialhelper_test

import (
	"context"
	"os"
	"testing"

	"github.com/EngFlow/credential-helper-go"

	"github.com/stretchr/testify/assert"
)

func runCredentialHelper(path string) (*credentialhelper.GetCredentialsResponse, error) {
	client, err := credentialhelper.NewClient(path)
	if err != nil {
		return nil, err
	}

	return client.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/foo",
		})
}

func TestClient_WithHeaders(t *testing.T) {
	response, err := runCredentialHelper("testdata/with-headers.sh")
	assert.NoError(t, err)
	assert.Equal(
		t,
		&credentialhelper.GetCredentialsResponse{
			Headers: map[string][]string{
				"foo": {"bar", "baz"},
				"bar": {"hello", "world"},
			},
		},
		response)
}

func TestClient_DoesNotReadStdin(t *testing.T) {
	response, err := runCredentialHelper("testdata/does-not-read-stdin.sh")
	assert.NoError(t, err)
	assert.Equal(
		t,
		&credentialhelper.GetCredentialsResponse{},
		response)
}

func TestClient_ReadStdin(t *testing.T) {
	response, err := runCredentialHelper("testdata/read-stdin.sh")
	assert.NoError(t, err)
	assert.Equal(
		t,
		&credentialhelper.GetCredentialsResponse{},
		response)
}

func TestClient_NoResponse(t *testing.T) {
	response, err := runCredentialHelper("testdata/no-response.sh")
	assert.ErrorContains(t, err, "could not read response from credential helper")
	assert.Nil(t, response)
}

func TestClient_NotExecutable(t *testing.T) {
	// Unfortunately, the source file ends up being executable when running on
	// RE. So we copy it first.
	data, err := os.ReadFile("testdata/not-executable.sh")
	if err != nil {
		t.Fatal(err)
		return
	}

	if err = os.WriteFile("testdata/not-executable-copy.sh", data, 0644); err != nil {
		t.Fatal(err)
		return
	}

	response, err := runCredentialHelper("testdata/not-executable-copy.sh")
	assert.ErrorContains(t, err, "could not lookup credential helper")
	assert.ErrorContains(t, err, "testdata/not-executable-copy.sh")
	assert.ErrorContains(t, err, "permission denied")
	assert.Nil(t, response)
}

func TestClient_WrongExitCode(t *testing.T) {
	response, err := runCredentialHelper("testdata/wrong-exit-code.sh")
	assert.ErrorContains(t, err, "error running credential helper")
	assert.ErrorContains(t, err, "exit status 1")
	assert.Nil(t, response)
}

func TestClient_WrongExitCodeWithStderr(t *testing.T) {
	response, err := runCredentialHelper("testdata/wrong-exit-code-with-stderr.sh")
	assert.ErrorContains(t, err, "error running credential helper")
	assert.ErrorContains(t, err, "exit status 1")
	assert.Nil(t, response)
}

func TestClient_InvalidResponse(t *testing.T) {
	response, err := runCredentialHelper("testdata/invalid-response.sh")
	assert.ErrorContains(t, err, "could not read response from credential helper")
	assert.Nil(t, response)
}
