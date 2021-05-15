package decode

import (
	"errors"
	"io"
	"reflect"
)

// ServerPacket ...
type ServerPacket interface {
	Parse(io.Reader) error
}

// Fill ....
// Дописать элементарные проверки!!!
func Fill(packet ServerPacket, buffer io.Reader) error {
	//log.Println(arr)
	//buffer := bytes.NewReader(arr)
	obj := reflect.ValueOf(packet).Elem()

	for i := 0; i < obj.NumField(); i++ {

		field := obj.Field(i)
		itemKind := field.Type().Kind()

		switch itemKind {
		case reflect.Uint8:
			v, err := ReadVarUint(buffer, 8)
			if err != nil {
				return err
			}

			field.SetUint(uint64(v))
		case reflect.Int16:
			v, err := ReadVarUint(buffer, 16)
			if err != nil {
				return err
			}

			field.SetInt(int64(v))
		case reflect.Int, reflect.Int32:
			v, err := ReadVarUint(buffer, 32)
			if err != nil {
				return err
			}

			field.SetInt(int64(v))
		case reflect.Int64:
			v, err := ReadVarUint(buffer, 64)
			if err != nil {
				return err
			}

			field.SetInt(int64(v))
		case reflect.Uint64:
			v, err := ReadVarUint(buffer, 64)
			if err != nil {
				return err
			}

			field.SetUint(uint64(v))
		case reflect.String:

			v, err := ReadVarUint(buffer, 16)
			if err != nil {
				return err
			}

			var arr = make([]byte, v)
			buffer.Read(arr)
			str := string(arr)
			field.SetString(str)
		}

	}

	return nil
}

func ReadVarUint(r io.Reader, n uint) (uint64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	p := make([]byte, 1)
	var res uint64
	var shift uint
	for {
		_, err := io.ReadFull(r, p)
		if err != nil {
			return 0, err
		}
		b := uint64(p[0])
		switch {
		case b < 1<<7 && b < 1<<n:
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid uint")
		}
	}
}

func ReadVarint(r io.Reader, n uint) (int64, error) {
	if n > 64 {
		panic(errors.New("leb128: n must <= 64"))
	}
	p := make([]byte, 1)
	var res int64
	var shift uint
	for {
		_, err := io.ReadFull(r, p)
		if err != nil {
			return 0, err
		}
		b := int64(p[0])
		switch {
		case b < 1<<6 && uint64(b) < uint64(1<<(n-1)):
			res += (1 << shift) * b
			return res, nil
		case b >= 1<<6 && b < 1<<7 && uint64(b)+1<<(n-1) >= 1<<7:
			res += (1 << shift) * (b - 1<<7)
			return res, nil
		case b >= 1<<7 && n > 7:
			res += (1 << shift) * (b - 1<<7)
			shift += 7
			n -= 7
		default:
			return 0, errors.New("leb128: invalid int")
		}
	}
}
