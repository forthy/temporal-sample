package main

import (
	"context"
	"fmt"
	"log"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/workflow"

	WF "temporal-sample/internal/workflow"
	M "temporal-sample/internal/workflow/shared"
)

func withTemporalClient[T any](fn func(client client.Client) E.Either[error, T]) E.Either[error, T] {
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

func startWorkflow(c client.Client) E.Either[error, *string] {
	input := M.PaymentDetails{
		SourceAccount: "85-150",
		TargetAccount: "43-812",
		Amount:        250,
		ReferenceID:   "12345",
	}

	options := client.StartWorkflowOptions{
		ID:        "pay-invoice-701",
		TaskQueue: M.MoneyTransferTaskQueueName,
	}

	log.Printf("Starting transfer from account %s to account %s for %d", input.SourceAccount, input.TargetAccount, input.Amount)

	return E.Chain(func(we client.WorkflowRun) E.Either[error, *string] {
		// DEBUG
		log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

		var result string

		return E.FromError(func(v *string) error {
			return we.Get(context.Background(), v)
		})(&result)
	})(E.Eitherize4(func(
		ctx context.Context,
		options client.StartWorkflowOptions,
		wf func(workflow.Context, M.PaymentDetails) (string, error),
		input M.PaymentDetails,
	) (client.WorkflowRun, error) {
		return c.ExecuteWorkflow(ctx, options, wf, input)
	})(context.Background(), options, WF.MoneyTransfer, input))
}

func main() {
	result := E.Fold(
		func(err error) string {
			return fmt.Sprintf("Error: %v", err)
		},
		func(result *string) string {
			return fmt.Sprintf("Workflow result: %s", *result)
		},
	)(withTemporalClient[*string](startWorkflow))

	log.Println(result)
}
