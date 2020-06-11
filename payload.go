package flip

type DisbursementRequest struct {
	AccountNumber string `json:"account_number" url:"account_number" valid:"required"`
	BankCode      string `json:"bank_code" url:"bank_code" valid:"required"`
	Amount        uint32 `json:"amount" url:"amount" valid:"required"`
	Remark        string `json:"remark" url:"remark"`
	RecipientCity uint32 `json:"recipient_city" url:"recipient_city"`
}
