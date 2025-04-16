package activities

import (
	"context"
	"fmt"
	"log"

	E "github.com/IBM/fp-go/either"

	S "temporal-sample/internal/services"
	M "temporal-sample/internal/workflow/shared"
)

func Withdraw(ctx context.Context, data M.PaymentDetails) (string, error) {
	log.Printf("Withdrawing $%d from account %s.\n\n",
		data.Amount,
		data.SourceAccount,
	)

	referenceID := fmt.Sprintf("%s-withdrawal", data.ReferenceID)
	bank := S.BankingServiceOf("bank-api.example.com")
	return E.Unwrap(bank.Withdraw(data.SourceAccount, data.Amount, referenceID))
}
