package flip

type Banks struct {
	BankCode string  `json:"bank_code"`
	Name     string  `json:"name"`
	Fee      float32 `json:"fee"`
	Queue    int     `json:"queue"`
	Status   string  `json:"status"`
}
