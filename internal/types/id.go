package types

import "strconv"

type Id uint

func IdFromUint(id uint8) Id {
	return Id(id)
}
func (i Id) ToUint() uint {
	return uint(i)
}

func (i Id) String() string {
	return strconv.Itoa(int(i))
}
