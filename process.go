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
	"os"
	"path"
)

// StartCredentialHelper is a util for turning the current process as credential helper.
//
// This function never returns.
func StartCredentialHelper(helper CredentialHelper) {
	os.Exit(startCredentialHelper(os.Stdin, os.Stdout, os.Stderr, os.Args, helper))
}

func startCredentialHelper(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args []string, helper CredentialHelper) int {
	if len(args) < 2 {
		printHelp(stderr, args[0])
		return 1
	}

	switch args[1] {
	case "get":
		return runGetCommand(stdin, stdout, stderr, args, helper)

	default:
		fmt.Fprintln(stderr, "Unknown command '"+args[1]+"'")
		fmt.Fprintln(stderr, "")
		printHelp(stderr, args[0])
		return 1
	}
}

func printHelp(w io.Writer, procName string) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  "+path.Base(procName)+" <command>")
}

func runGetCommand(stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, args []string, helper CredentialHelper) int {
	if len(args) != 2 {
		printHelp(stderr, args[0])
		return 1
	}

	var request GetCredentialsRequest
	if err := json.NewDecoder(stdin).Decode(&request); err != nil {
		fmt.Fprintln(stderr, "Invalid request for command 'get':")
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	response, err := helper.GetCredentials(context.Background(), &request)
	if err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	if err := json.NewEncoder(stdout).Encode(response); err != nil {
		fmt.Fprintln(stderr, err.Error())
		return 1
	}

	return 0
}
