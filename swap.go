package uniswap_core

import (
	"math/big"
)

type SwapCache struct {
	// the protocol fee for the input token
	feeProtocol *big.Int
	// liquidity at the beginning of the swap
	liquidityStart *big.Int
	// the current value of the tick accumulator, computed only if we cross an initialized tick
	tickCumulative *big.Int
	// the current value of seconds per liquidity accumulator, computed only if we cross an initialized tick
	secondsPerLiquidityCumulativeX128 *big.Int
	// whether we've computed and cached the above two accumulators
	computedLatestObservation bool
}

func NewSwapCache(zeroForOne bool, slot0 *Slot0) *SwapCache {
	feeProtocol := big.NewInt(0)
	if zeroForOne {
		feeProtocol.Mod(slot0.FeeProtocol, big.NewInt(16))
	} else {
		feeProtocol.Rsh(slot0.FeeProtocol, 4)
	}

	return &SwapCache{
		liquidityStart:                    big.NewInt(0).Set(slot0.Liquidity),
		feeProtocol:                       feeProtocol,
		secondsPerLiquidityCumulativeX128: big.NewInt(0),
		tickCumulative:                    big.NewInt(0),
		computedLatestObservation:         false}
}

// the top level state of the swap, the results of which are recorded in storage at the end
type SwapState struct {
	// the amount remaining to be swapped in/out of the input/output asset
	amountSpecifiedRemaining *big.Int
	// the amount already swapped out/in of the output/input asset
	amountCalculated *big.Int
	// current sqrt(price)
	sqrtPriceX96 *big.Int
	// the tick associated with the current price
	tick *big.Int
	// the global fee growth of the input token
	feeGrowthGlobalX128 *big.Int
	// amount of input token paid as protocol fee
	protocolFee *big.Int
	// the current liquidity in range
	liquidity *big.Int
}

func (state *SwapState) UpdateTickLiquidity(zeroForOne bool, step *StepComputations, ticker TickReader) {
	if state.sqrtPriceX96.Cmp(step.sqrtPriceNextX96) == 0 {
		if step.initialized {
			liquidityNet := ticker.GetLiquidityNet(step.tickNext)

			if zeroForOne {
				liquidityNet.Neg(liquidityNet)
			}

			state.liquidity.Set(AddLiquidityDelta(state.liquidity, liquidityNet))
		}

		if zeroForOne {
			state.tick.Sub(step.tickNext, ONE_UINT_256)
		} else {
			state.tick.Set(step.tickNext)
		}

	} else if state.sqrtPriceX96.Cmp(step.sqrtPriceStartX96) != 0 {
		state.tick.Set(GetTickAtSqrtRatio(state.sqrtPriceX96))
	}
}

func (state *SwapState) UpdateAmount(exactInput bool, step *StepComputations) {
	if exactInput {
		state.amountSpecifiedRemaining.Sub(state.amountSpecifiedRemaining, step.amountIn)
		state.amountSpecifiedRemaining.Sub(state.amountSpecifiedRemaining, step.feeAmount)
		state.amountCalculated.Sub(state.amountCalculated, step.amountOut)
	} else {
		state.amountSpecifiedRemaining.Add(state.amountSpecifiedRemaining, step.amountOut)
		state.amountCalculated.Add(state.amountCalculated, step.amountIn)
		state.amountCalculated.Add(state.amountCalculated, step.feeAmount)
	}
}

func NewSwapState(amountSpecified *big.Int, slot0 *Slot0, cache *SwapCache) *SwapState {
	return &SwapState{
		amountSpecifiedRemaining: big.NewInt(0).Set(amountSpecified),
		amountCalculated:         big.NewInt(0),
		sqrtPriceX96:             big.NewInt(0).Set(slot0.SqrtPriceX96),
		tick:                     big.NewInt(0).Set(slot0.TickCurrent),
		feeGrowthGlobalX128:      big.NewInt(0),
		protocolFee:              big.NewInt(0),
		liquidity:                big.NewInt(0).Set(cache.liquidityStart)}
}

type StepComputations struct {
	// the price at the beginning of the step
	sqrtPriceStartX96 *big.Int
	// the next tick to swap to from the current tick in the swap direction
	tickNext *big.Int
	// whether tickNext is initialized or not
	initialized bool
	// sqrt(price) for the next tick (1/0)
	sqrtPriceNextX96 *big.Int
	// how much is being swapped in in this step
	amountIn *big.Int
	// how much is being swapped out
	amountOut *big.Int
	// how much fee is being paid in
	feeAmount *big.Int
}

