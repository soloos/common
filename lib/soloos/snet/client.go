package snet

type Client struct {
	Servers map[int]int
}

func (p *Client) Init() error {
	return nil
}

func (p *Client) Write() {
}
