package main

var alphabetic = map[rune]bool{
	'a': true, 'A': true,
	'b': true, 'B': true,
	'c': true, 'C': true,
	'd': true, 'D': true,
	'e': true, 'E': true,
	'f': true, 'F': true,
	'g': true, 'G': true,
	'h': true, 'H': true,
	'i': true, 'I': true,
	'j': true, 'J': true,
	'k': true, 'K': true,
	'l': true, 'L': true,
	'm': true, 'M': true,
	'n': true, 'N': true,
	'o': true, 'O': true,
	'p': true, 'P': true,
	'q': true, 'Q': true,
	'r': true, 'R': true,
	's': true, 'S': true,
	't': true, 'T': true,
	'u': true, 'U': true,
	'v': true, 'V': true,
	'w': true, 'W': true,
	'x': true, 'X': true,
	'y': true, 'Y': true,
	'z': true, 'Z': true,
}

func validWord(candidate string) bool {
	if len(candidate) < 6 || candidate[0] == '#' || candidate[0] == ' ' {
		return false
	}

	for _, c := range candidate {
		if !alphabetic[c] {
			return false
		}
	}

	return true
}
