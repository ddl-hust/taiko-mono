package integration_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/ethereum/hive/hivesim"
	"github.com/stretchr/testify/assert"
)

func TestHiveHandler(t *testing.T) {
	baseDir := os.Getenv("HIVE_DIR")
	if baseDir == "" {
		t.SkipNow()
	}

	clientGroups := [][]string{
		{
			"taiko/anvil",
			"taiko/taiko-geth",
			"taiko/driver",
			"taiko/proposer",
			"taiko/prover",
		},
		{
			"taiko/taiko-geth",
			"taiko/driver",
		},
		{
			"taiko/taiko-geth",
			"taiko/driver",
		},
	}

	// Single cluster test.
	t.Run("taiko-genesis/l2-full-sync/clusters(1)", func(t *testing.T) {
		testDenebGenesis(t, "taiko-genesis/l2-full-sync", [][]string{clientGroups[0]})
	})
	t.Run("taiko-reorg/l2-full-sync/clusters(1)", func(t *testing.T) {
		testDenebReorg(t, "taiko-reorg/l2-full-sync", [][]string{clientGroups[0]})
	})

	// Multi clusters test.
	t.Run(fmt.Sprintf("taiko-genesis/l2-snap-sync/clusters(%d)", len(clientGroups)), func(t *testing.T) {
		testDenebGenesis(t, "taiko-genesis/l2-snap-sync", clientGroups)
	})
	t.Run(fmt.Sprintf("taiko-reorg/l2-snap-sync/clusters(%d)", len(clientGroups)), func(t *testing.T) {
		testDenebReorg(t, "taiko-reorg/l2-snap-sync", clientGroups)
	})
}

func testDenebGenesis(t *testing.T, simPattern string, clientGroups [][]string) {
	handler, err := hivesim.NewHiveFramework(&hivesim.HiveConfig{
		BuildOutput:     false,
		ContainerOutput: true,
		BaseDir:         os.Getenv("HIVE_DIR"),
		SimPattern:      "taiko",
		SimTestPattern:  simPattern,
		ClientGroups:    clientGroups,
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
}

func testDenebReorg(t *testing.T, simPattern string, clientGroups [][]string) {
	handler, err := hivesim.NewHiveFramework(&hivesim.HiveConfig{
		BuildOutput:     false,
		ContainerOutput: true,
		BaseDir:         os.Getenv("HIVE_DIR"),
		SimPattern:      "taiko",
		SimTestPattern:  simPattern,
		ClientGroups:    clientGroups,
	})
	assert.NoError(t, err)

	failedCount, err := handler.Run(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, failedCount)
}
