package object

import ()

type (
	Relationship struct {
		// The internal ID of the status
		ID AccountID `json:"-"`

		// The internal following_id of the account
		FollowingId AccountID `json:"following_id" db:"following_id"`

		// The internal follower_id of the account
		FollowerId AccountID `json:"follower_id" db:"follower_id"`

		// The time the relationship was created
		CreateAt DateTime `json:"create_at,omitempty" db:"create_at"`
	}
)
