// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package command defines the command line interface for rvcs
package command

import (
	"context"
	"flag"
	"fmt"
	"net/url"

	"github.com/google/recursive-version-control-system/config"
	"github.com/google/recursive-version-control-system/snapshot"
	"github.com/google/recursive-version-control-system/storage"
)

const removeMirrorUsage = `Usage: %s remove-mirror [<IDENTITY>] <MIRROR_URL>

Where <IDENTITY> is the optional identity to mirror (omit to apply to all identities), and <MIRROR_URL> is the URL of the mirror.
`

func removeMirrorCommand(ctx context.Context, s *storage.LocalFiles, cmd string, args []string) (int, error) {
	if len(args) < 1 {
		fmt.Fprintf(flag.CommandLine.Output(), removeMirrorUsage, cmd)
		return 1, nil
	}
	var id *snapshot.Identity
	var err error
	if len(args) > 1 {
		id, err = snapshot.ParseIdentity(args[0])
		if err != nil {
			return 1, fmt.Errorf("failure parsing the identity %q: %v", args[0], err)
		}
		args = args[1:]
	}
	mirrorURL, err := url.Parse(args[0])
	if err != nil {
		return 1, fmt.Errorf("failure parsing the mirror URL %q: %v", args[0], err)
	}
	settings, err := config.Read()
	if err != nil {
		return 1, fmt.Errorf("failure reading the existing config settings: %v", err)
	}
	if id == nil {
		settings = settings.WithoutAdditionalMirror(mirrorURL)
	} else {
		settings = settings.WithoutMirrorForIdentity(id.String(), mirrorURL)
	}
	if err := settings.Write(); err != nil {
		return 1, fmt.Errorf("failure writing the updated config settings: %v", err)
	}
	return 0, nil
}
