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
package main

import (
	"context"
	"flag"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ThomasHabets/firewalls-at-the-source/pkg/cisco"
	"github.com/ThomasHabets/firewalls-at-the-source/pkg/linux"
	"github.com/ThomasHabets/firewalls-at-the-source/pkg/rules"
)

func run(ctx context.Context, blocker rules.Blocker) error {
	defer blocker.Clear(ctx)

	// Create a rule to stop NTP-based DDoS of 8.8.8.8.
	r := &rules.Rule{
		IPVersion:   rules.IPv4,
		Protocol:    rules.UDP,
		Source:      "192.0.2.0/24",
		Destination: "8.8.8.8/32",
		SourcePort:  123,
	}

	if err := blocker.Add(ctx, r); err != nil {
		return err
	}
	time.Sleep(10 * time.Second)
	return nil
}

func main() {
	flag.Parse()

	ctx := context.Background()
	blockers := []rules.Blocker{
		linux.NewBlocker("sourcefirewall"),
		cisco.NewBlocker("gw.example.net", "aclname"),
	}
	if err := run(ctx, rules.RuleList(blockers)); err != nil {
		log.Fatal(err)
	}
}
