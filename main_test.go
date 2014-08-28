package winformat

import (
	. "github.com/r7kamura/gospel"
	"os"
	"testing"
)

func ReadBinaryFile(filename string) []byte {
	// ファイルを開く
	file, _ := os.Open(filename)
	// ファイルから1バイト読み出し
	b := make([]byte, 1400)
	file.Read(b)
	return b
}

func TestDescribe(t *testing.T) {
	Describe(t, "Parse WinFormat Data", func() {
		Context("when SamplingSize is 4", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/4.bin"))

			It("Sequence is 0 ", func() {
				Expect(winFormat.GetSequence()).To(Equal, 0)
			})

			It("SubSequence is 0", func() {
				Expect(winFormat.GetSubSequence()).To(Equal, 0)
			})

			It("channel is 0", func() {
				Expect(winFormat.GetChannel()).To(Equal, 0)
			})

			It("SamplingRate is 200", func() {
				Expect(winFormat.GetSamplingRate()).To(Equal, 200)
			})

			It("SamplingSize is 4", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 4)
			})

			It("datetime is 1409050412", func() {
				Expect(winFormat.GetUnixDateTime()).To(Equal, int64(1409050412))
			})

			It("firstsample is ", func() {
				Expect(winFormat.GetFirstSample()).To(Equal, 447712320)
			})

		})

		Context("when SamplingSize is 3", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/3.bin"))
			It("Sequence is 2 ", func() {
				Expect(winFormat.GetSequence()).To(Equal, 2)
			})

			It("SubSequence is 2", func() {
				Expect(winFormat.GetSubSequence()).To(Equal, 2)
			})

			It("channel is 1", func() {
				Expect(winFormat.GetChannel()).To(Equal, 1)
			})

			It("SamplingRate is 100", func() {
				Expect(winFormat.GetSamplingRate()).To(Equal, 100)
			})

			It("SamplingSize is 3", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 3)
			})

			It("datetime is 1409111848", func() {
				Expect(winFormat.GetUnixDateTime()).To(Equal, int64(1409111848))
			})
		})

		Context("when SamplingSize is 2", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/2.bin"))
			It("SamplingSize is 2", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 2)
			})
		})

		Context("when SamplingSize is 1", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/1.bin"))
			It("SamplingSize is 1", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 1)
			})
		})

		Context("when SamplingSize is 0", func() {
			winFormat := Parse(ReadBinaryFile("./testdata/0.bin"))
			It("SamplingSize is 0", func() {
				Expect(winFormat.GetSamplingSize()).To(Equal, 0)
			})
		})

	})

}
