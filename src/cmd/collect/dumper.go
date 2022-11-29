package main

import (
	"flag"
	"prometheus_dump_tools/internal/config"
	"prometheus_dump_tools/internal/dumpType"
	"prometheus_dump_tools/internal/prometheus_controller"
)

func main() {
	configPath := flag.String("config_path", "/config.yaml", "config yaml path")
	flag.Parse()
	conf := config.NewConfig(*configPath)

	var mc *prometheus_controller.ManageController
	mc = prometheus_controller.InitMC(conf)

	mc.RegisterMetricController("cpu_util", &dumpType.ContainerUtilization{
		TableName: "container_table",
	})

	//c := dumpType.ContainerUtilization{}
	//c.Dump(context.Background(), conf, time.Now().Unix())
	mc.Start()

}
