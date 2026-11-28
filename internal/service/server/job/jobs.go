package job

type CreateMinigameServer struct {
	ServerID string `json:"server_id"`
}

func (CreateMinigameServer) Kind() string {
	return "create_minigame_server"
}
