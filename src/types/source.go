package types

type AccessRecord struct {
	ExternalId int    `json:"externalId"`
	Run        string `json:"run"`
	FullName   string `json:"fullName"`
	EntryAt    string `json:"entryAt"`
	ExitAt     string `json:"exitAt"`
	Activity   string `json:"activity"`
	Location   string `json:"location"`
}

type AccessResponse struct {
	Message string `json:"message"`
	Data    struct {
		Records []AccessRecord `json:"records"`
	} `json:"data"`
}
