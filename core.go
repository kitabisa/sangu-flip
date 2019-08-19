package flip

import (
	"io"
	"net/url"
	"strings"
)

const (
	BankListURL = "/general/banks"
	BankAccountInquiryURL= "/disbursement/bank-account-inquiry"
)

type CoreGateway struct {
	Client Client
}

func (gateway *CoreGateway) Call(method, path string, header map[string]string, body io.Reader, v interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.BaseURL + path

	return gateway.Client.Call(method, path, header, body, v)
}

func (gateway *CoreGateway) GetBanks(bankCode string) (resp []Banks, err error) {
	data := url.Values{}
	if bankCode != "" {
		data.Set("code", bankCode)
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	err = gateway.Call("GET", BankListURL, headers, strings.NewReader(data.Encode()), &resp)
	if err != nil {
		return
	}

	return
}

func (gateway *CoreGateway) GetBankAccountInquiry(bankCode string, accountNumber string) (resp BankAccountInquiry, err error) {
	data := url.Values{}
	if bankCode != "" {
		data.Set("bank_code", bankCode)
	}

	if accountNumber != "" {
		data.Set("account_number", accountNumber)
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	err = gateway.Call("POST", BankAccountInquiryURL, headers, strings.NewReader(data.Encode()), &resp)
	if err != nil {
		return
	}

	return
}
