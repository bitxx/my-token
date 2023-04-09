package sui

import (
	"bytes"
	"fmt"
	"io"
	"mytoken/token/sui/types"
	"net/http"
)

func FaucetFundAccount(address string, faucetUrl string) ([]byte, error) {
	_, err := types.NewAddressFromHex(address)
	if err != nil {
		return nil, err
	}

	paramJson := fmt.Sprintf(`{"FixedAmountRequest":{"recipient":"%v"}}`, address)
	request, err := http.NewRequest(http.MethodPost, faucetUrl, bytes.NewBuffer([]byte(paramJson)))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		return nil, fmt.Errorf("post %v response code = %v", faucetUrl, res.Status)
	}
	defer res.Body.Close()

	resByte, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return resByte, nil
}
