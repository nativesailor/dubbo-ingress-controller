package main

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"flag"
	"fmt"
	"github/monkeyWie/dubbo-ingress-controller/pkg/controller"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	"os"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK\n"))
	if err != nil {
		return
	}
}

func main() {

	//创建一个http接口作为探针接口
	go func() {
		http.HandleFunc("/health", healthHandler)
		fmt.Print("Listening on 8080")
		err := http.ListenAndServe("8080", nil)
		if err != nil {
			return
		}
	}()

	host := os.Getenv("KUBERNETES_SERVICE_HOST")
	port := os.Getenv("KUBERNETES_SERVICE_PORT")

	var cfg *rest.Config

	if len(host) == 0 || len(port) == 0 {
		// 通过命令行参数获取kubeconfig文件路径
		kubeconfig := flag.String("kubeconfig", os.Getenv("HOME")+"/.kube/config", "path to the kubeconfig file")
		flag.Parse()

		// 使用clientcmd来加载kubeconfig文件并返回一个rest.Config对象
		var err error
		cfg, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			fmt.Printf("Error building kubeconfig: %v\n", err)
			os.Exit(1)
		}

	} else {
		var err error
		cfg, err = rest.InClusterConfig()
		if err != nil {
			logger.Fatal(err)
		}
	}

	controller, err := controller.NewController(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	if err := controller.Start(); err != nil {
		logger.Fatal(err)
	}

	// 等待程序终止
	select {}
}
