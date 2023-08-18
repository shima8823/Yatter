package object

import ()

type (
	Status struct {
		// The internal ID of the status
		ID AccountID `json:"-"`

		// The internal ID of the account
		AccountId AccountID `json:"account_id" db:"account_id"`

		// The content of the status
		Content string `json:"content"`

		// The time the status was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}
)
