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

package credentialhelper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_InvalidPath(t *testing.T) {
	c := &client{
		credentialHelperPath: "testdata/does-definitely-not-exist",
	}

	response, err := c.GetCredentials(
		context.Background(),
		&GetCredentialsRequest{
			URI: "http://example.com/foo/bar",
		})
	assert.ErrorContains(t, err, "could not start credential helper")
	assert.Nil(t, response)
}
