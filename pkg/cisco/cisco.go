// Copyright 2020 Google Inc.
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
package cisco

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/ThomasHabets/firewalls-at-the-source/pkg/rules"
)

var ()

type Blocker struct {
	host    string
	aclname string
}

func NewBlocker(host, aclname string) *Blocker {
	return &Blocker{
		host:    host,
		aclname: aclname,
	}
}

func (b *Blocker) Add(ctx context.Context, r *rules.Rule) error {
	// TODO
	log.Infof("Adding rule %+v to cisco router %q aclname %q", r, b.host, b.aclname)
	return nil
}

func (b *Blocker) Clear(ctx context.Context) {
	// TODO
	log.Infof("Flushing all rules from cisco router %q aclname %q", b.host, b.aclname)
}
