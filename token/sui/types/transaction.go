package types

type SignatureSchemeSerialized byte

const (
	TxnRequestTypeWaitForLocalExecution = "WaitForLocalExecution"
)

type TransactionFilter struct {
	Checkpoint   uint64 `json:"Checkpoint,omitempty"`
	MoveFunction *struct {
		Package  ObjectId `json:"package"`
		Module   string   `json:"module,omitempty"`
		Function string   `json:"function,omitempty"`
	} `json:"MoveFunction,omitempty"`
	InputObject      *ObjectId `json:"InputObject,omitempty"`
	ChangedObject    *ObjectId `json:"ChangedObject,omitempty"`
	FromAddress      *Address  `json:"FromAddress,omitempty"`
	ToAddress        *Address  `json:"ToAddress,omitempty"`
	FromAndToAddress *struct {
		From *Address `json:"from"`
		To   *Address `json:"to"`
	} `json:"FromAndToAddress,omitempty"`
	TransactionKind *string `json:"TransactionKind,omitempty"`
}

type SuiTransactionBlockResponseOptions struct {
	/* Whether to show transaction input data. Default to be false. */
	ShowInput bool `json:"showInput,omitempty"`
	/* Whether to show transaction effects. Default to be false. */
	ShowEffects bool `json:"showEffects,omitempty"`
	/* Whether to show transaction events. Default to be false. */
	ShowEvents bool `json:"showEvents,omitempty"`
	/* Whether to show object changes. Default to be false. */
	ShowObjectChanges bool `json:"showObjectChanges,omitempty"`
	/* Whether to show coin balance changes. Default to be false. */
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
}

type SuiTransactionBlockResponseQuery struct {
	Filter  *TransactionFilter                  `json:"filter,omitempty"`
	Options *SuiTransactionBlockResponseOptions `json:"options,omitempty"`
}

// TransactionBlockOptions 请求选项，确定要返回哪些内容
type TransactionBlockOptions struct {
	/* Whether to show transaction input data. Default to be false. */
	ShowInput bool `json:"showInput,omitempty"`
	/* Whether to show transaction effects. Default to be false. */
	ShowEffects bool `json:"showEffects,omitempty"`
	/* Whether to show transaction events. Default to be false. */
	ShowEvents bool `json:"showEvents,omitempty"`
	/* Whether to show object changes. Default to be false. */
	ShowObjectChanges bool `json:"showObjectChanges,omitempty"`
	/* Whether to show coin balance changes. Default to be false. */
	ShowBalanceChanges bool `json:"showBalanceChanges,omitempty"`
}
