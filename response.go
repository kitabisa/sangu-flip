package flip

type Balance struct {
	Balance int64 `json:"balance"`
}

type CheckIsOperational struct {
	Operational bool `json:"operational"`
}

type CheckIsMaintenance struct {
	Maintenance bool `json:"maintenance"`
}

type Bank struct {
	BankCode string `json:"bank_code"`
	Name     string `json:"name"`
	Fee      uint32 `json:"fee"`
	Queue    int    `json:"queue"`
	Status   string `json:"status"` // OPERATIONAL, DISTURBED, HEAVILY_DISTURBED
}

type BankAccountInquiry struct {
	BankCode      string `json:"bank_code"`
	AccountNumber string `json:"account_number"`
	AccountHolder string `json:"account_holder"`
	Status        string `json:"status"`
}

type Disbursement struct {
	ID            uint32  `json:"id"`
	UserID        uint32  `json:"user_id"`
	Amount        uint32  `json:"amount"`
	Status        string  `json:"status"`    // PENDING, CANCELLED, DONE
	Timestamp     string  `json:"timestamp"` // (Format: YYYY-MM-DD H:i:s)
	BankCode      string  `json:"bank_code"`
	AccountNumber string  `json:"account_number"`
	RecipientName string  `json:"recipient_name"`
	SenderBank    string  `json:"sender_bank"`
	Remark        string  `json:"remark"`
	Receipt       string  `json:"receipt"`
	TimeServed    string  `json:"time_served"` // Will only show valid value when the status is DONE (Format: YYYY-MM-DD H:i:s)
	BundleID      uint32  `json:"bundle_id"`
	CompanyID     uint32  `json:"company_id"`
	RecipientCity uint32  `json:"recipient_city"`
	CreatedFrom   string  `json:"created_from"`
	Direction     string  `json:"direction"`
	Sender        *Sender `json:"sender"`
	Fee           uint32  `json:"fee"`
}

type DisbursementError struct {
	Code   string            `json:"code"`
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Attribute string `json:"attribute"`
	Code      uint16 `json:"code"`
	Message   string `json:"message"`
}

type DisbursementQueue struct {
	Queue uint32 `json:"queue"`
}

type Sender struct {
	SenderName           string `json:"sender_name"`
	PlaceOfBirth         uint32 `json:"place_of_birth"`
	DateOfBirth          string `json:"date_of_birth"` // (Format: YYYY-MM-DD)
	Address              string `json:"address"`
	SenderIdentityType   string `json:"sender_identity_type"`
	SenderIdentityNumber string `json:"sender_identity_number"`
	SenderCountry        int16  `json:"sender_country"`
	Job                  string `json:"babu"`
}

type GetAllDisbursement struct {
	TotalData   uint32         `json:"total_data"`
	DataPerPage uint32         `json:"data_per_page"`
	TotalPage   uint32         `json:"total_page"`
	Page        uint32         `json:"page"`
	Data        []Disbursement `json:"data"`
}

type ErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

type GeneralErrorResponse struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	Type    string `json:"type"`
}
