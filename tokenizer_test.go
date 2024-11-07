package tokenizers_test

import (
	_ "embed"
	"github.com/daulet/tokenizers"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

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
		name                  string
		str                   string
		addSpecial            bool
		wantIDs               []uint32
		wantTypeIDs           []uint32
		wantTokens            []string
		wantSpecialTokensMask []uint32
		wantAttentionMask     []uint32
		wantOffsets           []tokenizers.Offset
	}{
		{
			name:                  "without special tokens",
			str:                   "brown fox jumps over the lazy dog",
			addSpecial:            false,
			wantIDs:               []uint32{0xca3f, 0x2f304, 0x5185b, 0x3c54, 0x3a89, 0x35fc3, 0x57b4},
			wantTypeIDs:           []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantTokens:            []string{"brown", "fox", "jumps", "over", "the", "lazy", "dog"},
			wantSpecialTokensMask: []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantAttentionMask:     []uint32{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			wantOffsets:           []tokenizers.Offset{{0x0, 0x5}, {0x6, 0x9}, {0xa, 0xf}, {0x10, 0x14}, {0x15, 0x18}, {0x19, 0x1d}, {0x1e, 0x21}},
		},
		{
			name:                  "with special tokens",
			str:                   "brown fox jumps over the lazy dog",
			addSpecial:            true,
			wantIDs:               []uint32{0x65, 0xca3f, 0x2f304, 0x5185b, 0x3c54, 0x3a89, 0x35fc3, 0x57b4, 0x66},
			wantTypeIDs:           []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantTokens:            []string{"[CLS]", "brown", "fox", "jumps", "over", "the", "lazy", "dog", "[SEP]"},
			wantSpecialTokensMask: []uint32{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantAttentionMask:     []uint32{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			wantOffsets:           []tokenizers.Offset{{0x0, 0x0}, {0x0, 0x5}, {0x6, 0x9}, {0xa, 0xf}, {0x10, 0x14}, {0x15, 0x18}, {0x19, 0x1d}, {0x1e, 0x21}, {0x0, 0x0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoding := tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnAllAttributes())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, tt.wantTypeIDs, encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, tt.wantTokens, encoding.Tokens, "wrong tokens")
			assert.Equal(t, tt.wantSpecialTokensMask, encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, tt.wantAttentionMask, encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, tt.wantOffsets, encoding.Offsets, "wrong offsets")

			ids, tokens := tk.Encode(tt.str, tt.addSpecial)
			assert.Equal(t, tt.wantIDs, ids, "wrong ids")
			assert.Equal(t, tt.wantTokens, tokens, "wrong tokens")
		})
	}
}

func TestEncodeWithAndWithoutOptions(t *testing.T) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(t, err)
	defer tk.Close()
	tests := []struct {
		name                  string
		str                   string
		addSpecial            bool
		wantIDs               []uint32
		wantTypeIDs           []uint32
		wantTokens            []string
		wantSpecialTokensMask []uint32
		wantAttentionMask     []uint32
		wantOffsets           []tokenizers.Offset
	}{
		{
			name:                  "without special tokens",
			str:                   "brown fox jumps over the lazy dog",
			addSpecial:            false,
			wantIDs:               []uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899},
			wantTypeIDs:           []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantTokens:            []string{"brown", "fox", "jumps", "over", "the", "lazy", "dog"},
			wantSpecialTokensMask: []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantAttentionMask:     []uint32{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			wantOffsets:           []tokenizers.Offset{{0x0, 0x5}, {0x6, 0x9}, {0xa, 0xf}, {0x10, 0x14}, {0x15, 0x18}, {0x19, 0x1d}, {0x1e, 0x21}},
		},
		{
			name:                  "with special tokens",
			str:                   "brown fox jumps over the lazy dog",
			addSpecial:            true,
			wantIDs:               []uint32{101, 2829, 4419, 14523, 2058, 1996, 13971, 3899, 102},
			wantTypeIDs:           []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantTokens:            []string{"[CLS]", "brown", "fox", "jumps", "over", "the", "lazy", "dog", "[SEP]"},
			wantSpecialTokensMask: []uint32{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1},
			wantAttentionMask:     []uint32{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			wantOffsets:           []tokenizers.Offset{{0x0, 0x0}, {0x0, 0x5}, {0x6, 0x9}, {0xa, 0xf}, {0x10, 0x14}, {0x15, 0x18}, {0x19, 0x1d}, {0x1e, 0x21}, {0x0, 0x0}},
		},
		{
			name:       "empty string",
			str:        "",
			addSpecial: false,
		},
		{
			name:                  "empty string with special tokens",
			str:                   "",
			addSpecial:            true,
			wantTypeIDs:           []uint32{0x0, 0x0},
			wantSpecialTokensMask: []uint32{0x1, 0x1},
			wantAttentionMask:     []uint32{0x1, 0x1},
			wantIDs:               []uint32{101, 102},
			wantTokens:            []string{"[CLS]", "[SEP]"},
			wantOffsets:           []tokenizers.Offset{{0x0, 0x0}, {0x0, 0x0}},
		},
		{
			name:       "invalid utf8 string",
			str:        "\x91D",
			addSpecial: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoding := tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnAllAttributes())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, tt.wantTypeIDs, encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, tt.wantTokens, encoding.Tokens, "wrong tokens")
			assert.Equal(t, tt.wantSpecialTokensMask, encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, tt.wantAttentionMask, encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, tt.wantOffsets, encoding.Offsets, "wrong offsets mask")

			ids, tokens := tk.Encode(tt.str, tt.addSpecial)
			assert.Equal(t, tt.wantIDs, ids, "wrong ids")
			assert.Equal(t, tt.wantTokens, tokens, "wrong tokens")
		})
	}
}

