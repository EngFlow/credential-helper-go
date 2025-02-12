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
	"context"
	"fmt"
	"os"

	"github.com/EngFlow/credential-helper-go"
)

const (
	CredentialHelperEnvironmentVariable = "REAL_CREDENTIAL_HELPER"
)

func Example_client() {
	credentialHelperPath := os.Getenv(CredentialHelperEnvironmentVariable)
	if credentialHelperPath == "" {
		fmt.Fprintln(os.Stderr, CredentialHelperEnvironmentVariable+" not set")
		return
	}

	helper, err := credentialhelper.NewClient(credentialHelperPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating credential helper: %v", err)
		return
	}

	response, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "grpcs://example.com",
		})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error fetching credentials: %v", err)
		return
	}

	for name, values := range response.Headers {
		for _, value := range values {
			fmt.Fprintf(os.Stdout, "%s: %s", name, value)
		}
	}
}
