package uniswap_core

import (
	"math/big"
	"testing"
)

func TestGetAmount0DeltaRoundingUp(t *testing.T) {

	sqrtRatioAX96 := big.NewInt(100)
	sqrtRatioAX96.Lsh(sqrtRatioAX96, 96)
	sqrtRatioBX96 := big.NewInt(300)
	sqrtRatioBX96.Lsh(sqrtRatioBX96, 96)
	liquidity := big.NewInt(10000)
	roundUp := false

	amount0 := GetAmount0DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, liquidity, roundUp)
	refRes := big.NewInt(66)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	}

	roundUp = true

	amount0 = GetAmount0DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, liquidity, roundUp)
	refRes = big.NewInt(67)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount0DeltaRoundingUp() = %d; want %d", amount0, refRes)
	}
}

func TestGetAmount1DeltaRoundingUp(t *testing.T) {

	sqrtRatioAX96 := big.NewInt(100)
	sqrtRatioAX96.Lsh(sqrtRatioAX96, 96)
	sqrtRatioBX96 := big.NewInt(300)
	sqrtRatioBX96.Lsh(sqrtRatioBX96, 96)
	liquidity := big.NewInt(10000)
	roundUp := false

	amount0 := GetAmount1DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, liquidity, roundUp)
	refRes := big.NewInt(2000000)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount1DeltaRoundingUp() = %d; want %d", amount0, refRes)
	}

	roundUp = true
	sqrtRatioBX96.Add(sqrtRatioBX96, big.NewInt(300))

	amount0 = GetAmount1DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, liquidity, roundUp)
	refRes = big.NewInt(2000001)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount1DeltaRoundingUp() = %d; want %d", amount0, refRes)
	}

	sqrtRatioAX96, sqrtRatioBX96 = sqrtRatioBX96, sqrtRatioAX96
	amount0 = GetAmount1DeltaRoundingUp(sqrtRatioAX96, sqrtRatioBX96, liquidity, roundUp)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount1DeltaRoundingUp() = %d; want %d", amount0, refRes)
	}
}

func TestGetAmount0Delta(t *testing.T) {
	sqrtRatioAX96 := big.NewInt(100)
	sqrtRatioAX96.Lsh(sqrtRatioAX96, 96)
	sqrtRatioBX96 := big.NewInt(300)
	sqrtRatioBX96.Lsh(sqrtRatioBX96, 96)
	liquidity := big.NewInt(10000)

	amount0 := GetAmount0Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	refRes := big.NewInt(67)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount0Delta() = %d; want %d", amount0, refRes)
	}

	sqrtRatioBX96.Add(sqrtRatioBX96, big.NewInt(300))
	liquidity = big.NewInt(-10000)

	amount0 = GetAmount0Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	refRes = big.NewInt(-66)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount0Delta() = %d; want %d", amount0, refRes)
	}
}

func TestGetAmount1Delta(t *testing.T) {
	sqrtRatioAX96 := big.NewInt(500)
	sqrtRatioAX96.Lsh(sqrtRatioAX96, 96)
	sqrtRatioBX96 := big.NewInt(400)
	sqrtRatioBX96.Lsh(sqrtRatioBX96, 96)
	liquidity := big.NewInt(10000)

	amount0 := GetAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	refRes := big.NewInt(1000000)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount1Delta() = %d; want %d", amount0, refRes)
	}

	sqrtRatioBX96.Add(sqrtRatioBX96, big.NewInt(300))
	liquidity = big.NewInt(-10000)

	amount0 = GetAmount1Delta(sqrtRatioAX96, sqrtRatioBX96, liquidity)
	refRes = big.NewInt(-999999)

	if amount0.Cmp(refRes) != 0 {
		t.Errorf("GetAmount1Delta() = %d; want %d", amount0, refRes)
	}
}

func TestGetNextSqrtPriceFromAmount0RoundingUp(t *testing.T) {
	sqrtPX96 := big.NewInt(300)
	sqrtPX96.Lsh(sqrtPX96, 96)
	liquidity := big.NewInt(12345)
	amount := big.NewInt(100)
	add := true

	result := getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amount, add)
	result.Rsh(result, 96)

	refResFloat := 12345.0 * 300.0 / (12345.0 + 100.0*300.0)
	refRes := big.NewInt(int64(refResFloat))

	if result.Cmp(refRes) != 0 {
		t.Errorf("getNextSqrtPriceFromAmount0RoundingUp() = %d; want %d", result, refRes)
	}

	liquidity = big.NewInt(300*100 + 300)
	add = false

	result = getNextSqrtPriceFromAmount0RoundingUp(sqrtPX96, liquidity, amount, add)
	result.Rsh(result, 96)

	refResFloat = 30300.0 * 300.0 / (30300.0 - 100.0*300.0)
	refRes = big.NewInt(int64(refResFloat))

	if result.Cmp(refRes) != 0 {
		t.Errorf("getNextSqrtPriceFromAmount0RoundingUp() = %d; want %d", result, refRes)
	}
}

func TestGetNextSqrtPriceFromAmount1RoundingDown(t *testing.T) {
	sqrtPX96 := big.NewInt(300)
	sqrtPX96.Lsh(sqrtPX96, 96)
	liquidity := big.NewInt(12345)
	amount := big.NewInt(100)
	add := true

	result := getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amount, add)
	result.Rsh(result, 96)

	refResFloat := 300.0 + (100.0 / 12345.0)
	refRes := big.NewInt(int64(refResFloat))

	if result.Cmp(refRes) != 0 {
		t.Errorf("getNextSqrtPriceFromAmount1RoundingDown() = %d; want %d", result, refRes)
	}

	add = false

	result = getNextSqrtPriceFromAmount1RoundingDown(sqrtPX96, liquidity, amount, add)
	result.Rsh(result, 96)

	refResFloat = 300.0 - (100.0 / 12345.0)
	refRes = big.NewInt(int64(refResFloat))

	if result.Cmp(refRes) != 0 {
		t.Errorf("getNextSqrtPriceFromAmount1RoundingDown() = %d; want %d", result, refRes)
	}
}
