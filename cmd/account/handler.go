package account

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"example.com/bear/cmd/httpserver"
)

type AccountHandler struct {
	dir string
}

func NewAccountHandler(path string) *AccountHandler {
	if path == "" {
		now, _ := os.Getwd()
		path = filepath.Join(now, "accounts")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}

	return &AccountHandler{
		dir: path,
	}
}

func (h *AccountHandler) Create(req *httpserver.AccountRequest) error {
	// 계좌 중복 확인
	_, err := os.Stat(filepath.Join(h.dir, req.Account+".json"))
	if err == nil {
		return fmt.Errorf("[Create] account already exists")
	}

	// 계좌 생성
	file, err := os.Create(filepath.Join(h.dir, req.Account+".json"))
	if err != nil {
		return fmt.Errorf("[Create] failed to create account: %w", err)
	}
	defer file.Close()

	// 계좌 생성 시 입금
	ai := httpserver.AccountInfo{
		Password: req.Password,
		Balance:  req.Amount,
		History:  []httpserver.History{},
	}
	if req.Amount > 0 {
		ai.History = append(ai.History, httpserver.History{
			Time:    time.Now().Format("2006-01-02 15:04:05"),
			Flow:    "Deposit",
			Amount:  req.Amount,
			Balance: req.Amount,
		})
	}

	// 계좌 정보 저장
	bytes, err := json.Marshal(ai)
	if err != nil {
		return fmt.Errorf("[Create] failed to marshal account info: %w", err)
	}
	err = os.WriteFile(filepath.Join(h.dir, req.Account+".json"), bytes, 0644)
	if err != nil {
		return fmt.Errorf("[Create] failed to write account info: %w", err)
	}

	return nil
}

func (h *AccountHandler) History(req *httpserver.AccountRequest) (*httpserver.AccountInfo, error) {
	// 계좌 정보 조회
	var ai httpserver.AccountInfo
	bytes, err := os.ReadFile(filepath.Join(h.dir, req.Account+".json"))
	if err != nil {
		return nil, fmt.Errorf("[History] failed to read account info: %w", err)
	}

	err = json.Unmarshal(bytes, &ai)
	if err != nil {
		return nil, fmt.Errorf("[History] failed to unmarshal account info: %w", err)
	}
	// 계좌 비밀번호 확인
	if ai.Password != req.Password {
		return nil, fmt.Errorf("[History] invalid password")
	}

	return &ai, nil
}

func (h *AccountHandler) Deposit(req *httpserver.AccountRequest) error {
	// 계좌 정보 조회
	ai, err := h.History(req)
	if err != nil {
		return fmt.Errorf("[Deposit] failed to get account info: %w", err)
	}

	// 입금
	ai.Balance += req.Amount
	ai.History = append(ai.History, httpserver.History{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Flow:    "Deposit",
		Amount:  req.Amount,
		Balance: ai.Balance,
	})

	bytes, err := json.Marshal(ai)
	if err != nil {
		return fmt.Errorf("[Deposit] failed to marshal account info: %w", err)
	}
	err = os.WriteFile(filepath.Join(h.dir, req.Account+".json"), bytes, 0644)
	if err != nil {
		return fmt.Errorf("[Deposit] failed to write account info: %w", err)
	}

	return nil
}

func (h *AccountHandler) Withdraw(req *httpserver.AccountRequest) error {
	// 계좌 정보 조회
	ai, err := h.History(req)
	if err != nil {
		return fmt.Errorf("[Withdraw] failed to get account info: %w", err)
	}

	if ai.Balance < req.Amount {
		return fmt.Errorf("[Withdraw] insufficient balance")
	}

	// 출금
	ai.Balance -= req.Amount
	ai.History = append(ai.History, httpserver.History{
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Flow:    "Withdraw",
		Amount:  req.Amount,
		Balance: ai.Balance,
	})

	bytes, err := json.Marshal(ai)
	if err != nil {
		return fmt.Errorf("[Withdraw] failed to marshal account info: %w", err)
	}
	err = os.WriteFile(filepath.Join(h.dir, req.Account+".json"), bytes, 0644)
	if err != nil {
		return fmt.Errorf("[Withdraw] failed to write account info: %w", err)
	}

	return nil
}
