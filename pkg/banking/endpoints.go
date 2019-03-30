package banking

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/shopspring/decimal"
)

func MakeCreateAccountEndpoint(svc BankingService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createAccountRequest)
		if req.Name == "" {
			return createAccountResponse{}, errBadRequest
		}
		account, err := svc.CreateAccount(ctx, req.Name)
		return createAccountResponse{Account: account}, err
	}
}

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
		if req.From.Name == "" || req.To.Name == "" || req.Amount.LessThanOrEqual(decimal.NewFromFloat(0)) {
			return nil, errBadRequest
		}
		err := svc.SendPayment(ctx, req.From, req.To, req.Amount)
		return map[string]interface{}{}, err
	}
}
