package bare_metal

import "fmt"

type TierPrice struct {
	Price
	Tier *int `json:"tier"`
}

type Price struct {
	Amount float32 `json:"amount"`
	Period string  `json:"period"`
}

func (p Price) String() string {
	return fmt.Sprintf("$%-.2f/%s", p.Amount, p.Period)
}
