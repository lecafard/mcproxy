package schemas

type ServerDescription struct {
	Text string `json:"text"`
}

type ServerPlayers struct {
	Maximum int `json:"max"`
	Online  int `json:"online"`
}

type ServerVersion struct {
	Name     string `json:"name"`
	Protocol int    `json:"protocol"`
}

type ServerStatus struct {
	Description ServerDescription `json:"description"`
	Players     ServerPlayers     `json:"players"`
	Version     ServerVersion     `json:"version"`
}
