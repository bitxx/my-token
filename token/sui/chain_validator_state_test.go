package sui

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGetGlobalActiveValidatorState(t *testing.T) {
	chain := TestnetChain()
	validatorService := NewValidatorState(chain)
	validatorState, err := validatorService.TotalActiveValidatorState(true)
	require.Nil(t, err)
	//t.Log(validatorState)
	for _, v := range validatorState.Validators {
		t.Logf("%-10v,Address = %s, APY = %v", v.Name, v.Address, v.APY)
	}
}

func TestGetValidator(t *testing.T) {
	v1Addr := "0x520289e77c838bae8501ae92b151b99a54407288fdd20dee6e5416bfe943eb7a"
	v2Addr := "0x89afed39dde1ce7d5f1c78f8c832e254b75e7c13ddf5158fa460e46416bd8f00"
	validatorStateService := NewValidatorState(TestnetChain())

	timeStart := time.Now()
	v1, err := validatorStateService.GetValidator(v1Addr, true)
	require.Nil(t, err)
	t.Logf("%-10v,Address = %s, APY = %v", v1.Name, v1.Address, v1.APY)
	timeMiddle := time.Now()
	fmt.Println("第一次执行耗时：", timeMiddle.Sub(timeStart))

	v2, err := validatorStateService.GetValidator(v2Addr, true)
	require.Nil(t, err)
	t.Logf("%-10v,Address = %s, APY = %v", v2.Name, v2.Address, v2.APY)
	fmt.Println("第二次执行耗时：", time.Now().Sub(timeMiddle))
}
