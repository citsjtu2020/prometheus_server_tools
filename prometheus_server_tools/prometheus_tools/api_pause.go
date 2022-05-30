package prometheus_tools

import (
	"bytes"
	"encoding/json"
	"strconv"

	//"bufio"
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

func check(err error) {
    if err != nil {
        panic(err.Error())
    }
}

func Get(url string) ([]byte,error){
	defer func() {
		if err := recover();err != nil{
			fmt.Println(err)
		}
	}()

	resp,err := http.Get(url)
	//Println(resp.StatusCode)
    if resp.StatusCode != 200 {
        panic(fmt.Sprintf("%v",resp.StatusCode))
    }
	check(err)
	//rb := bufio.NewReader(resp.Body)
	//rb.Read()
	return ioutil.ReadAll(resp.Body)
}

func Generate_address(ip string,port int) string{
	address := fmt.Sprintf("%s:%v",ip,port)
	address = strip(address)
	return address
}

func Gen_range_url(address string, metric string,start int64,end int64,unit string,step int) string{
	//cmd = "http://10.64.2.31:9090/api/v1/query_range?query=%s&start=%s&end=%s&step=3s"
	//% (k.strip(),start_date,end_date)
	url := 	fmt.Sprintf("http://%s/api/v1/query_range?query=%s&start=%v&end=%v&step=%v%s",address,metric,start,end,step,unit)
	return strip(url)
}



func strip(src string) string{
	src = strings.ToLower(src)
    re, _ := regexp.Compile(`<!doctype.*?>`)
    src = re.ReplaceAllString(src,"")
    re, _ = regexp.Compile(`<!--.*?-->`)
    src = re.ReplaceAllString(src, "")

    re, _ = regexp.Compile(`<script[\S\s]+?</script>`)
    src = re.ReplaceAllString(src, "")

    re, _ = regexp.Compile(`<style[\S\s]+?</style>`)
    src = re.ReplaceAllString(src, "")

    re, _ = regexp.Compile(`<.*?>`)
    src = re.ReplaceAllString(src, "")

    re, _ = regexp.Compile(`&.{1,5};|&#.{1,5};`)
    src = re.ReplaceAllString(src, "")

    src = strings.Replace(src, "\r\n", "\n", -1)
    src = strings.Replace(src, "\r", "\n", -1)
    return src
}

func Do(url string) string{
	body,err := Get(url)
	check(err)
	plainText := strip(string(body))
    return plainText
}

// query data from url
func Container_Raw_Metric(url string) []ContainerResult{
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	if (t > 0){
		results := make([]ContainerResult,0,t)
		tmpresult := ContainerResult{}
		for _,data := range datas{
			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)
			//data2.()
			tmpresult.container,_ = data2.Get("metric").Get("container").String()
			tmpresult.node,_ = data2.Get("metric").Get("instance").String()
			tmppod,_ := data2.Get("metric").Get("pod").String()
			tmpresult.pod = tmppod
			tmpresult.namespace,_ = data2.Get("metric").Get("namespace").String()
			tmpdeploys := strings.Split(tmppod,"-")
			if len(tmpdeploys) <= 3{
				tmpresult.deploy = tmpdeploys[0]
			} else{
				var tmpdeploy bytes.Buffer
				for k:=0;k<len(tmpdeploys)-3;k++{
				tmpdeploy.WriteString(tmpdeploys[k])
				tmpdeploy.WriteString("-")
				}
				tmpdeploy.WriteString(tmpdeploys[len(tmpdeploys)-3])
				tmpresult.deploy = tmpdeploy.String()
			}
			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}
			SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)
			results = append(results, tmpresult)
		}
		return results
	}
	return nil
}


