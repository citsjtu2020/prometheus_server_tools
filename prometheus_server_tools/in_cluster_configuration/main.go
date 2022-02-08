package main

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"time"
)

func main() {
	config,err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset,err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for{
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods,err := clientset.CoreV1().Pods("").List(context.TODO(),metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		// Examples for error handling:
		// - Use helper functions e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Messag
		namespace := "micro"
		pod := "geo-6667f69d8c-jghvb"
		_,err = clientset.CoreV1().Pods(namespace).Get(context.TODO(),pod,metav1.GetOptions{})
		if errors.IsNotFound(err){
			fmt.Printf("Pod %s in namespace %s not found\n", pod, namespace)
		} else if statusError,isStatus := err.(*errors.StatusError);isStatus{
			fmt.Printf("Error getting pod %s in namespace %s: %v\n",
				pod, namespace, statusError.ErrStatus.Message)
		}else if err != nil{
			panic(err.Error())
		}else{
			fmt.Printf("Found pod %s in namespace %s\n", pod, namespace)
		}

		time.Sleep(10 * time.Second)
	}
}
