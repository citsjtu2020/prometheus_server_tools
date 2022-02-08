package prometheus_tools

import "sort"

type Rawpoint struct{
	timestamp float64
	value float64
}

type Wrapper struct {
    points []Rawpoint
    by  func(p, q *Rawpoint) bool
}

type SortBy func(p, q *Rawpoint) bool

func (pw Wrapper) Len() int { // 重写 Len() 方法
    return len(pw.points)
}
func (pw Wrapper) Swap(i, j int) { // 重写 Swap() 方法
    pw.points[i], pw.points[j] = pw.points[j], pw.points[i]
}
func (pw Wrapper) Less(i, j int) bool { // 重写 Less() 方法
    return pw.by(&pw.points[i], &pw.points[j])
}

// 封装成 SortLog 方法
func SortPoint(points []Rawpoint, by SortBy) {
    sort.Sort(Wrapper{points, by})
}

//type ContainerAggre struct {
//	container string
//	pod       string
//	deploy    string
//	namespace string
//	node      string
//
//}

type ContainerResult struct {
	container string
	pod       string
	deploy    string
	namespace string
	node      string
	values    []Rawpoint
	//by func(p,q *Rawpoint) bool
}



//type Wrapper struct {
//    log []Log
//    by  func(p, q *Log) bool
//}

func (c *ContainerResult) GetContainer() string{
	return c.container
}

func (c *ContainerResult) SetContainer(container string){
	c.container = container
}

func (c *ContainerResult) GetPod() string{
	return c.pod
}

func (c *ContainerResult) SetPod(pod string){
	c.pod = pod
}

func (c *ContainerResult) GetDeploy() string {
	return c.deploy
}

func (c *ContainerResult) SetDeploy(deploy string) {
	c.deploy = deploy
}

func (c *ContainerResult) GetNamespace() string{
	return c.namespace
}

func (c *ContainerResult) SetNamespace(namespace string){
	c.namespace = namespace
}

func (c *ContainerResult) GetNode() string {
	return c.node
}

func (c *ContainerResult) SetNode(node string) {
	c.node = node
}


func (c *ContainerResult) GetValues() []Rawpoint {
	return c.values
}

func (c *ContainerResult) GetValuesi(i int) Rawpoint {
	return c.values[i]
}

func (c *ContainerResult) SetValues(values []Rawpoint) {
	//return c.values
	
	c.values = make([]Rawpoint,0,len(values))
	for _,v := range values{
		c.values  = append(c.values,v)
	}
}



func (c *ContainerResult) SetValuei(i int,r Rawpoint){
	//return c.values[i]
	if i >= 0{
		if i > cap(c.values){
			c.values = append(c.values,r)
		}else{
			c.values[i] = r
		}
	}
}

type ContainerNetworkResult struct {
	ContainerResult
	inter string
}

func (c *ContainerNetworkResult) GetContainer() string{
	return c.container
}

func (c *ContainerNetworkResult) SetContainer(container string){
	c.container = container
}

func (c *ContainerNetworkResult) GetInfer() string{
	return c.inter
}

func (c *ContainerNetworkResult) SetInfter(inter string){
	c.inter = inter
}

func (c *ContainerNetworkResult) GetPod() string{
	return c.pod
}

func (c *ContainerNetworkResult) SetPod(pod string){
	c.pod = pod
}

func (c *ContainerNetworkResult) GetDeploy() string {
	return c.deploy
}

func (c *ContainerNetworkResult) SetDeploy(deploy string) {
	c.deploy = deploy
}

func (c *ContainerNetworkResult) GetNamespace() string{
	return c.namespace
}

func (c *ContainerNetworkResult) SetNamespace(namespace string){
	c.namespace = namespace
}

func (c *ContainerNetworkResult) GetNode() string {
	return c.node
}

func (c *ContainerNetworkResult) SetNode(node string) {
	c.node = node
}


func (c *ContainerNetworkResult) GetValues() []Rawpoint {
	return c.values
}

func (c *ContainerNetworkResult) SetValues(values []Rawpoint) {
	//return c.values
	
	c.values = make([]Rawpoint,0,len(values))
	for _,v := range values{
		c.values  = append(c.values,v)
	}
}

