// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libp2p

import (
	"context"

	"github.com/ethersphere/bee/pkg/p2p"
)

type ProtocolService interface {
	Peer(ctx context.Context, peer p2p.Peer) error
	Start()
}
