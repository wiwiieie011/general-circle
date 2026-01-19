package dto

type PublicUserResponse struct {
	ID          uint   `json:"id"`
	FullName    string `json:"full_name"`
	IsOrganizer bool   `json:"is_organizer"`
}
