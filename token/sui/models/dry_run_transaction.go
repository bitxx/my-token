package models

import (
	"mytoken/core/lib/decimal"
	"mytoken/token/sui/types"
)

type DryRunTransaction struct {
	Effects struct {
		MessageVersion string `json:"messageVersion"`
		Status         struct {
			Status string `json:"status"`
			Error  string `json:"error"`
		} `json:"status"`
		ExecutedEpoch string `json:"executedEpoch"`
		GasUsed       struct {
			ComputationCost         string `json:"computationCost"`
			StorageCost             string `json:"storageCost"`
			StorageRebate           string `json:"storageRebate"`
			NonRefundableStorageFee string `json:"nonRefundableStorageFee"`
		} `json:"gasUsed"`
		ModifiedAtVersions []struct {
			ObjectID       types.ObjectId  `json:"objectId"`
			SequenceNumber decimal.Decimal `json:"sequenceNumber"`
		} `json:"modifiedAtVersions"`
		TransactionDigest string `json:"transactionDigest"`
		Mutated           []struct {
			Owner struct {
				AddressOwner string `json:"AddressOwner"`
			} `json:"owner"`
			Reference struct {
				ObjectID types.ObjectId `json:"objectId"`
				Version  int            `json:"version"`
				Digest   string         `json:"digest"`
			} `json:"reference"`
		} `json:"mutated"`
		GasObject struct {
			Owner struct {
				AddressOwner string `json:"AddressOwner"`
			} `json:"owner"`
			Reference struct {
				ObjectID types.ObjectId `json:"objectId"`
				Version  int            `json:"version"`
				Digest   string         `json:"digest"`
			} `json:"reference"`
		} `json:"gasObject"`
		Dependencies []string `json:"dependencies"`
	} `json:"effects"`
	Events        []interface{} `json:"events"`
	ObjectChanges []struct {
		Type   string `json:"type"`
		Sender string `json:"sender"`
		Owner  struct {
			AddressOwner string `json:"AddressOwner"`
		} `json:"owner"`
		ObjectType      string          `json:"objectType"`
		ObjectID        types.ObjectId  `json:"objectId"`
		Version         decimal.Decimal `json:"version"`
		PreviousVersion decimal.Decimal `json:"previousVersion"`
		Digest          string          `json:"digest"`
	} `json:"objectChanges"`
	BalanceChanges []struct {
		Owner struct {
			AddressOwner string `json:"AddressOwner"`
		} `json:"owner"`
		CoinType string `json:"coinType"`
		Amount   string `json:"amount"`
	} `json:"balanceChanges"`
}

func (te *DryRunTransaction) GasFee() (*decimal.Decimal, error) {
	gasUsed := te.Effects.GasUsed
	computationCost, err := decimal.NewFromString(gasUsed.ComputationCost)
	if err != nil {
		return nil, err
	}
	storageCost, err := decimal.NewFromString(gasUsed.StorageCost)
	if err != nil {
		return nil, err
	}
	storageRebate, err := decimal.NewFromString(gasUsed.StorageRebate)
	if err != nil {
		return nil, err
	}
	nonRefundableStorageFee, err := decimal.NewFromString(gasUsed.NonRefundableStorageFee)
	if err != nil {
		return nil, err
	}
	fee := computationCost.Add(storageCost).Sub(storageRebate).Sub(nonRefundableStorageFee)
	return &fee, nil
}
