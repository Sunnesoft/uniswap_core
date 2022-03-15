package uniswap_core

import (
	"github.com/machinebox/graphql"
	"testing"
)

func TestGetTicks(t *testing.T) {
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3")
	poolId := "0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8"
	ticks, err := GetTicks(client, poolId)

	if err != nil {
		t.Errorf("GetTicks(...): %s", err)
	}

	if len(ticks) == 0 {
		t.Errorf("GetTicks(...): len(tiks) < 1")
	}
}

func TestGetPool(t *testing.T) {
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3")
	poolId := "0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8"
	pool, err := GetPool(client, poolId)

	if err != nil {
		t.Errorf("GetPool(...): %s", err)
	}

	_ = pool
}

func TestGetSwap(t *testing.T) {
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3")
	swapId := "0x0000328cbe1ae7cba600d5c0cbc1387da983bb55ce49aa3519163334bf465430#115463"
	swap, err := GetSwap(client, swapId)

	if err != nil {
		t.Errorf("GetSwap(...): %s", err)
	}

	_ = swap
}
