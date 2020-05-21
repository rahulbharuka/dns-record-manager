package model

// DNSRecord is a model for a server's DNS record.
type DNSRecord struct {
	FQDN        string `json:"fqdn"`
	IP          string `json:"ip"`
	ServerName  string `json:"server_name"`
	ClusterName string `json:"cluster_name"`
}
