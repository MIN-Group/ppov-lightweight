package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"time"
	"unsafe"
)

//

func CompressIntSlice(a []int) []uint8{
	var b []uint8
	var tmp uint8 = 0
	for i := 1; i <= len(a); i++ {

		if a[i-1] == 0{
			tmp = tmp << 1
			tmp = tmp | 1
			tmp = tmp << 1
		}else if a[i-1] == 1{
			tmp = tmp << 2
			tmp = tmp | 1
		}else if a[i-1] == -1{
			tmp = tmp << 1
			tmp = tmp | 1
			tmp = tmp << 1
			tmp = tmp | 1
		}

		if i % 3 == 0 || i == len(a){
			b = append(b, tmp)
		}

		if i % 3 == 0{
			tmp = 0
		}
	}

	return b
}

func DecompressToIntSlice(b []uint8) []int{
	var c []int
	for i := 0; i < len(b); i++ {
		var s []int
		bt := b[i]
		for j := 0; j < 6; j+= 2 {
			tmp1, tmp2 := int(bt&1), int(bt&2>>1)
			if tmp1 == 1 && tmp2 == 0{
				s = append(s, 1)
			}else if tmp1 == 1 && tmp2 == 1{
				s = append(s, -1)
			}else if tmp1 == 0 && tmp2 == 1{
				s = append(s, 0)
			}
			bt = bt >> 2
		}
		for k := len(s)-1; k >= 0; k-- {
			c = append(c, s[k])
		}
	}
	return c
}

func SizeStruct(data interface{}) int {
	return sizeof(reflect.ValueOf(data))
}

func sizeof(v reflect.Value) int {
	switch v.Kind() {
	case reflect.Map:
		sum := 0
		keys := v.MapKeys()
		for i := 0; i < len(keys); i++ {
			mapkey := keys[i]
			s := sizeof(mapkey)
			if s < 0 {
				return -1
			}
			sum += s
			s = sizeof(v.MapIndex(mapkey))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum
	case reflect.Slice, reflect.Array:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.String:
		sum := 0
		for i, n := 0, v.Len(); i < n; i++ {
			s := sizeof(v.Index(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Ptr, reflect.Interface:
		p := (*[]byte)(unsafe.Pointer(v.Pointer()))
		if p == nil {
			return 0
		}
		return sizeof(v.Elem())
	case reflect.Struct:
		sum := 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeof(v.Field(i))
			if s < 0 {
				return -1
			}
			sum += s
		}
		return sum

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Int:
		return int(v.Type().Size())

	default:
		fmt.Println("t.Kind() no found:", v.Kind())
	}

	return 0
}

func Uint64ToString(i uint64) string {
	return strconv.FormatUint(i, 10)
}

func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func IntToBytes(i int) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func BytesToInt(buf []byte) int {
	return int(binary.BigEndian.Uint32(buf))
}
func BytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

func BytesToHex(data []byte) (dst string) {
	return hex.EncodeToString(data)
}
func HexToBytes(data string) (dst []byte, err error) {
	return hex.DecodeString(data)
}

//获取时间戳,单位为秒
func GetCurrentTime() float64 {
	return float64(time.Now().UTC().UnixNano()) / 1e9
}

//获取时间戳,单位为毫秒
func GetCurrentTimeMilli() float64 {
	return float64(time.Now().UTC().UnixNano()) / 1e6
}