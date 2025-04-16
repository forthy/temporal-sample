package main

import (
	"fmt"
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"

	WF "temporal-sample/internal/workflow"
	AC "temporal-sample/internal/workflow/activities"
	SR "temporal-sample/internal/workflow/shared"
)

func withClient[T any](fn func(client.Client) E.Either[error, T]) E.Either[error, T] {
	return E.WithResource[error, client.Client, T](
		func() E.Either[error, client.Client] {
			return E.Eitherize1(client.Dial)(client.Options{})
		},
		func(c client.Client) E.Either[error, any] {
			c.Close()

			return E.Right[error](F.ToAny("Client closed"))
		},
	)(fn)
}

func runWorkflow(c client.Client) E.Either[error, <-chan any] {
	w := worker.New(c, SR.MoneyTransferTaskQueueName, worker.Options{})

	// This worker hosts both Workflow and Activity functions.
	w.RegisterWorkflow(WF.MoneyTransfer)
	w.RegisterActivity(AC.Withdraw)
	w.RegisterActivity(AC.Deposit)
	w.RegisterActivity(AC.Refund)

	// Start listening to the Task Queue.
	return E.FromError(w.Run)(worker.InterruptCh())
}

func main() {
	result := E.Fold(
		func(err error) string {
			return fmt.Sprintf("unable to start Worker:[%v]", err)
		},
		func(_ <-chan any) string {
			return "Worker started"
		},
	)(withClient(runWorkflow))

	log.Print(result)
}
