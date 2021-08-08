package coinType

import "encoding/json"

type CoinType int

const (
	ETHEUREM CoinType = iota
	BITCOIN
	NOT_A_COIN
)

func GetSupportedCoinTypes() []CoinType{
	coinTypes := []CoinType{ETHEUREM}
	return coinTypes
}

func (c CoinType) String() string {
	return [...]string{"Etheurem", "Bitcoin" ,"NotACoin"}[c]
}

func (c CoinType) Repr() string {
	return [...]string{"60'", "0'" ,"-1"}[c]
}

func GetCoinType(repr string) CoinType {
	switch repr {
		case "60'":
			return ETHEUREM
		case "0'":
			return BITCOIN
		default:
			return NOT_A_COIN 
	}
}

func (c CoinType) MarshalJSON() ([]byte, error) {
    return []byte(`"` + c.Repr() + `"`), nil
}

func (this CoinType) UnmarshalJSON(b []byte) error {
	var c string
	err := json.Unmarshal(b, &c)
	if err != nil {
		return err
	}
	this = GetCoinType(c)
	return nil
}

type UnsupportedCoinError struct {
	Message string
}

func (e *UnsupportedCoinError) Error() string {
	return e.Message
}

func (c CoinType) CheckSupportCompatability() bool {
	for _, coinType := range GetSupportedCoinTypes() {
		if c == coinType {
			return true
		}
	}
	return false
}
