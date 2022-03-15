package uniswap_core

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

func GetTicks(client *graphql.Client, poolId string) ([]Tick, error) {
	numSkip := 0

	res := make([]Tick, 0)

	req := graphql.NewRequest(`
	query get_ticks($num_skip: Int, $pool_id: ID!) {
		ticks(skip: $num_skip, where: {pool: $pool_id}) {
			tickIdx
			liquidityGross
			liquidityNet
			feeGrowthOutside0X128
			feeGrowthOutside1X128
		}
	  }
    `)

	req.Var("pool_id", poolId)

	for {
		req.Var("num_skip", numSkip)

		var chunk struct {
			Ticks []Tick
		}

		if err := client.Run(context.Background(), req, &chunk); err != nil {
			return nil, err
		}

		n := len(chunk.Ticks)

		if n == 0 {
			break
		}

		numSkip += n
		res = append(res, chunk.Ticks...)
	}

	return res, nil
}

func GetPool(client *graphql.Client, poolId string) (*Pool, error) {
	req := graphql.NewRequest(`
		query get_pools($pool_id: ID!) {
			pools(where: {id: $pool_id}) {
			tick
			sqrtPrice
			liquidity
			feeTier
			feeGrowthGlobal0X128
			feeGrowthGlobal1X128
			token0 {
				symbol
				decimals
			}
			token1 {
				symbol
				decimals
			}
			swaps {
				id
			}
			ticks {
				id
				tickIdx
			}
			}
		}
	`)

	req.Var("pool_id", poolId)

	var res struct {
		Pools []Pool
	}

	if err := client.Run(context.Background(), req, &res); err != nil {
		return nil, err
	}

	if len(res.Pools) != 1 {
		return nil, fmt.Errorf("GetPool: incorrect count of results %d", len(res.Pools))
	}

	return &res.Pools[0], nil
}

func GetSwap(client *graphql.Client, swapId string) (*Swap, error) {
	req := graphql.NewRequest(`
		query get_swap($swap_id: ID!) {
			swaps(where: {id: $swap_id}) {
			sender
			recipient
			amount0
			amount1
			transaction {
			  id
			  blockNumber
			  gasUsed
			  gasPrice
			}
			timestamp
			sqrtPriceX96
			token0 {
			  id
			  symbol
			}
			token1 {
			  id
			  symbol
			}
		  }
		}`)

	req.Var("swap_id", swapId)

	var res struct {
		Swaps []Swap
	}

	if err := client.Run(context.Background(), req, &res); err != nil {
		return nil, err
	}

	if len(res.Swaps) != 1 {
		return nil, fmt.Errorf("GetSwap: incorrect count of results %d", len(res.Swaps))
	}

	return &res.Swaps[0], nil
}
