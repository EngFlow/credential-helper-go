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

// Package credentialhelper is a Go library for interacting with `Credential Helpers`.
//
// A `Credential Helper` is a tool for securely storing and
// retrieving credentials (e.g., for interacting with remote
// servers over `gRPC` or `HTTP(s)`).
//
// This library provides an easy and convenient way for
// implementing a credential helper as well as using a
// credential Helper from within an application to retrieve
// credentials.
//
// See https://github.com/bazelbuild/proposals/blob/main/designs/2022-06-07-bazel-credential-helpers.md#proposal
// for additional information about credential helpers.
package credentialhelper

import (
	"context"
	"errors"
)

// CredentialHelper provides an interface to implement a Credential Helper or
// communicate with one.
type CredentialHelper interface {
	// MustEmbedCredentialHelperBase is a private method so that this interface
	// can only be implemented by this package.
	mustEmbedCredentialHelperBase()

	// GetCredentials fetches credentials from the helper.
	GetCredentials(ctx context.Context, request *GetCredentialsRequest, extraParameters ...string) (*GetCredentialsResponse, error)
}

// CredentialHelperBase is the base for all implementations of
// [CredentialHelper]s.
type CredentialHelperBase struct{}

func (CredentialHelperBase) mustEmbedCredentialHelperBase() {}

func (CredentialHelperBase) GetCredentials(ctx context.Context, request *GetCredentialsRequest, extraParameters ...string) (*GetCredentialsResponse, error) {
	return nil, errors.New("credential Helper does not support command 'get'")
}

// Type assertions.
var (
	_ CredentialHelper = CredentialHelperBase{}
	_ CredentialHelper = &CredentialHelperBase{}
)
