package epc

import (
	"errors"
	"strings"

	"github.com/dundee/qrpay/base"
)

const EpcHeader = `BCD
002
1
SCT
`

type EpcPayment struct {
	*base.Payment
	Purpose string
}

func NewEpcPayment() *EpcPayment {
	return &EpcPayment{
		Payment: &base.Payment{
			Errors: make(map[string]error),
		},
	}
}

func (p *EpcPayment) SetPurpose(value string) {
	p.Purpose = value
}

func (p *EpcPayment) GenerateString() (string, error) {
	var sb strings.Builder
	sb.WriteString(EpcHeader)

	writeIfNotEmpty := func(value string, maxLength int) {
		if value != "" {
			sb.WriteString(base.TrimToLength(value, maxLength) + "\n")
		} else {
			sb.WriteString("\n")
		}
	}

	if p.BIC != "" {
		sb.WriteString(base.TrimToLength(p.BIC, 11) + "\n")
	} else {
		sb.WriteString("\n")
	}

	if p.Recipient == "" {
		return "", errors.New("name of the beneficiary is mandatory")
	}
	sb.WriteString(base.TrimToLength(p.Recipient, 70) + "\n")

	if p.IBAN == "" {
		return "", errors.New("IBAN is mandatory")
	}
	sb.WriteString(base.TrimToLength(p.IBAN, 34) + "\n")

	if p.Amount != "" {
		writeIfNotEmpty(p.Currency, 3)
		sb.WriteString(base.TrimToLength(p.Amount, 12) + "\n")
	} else {
		sb.WriteString("\n")
	}

	writeIfNotEmpty(p.Purpose, 4)
	writeIfNotEmpty(p.Reference, 140)
	writeIfNotEmpty(p.Msg, 70)

	return strings.TrimSpace(sb.String()), nil
}
