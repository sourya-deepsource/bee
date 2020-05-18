// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pullsync

import (
	"context"
	"io"

	"github.com/ethersphere/bee/pkg/logging"
	"github.com/ethersphere/bee/pkg/p2p"
	"github.com/ethersphere/bee/pkg/p2p/protobuf"
	"github.com/ethersphere/bee/pkg/pullsync/pb"
	"github.com/ethersphere/bee/pkg/pullsync/pullstorage"
	"github.com/ethersphere/bee/pkg/swarm"
)

const (
	protocolName    = "pullsync"
	protocolVersion = "1.0.0"
	streamName      = "pullsync"
)

type Syncer struct {
	streamer p2p.Streamer
	logger   logging.Logger

	storage pullstorage.Storer

	io.Closer
}

type Options struct {
	Streamer p2p.Streamer
	Logger   logging.Logger
}

func New(o Options) *Syncer {
	return &Syncer{
		streamer: o.Streamer,
		logger:   o.Logger,
	}
}

func (s *Syncer) Protocol() p2p.ProtocolSpec {
	return p2p.ProtocolSpec{
		Name:    protocolName,
		Version: protocolVersion,
		StreamSpecs: []p2p.StreamSpec{
			{
				Name:    streamName,
				Handler: s.handler,
			},
		},
	}
}

func (s *Syncer) handler(ctx context.Context, p p2p.Peer, stream p2p.Stream) error {
	w, r := protobuf.NewWriterAndReader(stream)
	defer stream.Close()

	// we are receiving a GetRange message
	var rn pb.GetRange
	if err := r.ReadMsg(&rn); err != nil {
		return nil, err
	}

	// we offer hashes in return
	offer := s.makeOffer(rn)
	if err := w.WriteMsg(&offer); err != nil {
		return nil, err
	}

	// we get a Want message as a response
	var want pb.Want
	if err := r.Read(&want); err != nil {
		return nil, err
	}

	for _, v := range s.processWant(offer, want) {
		var deliver pb.Delivery

		if err := w.WriteMsg(&deliver); err != nil {
			return nil, err
		}
	}
}

func (s *Syncer) makeOffer(rn pb.GetRange) (pb.Offer, error) {
	chs, top, err := s.storage.IntervalChunks(rn.Bin, rn.From, rn.To)
	///.... and so on, return an offer eventually

	return nil, nil
}

func (s *Syner) processWant(o pb.Offer, w pb.Want) []swarm.Chunk {

}

func (s *Syncer) Close() error {

}
