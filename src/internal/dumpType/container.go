package dumpType

import (
	"context"
	"errors"
	"fmt"
	clientInflux "github.com/influxdata/influxdb1-client/v2"
	"github.com/prometheus/client_golang/api"
	_ "github.com/prometheus/client_golang/api"
	apiV1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"log"
	"os"
	"prometheus_dump_tools/internal/config"
	"strconv"
	"strings"
	"time"
)

type ContainerUtilization struct {
	TableName string
}

//func (c *ContainerUtilization) ReQueryWithTimeout(ctx, query string)

func (c *ContainerUtilization) Query(ctx context.Context, v1api apiV1.API, query string, time2 time.Time, opt apiV1.Option) (model.Vector, error) {
	for i := 3; i > 0; i-- {
		podInfo, warnings, err := v1api.Query(ctx, query, time2, opt)
		if podInfo != nil && err == nil {
			if len(warnings) > 0 {
				log.Println("warning for query kube_pod_info", warnings)
			}
			result := podInfo.(model.Vector)
			if len(result) != 0 {
				return result, nil
			}
		}
		log.Println("Query Retry....")
		time.Sleep(time.Duration(500) * time.Millisecond)
		continue
	}
	return nil, errors.New("query failed")
}

func (c *ContainerUtilization) Dump(ctx context.Context, conf *config.Config, unixTime int64) {
	prometheusClient, err := api.NewClient(api.Config{
		Address: conf.Prometheus.Endpoints,
	})
	if err != nil {
		log.Fatal("connected to prometheus failed")
	}
	v1api := apiV1.NewAPI(prometheusClient)

	cInflux, err := clientInflux.NewHTTPClient(clientInflux.HTTPConfig{
		Addr:     conf.Influx.Endpoints,
		Username: conf.Influx.User,
		Password: conf.Influx.Pwd,
	})
	if err != nil {
		log.Fatal("connected to influxdb failed")
	}
	defer func(c clientInflux.Client) {
		err := c.Close()
		if err != nil {
			log.Fatal("close failed")
		}
	}(cInflux)

	fmt.Println("connect to prometheus and influxdb success")
	fmt.Println("start")

	var podInfo model.Vector

	interval := conf.Dumper.Interval

	podInfoMap := make(map[string]map[string]string)
	podInfo, err = c.Query(ctx, v1api, fmt.Sprintf("avg_over_time(kube_pod_info[%vs])", interval), time.Unix(unixTime, 0), apiV1.WithTimeout(20*time.Second))
	if err != nil {
		log.Println("query kube_pod_info failed", err)
		log.Println("ignore this timestamp")
		return
	}
	fmt.Println("pod info len:", len(podInfo))
	for _, v := range podInfo {
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
	var containerInfo model.Vector

	containerInfo, err = c.Query(ctx, v1api, fmt.Sprintf("avg_over_time(kube_pod_container_info[%vs])", interval), time.Unix(unixTime, 0), apiV1.WithTimeout(10*time.Second))

	if err != nil {
		log.Println("query kube_pod_container_info failed", err)
		log.Println("ignore this timestamp")
		return
	}
	fmt.Println("container info len:", len(containerInfo))
	for _, v := range containerInfo {
		container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
		namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
		pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
		deployName := "UNKNOWN"
		if len(pod) > 0 {
			sliceString := strings.Split(pod, "-")
			if len(sliceString) > 2 && len(sliceString[len(sliceString)-1]) == 5 {
				deployName = strings.Join(sliceString[0:len(sliceString)-2], "-")
			}
		}
		fieldMaps[container+pod+namespace] = map[string]interface{}{
			"container_name": container,
			"namespace":      namespace,
			"pod_name":       pod,
			"deploy_name":    deployName,
			//"pod_ip":               "",
			//"host_ip":              "",
			//"cpu_request":          0.0,
			//"cpu_limit":            0.0,
			//"cpu_utilization":      0.0,
			//"cpu_usage_per_second": 0.0,
			//"mem_request":          0.0,
			//"mem_limit":            0.0,
			//"mem_usage":            0.0,
			//"mem_utilization":      0.0,
			"memory_unit": "MB",
		}

		if value, ok := podInfoMap[pod+namespace]; ok {
			fieldMaps[container+pod+namespace]["host_ip"] = value["host_ip"]
			fieldMaps[container+pod+namespace]["pod_ip"] = value["pod_ip"]
		}

	}

	var containerRequest model.Vector
	containerRequest, err = c.Query(ctx, v1api, fmt.Sprintf("avg_over_time(kube_pod_container_resource_requests[%vs])", interval), time.Unix(unixTime, 0), apiV1.WithTimeout(10*time.Second))
	if err != nil {
		log.Printf("Error querying resource_requests: %v\n", err)
	} else {
		fmt.Printf("container request len: %v\n", len(containerRequest))
		for _, v := range containerRequest {
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
	}

	var containerLimit model.Vector
	containerLimit, err = c.Query(ctx, v1api, fmt.Sprintf("avg_over_time(kube_pod_container_resource_limits[%vs])", interval), time.Unix(unixTime, 0), apiV1.WithTimeout(10*time.Second))
	if err != nil {
		log.Printf("Error querying resource_limits: %v\n", err)
	} else {
		fmt.Printf("container limit len: %v\n", len(containerLimit))
		for _, v := range containerLimit {
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
	}

	cpuRate, err := c.Query(ctx, v1api, fmt.Sprintf("rate(container_cpu_usage_seconds_total[%vs])", interval), time.Unix(unixTime, 0), apiV1.WithTimeout(10*time.Second))
	if err != nil {
		log.Println("query container cpu usage seconds total failed")
	} else {
		fmt.Printf("cpu usage len: %v\n", len(cpuRate))
		for _, v := range cpuRate {
			container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
			namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
			pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
			key := container + pod + namespace
			value := float64(v.Value)
			if _, ok := fieldMaps[key]; ok {
				fieldMaps[key]["cpu_usage_per_second"] = value
				v, ok := fieldMaps[key]["cpu_request"]
				if ok && v.(float64) > 0.01 {
					fieldMaps[key]["cpu_utilization"], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value/(fieldMaps[key]["cpu_request"]).(float64)*100), 64)
				}
			}
		}
	}

	memUsage, err := c.Query(ctx, v1api, fmt.Sprintf("avg_over_time(container_memory_usage_bytes[%vs])", interval), time.Unix(unixTime, 0), apiV1.WithTimeout(10*time.Second))
	if err != nil {
		log.Println("query container memory usage bytes failed")
	} else {
		fmt.Printf("mem usage len: %v\n", len(memUsage))
		for _, v := range memUsage {
			container := strings.TrimSpace(string(v.Metric[model.LabelName("container")]))
			namespace := strings.TrimSpace(string(v.Metric[model.LabelName("namespace")]))
			pod := strings.TrimSpace(string(v.Metric[model.LabelName("pod")]))
			key := container + pod + namespace
			value := float64(v.Value)
			if _, ok := fieldMaps[key]; ok {
				fieldMaps[key]["mem_usage"] = value / 1024 / 1024

				v, ok := fieldMaps[key]["mem_request"]
				if ok && v.(float64) > 0.001 {
					fieldMaps[key]["mem_utilization"], _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value/1024/1024/(fieldMaps[key]["mem_request"]).(float64)*100), 64)
				}
			}
		}
	}

	bp, _ := clientInflux.NewBatchPoints(clientInflux.BatchPointsConfig{
		Database:  conf.Influx.Database,
		Precision: "s",
	})
	pointCnt := 0
	for _, v := range fieldMaps {
		field := v
		pt, _ := clientInflux.NewPoint(
			c.TableName,
			map[string]string{
				"pod":         v["pod_name"].(string),
				"container":   v["container_name"].(string),
				"namespace":   v["namespace"].(string),
				"deploy_name": v["deploy_name"].(string),
				"host_ip":     v["host_ip"].(string),
			},
			field,
			time.Unix(unixTime, 0),
		)
		//fmt.Println(pt)
		bp.AddPoint(pt)
		pointCnt += 1
	}
	if pointCnt == 0 {
		fmt.Println("Query Return None Data")
	} else {
		err = cInflux.Write(bp)
		if err != nil {
			fmt.Println("error write", err)
			os.Exit(1)
		}

		timeEnd := time.Now().Unix()
		fmt.Println("\nagg duration", timeEnd-time.Unix(unixTime, 0).Unix(), "s")
		fmt.Printf("time:[%v] datapoint num: %v\n", time.Unix(unixTime, 0).Format(time.UnixDate), pointCnt)
	}
}
