module github.com/Alexader/cosmos-proof

go 1.14

require (
	github.com/fatih/color v1.7.0
	github.com/gin-gonic/gin v1.6.3
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/iavl v0.14.0
	github.com/tendermint/tendermint v0.33.7
	github.com/urfave/cli/v2 v2.3.0
	xa.org/xablockchain/xchain-meta v0.0.0
)

replace (
	github.com/tendermint/tendermint v0.33.7 => github.com/FKmyCode/tendermint v0.33.7-0.20201207130517-07d76ab83293
	xa.org/xablockchain/xchain-meta => ../xchain-meta
)
