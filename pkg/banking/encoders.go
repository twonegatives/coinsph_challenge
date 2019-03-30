package banking

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

// paymentsJSONEncoder is a helper struct to convert database- and application-level
// entity.Payment objects into a desired JSON format.
// this format differs from entities.Payment by having Counterparty name
// translated to either to_account or from_account based on payment Direction
type paymentsJSONEncoder struct {
}

// encode is the method which takes application-level entity.Payment
// and converts it to JSON format of paymentsJSONEncoder
func (e *paymentsJSONEncoder) encode(payments []entities.Payment) ([]byte, error) {
	if payments == nil {
		return nil, errors.New("payments array should be initialized in order to encode it")
	}

	elements := make([]map[string]interface{}, len(payments))

	for index, payment := range payments {
		element := map[string]interface{}{
			"account":   payment.Account.Name,
			"amount":    payment.Amount,
			"direction": payment.Direction,
			"currency":  payment.Currency,
		}

		if payment.Direction == entities.Outgoing {
			element["to_account"] = payment.Counterparty.Name
		} else {
			element["from_account"] = payment.Counterparty.Name
		}

		elements[index] = element
	}

	result := map[string]interface{}{
		"payments": elements,
	}

	return json.Marshal(result)
}
