// Copyright (c) 2018 The VeChainThor developers

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package poa

import (
	"bytes"
	"errors"
	"sort"

	"github.com/vechain/thor/block"
	"github.com/vechain/thor/thor"
)

// SchedulerV2 to schedule the time when a proposer to produce a block.
// V2 is for post VIP-193 stage.
type SchedulerV2 struct {
	proposer        Proposer
	parentBlockTime uint64
	shuffled        []thor.Address
	activates       []thor.Address
}

var _ Scheduler = (*SchedulerV2)(nil)

// NewSchedulerV2 create a SchedulerV2 object.
// `addr` is the proposer to be scheduled.
// If `addr` is not listed in `proposers` or not active, an error returned.
func NewSchedulerV2(
	addr thor.Address,
	proposers []Proposer,
	parent *block.Block,
	seed []byte) (*SchedulerV2, error) {
	var (
		backers   = make(map[thor.Address]bool)
		activates []thor.Address
	)

	// handling parent block's backers in post VIP-193 stage, activate them when they backs
	bs, _, err := parent.Committee()
	if err != nil {
		return nil, err
	}
	for _, b := range bs {
		backers[b] = true
	}

	if canPropose := func() bool {
		for _, p := range proposers {
			if p.Address == addr && (p.Active || backers[p.Address]) {
				return true
			}
		}
		return false
	}(); !canPropose {
		return nil, errors.New("unauthorized or inactive block proposer")
	}

	var (
		proposer Proposer
		list     []struct {
			addr thor.Address
			hash thor.Bytes32
		}
	)

	for _, p := range proposers {
		if !p.Active && backers[p.Address] {
			activates = append(activates, p.Address)
		}
		if p.Active || backers[p.Address] {
			if p.Address == addr {
				proposer = p
			}
			list = append(list, struct {
				addr thor.Address
				hash thor.Bytes32
			}{
				p.Address,
				thor.Blake2b(seed, parent.Header().ID().Bytes()[:4], p.Address.Bytes()),
			})
		}
	}

	sort.Slice(list, func(i, j int) bool {
		return bytes.Compare(list[i].hash.Bytes(), list[j].hash.Bytes()) < 0
	})

	shuffled := make([]thor.Address, 0, len(list))
	for _, t := range list {
		shuffled = append(shuffled, t.addr)
	}

	return &SchedulerV2{
		proposer,
		parent.Header().Timestamp(),
		shuffled,
		activates,
	}, nil
}

// Schedule to determine time of the proposer to produce a block, according to `nowTime`.
// `newBlockTime` is promised to be >= nowTime and > parentBlockTime
func (s *SchedulerV2) Schedule(nowTime uint64) (newBlockTime uint64) {
	const T = thor.BlockInterval

	newBlockTime = s.parentBlockTime + T
	if nowTime > newBlockTime {
		// ensure T aligned, and >= nowTime
		newBlockTime += (nowTime - newBlockTime + T - 1) / T * T
	}

	offset := (newBlockTime-s.parentBlockTime)/T - 1
	for i := uint64(0); i < uint64(len(s.shuffled)); i++ {
		index := (i + offset) % uint64(len(s.shuffled))
		if s.shuffled[index] == s.proposer.Address {
			return newBlockTime + i*T
		}
	}

	// should never happen
	panic("something wrong with proposers list")
}

// IsTheTime returns if the newBlockTime is correct for the proposer.
func (s *SchedulerV2) IsTheTime(newBlockTime uint64) bool {
	return s.IsScheduled(newBlockTime, s.proposer.Address)
}

// IsScheduled returns if the schedule(proposer, blockTime) is correct.
func (s *SchedulerV2) IsScheduled(blockTime uint64, proposer thor.Address) bool {
	if s.parentBlockTime >= blockTime {
		// invalid block time
		return false
	}

	T := thor.BlockInterval
	if (blockTime-s.parentBlockTime)%T != 0 {
		// invalid block time
		return false
	}

	index := (blockTime - s.parentBlockTime - T) / T % uint64(len(s.shuffled))
	return s.shuffled[index] == proposer
}

// Updates returns proposers whose status are changed, and the score when new block time is assumed to be newBlockTime.
// In scheduler v2, Updates only deactivate proposers.
func (s *SchedulerV2) Updates(newBlockTime uint64) (updates []Proposer, score uint64) {
	T := thor.BlockInterval

	for _, a := range s.activates {
		updates = append(updates, Proposer{a, true})
	}

	for i := uint64(0); i < uint64(len(s.shuffled)); i++ {
		if s.parentBlockTime+i*T+T >= newBlockTime {
			break
		}
		if s.shuffled[i] != s.proposer.Address {
			updates = append(updates, Proposer{s.shuffled[i], false})
		}
	}

	score = uint64(len(s.shuffled) + len(s.activates) - len(updates))
	return
}
