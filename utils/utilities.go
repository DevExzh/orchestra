package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"unsafe"
)

func FetchFileIgnoreExt(path string, name string) (string, error) {
	var str string
	if wd, err := os.Getwd(); err != nil {
		return "", err
	} else {
		str = wd
	}
	if f, err := os.Open(str); err != nil {
		return "", err
	} else {
		if fs, e := f.ReadDir(-1); e != nil {
			return "", e
		} else {
			for _, entry := range fs {
				if i := strings.LastIndex(entry.Name(), "."); i != -1 {
					if strings.EqualFold(entry.Name()[:i], name) {
						if abs, er := filepath.Abs(entry.Name()); er != nil {
							return "", er
						} else {
							return abs, nil
						}
					}
				}
			}
		}
	}
	return "", os.ErrNotExist
}

func PrependString(slice []string, elements ...string) []string {
	var n []string
	n = append(n, elements...)
	n = append(n, slice...)
	return n
}

func Count(b []byte, sep byte) int {
	var c int = 0
	for _, y := range b {
		if y == sep {
			c++
		}
	}
	return c
}

func Index(b []byte, s byte) int {
	for i, y := range b {
		if y == s {
			return i
		}
	}
	return -1
}

func Equal(a []byte, b []byte) bool {
	if l := len(a); l != len(b) {
		return false
	} else {
		for i := 0; i < l; i++ {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	}
}

func Split(b []byte, sep byte) [][]byte {
	var n = Count(b, sep) + 1
	if n > len(b)+1 {
		n = len(b) + 1
	}
	a := make([][]byte, n)
	n--
	i := 0
	for i < n {
		m := Index(b, sep)
		if m < 0 {
			break
		}
		a[i] = b[:m]
		b = b[m+1:]
		i++
	}
	a[i] = b
	return a[:i+1]
}

func LastIndex(b []byte, str byte) int {
	for i := len(b); i > 0; i-- {
		if b[i-1] == str {
			return i - 1
		}
	}
	return -1
}

func HasPrefix(b []byte, prefix string) bool {
	var a = ConvertStringToByteSlice(prefix)
	for i, v := range b {
		if a[i] != v {
			return false
		}
		if i == len(a) {
			break
		}
	}
	return true
}

func ConvertByteSliceToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ConvertStringToByteSlice(str string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&str))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func AppendUnit(bytes uint64, bit bool, bin bool, precision int, addSpace bool) string {
	var ret = bytes
	var suffix string = "iB"
	var f uint64 = 1024
	if bit {
		ret <<= 3 // times eight to calculate how many bits
		if bin {
			f = 1000
			suffix = "ib"
		} else {
			suffix = "b"
		}
	} else {
		if !bin {
			suffix = "B"
		}
	}
	format := func(b uint64, div uint64, scale byte, suf string, prec int, addSpace bool) string {
		s := ConvertStringToByteSlice(strconv.FormatFloat(float64(bytes)/float64(div), 'f', prec, 64))
		if addSpace {
			s = append(s, ' ')
		}
		s = append(s, scale)
		s = append(s, suf...)
		return ConvertByteSliceToString(s)
	}
	if ret < f {
		return ConvertByteSliceToString(append(ConvertStringToByteSlice(strconv.FormatFloat(float64(bytes), 'f', -1, 64)), suffix...))
	} else { // Kilo
		if ret < f*f {
			return format(bytes, f, 'K', suffix, precision, addSpace)
		} else { // Mega
			if ret < f*f*f {
				return format(bytes, f*f, 'M', suffix, precision, addSpace)
			} else { // Giga
				if ret < f*f*f {
					return format(bytes, f*f*f, 'G', suffix, precision, addSpace)
				} else { // Tera
					return format(bytes, f*f*f*f, 'T', suffix, precision, addSpace)
				}
			}
		}
	}
}
