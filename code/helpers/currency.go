package helpers

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strings"
)

func ParseDecimalFromString(pS string) (decimal.Decimal, error) {
	removeBlank := strings.ReplaceAll(pS, " ", "")
	removeBRL := strings.ReplaceAll(removeBlank, "R$", "")
	switchComa := strings.ReplaceAll(removeBRL, ",", ".")

	fmt.Println(switchComa)
	value, err := decimal.NewFromString(switchComa)
	if err != nil {
		return decimal.Decimal{}, ErrInvalidCurrencyFormat
	}

	return value, err
}

func ParseDecimalToString(d decimal.Decimal) string {
	str := d.StringFixed(2)
	out := fmt.Sprintf("R$ %s", str)
	outF := strings.ReplaceAll(out, ".", ",")
	return outF
}
