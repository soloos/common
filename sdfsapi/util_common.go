package sdfsapi

import (
	"soloos/common/sdfsapitypes"
	"soloos/common/sdfsprotocol"
	"soloos/common/snettypes"
	"soloos/common/xerrors"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetCommonResponseCode(protocolBuilder *flatbuffers.Builder, code int) {
	sdfsprotocol.CommonResponseStart(protocolBuilder)
	sdfsprotocol.CommonResponseAddCode(protocolBuilder, int32(code))
	protocolBuilder.Finish(sdfsprotocol.CommonResponseEnd(protocolBuilder))
}

func CommonResponseToError(obj *sdfsprotocol.CommonResponse) error {
	switch obj.Code() {
	case snettypes.CODE_OK:
		return nil
	case snettypes.CODE_404:
		return sdfsapitypes.ErrObjectNotExists
	case snettypes.CODE_502:
		return sdfsapitypes.ErrRemoteService
	}

	return xerrors.New(string(obj.Error()))
	// return types.ErrRemoteService
}
