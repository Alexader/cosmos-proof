package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"xa.org/xablockchain/xchain-meta/relaychain"
)

// request of chaincode
type ProofInfo struct {
	Validators []byte           `json:"validators"`
	ChainID    string           `json:"chain_id"`
	Iccp       *relaychain.ICCP `json:"iccp"`
}

// response to send back to broker chaincode
type Response struct {
	Status bool   `json:"status"`
	Msg    string `json:"msg"`
}

func wrongRequest(g *gin.Context, err error) {
	g.JSON(http.StatusBadRequest, &Response{
		Status: false,
		Msg:    err.Error(),
	})
}

func verifyOK(g *gin.Context) {
	g.JSON(http.StatusOK, &Response{
		Status: true,
	})
}
