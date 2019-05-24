// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package sdfsprotocol

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type NetINodeSyncRequest struct {
	_tab flatbuffers.Table
}

func GetRootAsNetINodeSyncRequest(buf []byte, offset flatbuffers.UOffsetT) *NetINodeSyncRequest {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &NetINodeSyncRequest{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *NetINodeSyncRequest) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *NetINodeSyncRequest) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *NetINodeSyncRequest) NetINodeID() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func NetINodeSyncRequestStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func NetINodeSyncRequestAddNetINodeID(builder *flatbuffers.Builder, NetINodeID flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(NetINodeID), 0)
}
func NetINodeSyncRequestEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
