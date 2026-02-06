package game

type Question struct {
	ID     int    `json:"id"`
	Text   string `json:"text"`
	Answer string `json:"answer"`
	Hint   string `json:"hint"`
}

type Config struct {
	Questions    []Question `json:"questions"`
	FinalMessage string     `json:"final_message"`
	FinalHint    string     `json:"final_hint"`
}
