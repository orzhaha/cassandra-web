package main

import (
	"fmt"
	"regexp"

	"github.com/spf13/cast"
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
		case "bigint":
			val, err := cast.ToInt64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "boolean":
			val, err := cast.ToBoolE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "counter":
			val, err := cast.ToInt64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "date":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "decimal":
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "double":
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "float":
			val, err := cast.ToFloat64E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "inet":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "int":
			val, err := cast.ToInt32E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "smallint":
			val, err := cast.ToInt16E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "text":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "time":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "timestamp":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "timeuuid":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "tinyint":
			val, err := cast.ToInt8E(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "uuid":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "varchar":
			val, err := cast.ToStringE(v)
			if err != nil {
				return nil, nil, nil, err
			}
			itemData = append(itemData, val)
		case "varint":
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
	case "bigint":
		val, err := cast.ToInt64E(i)
		return val, err
	case "boolean":
		val, err := cast.ToBoolE(i)

		return val, err
	case "counter":
		val, err := cast.ToInt64E(i)

		return val, err
	case "date":
		val, err := cast.ToStringE(i)

		return val, err
	case "decimal":
		val, err := cast.ToFloat64E(i)

		return val, err
	case "double":
		val, err := cast.ToFloat64E(i)

		return val, err
	case "float":
		val, err := cast.ToFloat64E(i)

		return val, err
	case "inet":
		val, err := cast.ToStringE(i)

		return val, err
	case "int":
		val, err := cast.ToInt32E(i)

		return val, err
	case "smallint":
		val, err := cast.ToInt16E(i)

		return val, err
	case "text":
		val, err := cast.ToStringE(i)

		return val, err
	case "time":
		val, err := cast.ToStringE(i)

		return val, err
	case "timestamp":
		val, err := cast.ToStringE(i)

		return val, err
	case "timeuuid":
		val, err := cast.ToStringE(i)

		return val, err
	case "tinyint":
		val, err := cast.ToInt8E(i)

		return val, err
	case "uuid":
		val, err := cast.ToStringE(i)

		return val, err
	case "varchar":
		val, err := cast.ToStringE(i)

		return val, err
	case "varint":
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
