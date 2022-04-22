package utils

const (
	flagPrefix    = "-"
	flagSeparator = "="
)

//Chunks splits the xs slice into a slice of slices sized chunkSize
func Chunks(xs []int, chunkSize int) [][]int {
	if len(xs) == 0 {
		return nil
	}
	divided := make([][]int, (len(xs)+chunkSize-1)/chunkSize)
	prev := 0
	i := 0
	till := len(xs) - chunkSize
	for prev < till {
		next := prev + chunkSize
		divided[i] = xs[prev:next]
		prev = next
		i++
	}
	divided[i] = xs[prev:]
	return divided
}

//TrimPrefix Removes all the occurences of the character c on the beging of the string s
func TrimPrefix(s string, c string) string {
	tmp := s
	for len(c) > 0 && len(tmp) > 0 && tmp[0:len(c)] == c {
		tmp = tmp[len(c):]
	}
	return tmp
}
