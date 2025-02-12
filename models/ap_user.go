package models

type ActivityPubUser struct {
	ID          string              `json:"id"`
	AvatarURL   string              `json:"avatarUrl"`
	Host        string              `json:"host"`
	Name        string              `json:"name"`
	UserName    string              `json:"username"`
	Description string              `json:"description"`
	Instance    ActivityPubInstance `json:"instance"`
}
