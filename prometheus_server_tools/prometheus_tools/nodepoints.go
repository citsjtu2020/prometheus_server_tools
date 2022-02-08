package prometheus_tools

type NodeResult struct {
	node      string
	values    []Rawpoint
}

func (n *NodeResult) GetNode() string{
	return n.node
}

func (n *NodeResult) SetNode(node string){
	n.node = node
}

func (n *NodeResult) GetValues() []Rawpoint{
	return n.values
}

func (n *NodeResult) GetValuei(i int) Rawpoint{
	return n.values[i]
}

func (n *NodeResult) SetValues(rawpoints []Rawpoint){
	n.values = make([]Rawpoint,0,len(rawpoints))
	for _,v := range rawpoints{
		n.values = append(n.values,v)
	}
}

func (n *NodeResult) SetValuei(rawpoint Rawpoint, i int){
	if i < len(n.values){
		n.values[i] = rawpoint
	}else{
		n.values = append(n.values,rawpoint)
	}

}

type NodeCPUResult struct {
	NodeResult
	cpu int
	mode string
}

func (n *NodeCPUResult) GetCPU() int{
	return n.cpu
}

func (n *NodeCPUResult) SetCPU(cpu int){
	n.cpu = cpu
}

func (n *NodeCPUResult) GetNode() string{
	return n.node
}

func (n *NodeCPUResult) SetNode(node string){
	n.node = node
}

func (n *NodeCPUResult) GetMode() string{
	return n.mode
}

func (n *NodeCPUResult) SetMode(mode string){
	n.mode = mode
}

func (n *NodeCPUResult) GetValues() []Rawpoint{
	return n.values
}

func (n *NodeCPUResult) GetValuei(i int) Rawpoint{
	return n.values[i]
}

func (n *NodeCPUResult) SetValues(rawpoints []Rawpoint){
	n.values = make([]Rawpoint,0,len(rawpoints))
	for _,v := range rawpoints{
		n.values = append(n.values,v)
	}
}

func (n *NodeCPUResult) SetValuei(rawpoint Rawpoint, i int){
	if i < len(n.values){
		n.values[i] = rawpoint
	}else{
		n.values = append(n.values,rawpoint)
	}

}

type NodeNetworkResult struct {
	NodeResult
	inter string
}

func (n *NodeNetworkResult) GetNode() string{
	return n.node
}

func (n *NodeNetworkResult) SetNode(node string){
	n.node = node
}

func (n *NodeNetworkResult) GetInter() string{
	return n.inter
}

func (n *NodeNetworkResult) SetInter(inter string){
	n.inter = inter
}

func (n *NodeNetworkResult) GetValues() []Rawpoint{
	return n.values
}

func (n *NodeNetworkResult) GetValuei(i int) Rawpoint{
	return n.values[i]
}

func (n *NodeNetworkResult) SetValues(rawpoints []Rawpoint){
	n.values = make([]Rawpoint,0,len(rawpoints))
	for _,v := range rawpoints{
		n.values = append(n.values,v)
	}
}

func (n *NodeNetworkResult) SetValuei(rawpoint Rawpoint, i int){
	if i < len(n.values){
		n.values[i] = rawpoint
	}else{
		n.values = append(n.values,rawpoint)
	}

}

type NodeFSResult struct {
	NodeResult
	device string
}

func (n *NodeFSResult) GetNode() string{
	return n.node
}

func (n *NodeFSResult) SetNode(node string){
	n.node = node
}

func (n *NodeFSResult) GetDevice() string{
	return n.device
}

func (n *NodeFSResult) SetDevice(device string){
	n.device = device
}

func (n *NodeFSResult) GetValues() []Rawpoint{
	return n.values
}

func (n *NodeFSResult) GetValuei(i int) Rawpoint{
	return n.values[i]
}

func (n *NodeFSResult) SetValues(rawpoints []Rawpoint){
	n.values = make([]Rawpoint,0,len(rawpoints))
	for _,v := range rawpoints{
		n.values = append(n.values,v)
	}
}

func (n *NodeFSResult) SetValuei(rawpoint Rawpoint, i int){
	if i < len(n.values){
		n.values[i] = rawpoint
	}else{
		n.values = append(n.values,rawpoint)
	}

}
