package px

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func pow(i12, i13, i14, i15 int32) int32 {
	i16 := i15 % 10
	i17 := i14 % 10
	if i16 != 0 {
		i17 = i14 % i16
	}
	i18 := i12 * i12
	i19 := i13 * i13
	switch i17 {
	case 0:
		return i13 + i18
	case 1:
		break
	case 2:
		return i13 * i18
	case 3:
		return i13 ^ i12
	case 4:
		return i12 - i19
	case 5:
		i22 := i12 + 783
		i12 = i22 * i22
		break
	case 6:
		return i13 + (i12 ^ i13)
	case 7:
		return i18 - i19
	case 8:
		return i13 * i12
	case 9:
		return (i13 * i12) - i12
	default:
		return -1
	}
	return i12 + i19
}

// AppcInstruction populates the PX257, PX259, and PX256 keys in the payload. Input is the APPC 2 instruction
// ex: appc|2|1688967166574|42443399873c02eadc3ebd747567101f265e835563d18725f73ec88b7012eb8f,42443399873c02eadc3ebd747567101f265e835563d18725f73ec88b7012eb8f|374|3542|3311|1478|3717|352
func (p *Payload) AppcInstruction(appc string) error {
	p.T = "PX329"
	parts := strings.Split(appc, "|")
	if parts[1] != "2" || len(parts) != 10 {
		return errors.New("invalid appc")
	}
	p.D.Px256 = &parts[3]
	date, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return nil
	}
	p.D.Px259 = &date
	parts = parts[4:]
	numbers := []int32{}
	for _, str := range parts {
		num, err := strconv.ParseInt(str, 10, 32)
		if err != nil {
			return err
		}
		numbers = append(numbers, int32(num))
	}
	//aVar2.a(aVar2.a(aVar2.appcNumber3, aVar2.appcNumber4, aVar2.appcNumber1, aVar2.appcNumber6), aVar2.appcNumber5, aVar2.appcNumber2, aVar2.appcNumber6);
	mathOut := pow(pow(numbers[2], numbers[3], numbers[0], numbers[5]), numbers[4], numbers[1], numbers[5])
	bArr := []byte(p.D.Px320)
	res := int32(0)
	if len(bArr) >= 4 {
		binary.Read(bytes.NewBuffer(bArr), binary.BigEndian, &res)
	}
	out := fmt.Sprint(res ^ mathOut)
	p.D.Px257 = &out
	return nil
}