func NewStepComputations() *StepComputations {
	return &StepComputations{
		sqrtPriceStartX96: big.NewInt(0),
		tickNext:          big.NewInt(0),
		initialized:       false,
		sqrtPriceNextX96:  big.NewInt(0),
		amountIn:          big.NewInt(0),
		amountOut:         big.NewInt(0),
		feeAmount:         big.NewInt(0)}
}

func (step *StepComputations) ApplyTickLimits() {
	if step.tickNext.Cmp(MIN_TICK) < 0 {
		step.tickNext.Set(MIN_TICK)
	} else if step.tickNext.Cmp(MAX_TICK) > 0 {
		step.tickNext.Set(MAX_TICK)
	}
}

func (step *StepComputations) UpdateTickNext(zeroForOne bool, state *SwapState, ticker TickReader) {
	step.tickNext, step.initialized = ticker.NextInitializedTick(state.tick, zeroForOne)
}

func (step *StepComputations) UpdateSqrtPriceStartX96(state *SwapState) {
	step.sqrtPriceStartX96 = state.sqrtPriceX96
}

func (step *StepComputations) CalcSqrtPriceNextX96() {
	step.sqrtPriceNextX96.Set(GetSqrtRatioAtTick(step.tickNext))
}

func (step *StepComputations) GetSqrtRatioTargetX96(
	zeroForOne bool, sqrtPriceLimitX96 *big.Int) (sqrtRatioTargetX96 *big.Int) {
	sqrtRatioTargetX96 = step.sqrtPriceNextX96

	if zeroForOne {
		if sqrtRatioTargetX96.Cmp(sqrtPriceLimitX96) < 0 {
			sqrtRatioTargetX96 = sqrtPriceLimitX96
		}
	} else {
		if sqrtRatioTargetX96.Cmp(sqrtPriceLimitX96) > 0 {
			sqrtRatioTargetX96 = sqrtPriceLimitX96
		}
	}
	return
}

func (step *StepComputations) UpdateAmount(
	zeroForOne bool, sqrtPriceLimitX96 *big.Int, state *SwapState, slot0 *Slot0) (sqrtPriceX96 *big.Int) {
	sqrtRatioTargetX96 := step.GetSqrtRatioTargetX96(zeroForOne, sqrtPriceLimitX96)

	sqrtPriceX96, step.amountIn, step.amountOut, step.feeAmount = ComputeSwapStep(
		state.sqrtPriceX96, sqrtRatioTargetX96, state.liquidity, state.amountSpecifiedRemaining, slot0.Fee)

	return
}

func setDefaultSqrtPriceLimitX96(zeroForOne bool, sqrtPriceLimitX96 *big.Int) *big.Int {
	if sqrtPriceLimitX96.Cmp(ZERO_UINT_256) == 0 {
		if zeroForOne {
			sqrtPriceLimitX96.Set(MIN_SQRT_RATIO)
			sqrtPriceLimitX96.Add(sqrtPriceLimitX96, ONE_UINT_256)
		} else {
			sqrtPriceLimitX96.Set(MAX_SQRT_RATIO)
			sqrtPriceLimitX96.Sub(sqrtPriceLimitX96, ONE_UINT_256)
		}
	}

	return sqrtPriceLimitX96
}

// Swap token0 for token1, or token1 for token0
// zeroForOne	bool	The direction of the swap, true for token0 to token1, false for token1 to token0
// amountSpecified	big.Int	The amount of the swap, which implicitly configures the swap as exact input (positive), or exact output (negative)
// sqrtPriceLimitX96	big.Int	The Q64.96 sqrt price limit. If zero for one, the price cannot be less than this
// ticker TickReader	tick bitmap object
// slotReader PoolStateReader	Pool's state retriever object
// amount0	big.Int	The delta of the balance of token0 of the pool, exact when negative, minimum when positive
// amount1	big.Int	The delta of the balance of token1 of the pool, exact when negative, minimum when positive
func DoSwap(zeroForOne bool,
	amountSpecified *big.Int,
	sqrtPriceLimitX96 *big.Int,
	ticker TickReader,
	slotReader PoolStateReader) (amount0 *big.Int, amount1 *big.Int, feeTotal *big.Int) {

	sqrtPriceLimitX96 = setDefaultSqrtPriceLimitX96(zeroForOne, sqrtPriceLimitX96)

	feeTotal = big.NewInt(0)
	exactInput := amountSpecified.Cmp(ZERO_UINT_256) > 0

	slot0 := slotReader.CurrentState()
	cache := NewSwapCache(zeroForOne, slot0)
	state := NewSwapState(amountSpecified, slot0, cache)
	step := NewStepComputations()

	for state.amountSpecifiedRemaining.Cmp(ZERO_UINT_256) != 0 && state.sqrtPriceX96.Cmp(sqrtPriceLimitX96) != 0 {
		step.UpdateSqrtPriceStartX96(state)
		step.UpdateTickNext(zeroForOne, state, ticker)
		step.ApplyTickLimits()
		step.CalcSqrtPriceNextX96()
		state.sqrtPriceX96 = step.UpdateAmount(zeroForOne, sqrtPriceLimitX96, state, slot0)

		feeTotal.Add(feeTotal, step.feeAmount)

		state.UpdateAmount(exactInput, step)

		if cache.feeProtocol.Cmp(ZERO_UINT_256) > 0 {
			delta := big.NewInt(0)
			delta.Div(step.feeAmount, cache.feeProtocol)
			step.feeAmount.Sub(step.feeAmount, delta)
			state.protocolFee.Add(state.protocolFee, delta)
		}

		state.UpdateTickLiquidity(zeroForOne, step, ticker)
	}

	amount0 = big.NewInt(0)
	amount1 = big.NewInt(0)

	if zeroForOne == exactInput {
		amount0.Sub(amountSpecified, state.amountSpecifiedRemaining)
		amount1.Set(state.amountCalculated)
	} else {
		amount0.Set(state.amountCalculated)
		amount1.Sub(amountSpecified, state.amountSpecifiedRemaining)
	}

	return
}