func TestEncodeSpecialTokens(t *testing.T) {
	tk, err := tokenizers.FromBytes(embeddedBytes)
	require.NoError(t, err)
	// special tokens are not encoded by default,
	// meaning if input matches a special token, encoding will include the special token
	ids, _ := tk.Encode("[CLS]fox[SEP]", false)
	assert.Equal(t, []uint32{101, 193284, 102}, ids)
	tk.Close()

	tk, err = tokenizers.FromBytes(embeddedBytes, tokenizers.WithEncodeSpecialTokens())
	require.NoError(t, err)
	ids, _ = tk.Encode("[CLS]fox[SEP]", false)
	// assert that special tokens 101 and 102 are not present
	assert.Equal(t, []uint32{164, 304910, 166, 193284, 164, 211703, 166}, ids)
	tk.Close()
}

func TestEncodeOptions(t *testing.T) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(t, err)
	defer tk.Close()
	tests := []struct {
		name                  string
		str                   string
		addSpecial            bool
		wantIDs               []uint32
		wantTypeIDs           []uint32
		wantTokens            []string
		wantSpecialTokensMask []uint32
		wantAttentionMask     []uint32
		wantOffsets           []tokenizers.Offset
	}{
		{
			name:                  "without special tokens",
			str:                   "brown fox jumps over the lazy dog",
			addSpecial:            false,
			wantIDs:               []uint32{2829, 4419, 14523, 2058, 1996, 13971, 3899},
			wantTypeIDs:           []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantTokens:            []string{"brown", "fox", "jumps", "over", "the", "lazy", "dog"},
			wantSpecialTokensMask: []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantAttentionMask:     []uint32{0x1, 0x1, 0x1, 0x1, 0x1, 0x1, 0x1},
			wantOffsets:           []tokenizers.Offset{{0x0, 0x5}, {0x6, 0x9}, {0xa, 0xf}, {0x10, 0x14}, {0x15, 0x18}, {0x19, 0x1d}, {0x1e, 0x21}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoding := tk.EncodeWithOptions(tt.str, tt.addSpecial)
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, []uint32(nil), encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, []string(nil), encoding.Tokens, "wrong tokens")
			assert.Equal(t, []uint32(nil), encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, []uint32(nil), encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, []tokenizers.Offset(nil), encoding.Offsets, "wrong offsets")

			encoding = tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnTokens())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, []uint32(nil), encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, tt.wantTokens, encoding.Tokens, "wrong tokens")
			assert.Equal(t, []uint32(nil), encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, []uint32(nil), encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, []tokenizers.Offset(nil), encoding.Offsets, "wrong offsets")

			encoding = tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnTypeIDs())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, tt.wantTypeIDs, encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, []string(nil), encoding.Tokens, "wrong tokens")
			assert.Equal(t, []uint32(nil), encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, []uint32(nil), encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, []tokenizers.Offset(nil), encoding.Offsets, "wrong offsets")

			encoding = tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnSpecialTokensMask())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, []uint32(nil), encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, []string(nil), encoding.Tokens, "wrong tokens")
			assert.Equal(t, tt.wantSpecialTokensMask, encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, []uint32(nil), encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, []tokenizers.Offset(nil), encoding.Offsets, "wrong offsets")

			encoding = tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnAttentionMask())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, []uint32(nil), encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, []string(nil), encoding.Tokens, "wrong tokens")
			assert.Equal(t, []uint32(nil), encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, tt.wantAttentionMask, encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, []tokenizers.Offset(nil), encoding.Offsets, "wrong offsets")

			encoding = tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnOffsets())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, []uint32(nil), encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, []string(nil), encoding.Tokens, "wrong tokens")
			assert.Equal(t, []uint32(nil), encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, []uint32(nil), encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, tt.wantOffsets, encoding.Offsets, "wrong offsets")
		})
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

			ids, tokens := tk.Encode(tt.str, tt.addSpecial)
			assert.Equal(t, tt.wantIDs, ids, "wrong ids")
			assert.Equal(t, tt.wantTokens, tokens, "wrong tokens")
		})
	}
}

