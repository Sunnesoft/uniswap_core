package uniswap_core

import "math/big"

type Slot0 struct {
	TickSpacing          *big.Int
	Fee                  *big.Int
	Liquidity            *big.Int
	FeeGrowthGlobal0X128 *big.Int
	FeeGrowthGlobal1X128 *big.Int
	// the current price
	SqrtPriceX96 *big.Int
	// the current tick
	TickCurrent *big.Int
	// the most-recently updated index of the observations array
	ObservationIndex *big.Int
	// the current maximum number of observations that are being stored
	ObservationCardinality *big.Int
	// the next maximum number of observations to store, triggered in observations.write
	ObservationCardinalityNext *big.Int
	// the current protocol fee as a percentage of the swap fee taken on withdrawal
	// represented as an integer denominator (1/x)%
	FeeProtocol *big.Int
}

func NewSlot0() *Slot0 {
	s := new(Slot0)
	s.TickSpacing = big.NewInt(0)
	s.TickCurrent = big.NewInt(0)
	s.Fee = big.NewInt(0)
	s.Liquidity = big.NewInt(0)
	s.FeeGrowthGlobal0X128 = big.NewInt(0)
	s.FeeGrowthGlobal1X128 = big.NewInt(0)
	s.SqrtPriceX96 = big.NewInt(0)
	s.ObservationIndex = big.NewInt(0)
	s.ObservationCardinality = big.NewInt(0)
	s.ObservationCardinalityNext = big.NewInt(0)
	s.FeeProtocol = big.NewInt(0)
	return s
}

type PoolStateReader interface {
	CurrentState() *Slot0
}
