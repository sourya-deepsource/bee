package pullstorage

import (
	"github.com/ethersphere/bee/pkg/swarm"
)

type Storer interface {
	IntervalChunks(uint8, uint64, uint64) ([]swarm.Address, uint64, error)
}

type ps struct {
}

func (s *ps) IntervalChunks(bin uint8, from, to uint64) (chs []swarm.Address, topmost uint64, err error) {
	/*
		call iterator, iterate till upper bound reached or that limit was reached
	*/
	//return addresses, topmost is the topmost bin ID
	return chs, 0, nil
}
