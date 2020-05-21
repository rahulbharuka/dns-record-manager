package model

// Server is a model for a subdomain server.
type Server struct {
	ID              uint64 `json:"id"`
	Name            string `json:"name"`
	ClusterID       uint64 `json:"cluster_id"`
	IP              string `json:"ip"`
	AddedToRotation bool   `json:"added_to_rotation"`
	ClusterName     string `json:"cluster_name"`
	Subdomain       string `json:"subdomain"`
}
