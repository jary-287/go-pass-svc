package main

import (
	"log"
	"time"

	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/go-micro/plugins/v3/registry/consul"
	"github.com/go-micro/plugins/v3/wrapper/breaker/hystrix"
	limiter "github.com/go-micro/plugins/v3/wrapper/ratelimiter/uber"
	opentracing2 "github.com/go-micro/plugins/v3/wrapper/trace/opentracing"
	"github.com/jary-287/gopass-common/common"
	"github.com/jary-287/gopass-svc/handle"
	"github.com/jary-287/gopass-svc/model"
	"github.com/jary-287/gopass-svc/proto/svc"
	"github.com/jary-287/gopass-svc/service"
)

func main() {
	// var kubeconfig *string
	// if home := homedir.HomeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", path.Join(home, ".kube", "config"), "kubeconfig 位置")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "kubeconfig 位置")
	// }
	// flag.Parse()
	// log.Println(*kubeconfig)
	// //获取kubernets实例
	client, err := common.GetKubernetsClient("/Users/liujuwen/.kube/config")
	if err != nil {
		log.Fatal("get kubernetes client err:", err)
	}
	//注册中心
	consulRegister := consul.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{"192.168.0.19:8500"}
		o.Timeout = 20 * time.Second

	})
	t, io, err := common.NewTracer("service.pod", ":9333")
	if err != nil {
		log.Fatal(err)
	}

	defer io.Close()
	// 创建pod服务
	serv := micro.NewService(
		micro.Name("service.svc"),
		micro.Version("latest"),
		//注册中心
		micro.Address(":8888"),
		micro.Registry(consulRegister),
		//链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(t)),
		micro.WrapClient(opentracing2.NewClientWrapper(t)),
		//熔断
		micro.WrapClient(hystrix.NewClientWrapper()),
		micro.WrapHandler(limiter.NewHandlerWrapper(1000)),
	)
	serv.Init()
	if err := model.NewSvcRegistry(model.Db).InitTable(); err != nil {
		log.Fatal(err)
	}
	//注册句柄
	svcService := service.NewSvcService(model.NewSvcRegistry(model.Db), client)
	svc.RegisterSvcHandler(serv.Server(), &handle.SvcHandler{SvcService: svcService})
	if err := serv.Run(); err != nil {
		log.Fatal(err)
	}
}
