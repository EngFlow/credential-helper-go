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
	"errors"

	"github.com/EngFlow/credential-helper-go"
)

type exampleCredentialHelper struct {
	credentialhelper.CredentialHelperBase
}

func (e *exampleCredentialHelper) GetCredentials(ctx context.Context, request *credentialhelper.GetCredentialsRequest, extraParameters ...string) (*credentialhelper.GetCredentialsResponse, error) {
	return nil, errors.New("example does not provide credentials")
}

func Example_helperProcess() {
	credentialhelper.StartCredentialHelper(&exampleCredentialHelper{})

	panic("UNREACHED")
}
