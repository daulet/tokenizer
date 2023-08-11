package tokenizers_test

import (
	_ "embed"
	"math/rand"
	"testing"

	"github.com/daulet/tokenizers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//go:embed test/data/sentence-transformers-labse.json
var embeddedBytes []byte

// TODO test for leaks

func TestInvalidConfigPath(t *testing.T) {
	_, err := tokenizers.FromFile("./non-existent.json")
	require.Error(t, err)
}

func TestEmbeddingConfig(t *testing.T) {
	tk, err := tokenizers.FromBytes(embeddedBytes)
	require.NoError(t, err)
	defer tk.Close()

	tests := []struct {
		name       string
		str        string
		addSpecial bool
		wantIDs    []uint32
		wantTokens []string
	}{
		{
			name:       "without special tokens",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: false,
			wantIDs:    []uint32{0xca3f, 0x2f304, 0x5185b, 0x3c54, 0x3a89, 0x35fc3, 0x57b4},
			wantTokens: []string{"brown", "fox", "jumps", "over", "the", "lazy", "dog"},
		},
		{
			name:       "with special tokens",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: true,
			wantIDs:    []uint32{0x65, 0xca3f, 0x2f304, 0x5185b, 0x3c54, 0x3a89, 0x35fc3, 0x57b4, 0x66},
			wantTokens: []string{"[CLS]", "brown", "fox", "jumps", "over", "the", "lazy", "dog", "[SEP]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodeRes := tk.Encode(tt.str, tt.addSpecial, false, false)
			assert.Equal(t, tt.wantIDs, encodeRes.TokenIds)
			assert.Equal(t, tt.wantTokens, encodeRes.Tokens)
		})
	}
}

func TestEncode(t *testing.T) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(t, err)
	defer tk.Close()
	tests := []struct {
		name       string
		str        string
		addSpecial bool
		wantIDs    []uint32
		wantTokens []string
	}{
		{
			name:       "without special tokens",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: false,
			wantIDs:    []uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899},
			wantTokens: []string{"brown", "fox", "jumps", "over", "the", "lazy", "dog"},
		},
		{
			name:       "with special tokens",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: true,
			wantIDs:    []uint32{101, 2829, 4419, 14523, 2058, 1996, 13971, 3899, 102},
			wantTokens: []string{"[CLS]", "brown", "fox", "jumps", "over", "the", "lazy", "dog", "[SEP]"},
		},
		{
			name:       "empty string",
			str:        "",
			addSpecial: false,
		},
		{
			name:       "empty string with special tokens",
			str:        "",
			addSpecial: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encodeRes := tk.Encode(tt.str, tt.addSpecial, false, false)
			assert.Equal(t, tt.wantIDs, encodeRes.TokenIds)
			assert.Equal(t, tt.wantTokens, encodeRes.Tokens)
		})
	}
}

func TestEncodeBatch(t *testing.T) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(t, err)
	defer tk.Close()

	tests := []struct {
		name       string
		str        string
		addSpecial bool
		wantIDs    []uint32
		wantTokens []string
	}{
		{
			name:       "without special tokens-1",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: false,
			wantIDs:    []uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899},
			wantTokens: []string{"brown", "fox", "jumps", "over", "the", "lazy", "dog"},
		},
		{
			name:       "without special tokens-2",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: false,
			wantIDs:    []uint32{2829, 4419, 14523, 2058, 1996, 13971, 4937},
			wantTokens: []string{"brown", "fox", "jumps", "over", "the", "lazy", "cat"},
		},
	}

	for i, tt := range tk.EncodeBatch([]string{"brown fox jumps over the lazy dog", "brown fox jumps over the lazy cat"}, false, false, false) {
		assert.Equal(t, tests[i].wantIDs, tt.TokenIds)
		assert.Equal(t, tests[i].wantTokens, tt.Tokens)
	}
}

