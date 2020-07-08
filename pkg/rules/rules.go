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
package rules

import (
	"context"

	log "github.com/sirupsen/logrus"
)

type IPVersion int
type Protocol string

const (
	IPv4 IPVersion = 4
	IPv6 IPVersion = 6
	TCP  Protocol  = "tcp"
	UDP  Protocol  = "udp"
)

// TODO: turn into protobuf
type Rule struct {
	IPVersion   IPVersion
	Destination string
	Source      string
	Protocol    Protocol

	PortOptional    bool // If true and source can't block on port numbers, then drop whole address.
	SourcePort      int  // TODO: port range
	DestinationPort int
}

type Blocker interface {
	Add(context.Context, *Rule) error
	Clear(context.Context)
}

type RuleListBlocker struct {
	list []Blocker
}

func (b *RuleListBlocker) Add(ctx context.Context, r *Rule) error {
	var ret error
	for _, e := range b.list {
		if err := e.Add(ctx, r); err != nil {
			log.Error(err)
			// TODO: concat errors
			ret = err
		}
	}
	return ret
}

func (b *RuleListBlocker) Clear(ctx context.Context) {
	for _, e := range b.list {
		e.Clear(ctx)
	}
}

func RuleList(bs []Blocker) *RuleListBlocker {
	return &RuleListBlocker{
		list: bs,
	}
}
