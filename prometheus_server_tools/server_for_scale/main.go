package main

import (
	"context"
	"flag"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func homeDir() string{
	if h := os.Getenv("HOME"); h!= ""{
		return h
	}
	return os.Getenv("USERPROFILE") // windows

}

//func int2point(i int32) *int32{
//	return &i
//}
//sync.Ma


type MessageController struct {
	//dd sync.Map[net.Conn] chan string
	//connqueue []net.Conn
	serverReq map[net.Conn]chan string
	serverRes map[net.Conn]chan int
	stopChan map[net.Conn]chan struct{}
	//broadChan map[net.Conn]chan int
	quitconn chan net.Conn
	rlock *sync.RWMutex
}

//var GlobalStopChan chan struct{}

//func init() {
//	GlobalStopChan = make(chan struct{},1)
//}

//var Message map[net.Conn]chan struct{}
func HandleScaleRequest(c net.Conn,mc *MessageController){
	defer func() {
		if p := recover(); p!= nil{
			fmt.Printf("Recovered panic: %s\n",p)
		}
		fmt.Println("quit the handler")
	}()
	fmt.Println("Start the handler")
	buf := make([]byte,1024)
	//defer c.Close()
	//defer mc.rlock.RUnlock()
	Loop:
		for{
			mc.rlock.RLock()
			tmpStopChan,ok1 := mc.stopChan[c]
			tmpserverRes,ok2 := mc.serverRes[c]
			tmpserverReq,ok3 := mc.serverReq[c]
			//tmpboardChan,ok4 := mc.broadChan[c]
			if !ok1 || !ok2 || !ok3{
				mc.rlock.RUnlock()
				break Loop
			}
			select {
				case <-tmpStopChan:
					fmt.Println("server has ended.")
					mc.rlock.RUnlock()
					break Loop
				//case e,ok := <- tmpboardChan:
				//	if !ok{
				//		fmt.Fprintf(os.Stdout,"End.")
				//		mc.rlock.RUnlock()
				//		break Loop
				//	}
				//	if e == 1{
				//		mc.rlock.RUnlock()
				//		time.Sleep(2*time.Second)
				//	}
				case e,ok := <-tmpserverRes:
					fmt.Println("response case")
					if !ok{
						fmt.Fprintf(os.Stdout,"End.")
						mc.rlock.RUnlock()
						break Loop
					}
					fmt.Println("start to send response to client.")
					c.Write([]byte(strconv.Itoa(e)))
					fmt.Println("send back to client")
					mc.rlock.RUnlock()
				default:
					n,err := c.Read(buf)
					if err != nil{
						//res := "quit"
						//tmpserverReq
						close(mc.serverRes[c])
						fmt.Fprintf(os.Stderr,err.Error())
						mc.rlock.RUnlock()
						//break Loop
						//log.Fatal(err.Error())

					}
					line := strings.Trim(string(buf[0:n])," \r\n\t")
					if line == "@"{
						mc.rlock.RUnlock()
						continue
					}
					fmt.Println("default case")
					fmt.Println(line)
					lines := strings.Split(line,":")
					line2 := strings.ToLower(strings.Trim(lines[1]," \r\n\t"))
					if line2 == "exit" {
						fmt.Printf("user %s quit the client\n",lines[0])
						//tmpserverReq <- "quit"
						close(mc.serverRes[c])
						mc.rlock.RUnlock()
						//break Loop
					}
					tmpserverReq <- line
					mc.rlock.RUnlock()
			}
		}
}

func int2point(i int32) *int32{
	return &i
}

func HandleUpdate(config *rest.Config,mc *MessageController,stopCh chan struct{}){
	defer func() {
		if p := recover(); p!= nil{
			fmt.Printf("Recovered panic: %s\n",p)
		}
	}()
	clientset,err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err)
	}
	Loop:
		for{
		select {
			case <-stopCh:
				break Loop
			default:
				mc.rlock.RLock()
				for conn,_ := range mc.serverReq{
					select {
					case e,ok := <- mc.serverReq[conn]:
						if !ok{
							mc.serverRes[conn] <- 2
							close(mc.serverRes[conn])
							mc.quitconn <- conn
							continue
						}
						requests0 := strings.Split(e,":")
						user := strings.Trim(requests0[0]," \r\n\t")
						requests := strings.Split(requests0[1],"@")
						//if strings.Trim(requests[0]," \r\n\t") == "quit"{
						//							mc.quitconn <- conn
						//							mc.serverRes[conn] <- 2
						//							continue
						//						} else
						if strings.ToLower(strings.Trim(requests[0]," \r\n\t")) == "create"{
							namespace := strings.Trim(requests[1]," \r\n\t")
							name := strings.Trim(requests[2]," \r\n\t")
							deployment := &appsv1.Deployment{
								ObjectMeta: metav1.ObjectMeta{
								Name: name,
								},
								Spec: appsv1.DeploymentSpec{
										Replicas: int2point(2),
										Selector: &metav1.LabelSelector{
										MatchLabels: map[string]string{
											"app": name,
											},
										},
										Template: apiv1.PodTemplateSpec{
											ObjectMeta: metav1.ObjectMeta{
											Labels: map[string]string{
												"app": name,
												},
											},
											Spec: apiv1.PodSpec{
												Containers: []apiv1.Container{
												{
													Name: name,
													Image:  "nginx:1.12",
													Ports: []apiv1.ContainerPort{
													{
														Name: "http",
														Protocol: apiv1.ProtocolTCP,
														ContainerPort: 80,
													},
												},
											},
										},
									},
								},
							},
							}
							list,listErr := clientset.CoreV1().Namespaces().List(context.TODO(),metav1.ListOptions{})
							if listErr != nil{
								fmt.Printf("%s: Create failed\n",user)
								mc.serverRes[conn] <- -1
								continue
							}
							ifCreate := true
							for _,ns := range list.Items{
								if ns.Name == namespace{
									ifCreate = false
									//mc.serverRes[conn] <- -1
									break
								}
							}
							if ifCreate{
								createns := &apiv1.Namespace{
									ObjectMeta: metav1.ObjectMeta{
										Namespace: namespace,
										Name: namespace,
									},
								}
								_,createErr := clientset.CoreV1().Namespaces().Create(context.TODO(),createns,metav1.CreateOptions{})
								if createErr != nil{
									fmt.Printf("%s: Create failed\n",user)
									mc.serverRes[conn] <- -1
									continue
								}
								//createResult.Name
								//clientset.CoreV1().Namespaces().Create()
							}
							depolymentClient := clientset.AppsV1().Deployments(namespace)
							result,createErr := depolymentClient.Create(context.TODO(),deployment,metav1.CreateOptions{})
							if createErr != nil{
								fmt.Printf("%s: Create failed\n",user)
								mc.serverRes[conn] <- -1
								continue
							}
							fmt.Printf("%s: Created deployment %q.\n", user,result.GetObjectMeta().GetName())
							mc.serverRes[conn] <- 1

						} else if strings.ToLower(strings.Trim(requests[0]," \r\n\t")) == "update"{
							namespace := strings.Trim(requests[1]," \r\n\t")
							name := strings.Trim(requests[2]," \r\n\t")
							num := strings.Trim(requests[3]," \r\n\t")
							fmt.Println("start to update")
							depolymentClient := clientset.AppsV1().Deployments(namespace)
							retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
								result,getErr := depolymentClient.Get(context.TODO(),name,metav1.GetOptions{})
								if getErr != nil{
									panic(fmt.Errorf("%s: Failed to get latest version of Deployment: %v", user,getErr))
								}
								//strconv.Atoi(num).(int32)
								scalesize, scaleerr := strconv.Atoi(num)
								if scaleerr != nil{
									panic(fmt.Errorf("%s: updated error because of scale size"))
									//mc.serverRes[conn] <- 1
									//continue
									//return scaleerr
								}
								scalesize32 := int32(scalesize)
								result.Spec.Replicas = int2point(scalesize32)  // reduce replica count
								//result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
								_,UpdateErr := depolymentClient.Update(context.TODO(),result,metav1.UpdateOptions{})
								return UpdateErr
								//UpdateScale()
							})
							if retryErr != nil{
							//panic(fmt.Errorf("Update failed: %v", retryErr))
								fmt.Printf("%s: Update failed: %v\n", user,retryErr)
								mc.serverRes[conn] <- -1
								continue
							}
							fmt.Printf("%s: Updated deployment %s...\n",user,name)
							mc.serverRes[conn] <- 1
							fmt.Println("Update finished")
						}else if strings.ToLower(strings.Trim(requests[0]," \r\n\t")) == "delete"{
							deletePolicy := metav1.DeletePropagationForeground
							namespace := strings.Trim(requests[1]," \r\n\t")
							name := strings.Trim(requests[2]," \r\n\t")
							depolymentClient := clientset.AppsV1().Deployments(namespace)

							if err := depolymentClient.Delete(context.TODO(),name,metav1.DeleteOptions{
								PropagationPolicy: &deletePolicy,
							});err != nil{
								//panic(err)
								fmt.Printf("%s: Delete failed: %v\n", user,err)
								mc.serverRes[conn] <- -1
								continue
							}
							fmt.Printf("%s Deleted deployment %s\n.",user,name)
							mc.serverRes[conn] <- 1
						}else if strings.ToLower(strings.Trim(requests[0]," \r\n\t")) == "list"{
							namespace := strings.Trim(requests[1]," \r\n\t")
							depolymentClient := clientset.AppsV1().Deployments(namespace)
							list,listErr := depolymentClient.List(context.TODO(),metav1.ListOptions{})
							if listErr != nil{
								fmt.Printf("%s: List failed: %v\n", user,err)
								mc.serverRes[conn] <- -1
								continue
							}
							for _,li := range list.Items{
								fmt.Printf("%s: namespace: %s,found deploy: %s\n",user,li.Namespace,li.Name)
							}
							mc.serverRes[conn] <- 1
						}else if strings.ToLower(strings.Trim(requests[0]," \r\n\t")) == "get"{
							namespace := strings.Trim(requests[1]," \r\n\t")
							name := strings.Trim(requests[2]," \r\n\t")
							depolymentClient := clientset.AppsV1().Deployments(namespace)
							getResult,getErr := depolymentClient.Get(context.TODO(),name,metav1.GetOptions{})
							if getErr != nil{
								fmt.Printf("%s: Get failed: %v\n",user,getErr)
								mc.serverRes[conn] <- -1
							}
							fmt.Printf("%s: Find %s in %s with %d replicas\n",user,getResult.Name,getResult.Namespace,*getResult.Spec.Replicas)
							mc.serverRes[conn] <- 1
						} else{
							fmt.Printf("%s: Uncontrolled instruction\n",user)
							mc.serverRes[conn] <- -1
						}
					default:
						continue
					}
				}
				mc.rlock.RUnlock()

		}
	}
	//depolymentClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

}

