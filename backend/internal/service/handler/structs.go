package handler

type messageResponse struct {
	IsMine bool   `json:"isMine"`
	Text   string `json:"text"`
}

type send struct {
	Text string `json:"text"`
}

type auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
