package uniswap_core

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

var A [19]*big.Int
var B [19]*big.Int
var C [7]*big.Int

var ZERO_UINT_256 *big.Int
var ONE_UINT_256 *big.Int
var INIT0 *big.Int
var INIT1 *big.Int
var MAX_UINT_256 *big.Int
var MAX_UINT_160 *big.Int
var SQRT_10001 *big.Int
var LOWER_ERR_BOUND *big.Int
var UPPER_ERR_BOUND *big.Int

//The minimum tick that may be passed to #getSqrtRatioAtTick computed from log base 1.0001 of 2**-128
var MIN_TICK *big.Int

//The maximum tick that may be passed to #getSqrtRatioAtTick computed from log base 1.0001 of 2**128
var MAX_TICK *big.Int

//The minimum value that can be returned from #getSqrtRatioAtTick. Equivalent to getSqrtRatioAtTick(MIN_TICK)
var MIN_SQRT_RATIO *big.Int

//The maximum value that can be returned from #getSqrtRatioAtTick. Equivalent to getSqrtRatioAtTick(MAX_TICK)
var MAX_SQRT_RATIO *big.Int

func init() {
	A[0], _ = hexutil.DecodeBig("0xfff97272373d413259a46990580e213a")
	A[1], _ = hexutil.DecodeBig("0xfff2e50f5f656932ef12357cf3c7fdcc")
	A[2], _ = hexutil.DecodeBig("0xffe5caca7e10e4e61c3624eaa0941cd0")
	A[3], _ = hexutil.DecodeBig("0xffcb9843d60f6159c9db58835c926644")
	A[4], _ = hexutil.DecodeBig("0xff973b41fa98c081472e6896dfb254c0")
	A[5], _ = hexutil.DecodeBig("0xff2ea16466c96a3843ec78b326b52861")
	A[6], _ = hexutil.DecodeBig("0xfe5dee046a99a2a811c461f1969c3053")
	A[7], _ = hexutil.DecodeBig("0xfcbe86c7900a88aedcffc83b479aa3a4")
	A[8], _ = hexutil.DecodeBig("0xf987a7253ac413176f2b074cf7815e54")
	A[9], _ = hexutil.DecodeBig("0xf3392b0822b70005940c7a398e4b70f3")
	A[10], _ = hexutil.DecodeBig("0xe7159475a2c29b7443b29c7fa6e889d9")
	A[11], _ = hexutil.DecodeBig("0xd097f3bdfd2022b8845ad8f792aa5825")
	A[12], _ = hexutil.DecodeBig("0xa9f746462d870fdf8a65dc1f90e061e5")
	A[13], _ = hexutil.DecodeBig("0x70d869a156d2a1b890bb3df62baf32f7")
	A[14], _ = hexutil.DecodeBig("0x31be135f97d08fd981231505542fcfa6")
	A[15], _ = hexutil.DecodeBig("0x9aa508b5b7a84e1c677de54f3e99bc9")
	A[16], _ = hexutil.DecodeBig("0x5d6af8dedb81196699c329225ee604")
	A[17], _ = hexutil.DecodeBig("0x2216e584f5fa1ea926041bedfe98")
	A[18], _ = hexutil.DecodeBig("0x48a170391f7dc42444e8fa2")

	B[0], _ = hexutil.DecodeBig("0x2")
	B[1], _ = hexutil.DecodeBig("0x4")
	B[2], _ = hexutil.DecodeBig("0x8")
	B[3], _ = hexutil.DecodeBig("0x10")
	B[4], _ = hexutil.DecodeBig("0x20")
	B[5], _ = hexutil.DecodeBig("0x40")
	B[6], _ = hexutil.DecodeBig("0x80")
	B[7], _ = hexutil.DecodeBig("0x100")
	B[8], _ = hexutil.DecodeBig("0x200")
	B[9], _ = hexutil.DecodeBig("0x400")
	B[10], _ = hexutil.DecodeBig("0x800")
	B[11], _ = hexutil.DecodeBig("0x1000")
	B[12], _ = hexutil.DecodeBig("0x2000")
	B[13], _ = hexutil.DecodeBig("0x4000")
	B[14], _ = hexutil.DecodeBig("0x8000")
	B[15], _ = hexutil.DecodeBig("0x10000")
	B[16], _ = hexutil.DecodeBig("0x20000")
	B[17], _ = hexutil.DecodeBig("0x40000")
	B[18], _ = hexutil.DecodeBig("0x80000")

	C[0], _ = hexutil.DecodeBig("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	C[1], _ = hexutil.DecodeBig("0xFFFFFFFFFFFFFFFF")
	C[2], _ = hexutil.DecodeBig("0xFFFFFFFF")
	C[3], _ = hexutil.DecodeBig("0xFFFF")
	C[4], _ = hexutil.DecodeBig("0xFF")
	C[5], _ = hexutil.DecodeBig("0xF")
	C[6], _ = hexutil.DecodeBig("0x3")

	ZERO_UINT_256, _ = hexutil.DecodeBig("0x0")
	ONE_UINT_256, _ = hexutil.DecodeBig("0x1")
	INIT0, _ = hexutil.DecodeBig("0xfffcb933bd6fad37aa2d162d1a594001")
	INIT1, _ = hexutil.DecodeBig("0x100000000000000000000000000000000")
	MAX_UINT_256 = GetMaxValue(256)
	MAX_UINT_160 = GetMaxValue(160)
	SQRT_10001, _ = hexutil.DecodeBig("0x3627A301D71055774C85")
	LOWER_ERR_BOUND, _ = hexutil.DecodeBig("0x28F6481AB7F045A5AF012A19D003AAA")
	UPPER_ERR_BOUND, _ = hexutil.DecodeBig("0xDB2DF09E81959A81455E260799A0632F")

	MIN_TICK = big.NewInt(-887272)
	MAX_TICK = big.NewInt(-MIN_TICK.Int64())

	MIN_SQRT_RATIO = big.NewInt(4295128739)
	MAX_SQRT_RATIO, _ = hexutil.DecodeBig("0xFFFD8963EFD1FC6A506488495D951D5263988D26")
}

func GetMaxValue(bitsCount uint) *big.Int {
	return new(big.Int).Sub(new(big.Int).Lsh(common.Big1, bitsCount), common.Big1)
}

// Ported function getSqrtRatioAtTick(int24 tick) internal pure returns (uint160 sqrtPriceX96)
// Source: https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/TickMath.sol
// tick int24
// return int96 sqrtPriceX96
func GetSqrtRatioAtTick(tick *big.Int) *big.Int {
	absTick := big.NewInt(0)
	absTick.Abs(tick)

	if absTick.Cmp(MAX_TICK) > 0 {
		panic(fmt.Sprintf("ticks: tick %d out of interval [%d, %d]", tick, MIN_TICK, MAX_TICK))
	}

	var ratio *big.Int = big.NewInt(0)
	var mask *big.Int = big.NewInt(0)

	mask.And(absTick, ONE_UINT_256)

	if mask.Cmp(ZERO_UINT_256) != 0 {
		ratio.Set(INIT0)
	} else {
		ratio.Set(INIT1)
	}

	for i := 0; i < len(A); i++ {
		mask.And(absTick, B[i])
		if mask.Cmp(ZERO_UINT_256) != 0 {
			ratio.Mul(ratio, A[i])
			ratio.Rsh(ratio, 128)
		}
	}

	if tick.Sign() > 0 {
		ratio.Div(MAX_UINT_256, ratio)
	}

	mask.Lsh(ONE_UINT_256, 32)
	mask.Mod(ratio, mask)

	ratio.Rsh(ratio, 32)

	if mask.Cmp(ZERO_UINT_256) != 0 {
		ratio.Add(ratio, ONE_UINT_256)
	}

	return ratio
}

// Ported function getTickAtSqrtRatio(uint160 sqrtPriceX96) internal pure returns (int24 tick)
// Source: https://github.com/Uniswap/v3-core/blob/main/contracts/libraries/TickMath.sol
// sqrtPriceX96 int96
// return tick int24
func GetTickAtSqrtRatio(sqrtPriceX96 *big.Int) *big.Int {
	if sqrtPriceX96.Cmp(MIN_SQRT_RATIO) < 0 || sqrtPriceX96.Cmp(MAX_SQRT_RATIO) > 0 {
		panic(fmt.Sprintf("ticks: sqrtPriceX96 %d out of interval [%d, %d]", sqrtPriceX96, MIN_SQRT_RATIO, MAX_SQRT_RATIO))
	}

	ratio := big.NewInt(0)
	r := big.NewInt(0)

	var msb int = 0

	ratio.Lsh(sqrtPriceX96, 32)
	r.Set(ratio)

	for i := 0; i < 7; i++ {
		rel := r.Cmp(C[i])

		if rel > 0 {
			f := rel << (7 - i)
			msb |= f
			r.Rsh(r, uint(f))
		}
	}

	rel := r.Cmp(ONE_UINT_256)
	if rel > 0 {
		msb |= 1
	}

	if msb >= 128 {
		r.Rsh(ratio, uint(msb-127))
	} else {
		r.Lsh(ratio, uint(127-msb))
	}

	log_2 := big.NewInt(int64(msb - 128))
	log_2.Lsh(log_2, 64)

	for i := 0; i < 14; i++ {
		r.Mul(r, r)
		r.Rsh(r, 127)
		f := big.NewInt(0).Set(r)
		f.Rsh(f, 128)
		fSh := big.NewInt(0)
		fSh.Lsh(f, uint(63-i))
		log_2.Or(log_2, fSh)
		r.Rsh(r, uint(f.Uint64()))
	}

	log_2.Mul(log_2, SQRT_10001)

	tickLow := big.NewInt(0)
	tickHi := big.NewInt(0)

	tickLow.Sub(log_2, LOWER_ERR_BOUND)
	tickLow.Rsh(tickLow, 128)

	tickHi.Add(log_2, UPPER_ERR_BOUND)
	tickHi.Rsh(tickHi, 128)

	if tickLow.Cmp(tickHi) == 0 {
		return tickLow
	}

	if GetSqrtRatioAtTick(tickHi).Cmp(sqrtPriceX96) <= 0 {
		return tickHi
	}

	return tickLow
}
