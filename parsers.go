package swos_client

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func bootToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func fixJson(in string) string {
	re1 := regexp.MustCompile("([{,])([a-zA-Z][a-zA-Z0-9]+)")
	re2 := regexp.MustCompile("'")
	re3 := regexp.MustCompile("(0x[0-9a-zA-Z]+)")

	res := re1.ReplaceAllString(in, "$1\"$2\"")
	res = re2.ReplaceAllString(res, "\"")
	res = re3.ReplaceAllString(res, "\"$1\"")
	return res
}

func parseInt(in string) (int, error) {
	var val int
	_, err := fmt.Sscanf(in, "%v", &val)
	return val, err
}

func parseBool(in string) (bool, error) {
	v, er := parseInt(in)
	if er != nil {
		return false, er
	}
	return v != 0, nil
}

func parseInt64(in string) (int64, error) {
	var val int64
	_, err := fmt.Sscanf(in, "%v", &val)
	return val, err
}

func bitMaskToArray(in string, size int) ([]bool, error) {
	val, err := parseInt(in)
	if err != nil {
		return nil, err
	}
	result := make([]bool, size)
	for i := 0; i < size; i++ {
		result[i] = (val & (1 << i)) != 0
	}
	return result, nil
}

func arrayToBitMask(in []bool) int {
	result := 0

	for i, v := range in {
		if v {
			result |= 1 << i
		}
	}

	return result
}

func macFromMikrotik(in string) (net.HardwareAddr, error) {
	out := ""
	if len(in) != 12 {
		return nil, errors.New("invalid mac address")
	}
	for in != "" {
		hexv := in[0:2]
		in = in[2:]
		if out != "" {
			out = out + ":" + hexv
		} else {
			out = hexv
		}
	}
	return net.ParseMAC(out)
}

func macToMikrotik(addr net.HardwareAddr) string {
	out := ""
	in := addr.String()

	for in != "" {
		hexv := in[0:2]
		in = in[3:]
		out = fmt.Sprintf("%s:%s", out, hexv)
	}
	return out
}

func ipFromMikrotik(in string) (net.IP, error) {
	i, err := parseInt(in)
	if err != nil {
		return nil, err
	}
	return net.IPv4(byte(i), byte(i>>8), byte(i>>16), byte(i>>24)), nil
}

func ipToMikrotik(in net.IP) int {
	return int(in[0]) | int(in[1])<<8 | int(in[2])<<16 | int(in[3])<<24
}

func stringFromMikrotik(in string) (string, error) {
	out := ""
	if len(in)%2 != 0 {
		return "", errors.New("invalid port name")
	}
	for in != "" {
		hexv := in[0:2]
		in = in[2:]
		code, err := strconv.ParseInt(hexv, 16, 8)
		if err != nil {
			return "", err
		}
		out = out + string(rune(code))
	}
	return out, nil
}

func stringToMikrotik(in string) string {
	res := ""
	for _, c := range in {
		res = fmt.Sprintf("%s%02x", res, int(c))
	}
	return res
}

func arrayToMikrotik(value reflect.Value) string {
	var result = "["

	for j := 0; j < value.Len(); j++ {
		if j != 0 {
			result = result + ","
		}
		result = result + reflectToMikrotik(value.Index(j))
	}

	return result + "]"
}

func structToMikrotik(v reflect.Value) string {
	result := "{"
	for i := 0; i < v.NumField(); i++ {
		if i != 0 {
			result = result + ","
		}
		result = result + fmt.Sprintf("%s:", strings.ToLower(v.Type().Field(i).Name))

		result = result + reflectToMikrotik(v.Field(i))
	}

	return result + "}"
}

func toAnySlice(v reflect.Value) []any {
	res := make([]any, v.Len())
	for i := 0; i < v.Len(); i++ {
		res[i] = v.Index(i)
	}
	return res
}

func anyToMikrotik(value any) string {
	return reflectToMikrotik(reflect.ValueOf(value))
}

func reflectToMikrotik(v reflect.Value) string {
	if v.Type().Kind() == reflect.Slice {
		return arrayToMikrotik(v)
	}
	if v.Type().Kind() == reflect.Struct {
		return structToMikrotik(v)
	}
	if v.Type().Kind() == reflect.Int {
		return fmt.Sprintf("0x%02x", v.Int())
	}
	if v.Type().Kind() == reflect.Bool {
		iv := 0
		if v.Bool() {
			iv = 1
		}
		return fmt.Sprintf("0x%02x", iv)
	}
	if v.Type().Kind() == reflect.String {
		return fmt.Sprintf("'%s'", v.String())
	}

	panic("Unsupported type")
}
