package uniswap_core

import (
	"math/big"
	"testing"
)

func TestNewTickStorage(t *testing.T) {
	count := 20
	data := make([]Tick, count)

	tickSpacing := big.NewInt(200)

	for i := 0; i < count; i++ {
		data[i] = Tick{
			TickIdx:               BigInt{Val: big.NewInt(int64(i-count/2) * tickSpacing.Int64())},
			LiquidityGross:        BigInt{Val: big.NewInt(int64(i))},
			LiquidityNet:          BigInt{Val: big.NewInt(int64(-i))},
			FeeGrowthOutside0X128: BigInt{Val: big.NewInt(0)},
			FeeGrowthOutside1X128: BigInt{Val: big.NewInt(0)},
		}
	}

	ts := NewTickStorage(data, tickSpacing)

	if ts.TickSpacing.Cmp(tickSpacing) != 0 {
		t.Errorf("TickStorage.TickSpacing = %d; want %d", ts.TickSpacing, tickSpacing)
	}

	if len(ts.Ticks) != count {
		t.Errorf("len(TickStorage.Ticks) = %d; want %d", len(ts.Ticks), count)
	}

	key0 := data[0].TickIdx.Val.Int64()
	key1 := data[1].TickIdx.Val.Int64()

	if ts.Ticks[key0].TickIdx.Val.Cmp(ts.Ticks[key1].TickIdx.Val) == 0 {
		t.Errorf("Incorrect assign map elements")
	}
}

func TestGetLiquidityNet(t *testing.T) {
	count := 1
	data := make([]Tick, count)

	tickSpacing := big.NewInt(200)
	ref := big.NewInt(int64(-100))

	data[0] = Tick{
		TickIdx:               BigInt{Val: big.NewInt(int64(count/2) * tickSpacing.Int64())},
		LiquidityGross:        BigInt{Val: big.NewInt(int64(0))},
		LiquidityNet:          BigInt{Val: ref},
		FeeGrowthOutside0X128: BigInt{Val: big.NewInt(0)},
		FeeGrowthOutside1X128: BigInt{Val: big.NewInt(0)},
	}

	ts := NewTickStorage(data, tickSpacing)
	tickIdx := data[0].TickIdx.Val
	liqNet := ts.GetLiquidityNet(tickIdx)

	if liqNet.Cmp(ref) != 0 {
		t.Errorf("GetLiquidityNet(%d) = %d; want %d", tickIdx, liqNet, ref)
	}
}

func TestNextInitializedTick(t *testing.T) {
	count := 20
	data := make([]Tick, count)

	tickSpacing := big.NewInt(200)

	for i := 0; i < count; i++ {
		data[i] = Tick{
			TickIdx:               BigInt{Val: big.NewInt(int64(i-count/2) * tickSpacing.Int64())},
			LiquidityGross:        BigInt{Val: big.NewInt(int64(i - 1))},
			LiquidityNet:          BigInt{Val: big.NewInt(int64(-i))},
			FeeGrowthOutside0X128: BigInt{Val: big.NewInt(0)},
			FeeGrowthOutside1X128: BigInt{Val: big.NewInt(0)},
		}
	}

	ts := NewTickStorage(data, tickSpacing)

	tickNext, initialized := ts.NextInitializedTick(data[0].TickIdx.Val, false)

	if tickNext.Cmp(data[2].TickIdx.Val) != 0 {
		t.Errorf("tickNext = %d; want %d", tickNext, data[2].TickIdx.Val)
	}

	if initialized != true {
		t.Errorf("initialized = %t; want %t", initialized, true)
	}
}
