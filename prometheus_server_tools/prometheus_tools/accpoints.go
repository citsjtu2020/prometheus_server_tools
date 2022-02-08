package prometheus_tools

type ContainerAccResult struct {
	container string
	pod       string
	deploy    string
	namespace string
	node      string
	//values    []Rawpoint
	timestamp int64
	result map[string]float64
	//by func(p,q *Rawpoint) bool
}

type ContainerNetworkAccResult struct {
	ContainerAccResult
	inter string
}


func (c *ContainerAccResult ) GetContainer() string{
	return c.container
}

func (c *ContainerAccResult ) SetContainer(container string){
	c.container = container
}

func (c *ContainerAccResult ) GetPod() string{
	return c.pod
}

func (c *ContainerAccResult ) SetPod(pod string){
	c.pod = pod
}

func (c *ContainerAccResult ) GetDeploy() string {
	return c.deploy
}

func (c *ContainerAccResult) SetDeploy(deploy string) {
	c.deploy = deploy
}

func (c *ContainerAccResult) GetNamespace() string{
	return c.namespace
}

func (c *ContainerAccResult) SetNamespace(namespace string){
	c.namespace = namespace
}

func (c *ContainerAccResult) GetNode() string {
	return c.node
}

func (c *ContainerAccResult) SetNode(node string) {
	c.node = node
}

func (c *ContainerAccResult) GetTimestamp () int64{
	return c.timestamp
}

func (c *ContainerAccResult) SetTimestamp(timestamp int64){
	c.timestamp = timestamp
}

func (c *ContainerAccResult) GetResults() map[string]float64 {
	return c.result
}

func (c *ContainerAccResult) GetResultItem(key string) (float64,bool) {
	res,ok := c.result[key]
	return res,ok
}

func (c *ContainerAccResult) SetResults(t map[string]float64) {
	if c.result == nil{
		c.result = make(map[string]float64)
	}
	if t != nil{
		for k,v := range t{
			c.result[k] = v
		}
	}
}



func (c *ContainerAccResult) SetResultItem(key string,value float64){
	//return c.values[i]
	if c.result == nil{
		c.result = make(map[string]float64)
	}
	c.result[key] = value
}

func (c *ContainerNetworkAccResult) GetContainer() string{
	return c.container
}

func (c *ContainerNetworkAccResult) SetContainer(container string){
	c.container = container
}

func (c *ContainerNetworkAccResult) GetPod() string{
	return c.pod
}

func (c *ContainerNetworkAccResult) SetPod(pod string){
	c.pod = pod
}

func (c *ContainerNetworkAccResult) GetDeploy() string {
	return c.deploy
}

func (c *ContainerNetworkAccResult) SetDeploy(deploy string) {
	c.deploy = deploy
}

func (c *ContainerNetworkAccResult) GetNamespace() string{
	return c.namespace
}

func (c *ContainerNetworkAccResult) SetNamespace(namespace string){
	c.namespace = namespace
}

func (c *ContainerNetworkAccResult) GetNode() string {
	return c.node
}

func (c *ContainerNetworkAccResult) SetNode(node string) {
	c.node = node
}

func (c *ContainerNetworkAccResult) GetInter() string {
	return c.inter
}

func (c *ContainerNetworkAccResult) SetInter(inter string) {
	c.inter = inter
}

func (c *ContainerNetworkAccResult) GetTimestamp () int64{
	return c.timestamp
}

func (c *ContainerNetworkAccResult) SetTimestamp(timestamp int64){
	c.timestamp = timestamp
}

func (c *ContainerNetworkAccResult) GetResults() map[string]float64 {
	return c.result
}

func (c *ContainerNetworkAccResult) GetResultItem(key string) (float64,bool) {
	res,ok := c.result[key]
	return res,ok
}

func (c *ContainerNetworkAccResult) SetResults(t map[string]float64) {
	if c.result == nil{
		c.result = make(map[string]float64)
	}
	if t != nil{
		for k,v := range t{
			c.result[k] = v
		}
	}
}



func (c *ContainerNetworkAccResult) SetResultItem(key string,value float64){
	//return c.values[i]
	if c.result == nil{
		c.result = make(map[string]float64)
	}
	c.result[key] = value
}

