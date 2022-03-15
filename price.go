package uniswap_core

import (
	"fmt"
	"math/big"
)

// Calculates liquidity / sqrt(lower) - liquidity / sqrt(upper),
// i.e. liquidity * (sqrt(upper) - sqrt(lower)) / (sqrt(upper) * sqrt(lower))
// http://atiselsts.github.io/pdfs/uniswap-v3-liquidity-math.pdf
func GetAmount0DeltaRoundingUp(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
	roundUp bool) (amount0 *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	if sqrtRatioAX96.Cmp(ZERO_UINT_256) < 0 {
		panic(fmt.Sprintf("price: sqrtRatioAX96=%d must be non-negative", sqrtRatioAX96))
	}

	numerator1 := big.NewInt(0)
	numerator1.Lsh(liquidity, 96)

	numerator2 := big.NewInt(0)
	numerator2.Sub(sqrtRatioBX96, sqrtRatioAX96)

	if roundUp {
		amount0 = MulDivRoundingUp(numerator1, numerator2, sqrtRatioBX96)
		amount0 = DivRoundingUp(amount0, sqrtRatioAX96)
		return
	}

	amount0 = MulDiv(numerator1, numerator2, sqrtRatioBX96)
	amount0.Div(amount0, sqrtRatioAX96)
	return
}

// Calculates liquidity * (sqrt(upper) - sqrt(lower))
// http://atiselsts.github.io/pdfs/uniswap-v3-liquidity-math.pdf
func GetAmount1DeltaRoundingUp(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
	roundUp bool) (amount1 *big.Int) {
	if sqrtRatioAX96.Cmp(sqrtRatioBX96) > 0 {
		sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	}

	denumerator := big.NewInt(1)
	denumerator.Lsh(denumerator, 96)

	numerator2 := big.NewInt(0)
	numerator2.Sub(sqrtRatioBX96, sqrtRatioAX96)

	if roundUp {
		amount1 = MulDivRoundingUp(liquidity, numerator2, denumerator)
		return
	}

	amount1 = MulDiv(liquidity, numerator2, denumerator)
	return
}

func GetAmount0Delta(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int) (amount0 *big.Int) {
	if liquidity.Cmp(ZERO_UINT_256) < 0 {
		absLiq := big.NewInt(0)
		absLiq.Neg(liquidity)
		amount0 = GetAmount0DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, absLiq, false)
		amount0.Neg(amount0)
		return
	}

	amount0 = GetAmount0DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, liquidity, true)
	return
}

func GetAmount1Delta(
	sqrtRatioAX96 *big.Int,
	sqrtRatioBX96 *big.Int,
	liquidity *big.Int) (amount1 *big.Int) {
	if liquidity.Cmp(ZERO_UINT_256) < 0 {
		absLiq := big.NewInt(0)
		absLiq.Neg(liquidity)
		amount1 = GetAmount1DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, absLiq, false)
		amount1.Neg(amount1)
		return
	}

	amount1 = GetAmount1DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, liquidity, true)
	return
}

// The most precise formula for this is liquidity * sqrtPX96 / (liquidity +- amount * sqrtPX96),
// if this is impossible because of overflow, we calculate liquidity / (liquidity / sqrtPX96 +- amount)
// http://atiselsts.github.io/pdfs/uniswap-v3-liquidity-math.pdf
// add Whether to add, or remove, the amount of token0
func getNextSqrtPriceFromAmount0RoundingUp(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amount *big.Int,
	add bool) (result *big.Int) {
	if amount.Cmp(ZERO_UINT_256) == 0 {
		result = big.NewInt(0)
		result.Set(sqrtPX96)
		return
	}

	numerator1 := big.NewInt(0)
	numerator1.Lsh(liquidity, 96)

	product := big.NewInt(0)
	product.Mul(amount, sqrtPX96)

	denominator := big.NewInt(0)
	denominator.Div(product, amount)

	if add {
		if denominator.Cmp(sqrtPX96) == 0 {
			denominator.Add(numerator1, product)

			if denominator.Cmp(numerator1) >= 0 {
				result = MulDivRoundingUp(numerator1, sqrtPX96, denominator)
				return
			}
		}

		denominator.Div(numerator1, sqrtPX96)
		denominator.Add(denominator, amount)

		result = DivRoundingUp(numerator1, denominator)
		return
	}

	if denominator.Cmp(sqrtPX96) != 0 || numerator1.Cmp(product) <= 0 {
		panic(fmt.Sprintf("price: amount * sqrtPX96 / amount {%d} != sqrtPX96 or numerator1 {%d} <= 0", denominator, numerator1))
	}

	denominator.Sub(numerator1, product)
	result = MulDivRoundingUp(numerator1, sqrtPX96, denominator)
	return
}

// The formula we compute is within <1 wei of the lossless version: sqrtPX96 +- amount / liquidity
// http://atiselsts.github.io/pdfs/uniswap-v3-liquidity-math.pdf
// add Whether to add, or remove, the amount of token1
func getNextSqrtPriceFromAmount1RoundingDown(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amount *big.Int,
	add bool) (result *big.Int) {

	var quotient *big.Int
	result = big.NewInt(0)

	if add {
		if amount.Cmp(MAX_UINT_160) <= 0 {
			quotient = big.NewInt(0)
			quotient.Lsh(amount, 96)
			quotient.Div(quotient, liquidity)
		} else {
			quotient = big.NewInt(1)
			quotient.Lsh(quotient, 96)
			quotient = MulDiv(amount, quotient, liquidity)
		}

		result.Add(sqrtPX96, quotient)
		return
	}

	if amount.Cmp(MAX_UINT_160) <= 0 {
		quotient = big.NewInt(0)
		quotient.Lsh(amount, 96)
		quotient = DivRoundingUp(quotient, liquidity)
	} else {
		quotient = big.NewInt(1)
		quotient.Lsh(quotient, 96)
		quotient = MulDivRoundingUp(amount, quotient, liquidity)
	}

	if sqrtPX96.Cmp(quotient) <= 0 {
		panic(fmt.Sprintf("price: sqrtPX96 {%d} <= quotient {%d}", sqrtPX96, quotient))
	}

	result.Sub(sqrtPX96, quotient)
	return
}

func GetNextSqrtPriceFromInput(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amountIn *big.Int,
	zeroForOne bool) (sqrtQX96 *big.Int) {

	if sqrtPX96.Cmp(ZERO_UINT_256) <= 0 || liquidity.Cmp(ZERO_UINT_256) <= 0 {
		panic(fmt.Sprintf("price: sqrtPX96 {%d} or liquidity {%d} less than 0", sqrtPX96, liquidity))
	}

	if zeroForOne {
		sqrtQX96 = getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountIn, true)
		return
	}

	sqrtQX96 = getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountIn, true)
	return
}

func GetNextSqrtPriceFromOutput(
	sqrtPX96 *big.Int,
	liquidity *big.Int,
	amountOut *big.Int,
	zeroForOne bool) (sqrtQX96 *big.Int) {
	if sqrtPX96.Cmp(ZERO_UINT_256) <= 0 || liquidity.Cmp(ZERO_UINT_256) <= 0 {
		panic(fmt.Sprintf("price: sqrtPX96 {%d} or liquidity {%d} less than 0", sqrtPX96, liquidity))
	}

	if zeroForOne {
		sqrtQX96 = getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amountOut, false)
		return
	}

	sqrtQX96 = getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amountOut, false)
	return
}
