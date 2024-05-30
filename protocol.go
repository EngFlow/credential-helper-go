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
	"encoding/json"
	"time"
)

// GetCredentialsRequest represents the request for the `get` command of the Helper Protocol.
type GetCredentialsRequest struct {
	// The URI to get credentials for.
	URI string `json:"uri"`
}

// GetCredentialsResponse represents the response for the `get` command of the Helper Protocol.
type GetCredentialsResponse struct {
	// The headers containing credentials to add to all requests to the URI.
	Headers map[string][]string `json:"headers"`

	// The time the credentials expire and stop being valid for new requests,
	// formatted following [RFC 3339](https://www.rfc-editor.org/rfc/rfc3339.html).
	Expires *time.Time `json:"expires"`
}

func (resp GetCredentialsResponse) MarshalJSON() ([]byte, error) {
	// By default, time.Time is marshaled to a string with time.RFC3339Nano
	// instead of RFC3339, and Bazel rejects that format. We implement
	// json.Marshaler here to override that.
	v := struct {
		Headers map[string][]string `json:"headers"`
		Expires *string             `json:"expires"`
	}{
		Headers: resp.Headers,
	}
	if resp.Expires != nil {
		expires := resp.Expires.Format(time.RFC3339)
		v.Expires = &expires
	}
	return json.Marshal(v)
}
