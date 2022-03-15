package uniswap_core

import (
	"math/big"
	"testing"
)

func TestGetTickAtSqrtRatio(t *testing.T) {

	tick := big.NewInt(-887272)
	sqrtPrice := GetSqrtRatioAtTick(tick)

	newTick := GetTickAtSqrtRatio(sqrtPrice)

	if newTick.Cmp(tick) != 0 {
		t.Errorf("GetTickAtSqrtRatio(%d) = %d; want %d", sqrtPrice, newTick, tick)
	}
}

func TestGetSqrtRatioAtTick(t *testing.T) {

	cases := []int64{1, 2, 3, 4, -1, -2, -3, -5}
	trues := make([]*big.Int, 8)
	trues[0], _ = big.NewInt(0).SetString("79232123823359799118286999568", 10)
	trues[1], _ = big.NewInt(0).SetString("79236085330515764027303304732", 10)
	trues[2], _ = big.NewInt(0).SetString("79240047035742135098198828268", 10)
	trues[3], _ = big.NewInt(0).SetString("79244008939048815603706035062", 10)
	trues[4], _ = big.NewInt(0).SetString("79224201403219477170569942574", 10)
	trues[5], _ = big.NewInt(0).SetString("79220240490215316061937756561", 10)
	trues[6], _ = big.NewInt(0).SetString("79216279775241952975272415332", 10)
	trues[7], _ = big.NewInt(0).SetString("79208358939348018173455069825", 10)

	for i, v := range cases {
		tick := big.NewInt(v)
		sqrtPrice := GetSqrtRatioAtTick(tick)

		if sqrtPrice.Cmp(trues[i]) != 0 {
			t.Errorf("GetSqrtRatioAtTick(%d) = %d; want %d", tick, sqrtPrice, trues[i])
		}
	}
}
