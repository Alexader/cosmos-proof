package main

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/lite"
	"github.com/tendermint/tendermint/types"
	"github.com/urfave/cli/v2"
	"xa.org/xablockchain/xchain-meta/relaychain"
)

func startCMD() *cli.Command {
	return &cli.Command{
		Name:   "start",
		Usage:  "Start a long-running cosmos proof service process",
		Action: start,
	}
}

func start(ctx *cli.Context) error {
	r := gin.Default()
	r.POST("/verify", verify)
	return r.Run()
}

func verify(g *gin.Context) {
	proofInfo := &ProofInfo{}
	if err := g.BindJSON(proofInfo); err != nil {
		wrongRequest(g, err)
		return
	}

	validators := &types.ValidatorSet{}
	cdc := relaychain.ModuleCdc
	if err := cdc.UnmarshalJSON(proofInfo.Validators, validators); err != nil {
		wrongRequest(g, err)
		return
	}

	if err := verifyICCP(validators, proofInfo.ChainID, proofInfo.Iccp); err != nil {
		wrongRequest(g, err)
		return
	}
	verifyOK(g)
}

func verifyICCP(validatorSet *types.ValidatorSet, chainID string, iccp *relaychain.ICCP) error {
	cp := &relaychain.CommitProof{}
	if err := cp.UnmarshalJsonProof(iccp.Proof); err != nil {
		return fmt.Errorf("unmarshal proof error: %w", err)
	}

	if iccp.Height+1 > cp.NextSignedHeader.Height {
		return fmt.Errorf("Iccp height is not match with signed header height, Iccp height: %d, next signed header height: %d",
			iccp.Height, cp.NextSignedHeader.Height)
	}

	cert := lite.NewBaseVerifier(chainID, cp.NextSignedHeader.Height, validatorSet)
	err := cert.Verify(cp.NextSignedHeader)
	if err != nil {
		return fmt.Errorf("verify signed header error: %w", err)
	}

	// verify merkle proof
	prt := DefaultProofRuntime()
	kp := merkle.KeyPath{}
	kp = kp.AppendKey([]byte(relaychainKeyStore), merkle.KeyEncodingURL)
	kp = kp.AppendKey([]byte(iccp.ID()), merkle.KeyEncodingURL)

	iccpBytes, err := iccp.MarshalJSONWithoutProof()
	if err != nil {
		return fmt.Errorf("msrshal Iccp without proof: %w", err)
	}

	err = prt.VerifyValue(cp.Proof, cp.NextSignedHeader.AppHash, kp.String(), iccpBytes)
	if err != nil {
		return fmt.Errorf("verify proof error: %w", err)
	}
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Proof"
	app.Usage = "Start a cosmos proof verify service"
	app.Compiled = time.Now()

	// global flags
	app.Commands = cli.Commands{
		startCMD(),
	}
	err := app.Run(os.Args)
	if err != nil {
		color.Red(err.Error())
		os.Exit(-1)
	}
}