func TestEncodeWithTruncation(t *testing.T) {
	tests := []struct {
		name       string
		str        string
		addSpecial bool
		maxLen     int
		dir        tokenizers.TruncationDirection
		wantIDs    []uint32
		wantTokens []string
	}{
		{
			name:       "without special tokens, left truncation",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: false,
			maxLen:     5,
			dir:        tokenizers.TruncationDirectionLeft,
			wantIDs:    []uint32{0x5185b, 0x3c54, 0x3a89, 0x35fc3, 0x57b4},
			wantTokens: []string{"jumps", "over", "the", "lazy", "dog"},
		},
		{
			name:       "without special tokens, right truncation",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: false,
			maxLen:     5,
			dir:        tokenizers.TruncationDirectionRight,
			wantIDs:    []uint32{0xca3f, 0x2f304, 0x5185b, 0x3c54, 0x3a89},
			wantTokens: []string{"brown", "fox", "jumps", "over", "the"},
		},
		{
			name:       "with special tokens, left truncation",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: true,
			maxLen:     5,
			dir:        tokenizers.TruncationDirectionLeft,
			wantIDs:    []uint32{0x65, 0x3a89, 0x35fc3, 0x57b4, 0x66},
			wantTokens: []string{"[CLS]", "the", "lazy", "dog", "[SEP]"},
		},
		{
			name:       "with special tokens, right truncation",
			str:        "brown fox jumps over the lazy dog",
			addSpecial: true,
			maxLen:     5,
			dir:        tokenizers.TruncationDirectionRight,
			wantIDs:    []uint32{0x65, 0xca3f, 0x2f304, 0x5185b, 0x66},
			wantTokens: []string{"[CLS]", "brown", "fox", "jumps", "[SEP]"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tk, err := tokenizers.FromBytesWithTruncation(embeddedBytes, uint32(tt.maxLen), tt.dir)
			require.NoError(t, err)
			defer tk.Close()

			tk.Encode(tt.str, tt.addSpecial, false, false)
			encodeRes := tk.Encode(tt.str, tt.addSpecial, false, false)
			assert.Equal(t, tt.wantIDs, encodeRes.TokenIds)
			assert.Equal(t, tt.wantTokens, encodeRes.Tokens)
		})
	}
}

func TestDecode(t *testing.T) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(t, err)
	defer tk.Close()
	tests := []struct {
		name        string
		tokens      []uint32
		skipSpecial bool
		want        string
	}{
		{
			name:        "without special tokens, skip special tokens",
			tokens:      []uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899},
			skipSpecial: true,
			want:        "brown fox jumps over the lazy dog",
		},
		{
			name:        "with special tokens, skip special tokens",
			tokens:      []uint32{101, 2829, 4419, 14523, 2058, 1996, 13971, 3899, 102},
			skipSpecial: true,
			want:        "brown fox jumps over the lazy dog",
		},
		{
			name:        "without special tokens, don't skip special tokens",
			tokens:      []uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899},
			skipSpecial: false,
			want:        "brown fox jumps over the lazy dog",
		},
		{
			name:        "with special tokens, don't skip special tokens",
			tokens:      []uint32{101, 2829, 4419, 14523, 2058, 1996, 13971, 3899, 102},
			skipSpecial: false,
			want:        "[CLS] brown fox jumps over the lazy dog [SEP]",
		},
		{
			name:        "no tokens",
			tokens:      []uint32{},
			skipSpecial: false,
			want:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tk.Decode(tt.tokens, tt.skipSpecial)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestVocabSize(t *testing.T) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(t, err)
	defer tk.Close()
	assert.Equal(t, uint32(30522), tk.VocabSize())
}

func BenchmarkEncodeNTimes(b *testing.B) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(b, err)
	defer tk.Close()
	expected := []uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encodeRes := tk.Encode("brown fox jumps over the lazy dog", false, false, false)
		assert.Equal(b, expected, encodeRes.TokenIds)
	}
}

func BenchmarkEncodeNChars(b *testing.B) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(b, err)
	defer tk.Close()
	input := make([]rune, 0, b.N)
	for i := 0; i < b.N; i++ {
		input = append(input, rune(rand.Uint32()%tk.VocabSize()))
	}
	str := string(input)
	b.ResetTimer()
	encodeRes := tk.Encode(str, false, false, false)
	assert.Greater(b, len(encodeRes.TokenIds), 0)
}

func BenchmarkDecodeNTimes(b *testing.B) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(b, err)
	defer tk.Close()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		str := tk.Decode([]uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899}, true)
		assert.Equal(b, "brown fox jumps over the lazy dog", str)
	}
}

func BenchmarkDecodeNTokens(b *testing.B) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(b, err)
	defer tk.Close()
	input := make([]uint32, 0, b.N)
	for i := 0; i < b.N; i++ {
		input = append(input, rand.Uint32()%tk.VocabSize())
	}
	b.ResetTimer()
	text := tk.Decode(input, true)
	// a token is one or more characters
	assert.Greater(b, len(text), b.N)
}