func Container_Network_Metric(url string) []ContainerNetworkResult{
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	//var results []interface{}
	if (t>0){
		//datas2,_ := datas.([]interface{})
		//results
		//ifcontainer := false
		//var metric string
		//for i,data := range datas{
		//	data2,_ := data.(simplejson.Json)
		//	//data2.()
		//	if i == 0{
		//		//simplejson.NewJson([]byte(data))
		//		metric,_ = data2.Get("metric").Get("__name__").String()
		//		ifcontainer = strings.Contains(metric,"container")
		//		break
		//	}
		//}
		//if strings.Contains(metric,"network"){
		results := make([]ContainerNetworkResult,0,t)
		for _,data := range datas {
			tmpresult := ContainerNetworkResult{}
			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)
			//data2.()
			tmpresult.container, _ = data2.Get("metric").Get("container").String()
			tmpresult.node, _ = data2.Get("metric").Get("instance").String()
			tmppod, _ := data2.Get("metric").Get("pod").String()
			tmpresult.pod = tmppod
			tmpresult.namespace, _ = data2.Get("metric").Get("namespace").String()
			tmpdeploys := strings.Split(tmppod, "-")
			if len(tmpdeploys) <= 3{
				tmpresult.deploy = tmpdeploys[0]
			} else{
				var tmpdeploy bytes.Buffer
				for k:=0;k<len(tmpdeploys)-3;k++{
				tmpdeploy.WriteString(tmpdeploys[k])
				tmpdeploy.WriteString("-")
				}
				tmpdeploy.WriteString(tmpdeploys[len(tmpdeploys)-3])
				tmpresult.deploy = tmpdeploy.String()
			}
			tmpresult.inter, _ = data2.Get("metric").Get("interface").String()

			//fmt.Println(tmpresult.deploy)
			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}

	//		sort.Sort(Wrapper{log, func(p, q *Log) bool {
    //    return q.Num < p.Num // Num 递减排序
    //}})

    //fmt.Println(log)
    //间接封装
    		SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)

			results = append(results, tmpresult)
		}
		return results
	}
	return nil
			//}
//	{'status': 'success',
//	'data': {
//	'resultType': 'matrix',
//	'result': [
//	{'metric': {'__name__': 'container_cpu_cfs_periods_total',
//	'container': 'consul',
//	'instance': 'k8smaster-node02',
//	'job': 'pods',
//	'namespace': 'micro',
//	'pod': 'consul-79fb49d755-msncr'},
//	'values': [
//	[1643795060, '74682'], [1643795063.001, '74720'],
//	[1643795066.002, '74749'], [1643795069.003, '74780'],
//	[1643795072.004, '74813'], [1643795075.005, '74842'],
//	[1643795078.006, '74892'], [1643795081.007, '74892'],
//	[1643795084.008, '74956'], [1643795087.009, '74956'],
//	[1643795090.01, '75014'], [1643795093.011, '75014'],
//	[1643795096.012, '75053'], [1643795099.013, '75099'],
//	[1643795102.014, '75099'], [1643795105.015, '75143'],
//	[1643795108.016, '75173'], [1643795111.017, '75203'],
//	[1643795114.018, '75250'], [1643795117.019, '75280'],
//	[1643795120.02, '75280'], [1643795123.021, '75347'],
//	[1643795126.022, '75347'], [1643795129.023, '75383'],
//	[1643795132.024, '75408'], [1643795135.025, '75459'],
//	[1643795138.026, '75493'], [1643795141.027, '75493'],
//	[1643795144.028, '75532'], [1643795147.029, '75578'],
//	[1643795150.03, '75578'], [1643795153.031, '75620'],
//	[1643795156.032, '75650'], [1643795159.033, '75705']]},
}

func Container_FS_Metric(url string) []ContainerFSResult{
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	//var results []interface{}
	if (t>0) {
		results := make([]ContainerFSResult,0, t)
		for _, data := range datas {
			tmpresult := ContainerFSResult{}
			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)
			//data2.()
			tmpresult.container, _ = data2.Get("metric").Get("container").String()
			tmpresult.node, _ = data2.Get("metric").Get("instance").String()
			tmppod, _ := data2.Get("metric").Get("pod").String()
			tmpresult.pod = tmppod
			tmpresult.namespace, _ = data2.Get("metric").Get("namespace").String()
			tmpdeploys := strings.Split(tmppod, "-")
			if len(tmpdeploys) <= 3{
				tmpresult.deploy = tmpdeploys[0]
			} else{
				var tmpdeploy bytes.Buffer
				for k:=0;k<len(tmpdeploys)-3;k++{
				tmpdeploy.WriteString(tmpdeploys[k])
				tmpdeploy.WriteString("-")
				}
				tmpdeploy.WriteString(tmpdeploys[len(tmpdeploys)-3])
				tmpresult.deploy = tmpdeploy.String()
			}
			tmpresult.device, _ = data2.Get("metric").Get("device").String()

			//fmt.Println(tmpresult.deploy)
			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}
			SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)

			results = append(results, tmpresult)
		}
		return results
	}
	return nil
}

