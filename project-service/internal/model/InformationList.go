package model

type InformationList struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Key       string `json:"key"`
	Value     string `json:"value"`
}