func HandleQuit(mc *MessageController,stopChan chan struct{}){
	defer func() {
		if p := recover(); p!= nil{
			fmt.Printf("Recovered panic: %s\n",p)
		}
	}()
	Loop:
		for{
			select {
			case <-stopChan:
				fmt.Println("quit the HandleQuit")
				break Loop
			case conn,ok := <-mc.quitconn:
				if !ok{
					fmt.Println("Controller is closed!")
					break Loop
				}
				time.Sleep(2*time.Second)
				mc.rlock.Lock()
				//delete(mc.serverRes, )
				_,ok2 := mc.serverRes[conn]
				if ok2{
					//close(chan2)
					delete(mc.serverRes,conn)
				}
				_,ok3 := mc.serverReq[conn]
				if ok3{
					delete(mc.serverReq,conn)
				}
				_,ok4 := mc.stopChan[conn]
				if ok4{
					close(mc.stopChan[conn])
					delete(mc.stopChan,conn)
				}
				mc.rlock.Unlock()
				conn.Close()
			}
		}
}

func InitMC() *MessageController{
	return &MessageController{
		quitconn: make(chan net.Conn,1),
		rlock: new(sync.RWMutex),
		//connqueue: make(),
		serverReq: make(map[net.Conn]chan string),
		serverRes: make(map[net.Conn]chan int),
		stopChan: make(map[net.Conn]chan struct{}),
	}
}



