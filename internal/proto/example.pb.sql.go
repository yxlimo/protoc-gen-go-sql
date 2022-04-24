package proto

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
)

var _ sql.Scanner
var _ driver.Valuer
var _ = cast.ToString

// Scan implements sql.Scanner
func (msg *Foo) Scan(src interface{}) error {
	if msg == nil {
		return fmt.Errorf("scan into nil Foo")
	}
	var value []byte
	switch v := src.(type) {
	case []byte:
		value = v
	case string:
		value = []byte(v)
	default:
		return fmt.Errorf("can't convert %v to Foo, unsupported type %T", src, src)
	}

	err := protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: true,
	}.Unmarshal(value, msg)
	if err != nil {
		return fmt.Errorf("can't unmarshal Foo: %w", err)
	}
	return nil
}

// Value implements driver.Valuer
func (msg Foo) Value() (driver.Value, error) {
	value, err := protojson.MarshalOptions{
		Multiline:       false,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		UseProtoNames:   true,
		AllowPartial:    false,
	}.Marshal(&msg)
	if err != nil {
		return nil, fmt.Errorf("can't marshal Foo: %w", err)
	}
	return value, nil
}

// Scan implements sql.Scanner
func (msg *Bar) Scan(src interface{}) error {
	if msg == nil {
		return fmt.Errorf("scan into nil Bar")
	}
	var value []byte
	switch v := src.(type) {
	case []byte:
		value = v
	case string:
		value = []byte(v)
	default:
		return fmt.Errorf("can't convert %v to Bar, unsupported type %T", src, src)
	}

	err := protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: true,
	}.Unmarshal(value, msg)
	if err != nil {
		return fmt.Errorf("can't unmarshal Bar: %w", err)
	}
	return nil
}

// Value implements driver.Valuer
func (msg Bar) Value() (driver.Value, error) {
	value, err := protojson.MarshalOptions{
		Multiline:       false,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		UseProtoNames:   true,
		AllowPartial:    false,
	}.Marshal(&msg)
	if err != nil {
		return nil, fmt.Errorf("can't marshal Bar: %w", err)
	}
	return value, nil
}

// Scan implements sql.Scanner
func (msg *BarEnum) Scan(src interface{}) error {
	if msg == nil {
		return fmt.Errorf("scan into nil BarEnum")
	}
	var value []byte
	switch v := src.(type) {
	case []byte:
		value = v
	case string:
		value = []byte(v)
	default:
		return fmt.Errorf("can't convert %v to BarEnum, unsupported type %T", src, src)
	}

	err := protojson.UnmarshalOptions{
		AllowPartial:   false,
		DiscardUnknown: true,
	}.Unmarshal(value, msg)
	if err != nil {
		return fmt.Errorf("can't unmarshal BarEnum: %w", err)
	}
	return nil
}

// Value implements driver.Valuer
func (msg BarEnum) Value() (driver.Value, error) {
	value, err := protojson.MarshalOptions{
		Multiline:       false,
		UseEnumNumbers:  false,
		EmitUnpopulated: false,
		UseProtoNames:   true,
		AllowPartial:    false,
	}.Marshal(&msg)
	if err != nil {
		return nil, fmt.Errorf("can't marshal BarEnum: %w", err)
	}
	return value, nil
}

func (x *BarEnum_Enum) Scan(src any) error {

	switch v := src.(type) {
	case []byte, string:
		*x = BarEnum_Enum(BarEnum_Enum_value[cast.ToString(v)])
	case int, int8, int16, int32, int64:
		*x = BarEnum_Enum(cast.ToInt32(v))
	default:
		return fmt.Errorf("cannot scan type %T into BarEnum_Enum: %v", src, src)
	}
	return nil
}

func (x BarEnum_Enum) Value() (driver.Value, error) {
	return x.String(), nil
}
