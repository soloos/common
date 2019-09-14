package iron

func (p *Proxy) InitStandAloneWebServer(prefix string, options Options) error {
	var err error
	p.WebRouterPrefix = prefix
	err = p.StandAloneWebServer.Init(options)
	if err != nil {
		return err
	}

	p.StandAloneWebServer.Router(prefix+"/*", p.WebServe)
	return nil
}

func (p *Proxy) InitAttachModeWebServer(prefix string, webServer *Server) error {
	p.WebRouterPrefix = prefix
	p.AttachModeWebServer = webServer
	p.AttachModeWebServer.Router(prefix+"/*", p.WebServe)
	return nil
}

func (p *Proxy) StandAloneWebServerServe() error {
	return p.StandAloneWebServer.Serve()
}

func (p *Proxy) WebServe(ir *Request) {
	var path = ir.R.URL.Path[len(p.WebRouterPrefix):]
	var reqCtx RequestContext
	var resp, err = p.DispatchWithIronRequest(path, &reqCtx, ir).Resolve()
	if err != nil {
		ir.ApiOutput(nil, CODE_ERR, err.Error())
		return
	}
	ir.ApiOutput(resp, CODE_OK, "")
}
