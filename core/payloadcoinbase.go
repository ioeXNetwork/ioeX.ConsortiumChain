package core

import (
	"io"

	. "github.com/ioeX/ioeX.Utility/common"
)

const PayloadCoinBaseVersion byte = 0x04

type PayloadCoinBase struct {
	CoinbaseData []byte
}

func (a *PayloadCoinBase) Data(version byte) []byte {
	return a.CoinbaseData
}

func (a *PayloadCoinBase) Serialize(w io.Writer, version byte) error {
	return WriteVarBytes(w, a.CoinbaseData)
}

func (a *PayloadCoinBase) Deserialize(r io.Reader, version byte) error {
	temp, err := ReadVarBytes(r)
	a.CoinbaseData = temp
	return err
}
