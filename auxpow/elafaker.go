package auxpow

import (
	"time"

	. "github.com/ioeX/ioeX.Utility/common"
	"github.com/ioeX/ioeX.MainChain/auxpow"
	ioex "github.com/ioeX/ioeX.MainChain/core"
)

func getSideChainPowTx(msgBlockHash Uint256, genesisHash Uint256) *ioex.Transaction {

	txPayload := &ioex.PayloadSideChainPow{
		SideBlockHash:   msgBlockHash,
		SideGenesisHash: genesisHash,
	}

	sideChainPowTx := NewSideChainPowTx(txPayload, 0)

	return sideChainPowTx
}

func GenerateSideAuxPow(msgBlockHash Uint256, genesisHash Uint256) *SideAuxPow {
	sideAuxMerkleBranch := make([]Uint256, 0)
	sideAuxMerkleIndex := 0
	sideAuxBlockTx := getSideChainPowTx(msgBlockHash, genesisHash)
	elaBlockHeader := ioex.Header{
		Version:    0x7fffffff,
		Previous:   EmptyHash,
		MerkleRoot: sideAuxBlockTx.Hash(),
		Timestamp:  uint32(time.Now().Unix()),
		Bits:       0,
		Nonce:      0,
		Height:     0,
	}

	elahash := elaBlockHeader.Hash()
	// fake a btc blockheader and coinbase
	newAuxPow := auxpow.GenerateAuxPow(elahash)
	elaBlockHeader.AuxPow = *newAuxPow

	sideAuxPow := NewSideAuxPow(
		sideAuxMerkleBranch,
		sideAuxMerkleIndex,
		*sideAuxBlockTx,
		elaBlockHeader,
	)

	return sideAuxPow
}

func NewSideChainPowTx(payload *ioex.PayloadSideChainPow, currentHeight uint32) *ioex.Transaction {
	return &ioex.Transaction{
		TxType:  ioex.SideChainPow,
		Payload: payload,
		Inputs: []*ioex.Input{
			{
				Previous: ioex.OutPoint{
					TxID:  EmptyHash,
					Index: 0x0000,
				},
				Sequence: 0x00000000,
			},
		},
		Attributes: []*ioex.Attribute{},
		LockTime:   currentHeight,
		Programs:   []*ioex.Program{},
	}
}
