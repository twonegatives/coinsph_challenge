package banking

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

func MakeGetAccountsEndpoint(svc BankingService) endpoint.Endpoint {
	return func(ctx context.Context, _request interface{}) (interface{}, error) {
		accounts, err := svc.GetAccountsList(ctx)
		return getAccountsResponse{Accounts: accounts}, err
	}
}

func MakeGetPaymentsEndpoint(svc BankingService) endpoint.Endpoint {
	return func(ctx context.Context, _request interface{}) (interface{}, error) {
		payments, err := svc.GetPaymentsList(ctx)
		return getPaymentsResponse{Payments: payments}, err
	}
}

func MakeSendPaymentEndpoint(svc BankingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(sendPaymentRequest)
		err := svc.SendPayment(ctx, req.From, req.To, req.Amount)
		return map[string]interface{}{}, err
	}
}
