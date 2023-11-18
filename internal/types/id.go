package types

import "strconv"

type Id uint

func IdFromUint(id uint8) Id {
	return Id(id)
}
func (i Id) ToUint() uint {
	return uint(i)
}

func IdFromString(id string) (Id, error) {
	val, err := strconv.ParseUint(id, 10, 32)
	return Id(uint(val)), err
}

func (i Id) String() string {
	return strconv.Itoa(int(i))
}
