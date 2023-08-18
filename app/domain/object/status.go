package object

import ()

type (
	Status struct {
		// The internal ID of the account
		AccountId AccountID `json:"-"`

		// The content of the status
		Content string `json:"content"`

		// The time the status was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}
)
