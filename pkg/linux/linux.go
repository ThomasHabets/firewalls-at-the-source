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
// This implementation is for linux hosts. Mostly useful for a cloud
// VM to prevent costly outgoing reflection traffic.
package linux

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/ThomasHabets/firewalls-at-the-source/pkg/rules"
)

var (
	iptables  = "echo"
	ip6tables = "echo"
)

type Blocker struct {
	chain string
}

func NewBlocker(chain string) *Blocker {
	return &Blocker{
		chain: chain,
	}
}

func (b *Blocker) Add(ctx context.Context, r *rules.Rule) error {
	args := []string{
		"-A", b.chain, "-d", r.Destination,
	}
	var bin string
	switch r.IPVersion {
	case rules.IPv4:
		bin = iptables
	case rules.IPv6:
		bin = ip6tables
	default:
		return fmt.Errorf("unknown IP version %q", r.IPVersion)
	}
	if len(r.Protocol) > 0 {
		args = append(args, "-p", string(r.Protocol))
	}
	if len(r.Source) > 0 {
		args = append(args, "-s", string(r.Source))
	}
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (b *Blocker) Clear(ctx context.Context) {
	// TODO: work even with expired context
	if err := exec.CommandContext(ctx, iptables, "-F", b.chain); err != nil {
		log.Errorf("Failed to flush IPv4 rules on chain %q: %v", b.chain, err)
	}
	if err := exec.CommandContext(ctx, ip6tables, "-F", b.chain); err != nil {
		log.Errorf("Failed to flush IPv6 rules on chain %q: %v", b.chain, err)
	}
}
