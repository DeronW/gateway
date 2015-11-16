package protocol

type Command interface {
	GetOp() string
}

type CmdLogin struct {
	op string
}

func (c *CmdLogin) GetOp() string {
	return c.op
}

type CmdCommon struct {
	op string
}

func (c *CmdCommon) GetOp() string {
	return c.op
}
