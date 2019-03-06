package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cast"
)

const (
	AsciiType     = "ascii"
	BigintType    = "bigint"
	BlobType      = "blob"
	BooleanType   = "boolean"
	CounterType   = "counter"
	DateType      = "date"
	DecimalType   = "decimal"
	DoubleType    = "double"
	FloatType     = "float"
	FrozenType    = "frozen"
	InetType      = "inet"
	IntType       = "int"
	ListType      = "list"
	MapType       = "map"
	SetType       = "set"
	SmallintType  = "smallint"
	TextType      = "text"
	TimeType      = "time"
	TimestampType = "timestamp"
	TimeuuidType  = "timeuuid"
	TinyintType   = "tinyint"
	TupleType     = "tuple"
	UuidType      = "uuid"
	VarcharType   = "varchar"
	VarintType    = "varint"
)

// OutputTransformType Map Value 數值轉字串
func OutputTransformType(row map[string]interface{}) map[string]interface{} {
	for k, v := range row {
		switch v.(type) {
		case int64, float64, float32:
			row[k] = cast.ToString(v)
		case []int64:
			row[k] = cast.ToStringSlice(v)
		case map[string]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int64]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int32]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int16]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[int8]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[float64]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[float32]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		case map[bool]int64:
			val, err := cast.FToStringMapStringE(v)

			if err == nil {
				row[k] = val
			}
		}

	}

	return row
}

func cqlFormat(columnName string, columnType string, columnVal interface{}) string {
	switch columnType {
	case BigintType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case BlobType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case BooleanType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case CounterType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case DecimalType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case DoubleType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case FloatType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case FrozenType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case IntType:
		return fmt.Sprintf("%s=%v", columnName, columnVal)
	case ListType:
		b, _ := jsoni.Marshal(columnVal)
		return fmt.Sprintf("%s=%v", columnName, string(b))
	case MapType:
		b, _ := jsoni.Marshal(columnVal)
		return fmt.Sprintf("%s=%v", columnName, string(b))
	case SetType:
		b, _ := jsoni.Marshal(columnVal)
		return fmt.Sprintf("%s=%v", columnName, string(b))
	case AsciiType:
		return fmt.Sprintf("%s='%v'", columnName, columnVal)
	case TextType:
		return fmt.Sprintf("%s='%v'", columnName, columnVal)
	case TimeType:
		return fmt.Sprintf("%s='%v'", columnName, columnVal)
	case DateType:
		return fmt.Sprintf("%s='%v'", columnName, columnVal)
	case InetType:
		return fmt.Sprintf("%s='%v'", columnName, columnVal)
	case VarcharType:
		return fmt.Sprintf("%s='%v'", columnName, columnVal)
	default:
		return ""
	}
}

// InputTransformType 對應table schema型別作轉換
func InputTransformType(item map[string]interface{}, schema map[string]string) ([]string, []interface{}, []string, error) {
	var (
		itemKey         []string
		itemData        []interface{}
		itemPlaceholder []string
	)

	mapReg := regexp.MustCompile(`(?U)^map\<(.+),\s(.+)\>`)

	for k, v := range item {
		switch schema[k] {
		case BigintType:
			val, err := cast.ToInt64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case BooleanType:
			val, err := cast.ToBoolE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case CounterType:
			val, err := cast.ToInt64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case DateType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case DecimalType:
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case DoubleType:
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case FloatType:
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case InetType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case IntType:
			val, err := cast.ToInt32E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case SmallintType:
			val, err := cast.ToInt16E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case TextType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case TimeType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case TimestampType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case TimeuuidType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case TinyintType:
			val, err := cast.ToInt8E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case UuidType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case VarcharType:
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case VarintType:
			val, err := cast.ToInt32E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		default:
			var val interface{} = v
			var err error

			mapRet := mapReg.FindStringSubmatch(schema[k])

			if len(mapRet) == 3 {
				val, err = MapToCassandraMapType(v, mapRet[1], mapRet[2])

				if err != nil {
					return nil, nil, nil, err
				}
			}

			itemData = append(itemData, val)
		}

		itemKey = append(itemKey, k)
		itemPlaceholder = append(itemPlaceholder, "?")
	}

	return itemKey, itemData, itemPlaceholder, nil
}

// MapToCassandraMapType map對應cassandra map的型別
func MapToCassandraMapType(i interface{}, keyType string, valType string) (interface{}, error) {
	var m = map[interface{}]interface{}{}

	switch v := i.(type) {
	case map[string]interface{}:
		for k, val := range v {
			kRet, err := CassandraTypeToGoType(k, keyType)

			if err != nil {
				return nil, err
			}

			valRet, err := CassandraTypeToGoType(val, valType)

			if err != nil {
				return nil, err
			}

			m[kRet] = valRet
		}

		return m, nil
	case string:
		err := JsonStringToObject(v, &m)
		return m, err
	case nil:
		return m, nil
	default:
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]string", i, i)
	}
}

// CassandraTypeToGoType cassandra的型別轉Go型別
func CassandraTypeToGoType(i interface{}, t string) (interface{}, error) {
	switch t {
	case BigintType:
		val, err := cast.ToInt64E(i)
		return val, err
	case BooleanType:
		val, err := cast.ToBoolE(i)

		return val, err
	case CounterType:
		val, err := cast.ToInt64E(i)

		return val, err
	case DateType:
		val, err := cast.ToStringE(i)

		return val, err
	case DecimalType:
		val, err := cast.ToFloat64E(i)

		return val, err
	case DoubleType:
		val, err := cast.ToFloat64E(i)

		return val, err
	case FloatType:
		val, err := cast.ToFloat64E(i)

		return val, err
	case InetType:
		val, err := cast.ToStringE(i)

		return val, err
	case IntType:
		val, err := cast.ToInt32E(i)

		return val, err
	case SmallintType:
		val, err := cast.ToInt16E(i)

		return val, err
	case TextType:
		val, err := cast.ToStringE(i)

		return val, err
	case TimeType:
		val, err := cast.ToStringE(i)

		return val, err
	case TimestampType:
		val, err := cast.ToStringE(i)

		return val, err
	case TimeuuidType:
		val, err := cast.ToStringE(i)

		return val, err
	case TinyintType:
		val, err := cast.ToInt8E(i)

		return val, err
	case UuidType:
		val, err := cast.ToStringE(i)

		return val, err
	case VarcharType:
		val, err := cast.ToStringE(i)

		return val, err
	case VarintType:
		val, err := cast.ToInt32E(i)

		return val, err
	}

	return i, nil
}

// JsonStringToObject json轉obj
func JsonStringToObject(s string, v interface{}) error {
	data := []byte(s)
	return jsoni.Unmarshal(data, v)
}

// CreateTmpFile 建立暫存檔案
func CreateTmpFile(fpath string) error {
	f, err := os.Create(fpath)

	if err != nil {
		return err
	}

	f.Close()

	if err := os.Chmod(fpath, 0666); err != nil {
		return err
	}

	return nil
}
