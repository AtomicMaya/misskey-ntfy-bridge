package models

type ActivityPubNote struct {
	ID             string           `json:"id"`
	ContentWarning string           `json:"cw"`
	Text           string           `json:"text"`
	User           ActivityPubUser  `json:"user"`
	Reply          *ActivityPubNote `json:"reply,omitempty"`
	Renote         *ActivityPubNote `json:"renote,omitempty"`
}
