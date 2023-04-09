package models

import "mytoken/token/sui/types"

type Object struct {
	ObjectID types.ObjectId `json:"objectId"`
	Version  int            `json:"version"`
	Digest   string         `json:"digest"`
	Type     string         `json:"type"`
	Owner    struct {
		AddressOwner string `json:"AddressOwner"`
	} `json:"owner"`
	PreviousTransaction string `json:"previousTransaction"`
	StorageRebate       int    `json:"storageRebate"`
	Display             struct {
		Data  interface{} `json:"data"`
		Error interface{} `json:"error"`
	} `json:"display"`
	Content struct {
		DataType          string `json:"dataType"`
		Type              string `json:"type"`
		HasPublicTransfer bool   `json:"hasPublicTransfer"`
		Fields            struct {
			Balance string `json:"balance"`
			ID      struct {
				ID string `json:"id"`
			} `json:"id"`
		} `json:"fields"`
	} `json:"content"`
	Bcs struct {
		DataType          string `json:"dataType"`
		Type              string `json:"type"`
		HasPublicTransfer bool   `json:"hasPublicTransfer"`
		Version           int    `json:"version"`
		BcsBytes          string `json:"bcsBytes"`
	} `json:"bcs"`
}

type ObjectIndex struct {
	Data []struct {
		Object `json:"data"`
	} `json:"data"`
	NextCursor  string `json:"nextCursor"`
	HasNextPage bool   `json:"hasNextPage"`
}
