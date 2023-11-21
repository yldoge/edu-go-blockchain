package types

import (
	"encoding/hex"
	"fmt"
)

type Address [20]uint8

func (a Address) ToSlice() []byte {
	res := make([]byte, 20)
	for i := 0; i < 20; i++ {
		res[i] = a[i]
	}
	return res
}

func (a Address) String() string {
	return hex.EncodeToString(a.ToSlice())
}

func AddressFromBytes(b []byte) Address {
	if len(b) != 20 {
		msg := fmt.Sprintf("given bytes with length %d should be 20", len(b))
		panic(msg)
	}
	return Address(b)
}
