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

func cqlFormatValue(columnType string, columnVal interface{}) interface{} {
	mapReg := regexp.MustCompile(`(?U)^map\<(.+),\s(.+)\>`)

	switch columnType {
	case BigintType:
		val, err := cast.ToInt64E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case BooleanType:
		val, err := cast.ToBoolE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case CounterType:
		val, err := cast.ToInt64E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case DateType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case DecimalType:
		val, err := cast.ToFloat64E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case DoubleType:
		val, err := cast.ToFloat64E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case FloatType:
		val, err := cast.ToFloat64E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case InetType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case IntType:
		val, err := cast.ToInt32E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case SmallintType:
		val, err := cast.ToInt16E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case TextType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case TimeType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case TimestampType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case TimeuuidType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case TinyintType:
		val, err := cast.ToInt8E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case UuidType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case VarcharType:
		val, err := cast.ToStringE(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	case VarintType:
		val, err := cast.ToInt32E(columnVal)
		if err != nil {
			return columnVal
		}
		return val
	default:
		mapRet := mapReg.FindStringSubmatch(columnType)

		if len(mapRet) == 3 {
			val, err := MapToCassandraMapType(columnVal, mapRet[1], mapRet[2])

			if err != nil {
				return columnVal
			}

			return val
		}
	}

	return columnVal
}

func cqlFormatWhere(columnName string, operator string) string {
	return fmt.Sprintf("%s %s ?", columnName, operator)
}

var mapReg = regexp.MustCompile(`(?U)^map\<(.+),\s(.+)\>`)
var listReg = regexp.MustCompile(`(?U)^list\<(.+)>`)

// InputTransformType 對應table schema型別作轉換
func InputTransformType(item map[string]interface{}, schema map[string]string) ([]string, []interface{}, []string, error) {
	var (
		itemKey         []string
		itemData        []interface{}
		itemPlaceholder []string
	)

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

				itemData = append(itemData, val)
				break
			}

			listRet := listReg.FindStringSubmatch(schema[k])

			if len(listRet) == 2 {
				val, err = ListToCassandraListType(v, listRet[1])

				if err != nil {
					return nil, nil, nil, err
				}

				itemData = append(itemData, val)
				break
			}

			itemData = append(itemData, val)
		}

		itemKey = append(itemKey, k)
		itemPlaceholder = append(itemPlaceholder, "?")
	}

	return itemKey, itemData, itemPlaceholder, nil
}

// ListToCassandraListType list對應cassandra list的型別
func ListToCassandraListType(i interface{}, valType string) (interface{}, error) {
	var l = []interface{}{}

	switch v := i.(type) {
	case []interface{}:
		for _, val := range v {
			valRet, err := CassandraTypeToGoType(val, valType)

			if err != nil {
				return nil, err
			}

			l = append(l, valRet)
		}

		return l, nil
	case string:
		err := JsonStringToObject(v, &l)
		return l, err
	case nil:
		return l, nil
	default:
		return l, fmt.Errorf("unable to cast %#v of type %T to []interface{}", i, i)
	}
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
		return m, fmt.Errorf("unable to cast %#v of type %T to map[string]interface{}", i, i)
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
