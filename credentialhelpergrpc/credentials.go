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

package credentialhelpergrpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/credentials"

	"github.com/EngFlow/credential-helper-go"
)

// NewPerRPCCredentials creates a [credentials.PerRPCCredentials]
// using the provided [bazelcredentialhelper.CredentialHelper] to
// fetch credentials.
func NewPerRPCCredentials(helper credentialhelper.CredentialHelper) credentials.PerRPCCredentials {
	return &perRPCCredentials{
		helper: helper,
	}
}

type perRPCCredentials struct {
	helper credentialhelper.CredentialHelper
}

func (c *perRPCCredentials) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	if len(uri) != 1 {
		return nil, fmt.Errorf("must provide exactly one uri, got %v", len(uri))
	}

	response, err := c.helper.GetCredentials(
		ctx,
		&credentialhelper.GetCredentialsRequest{
			URI: uri[0],
		})
	if err != nil {
		return nil, fmt.Errorf("error fetching credentials from helper: %w", err)
	}

	metadata := make(map[string]string, len(response.Headers))
	for name, values := range response.Headers {
		switch len(values) {
		case 0:
			// Helper returned a header without value. Ignore.
			continue

		case 1:
			metadata[name] = values[0]
			break

		default:
			return nil, fmt.Errorf("helper returned more than one value for header %q", name)
		}
	}

	return metadata, nil
}

func (c *perRPCCredentials) RequireTransportSecurity() bool {
	return false
}
