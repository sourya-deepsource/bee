// Copyright 2020 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pullstorage

import (
	"context"
	"errors"

	"github.com/ethersphere/bee/pkg/storage"
	"github.com/ethersphere/bee/pkg/swarm"
)

var _ Storer = (*ps)(nil)
var ErrDbClosed = errors.New("db closed")

type Storer interface {
	IntervalChunks(ctx context.Context, bin uint8, from, to uint64, limit int) (chunks []swarm.Address, topmost uint64, err error)
	Get(ctx context.Context, mode storage.ModeGet, addrs ...swarm.Address) ([]swarm.Chunk, error)
	Put(ctx context.Context, mode storage.ModePut, chs ...swarm.Chunk) error
	Set(ctx context.Context, mode storage.ModeSet, addrs ...swarm.Address) error
	Has(ctx context.Context, addr swarm.Address) (bool, error)
}

// ps wraps storage.Storer
type ps struct {
	storage.Storer
}

func New(storer storage.Storer) Storer {
	return &ps{storer}
}

func (s *ps) IntervalChunks(ctx context.Context, bin uint8, from, to uint64, limit int) (chs []swarm.Address, topmost uint64, err error) {
	// call iterator, iterate till upper bound reached or that limit was reached
	// return addresses, topmost is the topmost bin ID

	ch, dbClosed, stop := s.SubscribePull(ctx, bin, from, to)
	defer stop()

	var nomore bool

LOOP:
	for limit > 0 {
		select {
		case v, ok := <-ch:
			if !ok {
				nomore = true
				break LOOP
			}
			chs = append(chs, v.Address)
			if v.BinID > topmost {
				topmost = v.BinID
			}
			limit--
		case <-ctx.Done():
			return nil, 0, ctx.Err()
		}
	}

	// this is needed for circumvent context cancellation or db close
	select {
	case <-ctx.Done():
		return nil, 0, ctx.Err()
	case <-dbClosed:
		return nil, 0, ErrDbClosed
	default:
	}

	if nomore {
		topmost = to
	}

	return chs, topmost, nil
}

func (s *ps) Get(ctx context.Context, mode storage.ModeGet, addrs ...swarm.Address) ([]swarm.Chunk, error) {
	return s.Storer.GetMulti(ctx, mode, addrs...)
}
func (s *ps) Put(ctx context.Context, mode storage.ModePut, chs ...swarm.Chunk) error {
	_, err := s.Storer.Put(ctx, mode, chs...)
	return err
}

func (s *ps) Set(ctx context.Context, mode storage.ModeSet, addrs ...swarm.Address) error {
	return s.Storer.Set(ctx, mode, addrs...)
}

func (s *ps) Has(ctx context.Context, addr swarm.Address) (bool, error) {
	return s.Storer.Has(ctx, addr)
}
