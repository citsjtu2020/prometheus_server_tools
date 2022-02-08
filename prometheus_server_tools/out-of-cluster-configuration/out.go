package main

import (
	"context"
	"flag"
	"fmt"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

func homeDir() string{
	if h := os.Getenv("HOME"); h!= ""{
		return h
	}
	return os.Getenv("USERPROFILE") // windows

}

func main() {
	var kubeconfig *string
	if home := homeDir();home != ""{
		kubeconfig = flag.String("kubeconfig",filepath.Join(home,".kube","config"), "(optional) absolute path to the kubeconfig file")
	}else{
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	// use the current context in kubeconfig
	config,err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err != nil{
		panic(err.Error())
	}

	// create the clientset
	clientset,err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err.Error())
	}
	for{
		pods,err := clientset.CoreV1().Pods("").List(context.TODO(),metav1.ListOptions{})
		if err != nil{
			panic(err.Error())
		}
		ns,err := clientset.CoreV1().Namespaces().List(context.TODO(),metav1.ListOptions{})
		if err != nil{
			panic(err.Error())
		}
		for _,n := range ns.Items{
			fmt.Printf("the namespace is: %s\n", n.Name)
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
		// Examples for error handling:
		// - Use helper functions like e.g. errors.IsNotFound()
		// - And/or cast to StatusError and use its properties like e.g. ErrStatus.Message
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
		time.Sleep(time.Second*10)



	}

}

