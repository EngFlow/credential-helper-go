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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
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
	stdin, err := json.Marshal(request)
	if err != nil {
		return err
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd := exec.CommandContext(ctx, credentialHelperPath, extraArgs...)
	cmd.Stdin = bytes.NewBuffer(stdin)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("could not start credential helper: %w", err)
	}
	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitErr.Stderr = stderr.Bytes()
		}
		return fmt.Errorf("error running credential helper: %w", err)
	}

	stdoutBytes := stdout.Bytes()
	if err := json.Unmarshal(stdoutBytes, response); err != nil {
		return fmt.Errorf("could not read response from credential helper: %w", err)
	}

	return nil
}
