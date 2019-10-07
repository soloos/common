package solofsapitypes

import (
	"soloos/common/soloosbase"
	"soloos/solodb/offheap"
)

type MockNetINodeTable struct {
	*soloosbase.SoloosEnv
	table offheap.LKVTableWithBytes64
}

func (p *MockNetINodeTable) Init(soloosEnv *soloosbase.SoloosEnv) error {
	var err error
	p.SoloosEnv = soloosEnv

	err = p.OffheapDriver.InitLKVTableWithBytes64(&p.table, "MockNetINode",
		int(NetINodeStructSize), -1, offheap.DefaultKVTableSharedCount, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *MockNetINodeTable) MustGetNetINode(netINodeID NetINodeID) (NetINodeUintptr, bool) {
	uObject, afterSetNewObj := p.table.MustGetObject(netINodeID)
	var loaded = afterSetNewObj == nil
	if afterSetNewObj != nil {
		afterSetNewObj()
	}
	uNetINode := (NetINodeUintptr)(uObject)
	return uNetINode, loaded
}

func (p *MockNetINodeTable) AllocNetINode(netBlockCap, memBlockCap int) NetINodeUintptr {
	var netINodeID NetINodeID
	InitTmpNetINodeID(&netINodeID)
	uNetINode, _ := p.MustGetNetINode(netINodeID)
	uNetINode.Ptr().ID = netINodeID
	uNetINode.Ptr().NetBlockCap = netBlockCap
	uNetINode.Ptr().MemBlockCap = memBlockCap
	return uNetINode
}

func (p *MockNetINodeTable) ReleaseNetINode(uNetINode NetINodeUintptr) {
	p.table.ReleaseObject(offheap.LKVTableObjectUPtrWithBytes64(uNetINode))
}