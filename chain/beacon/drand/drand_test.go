// stm: ignore
// Only tests external library behavior, therefore it should not be annotated
package drand

import (
	"bytes"
	"context"
	"os"
	"testing"

	dchain "github.com/drand/drand/chain"
	hclient "github.com/drand/drand/client/http"
	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/go-state-types/network"

	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/build/buildconstants"
)

func TestPrintGroupInfo(t *testing.T) {
	server := build.DrandConfigs[buildconstants.DrandTestnet].Servers[0]
	chainInfo := build.DrandConfigs[buildconstants.DrandTestnet].ChainInfoJSON

	drandChain, err := dchain.InfoFromJSON(bytes.NewReader([]byte(chainInfo)))
	assert.NoError(t, err)
	c, err := hclient.NewWithInfo(server, drandChain, nil)

	assert.NoError(t, err)
	cg := c.(interface {
		FetchChainInfo(ctx context.Context, groupHash []byte) (*dchain.Info, error)
	})
	chain, err := cg.FetchChainInfo(context.Background(), nil)
	assert.NoError(t, err)
	err = chain.ToJSON(os.Stdout, nil)
	assert.NoError(t, err)
}

func TestMaxBeaconRoundForEpoch(t *testing.T) {
	todayTs := uint64(1652222222)
	db, err := NewDrandBeacon(todayTs, build.BlockDelaySecs, nil, build.DrandConfigs[buildconstants.DrandTestnet])
	assert.NoError(t, err)
	assert.True(t, db.IsChained())
	mbr15 := db.MaxBeaconRoundForEpoch(network.Version15, 100)
	mbr16 := db.MaxBeaconRoundForEpoch(network.Version16, 100)
	assert.Equal(t, mbr15+1, mbr16)
}

func TestQuicknetIsChained(t *testing.T) {
	todayTs := uint64(1652222222)
	db, err := NewDrandBeacon(todayTs, build.BlockDelaySecs, nil, build.DrandConfigs[buildconstants.DrandQuicknet])
	assert.NoError(t, err)
	assert.False(t, db.IsChained())
}
