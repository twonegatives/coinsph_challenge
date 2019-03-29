package banking

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/twonegatives/coinsph_challenge/pkg/entities"
)

type paymentsJSONEncoder struct {
}

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
		}

		if payment.Direction == entities.Outgoing {
			element["to_account"] = payment.Participant.Name
		} else {
			element["from_account"] = payment.Participant.Name
		}

		elements[index] = element
	}

	result := map[string]interface{}{
		"payments": elements,
	}

	return json.Marshal(result)
}
