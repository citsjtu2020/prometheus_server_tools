package main

import (
	"context"
	"fmt"
	clientInflux "github.com/influxdata/influxdb1-client/v2"
	"github.com/prometheus/client_golang/api"
	apiV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	prometheusClient, err := api.NewClient(api.Config{
		Address: "http://localhost:30090",
	})
	if err != nil {
		log.Fatal("connected to prometheus failed")
	}
	v1api := apiV1.NewAPI(prometheusClient)
	ctx := context.Background()
	//r := apiV1.Range{
	//	Start: time.Unix(1669170480, 0),
	//	End:   time.Unix(1669171080, 0),
	//	Step:  time.Second * 5,
	//}

	c, err := clientInflux.NewHTTPClient(clientInflux.HTTPConfig{
		Addr:     "http://localhost:8086",
		Username: "root",
		Password: "kubernetes",
	})
	if err != nil {
		log.Fatal("connected to influxdb failed")
	}

	defer func(c clientInflux.Client) {
		err := c.Close()
		if err != nil {
			log.Fatal("close failed")
		}
	}(c)

	bp, _ := clientInflux.NewBatchPoints(clientInflux.BatchPointsConfig{
		Database:  "test",
		Precision: "s",
	})
	fmt.Println("start")
	timeNow := time.Now()
	var warnings apiV1.Warnings

	var podInfo model.Value
	podInfoMap := make(map[string]map[string]string)
	podInfo, warnings, err = v1api.Query(ctx, "avg_over_time(kube_pod_info[1m])", timeNow, apiV1.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatal("query kube_pod_info failed", err)
	}

	if len(warnings) > 0 {
		log.Println("warning for query kube_pod_info", warnings)
	}
	for _, v := range podInfo.(model.Vector) {
		namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
		pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
		hostIp := strings.TrimSpace(string(v.Metric[model.LabelName("host_ip")]))
		podIp := strings.TrimSpace(string(v.Metric[model.LabelName("pod_ip")]))
		podInfoMap[pod+namespace] = map[string]string{
			"namespace": namespace,
			"pod_name":  pod,
			"host_ip":   hostIp,
			"pod_ip":    podIp,
		}
	}

	fieldMaps := make(map[string]map[string]interface{}) // 用于存储points key=container+pod+namespace
	delete(fieldMaps, "")
	var containerInfo model.Value
	containerInfo, warnings, err = v1api.Query(ctx, "avg_over_time(kube_pod_container_info[1m])", timeNow, apiV1.WithTimeout(10*time.Second))
	if err != nil {
		log.Fatal("query kube_pod_container_info failed", err)
	}
	if len(warnings) > 0 {
		log.Println("warning for query kube_pod_container_info")
	}
	for _, v := range containerInfo.(model.Vector) {
		container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
		namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
		pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
		deployName := "UNKNOWN"
		if len(pod) > 0 {
			sliceString := strings.Split(pod, "-")
			if len(sliceString) > 2 && len(sliceString[len(sliceString)-1]) == 5 && len(sliceString[len(sliceString)-2]) == 9 {
				deployName = strings.Join(sliceString[0:len(sliceString)-2], "-")
			}
		}
		fieldMaps[container+pod+namespace] = map[string]interface{}{
			"container_name":       container,
			"namespace":            namespace,
			"pod_name":             pod,
			"deploy_name":          deployName,
			"pod_ip":               "",
			"host_ip":              "",
			"cpu_request":          -1.0,
			"cpu_limit":            -1.0,
			"cpu_utilization":      -1.0,
			"cpu_usage_per_second": -1.0,
			"mem_request":          -1.0,
			"mem_limit":            -1.0,
			"mem_usage":            -1.0,
			"mem_utilization":      -1.0,
			"memory_unit":          "MB",
		}

		if value, ok := podInfoMap[pod+namespace]; ok {
			fieldMaps[container+pod+namespace]["host_ip"] = value["host_ip"]
			fieldMaps[container+pod+namespace]["pod_ip"] = value["pod_ip"]
		}

	}

	var containerRequest model.Value
	containerRequest, warnings, err = v1api.Query(ctx, "avg_over_time(kube_pod_container_resource_requests[1m])", timeNow, apiV1.WithTimeout(10*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus Request Cores: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	for _, v := range containerRequest.(model.Vector) {
		container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
		namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
		pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
		resource := strings.TrimSpace(string(v.Metric[model.LabelName("resource")]))
		key := container + pod + namespace
		if _, ok := fieldMaps[key]; ok {
			if resource == "cpu" {
				fieldMaps[key]["cpu_request"] = float64(v.Value)
			} else if resource == "memory" {
				fieldMaps[key]["mem_request"] = float64(v.Value) / 1024 / 1024
			}

		}
	}

	var containerLimit model.Value
	containerLimit, warnings, err = v1api.Query(ctx, "avg_over_time(kube_pod_container_resource_limits[1m])", timeNow, apiV1.WithTimeout(10*time.Second))
	if err != nil {
		fmt.Printf("Error querying Prometheus Request Cores: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}

	for _, v := range containerLimit.(model.Vector) {
		container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
		namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
		pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
		resource := strings.TrimSpace(string(v.Metric[model.LabelName("resource")]))
		key := container + pod + namespace
		if _, ok := fieldMaps[key]; ok {
			if resource == "cpu" {
				fieldMaps[key]["cpu_limit"] = float64(v.Value)
			} else if resource == "memory" {
				fieldMaps[key]["mem_limit"] = float64(v.Value) / 1024 / 1024
			}

		}
	}

	cpuRate, _, _ := v1api.Query(ctx, "rate(container_cpu_usage_seconds_total[1m])", timeNow, apiV1.WithTimeout(10*time.Second))
	for _, v := range cpuRate.(model.Vector) {
		container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
		namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
		pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
		key := container + pod + namespace
		value := float64(v.Value)
		if _, ok := fieldMaps[key]; ok {
			fieldMaps[key]["cpu_usage_per_second"] = value

			if (fieldMaps[key]["cpu_request"]).(float64) > 0.001 {
				fieldMaps[key]["cpu_utilization"], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value/(fieldMaps[key]["cpu_request"]).(float64)*100), 64)
			}
		}
	}

	memUsage, _, _ := v1api.Query(ctx, "avg_over_time(container_memory_usage_bytes[1m])", timeNow, apiV1.WithTimeout(10*time.Second))
	for _, v := range memUsage.(model.Vector) {
		container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
		namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
		pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
		key := container + pod + namespace
		value := float64(v.Value)
		if _, ok := fieldMaps[key]; ok {
			fieldMaps[key]["mem_usage"] = value / 1024 / 1024

			if (fieldMaps[key]["mem_request"]).(float64) > 0.001 {
				fieldMaps[key]["mem_utilization"], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value/1024/1024/(fieldMaps[key]["mem_request"]).(float64)*100), 64)
			}
		}
	}

	fmt.Println("write to")

	for _, v := range fieldMaps {
		field := v
		pt, _ := clientInflux.NewPoint(
			"container_table",
			map[string]string{
				"pod":         v["pod_name"].(string),
				"container":   v["container_name"].(string),
				"namespace":   v["namespace"].(string),
				"deploy_name": v["deploy_name"].(string),
				"host_ip":     v["host_ip"].(string),
			},
			field,
			timeNow,
		)
		bp.AddPoint(pt)
	}
	err = c.Write(bp)
	if err != nil {
		fmt.Println("error write", err)
		os.Exit(1)
	}

	timeEnd := time.Now().Unix()
	fmt.Println(timeEnd - timeNow.Unix())
}