func main() {
	defer func() {
		if p := recover(); p!= nil{
			fmt.Printf("Recovered panic: %s\n",p)
		}
	}()
	var kubeconfig *string
	if home := homeDir(); home != ""{
		kubeconfig = flag.String("kubeconfig",filepath.Join(home,".kube","config"),"(optional) absolute path to the kubeconfig file")
	}else{
		kubeconfig = flag.String("kubeconfig","","absolute path to the kubeconfig file")
	}
	flag.Parse()
	l,err := net.Listen("tcp","127.0.0.1:28086")
	defer l.Close()
	if err != nil{
		fmt.Println("Start the server failed")
	}
	config,err := clientcmd.BuildConfigFromFlags("",*kubeconfig)
	if err != nil{
		panic(err)
	}
	stopQCH := make(chan struct{})
	defer close(stopQCH)
	stopHCH := make(chan struct{})
	defer close(stopHCH)

	//var rlock sync.RWMutex
	globalmc := InitMC()
	defer func() {
		close(globalmc.quitconn)
	}()
	defer func() {
		for conn,chans := range globalmc.serverReq{
			close(chans)
			chan2,ok := globalmc.serverRes[conn]
			if ok{
				select {
					case _,ee := <-chan2:
						if !ee{
							break
						}else{
							close(chan2)
							break
						}
				default:
					close(chan2)
				}
			}
			chan3,ok := globalmc.stopChan[conn]
			if ok{
				select {
					case _,ee := <-chan3:
						if !ee{
							break
						}else{
							close(chan3)
							break
						}
				default:
					close(chan3)
				}
			}
			_ = conn.Close()
		}
	}()



	go HandleUpdate(config,globalmc,stopHCH)
	go HandleQuit(globalmc,stopQCH)
	for{
		// need to do: receive and deal with the quit.
		conn,err := l.Accept()
		//&conn.Close()
		if err != nil{
			fmt.Println("Error accepting: ",err)
		}
		fmt.Printf("main server process: Received message %s -> %s \n",conn.RemoteAddr(),conn.LocalAddr())
		globalmc.rlock.Lock()
		globalmc.stopChan[conn] = make(chan struct{},1)
		globalmc.serverReq[conn] = make(chan string,1)
		globalmc.serverRes[conn] = make(chan int,1)
		globalmc.rlock.Unlock()
		fmt.Println("Add a client")
		go HandleScaleRequest(conn,globalmc)
	}
}
