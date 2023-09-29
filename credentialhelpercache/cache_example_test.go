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

package credentialhelpercache_test

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	credentialhelper "github.com/EngFlow/credential-helper-go"
	"github.com/EngFlow/credential-helper-go/credentialhelpercache"
)

type countingCredentialHelper struct {
	credentialhelper.CredentialHelperBase

	counter atomic.Int32
}

// GetCredentials invokes the specified credential helper to fetch credentials.
func (c *countingCredentialHelper) GetCredentials(ctx context.Context, request *credentialhelper.GetCredentialsRequest, extraParameters ...string) (*credentialhelper.GetCredentialsResponse, error) {
	count := c.counter.Add(1)

	return &credentialhelper.GetCredentialsResponse{
		Headers: map[string][]string{
			"uri":   {request.URI},
			"count": {fmt.Sprintf("%v", count)},
		},
	}, nil
}

func ExampleHelperProcess() {
	helper, err := credentialhelpercache.New(
		&countingCredentialHelper{},
		credentialhelpercache.Options{
			TTL: 5 * time.Second,
		})
	if err != nil {
		log.Fatalf("Error creating cache: %v", err)
		return
	}
	defer helper.Close()

	response1, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/foo",
		})
	if err != nil {
		log.Fatalf("Error reading credentials from cache: %v", err)
		return
	}
	fmt.Printf("response 1: %v\n", response1.Headers)

	response2, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/foo",
		})
	if err != nil {
		log.Fatalf("Error reading credentials from cache: %v", err)
		return
	}
	fmt.Printf("response 2: %v\n", response2.Headers)

	response3, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/bar",
		})
	if err != nil {
		log.Fatalf("Error reading credentials from cache: %v", err)
		return
	}
	fmt.Printf("response 3: %v\n", response3.Headers)

	time.Sleep(10 * time.Second)

	response4, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/foo",
		})
	if err != nil {
		log.Fatalf("Error reading credentials from cache: %v", err)
		return
	}
	fmt.Printf("response 4: %v\n", response4.Headers)

	response5, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/bar",
		})
	if err != nil {
		log.Fatalf("Error reading credentials from cache: %v", err)
		return
	}
	fmt.Printf("response 5: %v\n", response5.Headers)

	response6, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/bar",
		})
	if err != nil {
		log.Fatalf("Error reading credentials from cache: %v", err)
		return
	}
	fmt.Printf("response 6: %v\n", response6.Headers)

	response7, err := helper.GetCredentials(
		context.Background(),
		&credentialhelper.GetCredentialsRequest{
			URI: "https://example.com/foo",
		})
	if err != nil {
		log.Fatalf("Error reading credentials from cache: %v", err)
		return
	}
	fmt.Printf("response 7: %v\n", response7.Headers)

	// Output:
	// response 1: map[count:[1] uri:[https://example.com/foo]]
	// response 2: map[count:[1] uri:[https://example.com/foo]]
	// response 3: map[count:[2] uri:[https://example.com/bar]]
	// response 4: map[count:[3] uri:[https://example.com/foo]]
	// response 5: map[count:[4] uri:[https://example.com/bar]]
	// response 6: map[count:[4] uri:[https://example.com/bar]]
	// response 7: map[count:[3] uri:[https://example.com/foo]]
}
