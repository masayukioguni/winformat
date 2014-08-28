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

type WinPacket struct {
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
	firstSample int32
	Sampling    []int32
}

type WinFormat struct {
	Sequence    uint8
	SubSequence uint8
	A0          byte
	Length      uint16
	Datetime    int64
	Channel     uint16
	Rate        uint16
	Size        uint16
	FirstSample int32
	Sampling    []int32
}

func (winPacket *WinPacket) GetSize() uint16 {
	return uint16(winPacket.size >> 4)
}

func (winPacket *WinPacket) GetRate() uint16 {
	return uint16((winPacket.size&0x0f)<<8 | winPacket.rate&0xff)
}

func (winPacket *WinPacket) GetUnixDateTime() int64 {
	datetime := fmt.Sprintf("20%02d-%02d-%02dT%02d:%02d:%02d+09:00",
		bcd.BcdToInt(int(winPacket.year)),
		bcd.BcdToInt(int(winPacket.month)),
		bcd.BcdToInt(int(winPacket.day)),
		bcd.BcdToInt(int(winPacket.hour)),
		bcd.BcdToInt(int(winPacket.minute)),
		bcd.BcdToInt(int(winPacket.second)))
	t, err := time.Parse(
		time.RFC3339, datetime)
	if err != nil {
		log.Fatal(err)
	}
	return t.Unix()
}

func NewWinFormat(buffer []byte) *WinFormat {
	return Parse(buffer)
}

func Parse(buffer []byte) *WinFormat {
	winPacket := WinPacket{}
	buf := bytes.NewBuffer(buffer)
	binary.Read(buf, binary.BigEndian, &winPacket.sequence)
	binary.Read(buf, binary.BigEndian, &winPacket.subSequence)
	binary.Read(buf, binary.BigEndian, &winPacket.A0)
	binary.Read(buf, binary.BigEndian, &winPacket.length)
	binary.Read(buf, binary.BigEndian, &winPacket.year)
	binary.Read(buf, binary.BigEndian, &winPacket.month)
	binary.Read(buf, binary.BigEndian, &winPacket.day)
	binary.Read(buf, binary.BigEndian, &winPacket.hour)
	binary.Read(buf, binary.BigEndian, &winPacket.minute)
	binary.Read(buf, binary.BigEndian, &winPacket.second)
	binary.Read(buf, binary.BigEndian, &winPacket.channel)
	binary.Read(buf, binary.BigEndian, &winPacket.size)
	binary.Read(buf, binary.BigEndian, &winPacket.rate)
	binary.Read(buf, binary.BigEndian, &winPacket.firstSample)

	winFormat := WinFormat{
		Sequence:    winPacket.sequence,
		SubSequence: winPacket.subSequence,
		Length:      winPacket.length,
		Datetime:    winPacket.GetUnixDateTime(),
		Channel:     winPacket.channel,
		Rate:        winPacket.GetRate(),
		Size:        winPacket.GetSize(),
		FirstSample: winPacket.firstSample,
	}

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

	if winFormat.Size == 0 {
		winFormat.Rate = winFormat.Rate / 2
	}

	winFormat.Sampling = make([]int32, winFormat.Rate-1)

	for i := 0; i < int(winFormat.Rate)-1; i++ {

		switch winFormat.Size {
		case 4:
			var value int32
			binary.Read(buf, binary.BigEndian, &value)
			winFormat.Sampling[i] = value
		case 3:
			value := make([]byte, 3)
			binary.Read(buf, binary.BigEndian, &value)
			num, _ := strconv.ParseInt(fmt.Sprintf("%x%x%x", value[0], value[1], value[2]), 16, 32)
			winFormat.Sampling[i] = int32(num)
		case 2:
			var value int16
			binary.Read(buf, binary.BigEndian, &value)
			winFormat.Sampling[i] = int32(value)

		case 1, 0:
			var value int8
			binary.Read(buf, binary.BigEndian, &value)
			winFormat.Sampling[i] = int32(value)
		}
	}
	/*
			if winFormat.Size == 4 {
				var value int32
				binary.Read(buf, binary.BigEndian, &value)
				winFormat.Sampling[i] = value
			}
			if winFormat.Size == 3 {
				value := make([]byte, 3)
				binary.Read(buf, binary.BigEndian, &value)
				num, _ := strconv.ParseInt(fmt.Sprintf("%x%x%x", value[0], value[1], value[2]), 16, 32)
				winFormat.Sampling[i] = int32(num)
			}
			if winFormat.Size == 2 {
				var value int16
				binary.Read(buf, binary.BigEndian, &value)
				winFormat.Sampling[i] = int32(value)

			}
			if winFormat.Size == 1 {
				var value int8
				binary.Read(buf, binary.BigEndian, &value)
				winFormat.Sampling[i] = int32(value)

			}
			if winFormat.Size == 0 {
				var value int8
				binary.Read(buf, binary.BigEndian, &value)
				winFormat.Sampling[i] = int32(value)
			}
		}
	*/

	return &winFormat
}
