package activities

import (
	"context"
	"fmt"
	"log"

	E "github.com/IBM/fp-go/either"

	S "temporal-sample/internal/services"
	M "temporal-sample/internal/workflow/shared"
)

func Refund(ctx context.Context, data M.PaymentDetails) (string, error) {
	log.Printf("Refunding $%v back into account %v.\n\n",
		data.Amount,
		data.SourceAccount,
	)

	referenceID := fmt.Sprintf("%s-refund", data.ReferenceID)
	bank := S.BankingServiceOf("bank-api.example.com")
	return E.Unwrap(bank.Deposit(data.SourceAccount, data.Amount, referenceID))
}
