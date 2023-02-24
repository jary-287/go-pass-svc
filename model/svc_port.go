package model

type SvcPort struct {
	ID         uint64 `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	SvcID      uint64 `json:"svc_id"`
	Protocol   string `json:"protocol"`
	Port       int    `json:"port"`
	TargetPort int    `json:"target_port"`
	NodePort   int    `json:"node_port"`
}
