package snettypes

type APIRespCommonJSON struct {
	Errno  int    `json:"ErrNo"`
	ErrMsg string `json:"ErrMsg"`
}
