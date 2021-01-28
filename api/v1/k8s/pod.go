package k8s

import (
	"cmdb/utils"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"time"
)

var ClientSet *kubernetes.Clientset

func InitConfig(c *gin.Context) {
	var kubeconfig *string
	_, err := os.Lstat(utils.KubeFile)
	if err != nil {
		fmt.Sprintf("no such %s ", utils.KubeFile)
		panic("1")
	}
	kubeconfig = flag.String("kubeconfig", utils.KubeFile, "absolute path to the kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	ClientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

func ListPod(c *gin.Context) {
	for {
		// 通过实现 clientset 的 CoreV1Interface 接口列表中的 PodsGetter 接口方法 Pods(namespace string)返回 PodInterface
		// PodInterface 接口拥有操作 Pod 资源的方法，例如 Create、Update、Get、List 等方法
		// 注意：Pods() 方法中 namespace 不指定则获取 Cluster 所有 Pod 列表
		pods, err := ClientSet.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the k8s cluster\n", len(pods.Items))

		// 获取指定 namespace 中的 Pod 列表信息
		namespace := "default"
		pods, err = ClientSet.CoreV1().Pods(namespace).List(metav1.ListOptions{})
		if err != nil {
			panic(err)
		}
		fmt.Printf("\nThere are %d pods in namespaces %s\n", len(pods.Items), namespace)
		for _, pod := range pods.Items {
			fmt.Printf("Name: %s, Status: %s, CreateTime: %s\n", pod.ObjectMeta.Name, pod.Status.Phase, pod.ObjectMeta.CreationTimestamp)
		}
		time.Sleep(10 * time.Second)
	}
}

func ListDeployment(c *gin.Context) {

}
func ListStatefulSet(c *gin.Context) {

}
func ListCronJob(c *gin.Context) {

}
func ListJob(c *gin.Context) {

}