func TestEncodeWithPadding(t *testing.T) {
	tk, err := tokenizers.FromFile("./test/data/all-minilm-l6-v2.json")
	require.NoError(t, err)
	defer tk.Close()

	tests := []struct {
		name                  string
		str                   string
		addSpecial            bool
		wantIDs               []uint32
		wantTypeIDs           []uint32
		wantTokens            []string
		wantSpecialTokensMask []uint32
		wantAttentionMask     []uint32
		wantOffsets           []tokenizers.Offset
	}{
		{
			name:                  "sentence with padding",
			str:                   "this short sentence",
			addSpecial:            false,
			wantIDs:               []uint32{0x7e7, 0x99c, 0x186b, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantTypeIDs:           []uint32{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantTokens:            []string{"this", "short", "sentence", "[PAD]", "[PAD]", "[PAD]", "[PAD]", "[PAD]"},
			wantSpecialTokensMask: []uint32{0x0, 0x0, 0x0, 0x1, 0x1, 0x1, 0x1, 0x1},
			wantAttentionMask:     []uint32{0x1, 0x1, 0x1, 0x0, 0x0, 0x0, 0x0, 0x0},
			wantOffsets:           []tokenizers.Offset{{0x0, 0x4}, {0x5, 0xa}, {0xb, 0x13}, {0x0, 0x0}, {0x0, 0x0}, {0x0, 0x0}, {0x0, 0x0}, {0x0, 0x0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoding := tk.EncodeWithOptions(tt.str, tt.addSpecial, tokenizers.WithReturnAllAttributes())
			assert.Equal(t, tt.wantIDs, encoding.IDs, "wrong ids")
			assert.Equal(t, tt.wantTypeIDs, encoding.TypeIDs, "wrong type ids")
			assert.Equal(t, tt.wantTokens, encoding.Tokens, "wrong tokens")
			assert.Equal(t, tt.wantSpecialTokensMask, encoding.SpecialTokensMask, "wrong special tokens mask")
			assert.Equal(t, tt.wantAttentionMask, encoding.AttentionMask, "wrong attention mask")
			assert.Equal(t, tt.wantOffsets, encoding.Offsets, "wrong offsets")

			ids, tokens := tk.Encode(tt.str, tt.addSpecial)
			assert.Equal(t, tt.wantIDs, ids, "wrong ids")
			assert.Equal(t, tt.wantTokens, tokens, "wrong tokens")
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

func TestDecodeInvalidString(t *testing.T) {
	tk, err := tokenizers.FromFile("test/data/cohere-tokenizer.json")
	require.NoError(t, err)
	defer tk.Close()

	str := tk.Decode([]uint32{196}, true)
	assert.Empty(t, str)
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
		ids, _ := tk.Encode("brown fox jumps over the lazy dog", false)
		assert.Equal(b, expected, ids)
	}
}

func BenchmarkEncodeNChars(b *testing.B) {
	tk, err := tokenizers.FromFile("./test/data/bert-base-uncased.json")
	require.NoError(b, err)
	defer tk.Close()
	vocabSize := tk.VocabSize()
	input := make([]rune, 0, b.N)
	for i := 0; i < b.N; i++ {
		input = append(input, rune(rand.Uint32()%vocabSize))
	}
	str := string(input)
	b.ResetTimer()
	_, tokens := tk.Encode(str, false)
	assert.Greater(b, len(tokens), 0)
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
	vocabSize := tk.VocabSize()
	input := make([]uint32, 0, b.N)
	for i := 0; i < b.N; i++ {
		input = append(input, rand.Uint32()%vocabSize)
	}
	b.ResetTimer()
	text := tk.Decode(input, true)
	// a token is one or more characters
	assert.Greater(b, len(text), b.N)
}

func TestFromPretrained(t *testing.T) {
	tests := []struct {
		name          string
		modelID       string
		setupOpts     func(t *testing.T) ([]tokenizers.TokenizerConfigOption, string)
		expectedError bool
		expectedToken bool
		skipCache     bool
	}{
		{
			name:          "valid public model with cache dir",
			modelID:       "bert-base-uncased",
			expectedToken: true,
			setupOpts: func(t *testing.T) ([]tokenizers.TokenizerConfigOption, string) {
				tmpDir := t.TempDir()
				return []tokenizers.TokenizerConfigOption{
					tokenizers.WithCacheDir(tmpDir),
				}, tmpDir
			},
		},
		{
			name:          "valid public model without cache dir",
			modelID:       "bert-base-uncased",
			expectedToken: true,
			setupOpts: func(t *testing.T) ([]tokenizers.TokenizerConfigOption, string) {
				return nil, ""
			},
			skipCache: true,
		},
		{
			name:          "private model with valid auth token",
			modelID:       "bert-base-uncased",
			expectedToken: true,
			setupOpts: func(t *testing.T) ([]tokenizers.TokenizerConfigOption, string) {
				tmpDir := t.TempDir()
				return []tokenizers.TokenizerConfigOption{
					tokenizers.WithCacheDir(tmpDir),
					tokenizers.WithAuthToken("test-token"),
				}, tmpDir
			},
		},
		{
			name:          "private model with invalid auth token",
			modelID:       "private-model",
			expectedError: true,
			setupOpts: func(t *testing.T) ([]tokenizers.TokenizerConfigOption, string) {
				tmpDir := t.TempDir()
				return []tokenizers.TokenizerConfigOption{
					tokenizers.WithCacheDir(tmpDir),
					tokenizers.WithAuthToken("invalid-token"),
				}, tmpDir
			},
			skipCache: true,
		},
		{
			name:          "empty model ID",
			modelID:       "",
			expectedError: true,
			setupOpts: func(t *testing.T) ([]tokenizers.TokenizerConfigOption, string) {
				return nil, ""
			},
			skipCache: true,
		},
		{
			name:          "nonexistent model",
			modelID:       "nonexistent/model",
			expectedError: true,
			setupOpts: func(t *testing.T) ([]tokenizers.TokenizerConfigOption, string) {
				tmpDir := t.TempDir()
				return []tokenizers.TokenizerConfigOption{
					tokenizers.WithCacheDir(tmpDir),
				}, tmpDir
			},
			skipCache: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts, cacheDir := tt.setupOpts(t)
			tokenizer, err := tokenizers.FromPretrained(tt.modelID, opts...)

			if tt.expectedError && err == nil {
				t.Fatalf("expected error for case %s, got nil", tt.name)
			}
			if !tt.expectedError && err != nil {
				t.Fatalf("unexpected error for case %s: %v", tt.name, err)
			}
			if !tt.expectedToken && tokenizer != nil {
				t.Fatalf("expected nil tokenizer for case %s, got non-nil", tt.name)
			}
			if tt.expectedToken && tokenizer == nil {
				t.Fatalf("expected non-nil tokenizer for case %s", tt.name)
			}
			if tt.expectedError {
				return
			}
			if !tt.skipCache && cacheDir != "" {
				validateCache(t, cacheDir, tt.modelID)
			}
			if err := tokenizer.Close(); err != nil {
				t.Fatalf("error closing tokenizer: %v", err)
			}
		})
	}
}

func validateCache(t *testing.T, dir string, modelID string) {
	t.Helper()
	files := []string{"tokenizer.json", "vocab.txt"}
	for _, file := range files {
		path := filepath.Join(dir, modelID, file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Fatalf("expected file %s to exist in cache for model %s", file, modelID)
		}
	}
}
