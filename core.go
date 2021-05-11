package flip

import (
	"fmt"
	"io"
	"net/url"
	"sort"
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

func (gateway *CoreGateway) Call(method, path string, header map[string]string, body io.Reader, v interface{}) (err error, respErr ErrorResponse) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	path = gateway.Client.BaseURL + path
	err, respErr = gateway.Client.Call(method, path, header, body, v)
	return
}

func (gateway *CoreGateway) GetCurrentBalance() (resp Balance, respError ErrorResponse, err error) {
	err, respError = gateway.Call("GET", CurrentBalanceURL, nil, nil, &resp)
	return
}

func (gateway *CoreGateway) GetBanks(bankCode string) (resp []Bank, respError ErrorResponse, err error) {
	data := url.Values{}
	if bankCode != "" {
		data.Set("code", bankCode)
	}

	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}

	err, respError = gateway.Call("GET", fmt.Sprintf("%s?%s", BankListURL, data.Encode()), headers, nil, &resp)
	if err == nil && respError.Message == "" {
		sort.Slice(resp, func(i, j int) bool {
			return resp[i].Name < resp[j].Name
		})
	}

	return
}

func (gateway *CoreGateway) GetBankAccountInquiry(bankCode string, accountNumber string) (resp BankAccountInquiry, respError ErrorResponse, err error) {
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

	err, respError = gateway.Call("POST", BankAccountInquiryURL, headers, strings.NewReader(data.Encode()), &resp)
	return
}

func (gateway *CoreGateway) CheckIsMaintenance() (resp CheckIsMaintenance, respError ErrorResponse, err error) {
	err, respError = gateway.Call("GET", CheckIsMaintenanceURL, nil, nil, &resp)
	return
}

func (gateway *CoreGateway) CheckIsOperational() (resp CheckIsOperational, respError ErrorResponse, err error) {
	err, respError = gateway.Call("GET", CheckIsOperationalURL, nil, nil, &resp)
	return
}

func (gateway *CoreGateway) GetAllDisbursement(perPage int, page int) (resp GetAllDisbursement, respError ErrorResponse, err error) {
	data := url.Values{}
	strPerPage := strconv.Itoa(perPage)
	strPage := strconv.Itoa(page)
	data.Set("pagination", strPerPage)
	data.Set("page", strPage)
	data.Set("sort", "-id")

	err, respError = gateway.Call("GET", fmt.Sprintf("%s?%s", GetAllDisbursementURL, data.Encode()), nil, nil, &resp)
	return
}

func (gateway *CoreGateway) GetDisbursementInfo(trxId int) (resp Disbursement, respError ErrorResponse, err error) {

	err, respError = gateway.Call("GET", fmt.Sprintf("%s/%d", DisburseURL, trxId), nil, nil, &resp)
	return
}

func (gateway *CoreGateway) GetDisbursementQueue(trxId int) (resp DisbursementQueue, respError ErrorResponse, err error) {
	strTrxId := strconv.Itoa(trxId)
	url := strings.Replace(GetDisbursementQueueURL, "[trxId]", strTrxId, -1)

	err, respError = gateway.Call("GET", url, nil, nil, &resp)
	return
}

func (gateway *CoreGateway) Disburse(idempotencyKey string, payload DisbursementRequest) (resp Disbursement, respError ErrorResponse, err error) {
	headers := map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"idempotency-key": idempotencyKey,
	}

	params, err := query.Values(payload)
	if err != nil {
		return
	}

	err, respError = gateway.Call("POST", DisburseURL, headers, strings.NewReader(params.Encode()), &resp)
	return
}
