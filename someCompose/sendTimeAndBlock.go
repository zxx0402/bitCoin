package someCompose

import "time"

type SendTimeAndBlock struct {
	Time  time.Time
	Block *Block
}

func NewSendTimeAndBlock(b *Block, time time.Time) *SendTimeAndBlock {
	s := &SendTimeAndBlock{
		Time:  time,
		Block: b,
	}
	return s
}
func InitTime() (t time.Time) {
	t = time.Date(2024, time.January, 1, 0, 0, 0, 111111, time.UTC)
	return t
}
