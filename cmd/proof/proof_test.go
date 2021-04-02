package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/types"
	"xa.org/xablockchain/xchain-meta/relaychain"
)

func TestProof(t *testing.T) {
	// test gin context
	iccpBytes, err := ioutil.ReadFile("./testdata/cosmos-Iccp")
	require.Nil(t, err)
	validatorBytes, err := ioutil.ReadFile("./testdata/validator_set.txt")
	require.Nil(t, err)

	validators := &types.ValidatorSet{}
	cdc := relaychain.ModuleCdc
	require.Nil(t, cdc.UnmarshalJSON(validatorBytes, validators))

	iccp := &relaychain.ICCP{}
	require.Nil(t, json.Unmarshal(iccpBytes, iccp))

	proofInfo := &ProofInfo{
		Validators: validatorBytes,
		ChainID:    "appchain1",
		Iccp:       iccp,
	}
	proofInfoBytes, err := json.Marshal(proofInfo)
	require.Nil(t, err)

	buf := bytes.NewBuffer(proofInfoBytes)
	gin.SetMode(gin.TestMode)
	g, _ := gin.CreateTestContext(httptest.NewRecorder())

	r, err := http.NewRequest("POST", "http://localhost/verify", buf)
	require.Nil(t, err)
	g.Request = r
	verify(g)
	require.Equal(t, g.Writer.Status(), 200)
}
