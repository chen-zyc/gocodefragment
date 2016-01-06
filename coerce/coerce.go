package main

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	fmt.Println("-------- bool --------")
	fmt.Println(Coerce(false, true, ""))   // false <nil>
	fmt.Println(Coerce("false", true, "")) // false <nil>
	fmt.Println(Coerce(1, true, ""))       // false <nil>

	fmt.Println("-------- int64 --------")
	fmt.Println(Coerce(1234, 1, ""))   // 1234 <nil>
	fmt.Println(Coerce("5678", 1, "")) // 5678 <nil>

	fmt.Println("-------- duration --------")
	fmt.Println(Coerce("12", time.Second, "1s"))        // 12s <nil>
	fmt.Println(Coerce("12s", time.Second, ""))         // 12s <nil>
	fmt.Println(Coerce(12, time.Second, ""))            // 12ms <nil>
	fmt.Println(Coerce(2*time.Minute, time.Second, "")) // 2m0s <nil>

	fmt.Println("-------- string slice --------")
	fmt.Println(Coerce("a,b,c", []string{}, ""))                      // [a b c] <nil>
	fmt.Println(Coerce([]string{"a", "b", "c"}, []string{}, ""))      // [a b c] <nil>
	fmt.Println(Coerce([]interface{}{"a", "b", "c"}, []string{}, "")) // [a b c] <nil>

	fmt.Println("-------- float64 slice --------")
	fmt.Println(Coerce("1.2, 2.3, 3.4", []float64{}, ""))               // [1.2 2.3 3.4] <nil>
	fmt.Println(Coerce([]interface{}{1.2, 2.3, 3.4}, []float64{}, ""))  // [1.2 2.3 3.4] <nil>
	fmt.Println(Coerce([]string{"1.2", "2.3", "3.4"}, []float64{}, "")) // [1.2 2.3 3.4] <nil>
	fmt.Println(Coerce([]float64{1.2, 2.3, 3.4}, []float64{}, ""))      // [1.2 2.3 3.4] <nil>

	fmt.Println("-------- string --------")
	fmt.Println(Coerce("hello", "", "")) // hello <nil>
}

func Coerce(v interface{}, opt interface{}, arg string) (interface{}, error) {
	switch opt.(type) {
	case bool:
		return coerceBool(v)
	case int:
		i, err := coerceInt64(v)
		if err != nil {
			return nil, err
		}
		return int(i), nil
	case int16:
		i, err := coerceInt64(v)
		if err != nil {
			return nil, err
		}
		return int16(i), nil
	case uint16:
		i, err := coerceInt64(v)
		if err != nil {
			return nil, err
		}
		return uint16(i), nil
	case int32:
		i, err := coerceInt64(v)
		if err != nil {
			return nil, err
		}
		return int32(i), nil
	case uint32:
		i, err := coerceInt64(v)
		if err != nil {
			return nil, err
		}
		return uint32(i), nil
	case int64:
		return coerceInt64(v)
	case uint64:
		i, err := coerceInt64(v)
		if err != nil {
			return nil, err
		}
		return uint64(i), nil
	case string:
		return coerceString(v)
	case time.Duration:
		return coerceDuration(v, arg)
	case []string:
		return coerceStringSlice(v)
	case []float64:
		return coerceFloat64Slice(v)
	}
	return nil, errors.New("invalid type")
}

func coerceBool(v interface{}) (bool, error) {
	switch v.(type) {
	case bool:
		return v.(bool), nil
	case string:
		return strconv.ParseBool(v.(string))
	case int, int16, uint16, int32, uint32, int64, uint64:
		return reflect.ValueOf(v).Int() == 0, nil
	}
	return false, errors.New("invalid value type")
}

func coerceInt64(v interface{}) (int64, error) {
	switch v.(type) {
	case string:
		return strconv.ParseInt(v.(string), 10, 64)
	case int, int16, uint16, int32, uint32, int64, uint64:
		return reflect.ValueOf(v).Int(), nil
	}
	return 0, errors.New("invalid value type")
}

func coerceDuration(v interface{}, arg string) (time.Duration, error) {
	switch v.(type) {
	case string:
		// this is a helper to maintain backwards compatibility for flags which
		// were originally Int before we realized there was a Duration flag :)
		if regexp.MustCompile(`^[0-9]+$`).MatchString(v.(string)) {
			intVal, err := strconv.Atoi(v.(string))
			if err != nil {
				return 0, err
			}
			mult, err := time.ParseDuration(arg)
			if err != nil {
				return 0, err
			}
			return time.Duration(intVal) * mult, nil
		}
		return time.ParseDuration(v.(string))
	case int, int16, uint16, int32, uint32, int64, uint64:
		// treat like ms
		return time.Duration(reflect.ValueOf(v).Int()) * time.Millisecond, nil
	case time.Duration:
		return v.(time.Duration), nil
	}
	return 0, errors.New("invalid value type")
}

func coerceStringSlice(v interface{}) ([]string, error) {
	var tmp []string
	switch v.(type) {
	case string:
		for _, s := range strings.Split(v.(string), ",") {
			tmp = append(tmp, s)
		}
	case []interface{}:
		for _, si := range v.([]interface{}) {
			tmp = append(tmp, si.(string))
		}
	case []string:
		tmp = v.([]string)
	}
	return tmp, nil
}

func coerceFloat64Slice(v interface{}) ([]float64, error) {
	var tmp []float64
	switch v.(type) {
	case string:
		for _, s := range strings.Split(v.(string), ",") {
			f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
			if err != nil {
				return nil, err
			}
			tmp = append(tmp, f)
		}
	case []interface{}:
		for _, fi := range v.([]interface{}) {
			tmp = append(tmp, fi.(float64))
		}
	case []string:
		for _, s := range v.([]string) {
			f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
			if err != nil {
				return nil, err
			}
			tmp = append(tmp, f)
		}
	case []float64:
		tmp = v.([]float64)
	}
	return tmp, nil
}

func coerceString(v interface{}) (string, error) {
	switch v.(type) {
	case string:
		return v.(string), nil
	}
	return fmt.Sprintf("%s", v), nil
}