func (c *ContainerNetworkResult) GetValuesi(i int) Rawpoint {
	return c.values[i]
}

func (c *ContainerNetworkResult) SetValuei(i int,r Rawpoint){
	//return c.values[i]
	if i >= 0{
		if i > cap(c.values){
			c.values = append(c.values,r)
		}else{
			c.values[i] = r
		}
	}
}

type ContainerFSResult struct {
	ContainerResult
	device string
}

func (c *ContainerFSResult) GetContainer() string{
	return c.container
}

func (c *ContainerFSResult) SetContainer(container string){
	c.container = container
}

func (c *ContainerFSResult) GetDevice() string{
	return c.device
}

func (c *ContainerFSResult) SetInfter(device string){
	c.device = device
}

func (c *ContainerFSResult) GetPod() string{
	return c.pod
}

func (c *ContainerFSResult) SetPod(pod string){
	c.pod = pod
}

func (c *ContainerFSResult) GetDeploy() string {
	return c.deploy
}

func (c *ContainerFSResult) SetDeploy(deploy string) {
	c.deploy = deploy
}

func (c *ContainerFSResult) GetNamespace() string{
	return c.namespace
}

func (c *ContainerFSResult) SetNamespace(namespace string){
	c.namespace = namespace
}

func (c *ContainerFSResult) GetNode() string {
	return c.node
}

func (c *ContainerFSResult) SetNode(node string) {
	c.node = node
}


func (c *ContainerFSResult) GetValues() []Rawpoint {
	return c.values
}

func (c *ContainerFSResult) SetValues(values []Rawpoint) {
	//return c.values
	
	c.values = make([]Rawpoint,0,len(values))
	for _,v := range values{
		c.values  = append(c.values,v)
	}
}

func (c *ContainerFSResult) GetValuesi(i int) Rawpoint {
	return c.values[i]
}

func (c *ContainerFSResult) SetValuei(i int,r Rawpoint){
	//return c.values[i]
	if i >= 0{
		if i > cap(c.values){
			c.values = append(c.values,r)
		}else{
			c.values[i] = r
		}
	}
}

type ContainerMmemoryFailure struct {
	ContainerResult
	scope string
	failure_type string
}

func (c *ContainerMmemoryFailure) GetContainer() string{
	return c.container
}

func (c *ContainerMmemoryFailure) SetContainer(container string){
	c.container = container
}

func (c *ContainerMmemoryFailure) GetFailureType() string{
	return c.failure_type
}

func (c *ContainerMmemoryFailure) SetFailureType(failure_type string){
	c.failure_type = failure_type
}

func (c *ContainerMmemoryFailure) GetScope() string{
	return c.scope
}

func (c *ContainerMmemoryFailure) SetScope(scope string){
	c.scope = scope
}

func (c *ContainerMmemoryFailure) GetPod() string{
	return c.pod
}

func (c *ContainerMmemoryFailure) SetPod(pod string){
	c.pod = pod
}

func (c *ContainerMmemoryFailure) GetDeploy() string {
	return c.deploy
}

func (c *ContainerMmemoryFailure) SetDeploy(deploy string) {
	c.deploy = deploy
}

func (c *ContainerMmemoryFailure) GetNamespace() string{
	return c.namespace
}

func (c *ContainerMmemoryFailure) SetNamespace(namespace string){
	c.namespace = namespace
}

func (c *ContainerMmemoryFailure) GetNode() string {
	return c.node
}

func (c *ContainerMmemoryFailure) SetNode(node string) {
	c.node = node
}


func (c *ContainerMmemoryFailure) GetValues() []Rawpoint {
	return c.values
}

func (c *ContainerMmemoryFailure) SetValues(values []Rawpoint) {
	//return c.values
	
	c.values = make([]Rawpoint,0,len(values))
	for _,v := range values{
		c.values  = append(c.values,v)
	}
}

func (c *ContainerMmemoryFailure) GetValuesi(i int) Rawpoint {
	return c.values[i]
}

func (c *ContainerMmemoryFailure) SetValuei(i int,r Rawpoint){
	//return c.values[i]
	if i >= 0{
		if i > cap(c.values){
			c.values = append(c.values,r)
		}else{
			c.values[i] = r
		}
	}
}

