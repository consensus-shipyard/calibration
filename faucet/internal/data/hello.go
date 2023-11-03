package data

type LivenessResponse struct {
	LotusVersion         string `json:"lotus_version"`
	Build                string `json:"build"`
	Epoch                uint64 `json:"epoch"`
	Behind               uint64 `json:"behind"`
	PeerNumber           int    `json:"peer_number"`
	Host                 string `json:"host"`
	PeersToPublishMsgs   int    `json:"peers_to_publish_msgs"`
	PeersToPublishBlocks int    `json:"peers_to_publish_blocks"`
	PeerID               string `json:"peer_id"`
	ServiceVersion       string `json:"service_version"`
}
