package uniswap_core

import (
	"math/big"
)

type TickReader interface {
	NextInitializedTick(tick *big.Int, zeroForOne bool) (*big.Int, bool)
	GetLiquidityNet(tick *big.Int) *big.Int
}

type TickStorage struct {
	Ticks       map[int64]*Tick
	TickSpacing *big.Int
}

func (t TickStorage) IsInitialized(tickKey *big.Int) bool {
	key := tickKey.Int64()
	if tick, ok := t.Ticks[key]; ok {
		if tick.IsInitialized() && key%t.TickSpacing.Int64() == 0 {
			return true
		}
	}
	return false
}

func NewTickStorage(ticks []Tick, tickSpacing *big.Int) *TickStorage {
	t := new(TickStorage)
	t.TickSpacing = big.NewInt(tickSpacing.Int64())
	t.Ticks = make(map[int64]*Tick)

	for i, tick := range ticks {
		t.Ticks[tick.TickIdx.Val.Int64()] = &ticks[i]
	}
	return t
}

func (t TickStorage) NextInitializedTick(tick *big.Int, zeroForOne bool) (*big.Int, bool) {
	tickNext := big.NewInt(0).Set(tick)
	tickNext.Div(tickNext, t.TickSpacing)
	tickNext.Mul(tickNext, t.TickSpacing)

	if zeroForOne && tickNext.Cmp(tick) >= 0 {
		tickNext.Sub(tickNext, t.TickSpacing)
	} else if !zeroForOne && tickNext.Cmp(tick) <= 0 {
		tickNext.Add(tickNext, t.TickSpacing)
	}

	for {
		if t.IsInitialized(tickNext) {
			return tickNext, true
		}

		overBounds := false
		if tickNext.Cmp(MIN_TICK) < 0 {
			tickNext.Set(MIN_TICK)
			overBounds = true
		} else if tickNext.Cmp(MAX_TICK) > 0 {
			tickNext.Set(MAX_TICK)
			overBounds = true
		}

		if overBounds {
			return tickNext, t.IsInitialized(tickNext)
		}

		if zeroForOne {
			tickNext.Sub(tickNext, t.TickSpacing)
		} else {
			tickNext.Add(tickNext, t.TickSpacing)
		}
	}
}

func (t TickStorage) GetLiquidityNet(tick *big.Int) *big.Int {
	if data, ok := t.Ticks[tick.Int64()]; ok {
		return big.NewInt(0).Set(data.LiquidityNet.Val)
	}
	return nil
}
