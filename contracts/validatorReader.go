package contracts

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	slotValidatorMapping = map[string]uint64{
		"withdrawsState":         0,
		"validatorsState":        1,
		"voters":                 2,
		"candidates":             3,
		"candidateCount":         4,
		"minCandidateCap":        5,
		"minVoterCap":            6,
		"maxValidatorNumber":     7,
		"candidateWithdrawDelay": 8,
		"voterWithdrawDelay":     9,
	}
)

func GetCandidates(statedb *state.StateDB) []common.Address {
	slot := slotValidatorMapping["candidates"]
	slotHash := common.BigToHash(new(big.Int).SetUint64(slot))
	arrLength := statedb.GetState(common.HexToAddress(common.MasternodeVotingSMC), slotHash)
	keys := []common.Hash{}
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		key := state.GetLocDynamicArrAtElement(slotHash, i, 1)
		keys = append(keys, key)
	}
	rets := []common.Address{}
	for _, key := range keys {
		ret := statedb.GetState(common.HexToAddress(common.MasternodeVotingSMC), key)
		rets = append(rets, common.HexToAddress(ret.Hex()))
	}
	return rets
}

func GetCandidateOwner(statedb *state.StateDB, candidate common.Address) common.Address {
	slot := slotValidatorMapping["validatorsState"]
	// validatorsState[_candidate].owner;
	locValidatorsState := state.GetLocMappingAtKey(candidate.Hash(), slot)
	locCandidateOwner := locValidatorsState.Add(locValidatorsState, new(big.Int).SetUint64(uint64(0)))
	ret := statedb.GetState(common.HexToAddress(common.MasternodeVotingSMC), common.BigToHash(locCandidateOwner))
	return common.HexToAddress(ret.Hex())
}

func GetCandidateCap(statedb *state.StateDB, candidate common.Address) *big.Int {
	slot := slotValidatorMapping["validatorsState"]
	locValidatorsState := state.GetLocMappingAtKey(candidate.Hash(), slot)
	locCandidateCap := locValidatorsState.Add(locValidatorsState, new(big.Int).SetUint64(uint64(1)))
	ret := statedb.GetState(common.HexToAddress(common.MasternodeVotingSMC), common.BigToHash(locCandidateCap))
	return ret.Big()
}

func GetVoters(statedb *state.StateDB, candidate common.Address) []common.Address {
	//mapping(address => address[]) voters;
	slot := slotValidatorMapping["voters"]
	locVoters := state.GetLocMappingAtKey(candidate.Hash(), slot)
	arrLength := statedb.GetState(common.HexToAddress(common.MasternodeVotingSMC), common.BigToHash(locVoters))
	keys := []common.Hash{}
	for i := uint64(0); i < arrLength.Big().Uint64(); i++ {
		key := state.GetLocDynamicArrAtElement(common.BigToHash(locVoters), i, 1)
		keys = append(keys, key)
	}
	rets := []common.Address{}
	for _, key := range keys {
		ret := statedb.GetState(common.HexToAddress(common.MasternodeVotingSMC), key)
		rets = append(rets, common.HexToAddress(ret.Hex()))
	}

	return rets
}

func GetVoterCap(statedb *state.StateDB, candidate, voter common.Address) *big.Int {
	slot := slotValidatorMapping["validatorsState"]
	locValidatorsState := state.GetLocMappingAtKey(candidate.Hash(), slot)
	locCandidateVoters := locValidatorsState.Add(locValidatorsState, new(big.Int).SetUint64(uint64(2)))
	retByte := crypto.Keccak256(voter.Hash().Bytes(), common.BigToHash(locCandidateVoters).Bytes())
	ret := statedb.GetState(common.HexToAddress(common.MasternodeVotingSMC), common.BytesToHash(retByte))
	return ret.Big()
}
