package uniswap_core

import (
	"fmt"
	"math/big"
)

type BigInt struct {
	Val *big.Int
}

type BigDecimal struct {
	Val *big.Float
}

func (bi *BigInt) UnmarshalJSON(data []byte) error {
	bi.Val = big.NewInt(0)

	var ok bool
	strField := string(data[1 : len(data)-1])
	bi.Val, ok = bi.Val.SetString(strField, 10)

	if !ok {
		return fmt.Errorf("BigInt: UnmarshalJSON: something goes wrong with field %s", strField)
	}

	return nil
}

func (bi *BigDecimal) UnmarshalJSON(data []byte) error {
	bi.Val = big.NewFloat(0.0)

	var ok bool
	strField := string(data[1 : len(data)-1])
	bi.Val, ok = bi.Val.SetString(strField)

	if !ok {
		return fmt.Errorf("BigInt: UnmarshalJSON: something goes wrong with field %s", strField)
	}

	return nil
}

type FieldId struct {
	Id string
}

type Tick struct {
	Id                     string
	PoolAddress            string
	TickIdx                BigInt
	Pool                   FieldId
	LiquidityGross         BigInt
	LiquidityNet           BigInt
	Price0                 BigDecimal
	Price1                 BigDecimal
	VolumeToken0           BigDecimal
	VolumeToken1           BigDecimal
	VolumeUSD              BigDecimal
	UntrackedVolumeUSD     BigDecimal
	FeesUSD                BigDecimal
	CollectedFeesToken0    BigDecimal
	CollectedFeesToken1    BigDecimal
	CollectedFeesUSD       BigDecimal
	CreatedAtTimestamp     BigInt
	CreatedAtBlockNumber   BigInt
	LiquidityProviderCount BigInt
	FeeGrowthOutside0X128  BigInt
	FeeGrowthOutside1X128  BigInt
}

func (t Tick) IsInitialized() bool {
	return t.LiquidityGross.Val.Cmp(ZERO_UINT_256) != 0
}

type Token struct {
	Id                           string
	Symbol                       string
	Name                         string
	Decimals                     BigInt
	TotalSupply                  BigInt
	Volume                       BigDecimal
	VolumeUSD                    BigDecimal
	UntrackedVolumeUSD           BigDecimal
	FeesUSD                      BigDecimal
	TxCount                      BigInt
	PoolCount                    BigInt
	TotalValueLocked             BigDecimal
	TotalValueLockedUSD          BigDecimal
	TotalValueLockedUSDUntracked BigDecimal
	DerivedETH                   BigDecimal
}

type Tx struct {
	Id          string
	BlockNumber BigInt
	Timestamp   BigInt
	GasUsed     BigInt
	GasPrice    BigInt
}

type Swap struct {
	Id           string
	Transaction  Tx
	Timestamp    BigInt
	Pool         FieldId
	Token0       Token
	Token1       Token
	Sender       string
	Recipient    string
	Origin       string
	Amount0      BigDecimal
	Amount1      BigDecimal
	AmountUSD    BigDecimal
	SqrtPriceX96 BigInt
	Tick         BigInt
	LogIndex     BigInt
}

type Pool struct {
	Id                           string
	CreatedAtTimestamp           BigInt
	CreatedAtBlockNumber         BigInt
	Token0                       Token
	Token1                       Token
	FeeTier                      BigInt
	Liquidity                    BigInt
	SqrtPrice                    BigInt
	FeeGrowthGlobal0X128         BigInt
	FeeGrowthGlobal1X128         BigInt
	Token0Price                  BigDecimal
	Token1Price                  BigDecimal
	Tick                         BigInt
	ObservationIndex             BigInt
	VolumeToken0                 BigDecimal
	VolumeToken1                 BigDecimal
	VolumeUSD                    BigDecimal
	UntrackedVolumeUSD           BigDecimal
	FeesUSD                      BigDecimal
	TxCount                      BigInt
	CollectedFeesToken0          BigDecimal
	CollectedFeesToken1          BigDecimal
	CollectedFeesUSD             BigDecimal
	TotalValueLockedToken0       BigDecimal
	TotalValueLockedToken1       BigDecimal
	TotalValueLockedETH          BigDecimal
	TotalValueLockedUSD          BigDecimal
	TotalValueLockedUSDUntracked BigDecimal
	LiquidityProviderCount       BigInt
	Swaps                        []Swap
	Ticks                        []Tick
}

func (p Pool) FeerTierToTickSpacing() *big.Int {
	switch p.FeeTier.Val.Uint64() {
	case 10000:
		return big.NewInt(200)
	case 3000:
		return big.NewInt(60)
	case 500:
		return big.NewInt(10)
	case 100:
		return big.NewInt(1)
	}

	panic(fmt.Errorf("gql: Unexpected fee tier %d", p.FeeTier.Val))
}

func (p Pool) CurrentState() *Slot0 {
	slot0 := NewSlot0()
	slot0.TickSpacing.Set(p.FeerTierToTickSpacing())
	slot0.TickCurrent.Set(p.Tick.Val)
	slot0.Fee.Set(p.FeeTier.Val)
	slot0.Liquidity.Set(p.Liquidity.Val)
	slot0.FeeGrowthGlobal0X128.Set(p.FeeGrowthGlobal0X128.Val)
	slot0.FeeGrowthGlobal1X128.Set(p.FeeGrowthGlobal1X128.Val)
	slot0.SqrtPriceX96.Set(p.SqrtPrice.Val)
	slot0.FeeProtocol.Set(big.NewInt(0))

	return slot0
}
