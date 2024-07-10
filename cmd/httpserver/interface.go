package httpserver

type AccountRequest struct {
	Account  string `json:"account"`  // 계좌 번호
	Password string `json:"password"` // 계좌 비밀번호
	Amount   int    `json:"amount"`   // 금액
}

type AccountInfo struct {
	Password string    `json:"password,omitempty"`
	Balance  int       `json:"balance"`
	History  []History `json:"history"`
}

type History struct {
	Time    string `json:"time"`    // 거래 시간
	Flow    string `json:"flow"`    // 입출금 여부
	Amount  int    `json:"amount"`  // 거래 금액
	Balance int    `json:"balance"` // 거래 후 잔액
}

type Handler interface {
	Create(*AccountRequest) error                  // 계좌 생성
	History(*AccountRequest) (*AccountInfo, error) // 계좌 정보
	Deposit(*AccountRequest) error                 // 입금
	Withdraw(*AccountRequest) error                // 출금
}
