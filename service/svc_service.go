package service

import (
	"context"
	"fmt"

	"github.com/jary-287/gopass-svc/model"
	"github.com/jary-287/gopass-svc/proto/svc"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type ISvcService interface {
	GetAllSvc() ([]model.Svc, error)
	GetSvcById(uint64) (*model.Svc, error)
	AddSvc(*model.Svc) (uint64, error)
	UpdateSvc(*model.Svc) error
	DeleteSvc(uint64) error
	CreateToK8s(*svc.SvcInfo) error
	UpdateToK8s(*svc.SvcInfo) error
	DeleteFromK8s(*svc.SvcInfo) error
}

type SvcService struct {
	SvcRegistry model.ISvcRegistry
	K8sClient   *kubernetes.Clientset
	Service     *v1.Service
}

// AddSvc implements ISvcService
func (ss *SvcService) AddSvc(SvcInfo *model.Svc) (uint64, error) {
	return ss.SvcRegistry.CreateSvc(SvcInfo)
}

// CreateToK8s implements ISvcService
func (ss *SvcService) CreateToK8s(svcInfo *svc.SvcInfo) error {
	if _, err := ss.K8sClient.CoreV1().Services(svcInfo.SvcNamespace).
		Get(context.TODO(), svcInfo.SvcName, metav1.GetOptions{}); err != nil {
		s := ss.SetService(svcInfo)
		if _, err := ss.K8sClient.CoreV1().Services(s.Namespace).
			Create(context.TODO(), s, metav1.CreateOptions{}); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("svc have exsit:%s", svcInfo.SvcName)
	}
	return nil
}

// DeleteFromK8s implements ISvcService
func (ss *SvcService) DeleteFromK8s(svcInfo *svc.SvcInfo) error {
	if _, err := ss.K8sClient.CoreV1().Services(svcInfo.SvcNamespace).
		Get(context.TODO(), svcInfo.SvcName, metav1.GetOptions{}); err != nil {
		return fmt.Errorf("svc don`t exsit:%s", svcInfo.SvcName)
	} else {
		if err := ss.K8sClient.CoreV1().Services(svcInfo.SvcNamespace).
			Delete(context.TODO(), svcInfo.SvcName, metav1.DeleteOptions{}); err != nil {
			return err
		}
	}
	return nil
}

// DeleteSvc implements ISvcService
func (ss *SvcService) DeleteSvc(id uint64) error {
	return ss.SvcRegistry.DeleteSvc(id)
}

// GetAllSvc implements ISvcService
func (ss *SvcService) GetAllSvc() ([]model.Svc, error) {
	return ss.SvcRegistry.GetSvc()

}

// GetSvcById implements ISvcService
func (ss *SvcService) GetSvcById(id uint64) (*model.Svc, error) {
	return ss.SvcRegistry.GetSvcByID(id)
}

// UpdateSvc implements ISvcService
func (ss *SvcService) UpdateSvc(svcInfo *model.Svc) error {
	return ss.SvcRegistry.UpdateSvc(svcInfo)
}

// UpdateToK8s implements ISvcService
func (ss *SvcService) UpdateToK8s(svcInfo *svc.SvcInfo) error {
	if _, err := ss.K8sClient.CoreV1().Services(svcInfo.SvcNamespace).
		Get(context.TODO(), svcInfo.SvcName, metav1.GetOptions{}); err != nil {
		return fmt.Errorf("svc don`t exsit:%s", svcInfo.SvcName)
	} else {
		s := ss.SetService(svcInfo)
		if _, err := ss.K8sClient.CoreV1().Services(s.Namespace).
			Update(context.TODO(), s, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func NewSvcService(svcRegistry model.ISvcRegistry, client *kubernetes.Clientset) ISvcService {
	return &SvcService{
		SvcRegistry: svcRegistry,
		K8sClient:   client,
		Service:     &v1.Service{},
	}
}

//TODO
func (ss SvcService) SetService(svcInfo *svc.SvcInfo) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcInfo.SvcName,
			Namespace: svcInfo.SvcNamespace,
			Labels: map[string]string{
				"app": "pass-svc",
			},
			Annotations: map[string]string{
				"app":    "pass-svc",
				"author": "ljw",
			},
		},
		Spec: v1.ServiceSpec{
			Ports:          GetSvcPort(svcInfo),
			Selector:       svcInfo.Selector,
			Type:           v1.ServiceType(svcInfo.SvcType),
			LoadBalancerIP: svcInfo.GetLoadBanlancerIp(),
			ExternalName:   svcInfo.GetExternalName(),
			ClusterIP:      svcInfo.GetClusterIp(),
		},
	}
}

func GetSvcPort(info *svc.SvcInfo) (ports []v1.ServicePort) {
	for _, svcPort := range info.Ports {
		ports = append(ports, v1.ServicePort{
			Protocol:   v1.Protocol(svcPort.Protocol),
			Port:       svcPort.Port,
			TargetPort: intstr.FromInt(int(svcPort.TargetPort)),
			NodePort:   svcPort.NodePort,
		})
	}
	return
}
