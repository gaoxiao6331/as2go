package preprocess

import "regexp"

var numSuffixRe = regexp.MustCompile(`\b(\d+(?:\.\d+)?)(i8|i16|i32|i64|u8|u16|u32|u64|isize|usize|f32|f64)\b`)

func Process(src []byte) []byte {
    return numSuffixRe.ReplaceAll(src, []byte("($1 as $2)"))
}
