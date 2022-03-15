package uniswap_core

import (
	"math/big"
)

func MulDiv(
	a *big.Int,
	b *big.Int,
	denominator *big.Int) (result *big.Int) {
	result = big.NewInt(0)
	result.Mul(a, b)
	result.Div(result, denominator)
	return
}

// Deprecated: Obtain [a * b / denominator]
// Source: https://xn--2-umb.com/21/muldiv/index.html
// Origin: https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/FullMath.sol
func MulDivChn(
	a *big.Int,
	b *big.Int,
	denominator *big.Int) (result *big.Int) {
	result = big.NewInt(0)
	prod0 := big.NewInt(0) // Least significant 256 bits of the product
	prod1 := big.NewInt(0) // Most significant 256 bits of the product

	prod0.Mul(a, b)

	mm := big.NewInt(0)
	mm.Mod(prod0, MAX_UINT_256)

	prod1.Sub(mm, prod0)

	if mm.Cmp(prod0) < 0 {
		prod1.Sub(prod1, ONE_UINT_256)
	}

	if prod1.Cmp(ZERO_UINT_256) == 0 {
		result.Div(prod0, denominator)
		return
	}

	remainder := big.NewInt(0)
	remainder.Mod(prod0, denominator)

	if remainder.Cmp(prod0) > 0 {
		prod1.Sub(prod1, ONE_UINT_256)
		prod0.Sub(prod0, remainder)
	}

	twos := big.NewInt(0)
	twos.Neg(denominator)
	twos.And(twos, denominator)

	den := big.NewInt(0)
	den.Div(denominator, twos)

	prod0.Div(prod0, twos)
	mm.Sub(ZERO_UINT_256, twos)
	mm.Div(mm, twos)
	twos.Add(mm, ONE_UINT_256)

	mm.Mul(prod1, twos)
	prod0.Or(prod0, mm)

	inv := big.NewInt(3)
	inv.Mul(inv, denominator)
	inv.Mul(inv, inv)

	two := big.NewInt(2)

	for i := 0; i < 6; i++ {
		mm.Mul(denominator, inv)
		mm.Sub(two, mm)

		inv.Mul(inv, mm)
	}

	result.Mul(prod0, inv)
	return
}

// Origin: https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/FullMath.sol
func MulDivRoundingUp(
	a *big.Int,
	b *big.Int,
	denominator *big.Int) (result *big.Int) {
	result = MulDiv(a, b, denominator)

	prod0 := big.NewInt(0)

	prod0.Mul(a, b)
	prod0.Mod(prod0, denominator)

	if prod0.Cmp(ZERO_UINT_256) > 0 {
		result.Add(result, ONE_UINT_256)
	}
	return
}

// Origin: https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/UnsafeMath.sol
func DivRoundingUp(
	a *big.Int,
	b *big.Int) (result *big.Int) {

	result = big.NewInt(0)
	prod1 := big.NewInt(0)

	result.DivMod(a, b, prod1)

	if prod1.Cmp(ZERO_UINT_256) > 0 {
		result.Add(result, ONE_UINT_256)
	}

	return
}

func AddLiquidityDelta(x *big.Int, y *big.Int) (z *big.Int) {
	z = big.NewInt(0)
	z.Add(x, y)
	return
}
