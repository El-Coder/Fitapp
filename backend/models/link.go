package models

type Link struct {
	FitID    string `json:"fit_id"`
	ItemID   string `json:"item_id"`
	ItemName string `json:"item_name,omitempty"`
}

type LinkRequest struct {
	FitID  string `json:"fit_id"`
	ItemID string `json:"item_id"`
}
