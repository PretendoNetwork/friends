package types

type NEXToken struct {
	SystemType  uint8
	TokenType   uint8
	UserPID     uint32
	AccessLevel uint8
	TitleID     uint64
	ExpireTime  uint64
}
