package uniswap_core

import (
	"math/big"
	"testing"
)

func TestMulDiv(t *testing.T) {
	expCount := 4
	a := make([]*big.Int, expCount)
	b := make([]*big.Int, expCount)
	d := make([]*big.Int, expCount)
	refRes := make([]*big.Int, expCount)

	a[0], _ = big.NewInt(0).SetString("12324343523453453245551232354554393884940323134234", 10)
	b[0], _ = big.NewInt(0).SetString("12312343453039344948483291032493344853202340934233", 10)
	d[0], _ = big.NewInt(0).SetString("93383847573820304934", 10)
	refRes[0], _ = big.NewInt(0).SetString("1624922877310743438599627458588429791217975980864583739966581637877291143881239", 10)

	a[1] = big.NewInt(5)
	b[1] = big.NewInt(10)
	d[1] = big.NewInt(2)
	refRes[1] = big.NewInt(25)

	a[2] = big.NewInt(-5)
	b[2] = big.NewInt(10)
	d[2] = big.NewInt(2)
	refRes[2] = big.NewInt(-25)

	a[3], _ = big.NewInt(0).SetString("-99999999999999999999999999999999999999999999432134234", 10)
	b[3], _ = big.NewInt(0).SetString("123123434530393449484832910324933448532303949485302340934233", 10)
	d[3], _ = big.NewInt(0).SetString("93383847573820304934", 10)
	refRes[3], _ = big.NewInt(0).SetString("-131846607019553228557932428302470359356043248972505423022200290231622387584755406540021411092", 10)

	for i := 0; i < expCount; i++ {
		res := MulDiv(a[i], b[i], d[i])

		if res.Cmp(refRes[i]) != 0 {
			t.Errorf("MulDiv(%d, %d, %d) = %d; want %d", a[i], b[i], d[i], res, refRes[i])
		}
	}
}

func TestMulDivRoundingUp(t *testing.T) {
	expCount := 4
	a := make([]*big.Int, expCount)
	b := make([]*big.Int, expCount)
	d := make([]*big.Int, expCount)
	refRes := make([]*big.Int, expCount)

	a[0], _ = big.NewInt(0).SetString("12324343523453453245551232354554393884940323134234", 10)
	b[0], _ = big.NewInt(0).SetString("12312343453039344948483291032493344853202340934233", 10)
	d[0], _ = big.NewInt(0).SetString("93383847573820304934", 10)
	refRes[0], _ = big.NewInt(0).SetString("1624922877310743438599627458588429791217975980864583739966581637877291143881240", 10)

	a[1] = big.NewInt(5)
	b[1] = big.NewInt(10)
	d[1] = big.NewInt(3)
	refRes[1] = big.NewInt(17)

	a[2] = big.NewInt(-5)
	b[2] = big.NewInt(10)
	d[2] = big.NewInt(2)
	refRes[2] = big.NewInt(-25)

	a[3], _ = big.NewInt(0).SetString("-99999999999999999999999999999999999999999999432134234", 10)
	b[3], _ = big.NewInt(0).SetString("123123434530393449484832910324933448532303949485302340934233", 10)
	d[3], _ = big.NewInt(0).SetString("93383847573820304934", 10)
	refRes[3], _ = big.NewInt(0).SetString("-131846607019553228557932428302470359356043248972505423022200290231622387584755406540021411091", 10)

	for i := 0; i < expCount; i++ {
		res := MulDivRoundingUp(a[i], b[i], d[i])

		if res.Cmp(refRes[i]) != 0 {
			t.Errorf("MulDivRoundingUp(%d, %d, %d) = %d; want %d", a[i], b[i], d[i], res, refRes[i])
		}
	}
}

func TestDivRoundingUp(t *testing.T) {

	a := big.NewInt(5)
	b := big.NewInt(2)
	res := DivRoundingUp(a, b)
	org := big.NewInt(3)

	if res.Cmp(org) != 0 {
		t.Errorf("DivRoundingUp(%d, %d) = %d; want %d", a, b, res, org)
	}
}
