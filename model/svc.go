package model

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Svc struct {
	SvcID        uint64 `gorm:"primaryKey;AUTO_INCREMENT" json:"svc_id"`
	SvcName      string `gorm:"not null;unique" json:"svc_name"`
	SvcNamespace string `gorm:"default:'default'" json:"svc_namespace"`
	//type should be ClusterIP,ExternalName,LoadBalancer,NodePort
	SvcType        string         `gorm:"default:'ClusterIP'" json:"svc_type"`
	SvcTeamID      uint64         `json:"svc_team_id"`
	Ports          []SvcPort      `gorm:"foreignKey:svc_id;references:svc_id" json:"ports"`
	Selector       datatypes.JSON `json:"selector"`
	LoadBalancerIP string         `json:"load_balancer_ip"`
	ExternalName   string         `json:"external_name"`
	ClusterIP      string         `json:"cluster_ip"`
}

type ISvcRegistry interface {
	InitTable() error
	CreateSvc(*Svc) (uint64, error)
	UpdateSvc(*Svc) error
	DeleteSvc(id uint64) error
	GetSvc() ([]Svc, error)
	GetSvcByID(id uint64) (*Svc, error)
}

type SvcRegistry struct {
	db *gorm.DB
}

func (sr *SvcRegistry) CreateSvc(svcInfo *Svc) (uint64, error) {
	err := sr.db.Create(svcInfo).Error
	return svcInfo.SvcID, err
}

// DeleteSvc implements ISvcRegistry
func (sr *SvcRegistry) DeleteSvc(id uint64) error {
	tx := sr.db.Begin()
	tx.Where("svc_id=?", id).Delete(&SvcPort{})
	tx.Where("svc_id=?", id).Delete(&Svc{})
	if err := tx.Commit().Error; err != nil {
		tx.Callback()
		return err
	}
	return nil
}

// GetSvc implements ISvcRegistry
func (sr *SvcRegistry) GetSvc() (svcs []Svc, err error) {
	err = sr.db.Preload("SvcPorts").Find(&svcs).Error
	return
}

// GetSvcByID implements ISvcRegistry
func (sr *SvcRegistry) GetSvcByID(id uint64) (svcInfo *Svc, err error) {
	err = sr.db.Where("svc_id=?", id).First(svcInfo).Error
	return
}

func (sr *SvcRegistry) InitTable() error {
	return sr.db.AutoMigrate(&Svc{}, &SvcPort{})
}

// UpdateSvc implements ISvcRegistry
func (sr *SvcRegistry) UpdateSvc(svcInfo *Svc) error {
	tx := sr.db.Begin()
	tx.Where("svc_id = ?", svcInfo.SvcID).Save(&SvcPort{})
	tx.Save(svcInfo)
	if err := tx.Commit().Error; err != nil {
		tx.Callback()
		return err
	}
	return nil
}

func NewSvcRegistry(db *gorm.DB) ISvcRegistry {
	return &SvcRegistry{
		db: db,
	}
}
