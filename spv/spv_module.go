package spv

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"os"

	"github.com/ioeX/ioeX.SideChain/config"
	"github.com/ioeX/ioeX.SideChain/core"
	"github.com/ioeX/ioeX.SideChain/log"

	"github.com/ioeX/ioeX.SPV/interface"
	spvlog "github.com/ioeX/ioeX.SPV/log"
	. "github.com/ioeX/ioeX.MainChain/bloom"
	ioex "github.com/ioeX/ioeX.MainChain/core"
)

var spvService _interface.SPVService

func SpvInit() error {
	var err error
	spvlog.Init(config.Parameters.SpvPrintLevel)

	var id = make([]byte, 8)
	var clientId uint64
	rand.Read(id)
	binary.Read(bytes.NewReader(id), binary.LittleEndian, &clientId)

	spvService, err = _interface.NewSPVService(config.Parameters.SpvMagic, clientId,
		config.Parameters.SpvSeedList, config.Parameters.SpvMinOutbound, config.Parameters.SpvMaxConnections)
	if err != nil {
		return err
	}

	go func() {
		if err := spvService.Start(); err != nil {
			log.Info("Spv service start failed ：", err)
		}
		log.Info("Spv service stoped")
		os.Exit(-1)
	}()
	return nil
}

func VerifyTransaction(tx *core.Transaction) error {
	proof := new(MerkleProof)
	mainChainTransaction := new(ioex.Transaction)

	payloadObj, ok := tx.Payload.(*core.PayloadRechargeToSideChain)
	if !ok {
		return errors.New("Invalid payload core.PayloadRechargeToSideChain")
	}

	reader := bytes.NewReader(payloadObj.MerkleProof)
	if err := proof.Deserialize(reader); err != nil {
		return errors.New("RechargeToSideChain payload deserialize failed")
	}
	reader = bytes.NewReader(payloadObj.MainChainTransaction)
	if err := mainChainTransaction.Deserialize(reader); err != nil {
		return errors.New("RechargeToSideChain mainChainTransaction deserialize failed")
	}

	if err := spvService.VerifyTransaction(*proof, *mainChainTransaction); err != nil {
		return errors.New("SPV module verify transaction failed.")
	}

	return nil
}
