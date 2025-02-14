package models

type ActivityPubEvent struct {
	Server string         `json:"server"`
	HookID string         `json:"hookId"`
	UserID string         `json:"userId"`
	Type   string         `json:"type"`
	Body   map[string]any `json:"body"`
}
