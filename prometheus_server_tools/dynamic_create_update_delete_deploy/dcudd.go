package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"os"
	"path/filepath"
)

func homeDir() string{
	if h := os.Getenv("HOME"); h!= ""{
		return h
	}
	return os.Getenv("USERPROFILE") // windows

}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}

func main() {
	var kubeconfig *string
	if home := homeDir(); home != ""{
		kubeconfig = flag.String("kubeconfig",filepath.Join(home,".kube","config"),"(optional) absolute path to the kubeconfig file")
	}else{
		kubeconfig = flag.String("kubeconfig","","absolute path to the kubeconfig file")
	}
	flag.Parse()

	namespace := "default"
	config,err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err != nil{
		panic(err)
	}
	client,err := dynamic.NewForConfig(config)
	if err != nil{
		panic(err)
	}
	deploymentRes := schema.GroupVersionResource{Group: "apps",Version: "v1",Resource: "deployments"}
	deployment := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "apps/v1",
			"kind": "Deployment",
			"metadata": map[string]interface{}{
				"name":  "demo-deployment",
			},
			"spec": map[string]interface{}{
				"replicas": 2,
				"selector": map[string]interface{}{
					"matchLabels": map[string]interface{}{
						"app": "demo",
					},
				},
				"template": map[string]interface{}{
					"metadata": map[string]interface{}{
						"labels": map[string]interface{}{
							"app": "demo",
						},
					},
					"spec": map[string]interface{}{
						"containers": []map[string]interface{}{
							{
								"name":  "webdemo",
								"image": "nginx:1.12",
								"ports": []map[string]interface{}{
									{
										"name":          "http",
										"protocol":      "TCP",
										"containerPort": 80,
									},
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
	result,err := client.Resource(deploymentRes).Namespace(namespace).Create(context.TODO(),deployment,metav1.CreateOptions{})
	if err != nil{
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetName())
	// Update Deployment
	prompt()
	fmt.Println("Updating deployment...")
	//    You have two options to Update() this Deployment:
	//
	//    1. Modify the "deployment" variable and call: Update(deployment).
	//       This works like the "kubectl replace" command and it overwrites/loses changes
	//       made by other clients between you Create() and Update() the object.
	//    2. Modify the "result" returned by Get() and retry Update(result) until
	//       you no longer get a conflict error. This way, you can preserve changes made
	//       by other clients between Create() and Update(). This is implemented below
	//			 using the retry utility package included with client-go. (RECOMMENDED)
	//
	// More Info:
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result,getErr := client.Resource(deploymentRes).Namespace(namespace).Get(context.TODO(),"demo-deployment",metav1.GetOptions{})
		if getErr != nil{
			panic(fmt.Errorf("failed to get latest version of Deployment: %v", getErr))
		}

		// update replicas to 1
		if err := unstructured.SetNestedField(result.Object,int64(1),"spec","replicas");err != nil{
			panic(fmt.Errorf("failed to set replica value: %v", err))
		}

		// extract spec containers
		containers,found,err := unstructured.NestedSlice(result.Object,"spec","template","spec","containers")
		if err != nil || !found || containers == nil{
			panic(fmt.Errorf("deployment containers not found or error in spec: %v", err))
		}
		// update container[0] image
		if err := unstructured.SetNestedField(containers[0].(map[string]interface{}),"nginx:1.13","image");err != nil{
			panic(err)
		}
		if err := unstructured.SetNestedField(result.Object,containers,"spec","template","spec","containers");err != nil{
			panic(err)
		}
		_,updateErr := client.Resource(deploymentRes).Namespace(namespace).Update(context.TODO(),result,metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil{
		panic(fmt.Errorf("update failed: %v", retryErr))
	}
	fmt.Println("Updated deployment...")
	// List Deployments
	prompt()
	fmt.Printf("Listing deployments in namespace %q:\n", apiv1.NamespaceDefault)
	list,err := client.Resource(deploymentRes).Namespace(namespace).List(context.TODO(),metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _,d := range list.Items{
		replicas,found,err := unstructured.NestedInt64(d.Object,"spec","replicas")
		if err != nil || found{
			fmt.Printf("Replicas not found for deployment %s: error=%s", d.GetName(), err)
			continue
		}
		fmt.Printf(" * %s (%d replicas)\n", d.GetName(), replicas)
	}
	// Delete Deployment
	prompt()
	fmt.Println("Deleting deployment...")

	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := client.Resource(deploymentRes).Namespace(namespace).Delete(context.TODO(),"demo-deployment", deleteOptions);err != nil{
		panic(err)
	}
	fmt.Println("Deleted deployment.")
}
