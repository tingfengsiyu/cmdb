package k8s

import (
	"cmdb/utils"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

type clientconfig struct {
	clientset *kubernetes.Clientset
}

var Client = clientconfig{}

func int32Ptr(i int32) *int32 { return &i }

func Initk8s() {
	kubeconfig := utils.KubeFile

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	Client = clientconfig{
		clientset,
	}
	if err != nil {
		panic(err.Error())
	}
}

func (c *clientconfig) CreateDeployment(namespace, deploymentName, imageName string, replicas int32) {
	deploymentsClient := c.clientset.AppsV1().Deployments(namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "web",
							Image: imageName,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

func (c *clientconfig) DeleteDeployment(namespace, deploymentName string) {
	deploymentsClient := c.clientset.AppsV1().Deployments(namespace)
	fmt.Println("Deleting deployment...")
	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), deploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}

func (c *clientconfig) UpdateDeployment(namespace, deploymentName, imageName string, replicas int32) {
	deploymentsClient := c.clientset.AppsV1().Deployments(namespace)
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentsClient.Get(context.TODO(), deploymentName, metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(replicas)                 // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = imageName // change nginx version
		_, updateErr := deploymentsClient.Update(context.TODO(), result, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")

}

func (c *clientconfig) ListDeployment(namespace string) {
	deploymentList, _ := c.clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
	for _, deployment := range deploymentList.Items {
		fmt.Println(deployment.Name, deployment.Namespace, deployment.Status.Replicas, deployment.Spec.Template.Spec.Containers[0].Image)
	}
}

func (c *clientconfig) ListPod(namespace string) {
	pods, err := c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name, pod.Namespace, pod.Status.StartTime, pod.Status.Phase)
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func (c *clientconfig) ListNode() {

	nodeList, err := c.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	for _, node := range nodeList.Items {
		fmt.Printf("%s\t%s\t%s\t%s\n", node.Name, node.Status.Addresses, node.Status.NodeInfo.Architecture)
	}
}

func (c *clientconfig) ListNamespace() {
	namespaceList, _ := c.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	for _, namespace := range namespaceList.Items {
		fmt.Println(namespace.Name, namespace.GetCreationTimestamp())
	}

}

func (c *clientconfig) GetPod(namespace, pod string) {
	_, err := c.clientset.CoreV1().Pods(namespace).Get(context.TODO(), pod, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
	} else if statusError, isStatus := err.(*errors.StatusError); isStatus {
		fmt.Printf("Error getting pod %s in namespace %s: %v\n",
			pod, namespace, statusError.ErrStatus.Message)
	} else if err != nil {
		panic(err.Error())
	} else {
		fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
	}
}

func a(c *gin.Context) {
	Client.ListPod("namespace")
}

//func ListStatefulSet(c *gin.Context) {
//
//}
//func ListCronJob(c *gin.Context) {
//
//}
//func ListJob(c *gin.Context) {
//
//}
