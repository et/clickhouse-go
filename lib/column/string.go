package column

import (
	"fmt"

	"github.com/ClickHouse/clickhouse-go/v2/lib/binary"
)

type String []string

func (String) Type() Type {
	return "String"
}

func (col *String) Rows() int {
	return len(*col)
}

func (s *String) RowValue(row int) interface{} {
	value := *s
	return value[row]
}

func (s *String) ScanRow(dest interface{}, row int) error {
	v := *s
	switch d := dest.(type) {
	case *string:
		*d = v[row]
	case **string:
		*d = new(string)
		**d = v[row]
	default:
		return &ColumnConverterErr{
			op:   "ScanRow",
			to:   fmt.Sprintf("%T", dest),
			from: "String",
		}
	}
	return nil
}

func (s *String) Append(v interface{}) error {
	switch v := v.(type) {
	case []string:
		*s = append(*s, v...)
	default:
		return &ColumnConverterErr{
			op:   "Append",
			to:   "String",
			from: fmt.Sprintf("%T", v),
		}
	}
	return nil
}

func (s *String) AppendRow(v interface{}) error {
	switch v := v.(type) {
	case string:
		*s = append(*s, v)
	case *string:
		*s = append(*s, *v)
	case null:
		*s = append(*s, "")
	default:
		return &ColumnConverterErr{
			op:   "AppendRow",
			to:   "String",
			from: fmt.Sprintf("%T", v),
		}
	}
	return nil
}

func (s *String) Decode(decoder *binary.Decoder, rows int) error {
	for i := 0; i < int(rows); i++ {
		v, err := decoder.String()
		if err != nil {
			return err
		}
		*s = append(*s, v)
	}
	return nil
}

func (s *String) Encode(encoder *binary.Encoder) error {
	for _, v := range *s {
		if err := encoder.String(v); err != nil {
			return err
		}
	}
	return nil
}

var _ Interface = (*String)(nil)
