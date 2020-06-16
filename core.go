package flip

import (
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	BankListURL             = "/general/banks"
	BankAccountInquiryURL   = "/disbursement/bank-account-inquiry"
	CurrentBalanceURL       = "/general/balance"
	CheckIsOperationalURL   = "/general/operational"
	CheckIsMaintenanceURL   = "/general/maintenance"
	GetAllDisbursementURL   = "/disbursement"
	GetDisbursementQueueURL = "/disbursement/[trxId]/queue"
	DisburseURL             = "/disbursement"
)

type CoreGateway struct {
	Client Client
}

func (gateway *CoreGateway) Call(method, path string, header map[string]string, body io.Reader, v interface{}, x interface{}) error {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	path = gateway.Client.BaseURL + path

	return gateway.Client.Call(method, path, header, body, v, x)
}

func (gateway *CoreGateway) GetCurrentBalance() (resp Balance, err error) {
	err = gateway.Call("GET", CurrentBalanceURL, nil, nil, &resp, nil)
	return
}

func (gateway *CoreGateway) GetBanks(bankCode string) (resp []Bank, err error) {
	data := url.Values{}
	if bankCode != "" {
		data.Set("code", bankCode)
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	err = gateway.Call("GET", fmt.Sprintf("%s?%s", BankListURL, data.Encode()), headers, nil, &resp, nil)
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

	err = gateway.Call("POST", BankAccountInquiryURL, headers, strings.NewReader(data.Encode()), &resp, nil)
	return
}

func (gateway *CoreGateway) CheckIsMaintenance() (resp CheckIsMaintenance, err error) {
	err = gateway.Call("GET", CheckIsMaintenanceURL, nil, nil, &resp, nil)
	return
}

func (gateway *CoreGateway) CheckIsOperational() (resp CheckIsOperational, err error) {
	err = gateway.Call("GET", CheckIsOperationalURL, nil, nil, &resp, nil)
	return
}

func (gateway *CoreGateway) GetAllDisbursement(perPage int, page int) (resp GetAllDisbursement, err error) {
	data := url.Values{}
	strPerPage := strconv.Itoa(perPage)
	strPage := strconv.Itoa(page)
	data.Set("pagination", strPerPage)
	data.Set("page", strPage)
	data.Set("sort", "-id")

	err = gateway.Call("GET", fmt.Sprintf("%s?%s", GetAllDisbursementURL, data.Encode()), nil, nil, &resp, nil)
	return
}

func (gateway *CoreGateway) GetDisbursementInfo(trxId int) (resp Disbursement, respError ErrorResponse, err error) {

	err = gateway.Call("GET", fmt.Sprintf("%s/%d", DisburseURL, trxId), nil, nil, &resp, &respError)
	return
}

func (gateway *CoreGateway) GetDisbursementQueue(trxId int) (resp DisbursementQueue, respError ErrorResponse, err error) {
	strTrxId := strconv.Itoa(trxId)
	url := strings.Replace(GetDisbursementQueueURL, "[trxId]", strTrxId, -1)

	err = gateway.Call("GET", url, nil, nil, &resp, &respError)
	return
}

func (gateway *CoreGateway) Disburse(idempotencyKey string, payload DisbursementRequest) (resp Disbursement, respError DisbursementError, err error) {
	headers := map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"idempotency-key": idempotencyKey,
	}

	params, err := query.Values(payload)
	if err != nil {
		return
	}

	err = gateway.Call("POST", DisburseURL, headers, strings.NewReader(params.Encode()), &resp, &respError)
	return
}
