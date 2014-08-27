package winformat

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/masayukioguni/bcd"
	"log"
	"strconv"
	"time"
)

type WinFormat struct {
	sequence    uint8
	subSequence uint8
	A0          byte
	length      uint16
	year        byte
	month       byte
	day         byte
	hour        byte
	minute      byte
	second      byte
	channel     uint16
	rate        byte
	size        byte
	FirstSample int32
	Sampling    []int32
}

func (winFormat *WinFormat) GetSequence() int {
	return int(winFormat.sequence)
}

func (winFormat *WinFormat) GetSubSequence() int {
	return int(winFormat.subSequence)
}

func (winFormat *WinFormat) GetLength() int {
	return int(winFormat.length)
}

func (winFormat *WinFormat) GetChannel() int {
	return int(winFormat.channel)
}

func (winFormat *WinFormat) GetSamplingSize() int {
	return int(winFormat.size >> 4)
}

func (winFormat *WinFormat) GetSamplingRate() int {
	return int((winFormat.size&0x0f)<<8 | winFormat.rate&0xff)
}

func (winFormat *WinFormat) GetDateTime() string {
	datetime := fmt.Sprintf("20%02d-%02d-%02d %02d:%02d:%02d +09:00",
		bcd.BcdToInt(int(winFormat.year)),
		bcd.BcdToInt(int(winFormat.month)),
		bcd.BcdToInt(int(winFormat.day)),
		bcd.BcdToInt(int(winFormat.hour)),
		bcd.BcdToInt(int(winFormat.minute)),
		bcd.BcdToInt(int(winFormat.second)))
	t, err := time.Parse(
		"2006-01-02 15:04:05 +09:00", // スキャンフォーマット
		datetime)                     // パースしたい文字列
	if err != nil {
		log.Fatal(err)
	}

	return t.Format(time.RFC3339)
}

func Parse(buffer []byte) *WinFormat {
	winformat := WinFormat{}
	buf := bytes.NewBuffer(buffer)
	binary.Read(buf, binary.BigEndian, &winformat.sequence)
	binary.Read(buf, binary.BigEndian, &winformat.subSequence)
	binary.Read(buf, binary.BigEndian, &winformat.A0)
	binary.Read(buf, binary.BigEndian, &winformat.length)
	binary.Read(buf, binary.BigEndian, &winformat.year)
	binary.Read(buf, binary.BigEndian, &winformat.month)
	binary.Read(buf, binary.BigEndian, &winformat.day)
	binary.Read(buf, binary.BigEndian, &winformat.hour)
	binary.Read(buf, binary.BigEndian, &winformat.minute)
	binary.Read(buf, binary.BigEndian, &winformat.second)
	binary.Read(buf, binary.BigEndian, &winformat.channel)
	binary.Read(buf, binary.BigEndian, &winformat.size)
	binary.Read(buf, binary.BigEndian, &winformat.rate)
	binary.Read(buf, binary.BigEndian, &winformat.FirstSample)

	size := winformat.size >> 4
	rate := (winformat.size&0x0f)<<8 | winformat.rate&0xff

	winformat.Sampling = make([]int32, rate-1)

	/*
		fmt.Printf("%d %d %d %d %d %d\n", int(winformat.year),
			int(winformat.month),
			int(winformat.day),
			int(winformat.hour),
			int(winformat.minute),
			int(winformat.second))
		fmt.Printf("%0X%0X %X %04X %02d%02d%02d%02d%02d%02d %04X %02X%02X(%d %d) %08X\n", winformat.sequence,
			winformat.subSequence,
			int(winformat.A0),
			winformat.length,
			bcd.BcdToInt(int(winformat.year)),
			bcd.BcdToInt(int(winformat.month)),
			bcd.BcdToInt(int(winformat.day)),
			bcd.BcdToInt(int(winformat.hour)),
			bcd.BcdToInt(int(winformat.minute)),
			bcd.BcdToInt(int(winformat.second)),
			winformat.channel,
			winformat.size,
			winformat.rate,
			size,
			rate,
			winformat.FirstSample,
		)
	*/
	if size == 0 {
		rate = rate / 2

	}

	for i := 0; i < int(rate)-1; i++ {
		if size == 4 {
			var value int32
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = value
		}
		if size == 3 {
			value := make([]byte, 3)
			binary.Read(buf, binary.BigEndian, &value)
			num, _ := strconv.ParseInt(fmt.Sprintf("%x%x%x", value[0], value[1], value[2]), 16, 32)
			winformat.Sampling[i] = int32(num)
		}
		if size == 2 {
			var value int16
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = int32(value)

		}
		if size == 1 {
			var value int8
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = int32(value)

		}
		if size == 0 {
			var value int8
			binary.Read(buf, binary.BigEndian, &value)
			winformat.Sampling[i] = int32(value)
		}
	}

	return &winformat
}
