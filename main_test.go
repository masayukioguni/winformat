package winformat

import (
	//"fmt"
	. "github.com/r7kamura/gospel"
	"os"
	"testing"
)

func ReadBinaryFile(filename string) []byte {
	file, _ := os.Open(filename)
	b := make([]byte, 1400)
	file.Read(b)
	return b
}

func TestDescribe(t *testing.T) {
	Describe(t, "Parse WinFormat Data", func() {
		Context("when SamplingSize is 4", func() {
			winFormat := NewWinFormat(ReadBinaryFile("./testdata/4.bin"))

			It("Sequence is 0 ", func() {
				Expect(winFormat.Sequence).To(Equal, uint8(0))
			})

			It("SubSequence is 0", func() {
				Expect(winFormat.SubSequence).To(Equal, uint8(0))
			})

			It("channel is 0", func() {
				Expect(winFormat.Channel).To(Equal, uint16(0))
			})

			It("SamplingRate is 200", func() {
				Expect(winFormat.Rate).To(Equal, uint16(200))
			})

			It("Size is 4", func() {
				Expect(winFormat.Size).To(Equal, uint16(4))
			})

			It("datetime is 1409050412", func() {
				Expect(winFormat.Datetime).To(Equal, int64(1409050412))
			})

			It("firstsample is ", func() {
				Expect(winFormat.FirstSample).To(Equal, int32(447712320))
			})

			It("Sampleing Length is 199 ", func() {
				Expect(len(winFormat.Sampling)).To(Equal, 199)
			})
		})

		Context("when SamplingSize is 3", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/3.bin"))
			It("Sequence is 0 ", func() {
				Expect(winFormat.Sequence).To(Equal, uint8(2))
			})

			It("SubSequence is 0", func() {
				Expect(winFormat.SubSequence).To(Equal, uint8(2))
			})

			It("channel is 0", func() {
				Expect(winFormat.Channel).To(Equal, uint16(1))
			})

			It("SamplingRate is 200", func() {
				Expect(winFormat.Rate).To(Equal, uint16(100))
			})

			It("Size is 3", func() {
				Expect(winFormat.Size).To(Equal, uint16(3))
			})

			It("datetime is 1409050412", func() {
				Expect(winFormat.Datetime).To(Equal, int64(1409111848))
			})

			It("Sampleing Length is 99 ", func() {
				Expect(len(winFormat.Sampling)).To(Equal, 99)
			})

		})

		Context("when SamplingSize is 2", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/2.bin"))

			It("Size is 2", func() {
				Expect(winFormat.Size).To(Equal, uint16(2))
			})
		})

		Context("when SamplingSize is 1", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/1.bin"))

			It("Size is 1", func() {
				Expect(winFormat.Size).To(Equal, uint16(1))
			})
		})

		Context("when SamplingSize is 0", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/0.bin"))

			It("Size is 0", func() {
				Expect(winFormat.Size).To(Equal, uint16(0))
			})

			It("Sampleing Length is 99 ", func() {
				Expect(len(winFormat.Sampling)).To(Equal, 99)
			})
		})

	})

}
