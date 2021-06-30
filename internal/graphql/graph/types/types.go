package types

import (
	"encoding/binary"
	"fmt"
	"io"
)

type Int64 int64

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (i *Int64) UnmarshalGQL(v interface{}) error {
	value, ok := v.(int64)
	if !ok {
		return fmt.Errorf("invalid type for Int64")
	}

	*i = Int64(value)

	return nil
}

// MarshalGQL implements the graphql.Marshaler interface
func (i Int64) MarshalGQL(w io.Writer) {
	buf := make([]byte, binary.MaxVarintLen64)
	_ = binary.PutVarint(buf, int64(i))
	_, _ = w.Write(buf)
}
