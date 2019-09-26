package solofsapi

import (
	"soloos/common/solofsapitypes"
	"soloos/common/solofsprotocol"
	"soloos/common/snettypes"
	"soloos/common/xerrors"

	flatbuffers "github.com/google/flatbuffers/go"
)

func SetCommonResponseCode(protocolBuilder *flatbuffers.Builder, code int) {
	solofsprotocol.CommonResponseStart(protocolBuilder)
	solofsprotocol.CommonResponseAddCode(protocolBuilder, int32(code))
	protocolBuilder.Finish(solofsprotocol.CommonResponseEnd(protocolBuilder))
}

func CommonResponseToError(obj *solofsprotocol.CommonResponse) error {
	switch obj.Code() {
	case snettypes.CODE_OK:
		return nil
	case snettypes.CODE_404:
		return solofsapitypes.ErrObjectNotExists
	case snettypes.CODE_502:
		return solofsapitypes.ErrRemoteService
	}

	return xerrors.New(string(obj.Error()))
	// return types.ErrRemoteService
}
