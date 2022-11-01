package ByteCont

type ByteView struct {
	b []byte
}

func (b *ByteView) Len() int {
	return len(b.b)
}

func (b *ByteView) String() string {
	return string(b.b)
}

func (b *ByteView) GetSlice() []byte {
	return b.b
}

func (b *ByteView) GetCopySlice() []byte {
	cp := make([]byte, b.Len())
	copy(cp, b.b)
	return cp
}
