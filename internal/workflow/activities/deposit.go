package activities

import (
	"context"
	"fmt"
	"log"

	E "github.com/IBM/fp-go/either"

	S "temporal-sample/internal/services"
	M "temporal-sample/internal/workflow/shared"
)

func Deposit(ctx context.Context, data M.PaymentDetails) (string, error) {
	log.Printf("Depositing $%d into account %s.\n\n",
		data.Amount,
		data.TargetAccount,
	)

	referenceID := fmt.Sprintf("%s-deposit", data.ReferenceID)
	bank := S.BankingServiceOf("bank-api.example.com")
	// Uncomment the next line and comment the one after that to simulate an unknown failure
	// confirmation, err := bank.DepositThatFails(data.TargetAccount, data.Amount, referenceID)
	return E.Unwrap(bank.Deposit(data.TargetAccount, data.Amount, referenceID))
}