func ComputeSwapStep(
	sqrtRatioCurrentX96 *big.Int,
	sqrtRatioTargetX96 *big.Int,
	liquidity *big.Int,
	amountRemaining *big.Int,
	feePips *big.Int) (sqrtRatioNextX96 *big.Int, amountIn *big.Int, amountOut *big.Int, feeAmount *big.Int) {

	zeroForOne := sqrtRatioCurrentX96.Cmp(sqrtRatioTargetX96) >= 0
	exactIn := amountRemaining.Cmp(ZERO_UINT_256) >= 0

	onex6 := big.NewInt(1e6)

	absAmountRemaining := big.NewInt(0)
	absAmountRemaining.Abs(amountRemaining)

	oneSubFeePips := big.NewInt(0)
	oneSubFeePips.Sub(onex6, feePips)

	if exactIn {
		amountRemainingLessFee := MulDiv(amountRemaining, oneSubFeePips, onex6)

		if zeroForOne {
			amountIn = GetAmount0DeltaRoundingUp(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, true)
		} else {
			amountIn = GetAmount1DeltaRoundingUp(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, true)
		}

		if amountRemainingLessFee.Cmp(amountIn) >= 0 {
			sqrtRatioNextX96 = big.NewInt(0)
			sqrtRatioNextX96.Set(sqrtRatioTargetX96)
		} else {
			sqrtRatioNextX96 = GetNextSqrtPriceFromInput(
				sqrtRatioCurrentX96,
				liquidity,
				amountRemainingLessFee,
				zeroForOne)
		}
	} else {
		if zeroForOne {
			amountOut = GetAmount1DeltaRoundingUp(sqrtRatioTargetX96, sqrtRatioCurrentX96, liquidity, false)
		} else {
			amountOut = GetAmount0DeltaRoundingUp(sqrtRatioCurrentX96, sqrtRatioTargetX96, liquidity, false)
		}

		if absAmountRemaining.Cmp(amountOut) >= 0 {
			sqrtRatioNextX96 = big.NewInt(0)
			sqrtRatioNextX96.Set(sqrtRatioTargetX96)
		} else {
			sqrtRatioNextX96 = GetNextSqrtPriceFromOutput(
				sqrtRatioCurrentX96,
				liquidity,
				absAmountRemaining,
				zeroForOne)
		}
	}

	max := sqrtRatioTargetX96.Cmp(sqrtRatioNextX96) == 0

	// get the input/output amounts
	if zeroForOne {
		if !(max && exactIn) {
			amountIn = GetAmount0DeltaRoundingUp(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, true)
		}

		if !(max && !exactIn) {
			amountOut = GetAmount1DeltaRoundingUp(sqrtRatioNextX96, sqrtRatioCurrentX96, liquidity, false)
		}
	} else {
		if !(max && exactIn) {
			amountIn = GetAmount1DeltaRoundingUp(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, true)
		}

		if !(max && !exactIn) {
			amountOut = GetAmount0DeltaRoundingUp(sqrtRatioCurrentX96, sqrtRatioNextX96, liquidity, false)
		}
	}

	if !exactIn && amountOut.Cmp(absAmountRemaining) > 0 {
		amountOut.Set(absAmountRemaining)
	}

	if exactIn && sqrtRatioNextX96.Cmp(sqrtRatioTargetX96) != 0 {
		feeAmount = big.NewInt(0)
		feeAmount.Sub(amountRemaining, amountIn)
	} else {
		feeAmount = MulDivRoundingUp(amountIn, feePips, oneSubFeePips)
	}

	return
}
