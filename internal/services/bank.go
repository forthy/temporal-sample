package services

import (
	E "github.com/IBM/fp-go/either"
)

type BankingService struct {
	bankAPI string
}

func BankingServiceOf(bankAPI string) BankingService {
	return BankingService{bankAPI: bankAPI}
}

func (b *BankingService) Deposit(account string, amount int, referenceID string) E.Either[error, string] {
	return E.Right[error]("confirmation")
}

func (b *BankingService) Withdraw(account string, amount int, referenceID string) E.Either[error, string] {
	return E.Right[error]("confirmation")
}
