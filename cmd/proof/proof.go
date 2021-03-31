package main

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	relaycodec "xa.org/xablockchain/xchain-meta/relaychain"
)

const (
	relaychainKeyStore = "cross"
)

// MultiStoreProof defines a collection of store proofs in a multi-store
type MultiStoreProof struct {
	StoreInfos []storeInfo
}

// ComputeRootHash returns the root hash for a given multi-store proof.
func (proof *MultiStoreProof) ComputeRootHash() []byte {
	ci := commitInfo{
		Version:    -1, // TODO: Not needed; improve code.
		StoreInfos: proof.StoreInfos,
	}
	return ci.Hash()
}

//-----------------------------------------------------------------------------

var _ merkle.ProofOperator = MultiStoreProofOp{}

// the multi-store proof operation constant value
const ProofOpMultiStore = "multistore"

// TODO: document
type MultiStoreProofOp struct {
	// Encoded in ProofOp.Key
	key []byte

	// To encode in ProofOp.Data.
	Proof *MultiStoreProof `json:"proof"`
}

func NewMultiStoreProofOp(key []byte, proof *MultiStoreProof) MultiStoreProofOp {
	return MultiStoreProofOp{
		key:   key,
		Proof: proof,
	}
}

// MultiStoreProofOpDecoder returns a multi-store merkle proof operator from a
// given proof operation.
func MultiStoreProofOpDecoder(pop merkle.ProofOp) (merkle.ProofOperator, error) {
	if pop.Type != ProofOpMultiStore {
		return nil, fmt.Errorf("unexpected ProofOp.Type; got %v, want %v", pop.Type, ProofOpMultiStore)
	}

	// XXX: a bit strange as we'll discard this, but it works
	var op MultiStoreProofOp

	err := relaycodec.ModuleCdc.UnmarshalBinaryLengthPrefixed(pop.Data, &op)
	if err != nil {
		return nil, fmt.Errorf("decoding ProofOp.Data into MultiStoreProofOp: %w", err)
	}

	return NewMultiStoreProofOp(pop.Key, op.Proof), nil
}

// ProofOp return a merkle proof operation from a given multi-store proof
// operation.
func (op MultiStoreProofOp) ProofOp() merkle.ProofOp {
	bz := relaycodec.ModuleCdc.MustMarshalBinaryLengthPrefixed(op)
	return merkle.ProofOp{
		Type: ProofOpMultiStore,
		Key:  op.key,
		Data: bz,
	}
}

// String implements the Stringer interface for a mult-store proof operation.
func (op MultiStoreProofOp) String() string {
	return fmt.Sprintf("MultiStoreProofOp{%v}", op.GetKey())
}

// GetKey returns the key for a multi-store proof operation.
func (op MultiStoreProofOp) GetKey() []byte {
	return op.key
}

// Run executes a multi-store proof operation for a given value. It returns
// the root hash if the value matches all the store's commitID's hash or an
// error otherwise.
func (op MultiStoreProofOp) Run(args [][]byte) ([][]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("value size is not 1")
	}

	value := args[0]
	root := op.Proof.ComputeRootHash()

	for _, si := range op.Proof.StoreInfos {
		if si.Name == string(op.key) {
			if bytes.Equal(value, si.Core.CommitID.Hash) {
				return [][]byte{root}, nil
			}

			return nil, fmt.Errorf("hash mismatch for substore %v: %X vs %X", si.Name, si.Core.CommitID.Hash, value)
		}
	}

	return nil, fmt.Errorf("key %v not found in multistore proof", op.key)
}

//-----------------------------------------------------------------------------

// XXX: This should be managed by the rootMultiStore which may want to register
// more proof ops?
func DefaultProofRuntime() (prt *merkle.ProofRuntime) {
	prt = merkle.NewProofRuntime()
	prt.RegisterOpDecoder(merkle.ProofOpSimpleValue, merkle.SimpleValueOpDecoder)
	prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.ValueOpDecoder)
	prt.RegisterOpDecoder(iavl.ProofOpIAVLAbsence, iavl.AbsenceOpDecoder)
	prt.RegisterOpDecoder(ProofOpMultiStore, MultiStoreProofOpDecoder)
	return
}

//----------------------------------------
// commitInfo

// NOTE: Keep commitInfo a simple immutable struct.
type commitInfo struct {

	// Version
	Version int64

	// Store info for
	StoreInfos []storeInfo
}

// Hash returns the simple merkle root hash of the stores sorted by name.
func (ci commitInfo) Hash() []byte {
	// TODO: cache to ci.hash []byte
	m := make(map[string][]byte, len(ci.StoreInfos))
	for _, storeInfo := range ci.StoreInfos {
		m[storeInfo.Name] = storeInfo.Hash()
	}

	return merkle.SimpleHashFromMap(m)
}

func (ci commitInfo) CommitID() CommitID {
	return CommitID{
		Version: ci.Version,
		Hash:    ci.Hash(),
	}
}

//----------------------------------------
// storeInfo

// storeInfo contains the name and core reference for an
// underlying store.  It is the leaf of the Stores top
// level simple merkle tree.
type storeInfo struct {
	Name string
	Core storeCore
}

type storeCore struct {
	// StoreType StoreType
	CommitID CommitID
	// ... maybe add more state
}

// Implements merkle.Hasher.
func (si storeInfo) Hash() []byte {
	// Doesn't write Name, since merkle.SimpleHashFromMap() will
	// include them via the keys.
	bz := si.Core.CommitID.Hash
	hasher := tmhash.New()

	_, err := hasher.Write(bz)
	if err != nil {
		// TODO: Handle with #870
		panic(err)
	}

	return hasher.Sum(nil)
}

//----------------------------------------
// CommitID

// CommitID contains the tree version number and its merkle root.
type CommitID struct {
	Version int64
	Hash    []byte
}

func (cid CommitID) IsZero() bool { //nolint
	return cid.Version == 0 && len(cid.Hash) == 0
}

func (cid CommitID) String() string {
	return fmt.Sprintf("CommitID{%v:%X}", cid.Hash, cid.Version)
}
