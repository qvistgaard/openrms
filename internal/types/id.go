package types

import "strconv"

type CarId uint

func IdFromUint(id uint8) CarId {
	return CarId(id)
}
func (i CarId) ToUint() uint {
	return uint(i)
}

func IdFromString(id string) (CarId, error) {
	val, err := strconv.ParseUint(id, 10, 32)
	return CarId(uint(val)), err
}

func (i CarId) String() string {
	return strconv.Itoa(int(i))
}
