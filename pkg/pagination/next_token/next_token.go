package next_token

type NextToken struct {
	raw []byte
}

func (t *NextToken) Bytes() []byte {
	return t.raw
}
