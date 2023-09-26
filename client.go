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

package credentialhelper

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
)

// NewClient returns a new [CredentialHelper] invoking the provided path following
// the protocol for Bazel Credential Helpers.
func NewClient(credentialHelperPath string) (CredentialHelper, error) {
	path, err := exec.LookPath(credentialHelperPath)
	if err != nil {
		return nil, fmt.Errorf("could not lookup credential helper %q: %w", credentialHelperPath, err)
	}

	c := &client{
		credentialHelperPath: path,
	}
	return c, nil
}

type client struct {
	CredentialHelperBase

	credentialHelperPath string
}

// GetCredentials invokes the specified credential helper to fetch credentials.
func (c *client) GetCredentials(ctx context.Context, request *GetCredentialsRequest, extraParameters ...string) (*GetCredentialsResponse, error) {
	var response GetCredentialsResponse
	if err := invoke(ctx, c.credentialHelperPath, request, &response, append([]string{"get"}, extraParameters...)...); err != nil {
		return nil, err
	}
	return &response, nil
}

func invoke[RequestT any, ResponseT any](ctx context.Context, credentialHelperPath string, request *RequestT, response *ResponseT, extraArgs ...string) error {
	cmd := exec.CommandContext(ctx, credentialHelperPath, extraArgs...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("could not open stdin of credential helper: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		defer stdin.Close()
		return fmt.Errorf("could not open stdout of credential helper: %w", err)
	}
	defer stdout.Close()

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start credential helper: %w", err)
	}

	if err := writeRequest(stdin, request); err != nil {
		return fmt.Errorf("could not write request to credential helper: %w", err)
	}

	if err := json.NewDecoder(stdout).Decode(response); err != nil {
		return fmt.Errorf("could not read response from credential helper: %w", err)
	}

	return cmd.Wait()
}

func writeRequest(stdin io.WriteCloser, request any) error {
	defer stdin.Close()

	if err := json.NewEncoder(stdin).Encode(request); err != nil {
		// This can happen if the helper prints a static set of credentials without reading from
		// stdin (e.g., with a simple shell script running `echo "{...}"`). This is fine to
		// ignore.
	}

	return nil
}