func Container_Memory_Failure(url string) []ContainerMmemoryFailure {
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	//var results []interface{}
	if (t>0) {
		results := make([]ContainerMmemoryFailure,0, t)
		for _, data := range datas {
			tmpresult := ContainerMmemoryFailure{}
			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)
			//data2.()
			tmpresult.container, _ = data2.Get("metric").Get("container").String()
			tmpresult.node, _ = data2.Get("metric").Get("instance").String()
			tmppod, _ := data2.Get("metric").Get("pod").String()
			tmpresult.pod = tmppod
			tmpresult.namespace, _ = data2.Get("metric").Get("namespace").String()
			tmpdeploys := strings.Split(tmppod, "-")

			if len(tmpdeploys) <= 3{
				tmpresult.deploy = tmpdeploys[0]
			} else{
				var tmpdeploy bytes.Buffer
				for k:=0;k<len(tmpdeploys)-3;k++{
				tmpdeploy.WriteString(tmpdeploys[k])
				tmpdeploy.WriteString("-")
				}
				tmpdeploy.WriteString(tmpdeploys[len(tmpdeploys)-3])
				tmpresult.deploy = tmpdeploy.String()
			}

			//if len(tmpdeploys) > 3{
			//	tmpdeploy.WriteString(tmpdeploys[len(tmpdeploys)-2])
			//}

			tmpresult.failure_type, _ = data2.Get("metric").Get("failure_type").String()
			tmpresult.scope, _ = data2.Get("metric").Get("scope").String()

			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}
			SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)
			results = append(results, tmpresult)
		}
		return results
	}
	return nil
}

func Node_Raw_Metric(url string) []NodeResult{
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	if (t > 0){
		results := make([]NodeResult,0,t)
		tmpresult := NodeResult{}
		for _,data := range datas{
			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)
			
			tmpresult.node,_ = data2.Get("metric").Get("nodename").String()
			
			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}
			SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)
			results = append(results, tmpresult)
		}
		return results
	}
	return nil
}

func Node_FS_Metric(url string) []NodeFSResult{
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	if (t > 0){
		results := make([]NodeFSResult,0,t)
		tmpresult := NodeFSResult{}
		for _,data := range datas{
			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)
			
			tmpresult.node,_ = data2.Get("metric").Get("nodename").String()
			tmpresult.device,_ = data2.Get("metric").Get("device").String()
			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}
			SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)
			results = append(results, tmpresult)
		}
		return results
	}
	return nil
}

func Node_CPU_Metric(url string) []NodeCPUResult{
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	if (t > 0){
		results := make([]NodeCPUResult,0,t)
		tmpresult := NodeCPUResult{}
		for _,data := range datas{
			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)

			tmpresult.node,_ = data2.Get("metric").Get("nodename").String()
			tmpresult.cpu,_ = data2.Get("metric").Get("cpu").Int()
			tmpresult.mode,_ = data2.Get("metric").Get("mode").String()

			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}
			SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)
			results = append(results, tmpresult)
		}
		return results
	}
	return nil
}

func Node_Network_Metric(url string) []NodeNetworkResult{
	raws := Do(url)
	js,err := simplejson.NewJson([]byte(raws))
	check(err)
	datas,err := js.Get("data").Get("result").Array()
	check(err)
	t := len(datas)
	if (t > 0){
		results := make([]NodeNetworkResult,0,t)
		tmpresult := NodeNetworkResult{}
		var metric string
		ifdevice := false
		for i,data := range datas{

			bytedata,_ := json.Marshal(data)
			data2,_ := simplejson.NewJson(bytedata)
			if i == 0{
				metric,_ = data2.Get("metric").Get("__name__").String()
				ifdevice = strings.Contains(metric,"network_speed_bytes")
				break
			}
			tmpresult.node,_ = data2.Get("metric").Get("nodename").String()
			//network_speed_bytes
			//interface
			if ifdevice{
				tmpresult.inter,_ = data2.Get("metric").Get("interface").String()
			}else{
				tmpresult.inter,_ = data2.Get("metric").Get("device").String()
			}
			//fmt.Println(tmpresult.deploy)
			values,_ := data2.Get("values").Array()
			//fmt.Println(values)
			lv := len(values)
			rawpoints := make([]Rawpoint,0,lv)
			//fmt.Println(values)
			for _,v := range values{
				//timestamp :=
				value,_ := v.([]interface{})
				tmpraw := Rawpoint{}
				for k,u := range value{
					//fmt.Print(u)
					x,_ := json.Marshal(u)
					xstr := string(x)
					if strings.Contains(xstr,"\""){
						xstr = strings.Replace(xstr,"\"","",-1)
					}
					//fmt.Println(xstr)
					//json.
					y, _ := strconv.ParseFloat(xstr, 64) //将字符型号转化为float64
					//println()
					//fmt.Printf("%v:%v\n",k,y)
					switch k {
					case 0:
						tmpraw.timestamp = float64(y)
					case 1:
						tmpraw.value = float64(y)
					}
				}
				rawpoints = append(rawpoints,tmpraw)
			}
			SortPoint(rawpoints, func(p, q *Rawpoint) bool {
        		return p.timestamp < q.timestamp // CreateTime 递增排序
    		})
			tmpresult.SetValues(rawpoints)
			results = append(results, tmpresult)
		}
		return results
	}
	return nil
}


