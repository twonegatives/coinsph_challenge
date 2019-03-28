package banking

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeGetAccountsEndpoint(svc BankingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccountsRequest)
		accounts, err := svc.GetAccountsList(ctx, req)
		return accounts, err
	}
}

func MakeGetPaymentsEndpoint(svc BankingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetPaymentsRequest)
		payments, err := svc.GetPaymentsList(ctx, req)
		return payments, err
	}
}

func MakeSendPaymentEndpoint(svc BankingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SendPaymentRequest)
		err := svc.SendPayment(ctx, req)
		return nil, err
	}
}
