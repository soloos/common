package snet

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"soloos/common/snettypes"
	"strconv"

	"golang.org/x/xerrors"
)

type NetDriverWebClient struct {
	netDriver     *NetDriver
	webServerAddr string
}

func NewNetDriverWebClient(netDriver *NetDriver, webServerAddr string) (*NetDriverWebClient, error) {
	var (
		ret *NetDriverWebClient = new(NetDriverWebClient)
		err error
	)

	ret.Init(netDriver, webServerAddr)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func (p *NetDriverWebClient) Init(netDriver *NetDriver, webServerAddr string) error {
	p.webServerAddr = webServerAddr
	return nil
}

func (p *NetDriverWebClient) readResp(resp *http.Response, err error) (snettypes.Peer, error) {
	var ret snettypes.Peer
	var respJSON snettypes.GetPeerRespJSON
	if err != nil {
		return ret, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		goto DONE
	}

	err = json.Unmarshal(body, &respJSON)
	if err != nil {
		goto DONE
	}

DONE:
	resp.Body.Close()
	if err != nil {
		return ret, err
	}

	if respJSON.Errno != snettypes.CODE_OK {
		return ret, xerrors.New(respJSON.ErrMsg)
	}
	ret = snettypes.PeerJSONToPeer(respJSON.Data)

	return ret, nil
}

func (p *NetDriverWebClient) GetPeer(peerID snettypes.PeerID) (snettypes.Peer, error) {
	var urlPath = p.webServerAddr + "/Peer/Get"
	return p.readResp(http.PostForm(urlPath, url.Values{"PeerID": {peerID.Str()}}))
}

// MustGetPee return uPeer and peer is inited before
func (p *NetDriverWebClient) RegisterPeer(peerID snettypes.PeerID, addr string, protocol int) error {
	if protocol != snettypes.ProtocolSRPC {
		return nil
	}

	var err error
	var urlPath = p.webServerAddr + "/Peer/Register"
	var urlValues = url.Values{
		"PeerID":   {peerID.Str()},
		"Addr":     {addr},
		"Protocol": {strconv.Itoa(protocol)},
	}
	_, err = p.readResp(http.PostForm(urlPath, urlValues))
	return err
}
