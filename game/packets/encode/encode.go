package encode

import (
	"reflect"
)

// ClientPacket ...
type ClientPacket interface {
	Bytes() []byte
}

// Load ...
func Load(packet ClientPacket) []byte {

	buffer := make([]byte, 0)
	strct := reflect.ValueOf(packet).Elem()

	for i := 0; i < strct.NumField(); i++ {

		fieldType := strct.Field(i).Type()
		fieldValue := strct.Field(i).Interface()

		switch fieldType.Name() {
		case "byte", "uint8":
			buffer = WriteNumber(buffer, uint64(fieldValue.(uint8)))
		case "int16":
			buffer = WriteNumber(buffer, uint64(fieldValue.(int16)))
		case "uint16":
			buffer = WriteNumber(buffer, uint64(fieldValue.(uint16)))
		case "int":
			buffer = WriteNumber(buffer, uint64(fieldValue.(int)))
		case "uint":
			buffer = WriteNumber(buffer, uint64(fieldValue.(uint)))
		case "int32":
			buffer = WriteNumber(buffer, uint64(fieldValue.(int32)))
		case "uint32":
			buffer = WriteNumber(buffer, uint64(fieldValue.(uint32)))
		case "int64":
			buffer = WriteNumber(buffer, uint64(fieldValue.(int64)))
		case "uint64":
			buffer = WriteNumber(buffer, fieldValue.(uint64))
		case "string":
			buffer = WriteString(buffer, fieldValue.(string))
		}
	}

	return buffer
}

func WriteNumber(buffer []byte, v uint64) []byte {
	var b []byte

	for {
		c := uint8(v & 0x7f)
		v >>= 7
		if v != 0 {
			c |= 0x80
		}
		b = append(b, c)
		if c&0x80 == 0 {
			break
		}
	}

	return append(buffer, b...)
}

func WriteString(buffer []byte, str string) []byte {

	arr := []byte(str)

	buffer = WriteNumber(buffer, uint64(len(arr)))

	return append(buffer, arr...)
}
