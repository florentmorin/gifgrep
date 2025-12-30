package gifdecode

import "testing"

func BenchmarkDecodeSmall(b *testing.B) {
	data := makeTestGIF(4)
	opts := DefaultOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Decode(data, opts); err != nil {
			b.Fatalf("decode failed: %v", err)
		}
	}
}

func BenchmarkDecodeMedium(b *testing.B) {
	data := makeMediumGIF()
	opts := DefaultOptions()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := Decode(data, opts); err != nil {
			b.Fatalf("decode failed: %v", err)
		}
	}
}
