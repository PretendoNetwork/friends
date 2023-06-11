package types

type NEXToken struct {
	SystemType  uint8
	TokenType   uint8
	UserPID     uint32
	ExpireTime  uint64
	TitleID     uint64
	AccessLevel int8
}
