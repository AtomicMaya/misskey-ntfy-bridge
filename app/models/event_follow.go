package models

type ActivityPubFollowEvent struct {
	User ActivityPubUser `json:"user"`
}
