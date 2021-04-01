package main

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/types"
	"io/ioutil"
	"testing"
	"xa.org/xablockchain/xchain-meta/relaychain"
)

func TestProof(t *testing.T) {
	iccpBytes, err := ioutil.ReadFile("./testdata/cosmos-iccp")
	require.Nil(t, err)
	validatorBytes, err := ioutil.ReadFile("./testdata/validator_set.txt")
	require.Nil(t, err)
	validators := &types.ValidatorSet{}
	require.Nil(t, json.Unmarshal(validatorBytes, validators))
	iccp := &relaychain.ICCP{}
	require.Nil(t, json.Unmarshal(iccpBytes, iccp))
	require.Nil(t, verifyICCP(validators, "appchain1", iccp))
}
