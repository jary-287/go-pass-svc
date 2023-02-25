package handle

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/jary-287/gopass-svc/model"
	"github.com/jary-287/gopass-svc/proto/svc"
	"github.com/jary-287/gopass-svc/service"
)

//完成porto中的方法
type SvcHandler struct {
	SvcService service.ISvcService
}

//第二个值传指针
func Swap(source, target interface{}) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &target)
}

func (sh *SvcHandler) AddSvc(ctx context.Context, info *svc.SvcInfo, rsp *svc.Response) error {
	log.Println("start add svc:", info.SvcName)
	svcModel := &model.Svc{}
	if err := Swap(info, svcModel); err != nil {
		rsp.Msg = "swap to model svc failed"
		return err
	}
	log.Println("swap over:", info, svcModel)
	if err := sh.SvcService.CreateToK8s(info); err != nil {
		rsp.Msg = "crate to k8s failed"
		return err
	}
	log.Println("create svc to k8s success:", info.SvcName)
	if id, err := sh.SvcService.AddSvc(svcModel); err != nil {
		rsp.Msg = "add svc model failed"
		return err
	} else {
		rsp.Msg = strconv.Itoa(int(id))
	}
	log.Println("add svc success:", info.SvcName)
	return nil
}

func (sh *SvcHandler) FindAllSvc(ctx context.Context, findAll *svc.FindAll, allSvc *svc.AllSvc) error {
	svcs, err := sh.SvcService.GetAllSvc()
	if err != nil {
		return err
	}
	if err := Swap(svcs, &allSvc.SvcInfo); err != nil {
		return err
	}
	log.Println("find all svc success")
	return nil
}

func (sh *SvcHandler) DeleteSvc(ctx context.Context, SvcInfo *svc.SvcInfo, rsp *svc.Response) error {
	log.Println("start delete svc:", SvcInfo.SvcName)
	log.Println(SvcInfo)
	if err := sh.SvcService.DeleteFromK8s(SvcInfo); err != nil {
		rsp.Msg = "delete from  k8s failed"
		return err
	}
	log.Println("delete svc to k8s success:", SvcInfo.SvcName)
	if err := sh.SvcService.DeleteSvc(SvcInfo.SvcId); err != nil {
		rsp.Msg = "delete svc model failed"
		return err
	}
	log.Println("delete svc success:", SvcInfo.SvcName)
	return nil

}

func (sh *SvcHandler) UpdateSvc(ctx context.Context, SvcInfo *svc.SvcInfo, rsp *svc.Response) error {
	log.Println("start update svc:", SvcInfo.SvcName)
	svcModel := &model.Svc{}
	if err := Swap(SvcInfo, svcModel); err != nil {
		rsp.Msg = "swap to model svc failed"
		return err
	}
	if err := sh.SvcService.UpdateToK8s(SvcInfo); err != nil {
		rsp.Msg = "update from  k8s failed"
		return err
	}
	log.Println("update svc to k8s success:", SvcInfo.SvcName)
	if err := sh.SvcService.UpdateSvc(svcModel); err != nil {
		rsp.Msg = "update svc model failed"
		return err
	}
	log.Println("update svc success:", SvcInfo.SvcName)
	return nil
}
func (sh *SvcHandler) FindSvcById(ctx context.Context, SvcId *svc.SvcId, svcInfo *svc.SvcInfo) error {
	if svc, err := sh.SvcService.GetSvcById(SvcId.Id); err != nil {
		log.Println("find svc by id failed:", err)
		return err
	} else {

		if err := Swap(svc, svcInfo); err != nil {
			log.Println("find svc by id swap  failed:", err)
			return err
		}
	}
	return nil
}
