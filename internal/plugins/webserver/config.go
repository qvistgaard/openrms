package webserver

type Config struct {
	Webserver struct {
		Listen string
	}
}

/*func CreateFromConfig(ctx *application.Context) (*Server, error) {
	c := &Config{}
	err := mapstructure.Decode(ctx.Config, c)
	if err != nil {
		return nil, err
	}

	return NewServer(c.Webserver.Listen, ctx), nil
}
*/