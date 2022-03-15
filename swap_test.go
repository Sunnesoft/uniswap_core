package uniswap_core

import (
	"fmt"
	"github.com/machinebox/graphql"
	"math/big"
	"testing"
)

func TestComputeSwapStep(t *testing.T) {
	exSize := 8
	ex := []int64{
		0, 60, 10000, 300, 10000, 31, 29, 1,
		-60, 0, 10000, 300, 10000, 30, 30, 1,
		100, 200, 100000000, 3000, 3000, 2991, 2961, 9}

	for i := 0; i < len(ex)/exSize; i++ {
		tickCur := big.NewInt(ex[exSize*i+0])
		tickNext := big.NewInt(ex[exSize*i+1])
		sqrtPriceCurrentX96 := GetSqrtRatioAtTick(tickCur)
		sqrtPriceNextX96 := GetSqrtRatioAtTick(tickNext)
		liquidity := big.NewInt(ex[exSize*i+2])
		amountRemaining := big.NewInt(ex[exSize*i+3])
		feePips := big.NewInt(ex[exSize*i+4])

		_, amountIn, amountOut, feeAmount := ComputeSwapStep(
			sqrtPriceCurrentX96, sqrtPriceNextX96, liquidity, amountRemaining, feePips)

		// fmt.Println(sqrtRatioNextX96, amountIn, amountOut, feeAmount)

		if amountIn.Int64() != ex[exSize*i+5] ||
			amountOut.Int64() != ex[exSize*i+6] ||
			feeAmount.Int64() != ex[exSize*i+7] {
			t.Errorf("ComputeSwapStep(%v) = %d, %d, %d; want %v",
				ex[exSize*i:exSize*i+5], amountIn, amountOut, feeAmount, ex[exSize*i+5:exSize*i+8])
		}
	}
}

func TestDoSwap(t *testing.T) {
	zeroForOne := true

	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3")
	poolId := "0x8ad599c3a0ff1de082011efddc58f1908eb6e6d8"
	ticks, err := GetTicks(client, poolId)

	if err != nil {
		t.Errorf("TestDoSwap(...): %s", err)
	}

	pool, err := GetPool(client, poolId)

	if err != nil {
		t.Errorf("TestDoSwap(...): %s", err)
	}

	amountSpecified := big.NewInt(-100000000000000)
	sqrtPriceLimitX96 := big.NewInt(0)

	ticker := NewTickStorage(ticks, pool.FeerTierToTickSpacing())

	amount0, amount1, fee := DoSwap(zeroForOne, amountSpecified, sqrtPriceLimitX96, ticker, pool)
	fmt.Println(amount0, amount1, fee)
}
