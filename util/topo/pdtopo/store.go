// Copyright 2021 PingCAP, Inc. Licensed under Apache-2.0.

package pdtopo

import (
	"context"
	"strings"

	"github.com/pingcap/log"
	"go.uber.org/zap"

	"github.com/pingcap/tidb-dashboard/util/client/pdclient"
	"github.com/pingcap/tidb-dashboard/util/netutil"
	"github.com/pingcap/tidb-dashboard/util/topo"
)

// GetStoreInstances returns TiKV info and TiFlash info.
func GetStoreInstances(ctx context.Context, pdAPI *pdclient.APIClient) ([]topo.StoreInfo, []topo.StoreInfo, error) {
	stores, err := pdAPI.HLGetStores(ctx)
	if err != nil {
		return nil, nil, err
	}

	tiKVStores := make([]pdclient.GetStoresResponseStore, 0, len(stores))
	tiFlashStores := make([]pdclient.GetStoresResponseStore, 0, len(stores))
	for _, store := range stores {
		isTiFlash := false
		for _, label := range store.Labels {
			if label.Key == "engine" && label.Value == "tiflash" {
				isTiFlash = true
			}
		}
		if isTiFlash {
			tiFlashStores = append(tiFlashStores, store)
		} else {
			tiKVStores = append(tiKVStores, store)
		}
	}

	return buildStoreTopology(tiKVStores), buildStoreTopology(tiFlashStores), nil
}

func buildStoreTopology(stores []pdclient.GetStoresResponseStore) []topo.StoreInfo {
	nodes := make([]topo.StoreInfo, 0, len(stores))
	for _, v := range stores {
		hostname, port, err := netutil.ParseHostAndPortFromAddress(v.Address)
		if err != nil {
			log.Warn("Failed to parse store address", zap.Any("store", v))
			continue
		}
		_, statusPort, err := netutil.ParseHostAndPortFromAddress(v.StatusAddress)
		if err != nil {
			log.Warn("Failed to parse store status address", zap.Any("store", v))
			continue
		}
		// In current TiKV, it's version may not start with 'v',
		// so we may need to add a prefix 'v' for it.
		version := strings.Trim(v.Version, "\n ")
		if !strings.HasPrefix(version, "v") {
			version = "v" + version
		}
		node := topo.StoreInfo{
			Version:        version,
			IP:             hostname,
			Port:           port,
			GitHash:        v.GitHash,
			DeployPath:     v.DeployPath,
			Status:         parseStoreState(v.StateName),
			StatusPort:     statusPort,
			Labels:         map[string]string{},
			StartTimestamp: v.StartTimestamp,
		}
		for _, v := range v.Labels {
			node.Labels[v.Key] = v.Value
		}
		nodes = append(nodes, node)
	}

	return nodes
}

func parseStoreState(state string) topo.ComponentStatus {
	state = strings.Trim(strings.ToLower(state), "\n ")
	switch state {
	case "up":
		return topo.ComponentStatusUp
	case "tombstone":
		return topo.ComponentStatusTombstone
	case "offline":
		return topo.ComponentStatusOffline
	case "down":
		return topo.ComponentStatusDown
	case "disconnected":
		return topo.ComponentStatusUnreachable
	default:
		return topo.ComponentStatusUnreachable
	}
}
