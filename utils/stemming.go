package utils

import (
	"strings"
)

func Stem(word string) string {
	if len(word) > 2 {
		return step1a(strings.ToLower(word))
	}
	return strings.ToLower(word)
}

func isConsonant(word string, pos int) bool {
	switch word[pos] {
	case 'a', 'e', 'i', 'o', 'u':
		return false
	case 'y':
		if pos == 0 {
			return true
		}
		return !isConsonant(word, pos-1)
	}
	return true
}

func getMeasure(word string) int {
	form := ""
	for i := 0; i < len(word); i++ {
		if isConsonant(word, i) {
			form += "C"
		} else {
			form += "V"
		}
	}
	return strings.Count(form, "VC")
}

func containsVowel(word string) bool {
	for i := 0; i < len(word); i++ {
		if !isConsonant(word, i) {
			return true
		}
	}
	return false
}

func endsWithDoubleConsonant(word string) bool {
	return len(word) > 2 && word[len(word)-1] == word[len(word)-2] && isConsonant(word, len(word)-1)
}

func endsWithCVC(word string) bool {
	return len(word) >= 3 && isConsonant(word, len(word)-1) && !isConsonant(word, len(word)-2) &&
		isConsonant(word, len(word)-3) && word[len(word)-1] != 'x' && word[len(word)-1] != 'y' &&
		word[len(word)-1] != 'w'
}

func step1a(word string) string {
	if strings.HasSuffix(word, "sses") {
		word = strings.TrimSuffix(word, "es")
	} else if strings.HasSuffix(word, "ies") {
		word = strings.TrimSuffix(word, "es")
	} else if strings.HasSuffix(word, "s") {
		word = strings.TrimSuffix(word, "s")
	}
	return step1b(word)
}

func step1bExtra(word string) string {
	if strings.HasSuffix(word, "at") || strings.HasSuffix(word, "bl") ||
		strings.HasSuffix(word, "iz") {
		return word + "e"
	}
	if endsWithDoubleConsonant(word) && word[len(word)-1] != 'l' && word[len(word)-1] != 's' &&
		word[len(word)-1] != 'z' {
		return word[:len(word)-1]
	}
	if getMeasure(word) == 1 && endsWithCVC(word) {
		return word + "e"
	}
	return word
}

func step1b(word string) string {
	if strings.HasSuffix(word, "eed") && getMeasure(word[:len(word)-3]) > 0 {
		word = strings.TrimSuffix(word, "d")
	} else if strings.HasSuffix(word, "ed") && containsVowel(word[:len(word)-2]) {
		word = step1bExtra(word[:len(word)-2])
	} else if strings.HasSuffix(word, "ing") && containsVowel(word[:len(word)-3]){
		word = step1bExtra(word[:len(word)-3])
	}
	return step1c(word)
}

func step1c(word string) string {
	if strings.HasSuffix(word, "y") && containsVowel(word[:len(word)-1]){
		word = strings.TrimSuffix(word, "y") + "i"
	}
	return step2(word)
}

func step2(word string) string {
	if strings.HasSuffix(word, "ational") && getMeasure(word[:len(word)-7]) > 0 {
		word = strings.TrimSuffix(word, "ational") + "ate"
	} else if strings.HasSuffix(word, "tional") && getMeasure(word[:len(word)-6]) > 0 {
		word = strings.TrimSuffix(word, "al")
	} else if strings.HasSuffix(word, "enci") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "enci") + "ence"
	} else if strings.HasSuffix(word, "anci") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "anci") + "ance"
	} else if strings.HasSuffix(word, "izer") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "r")
	} else if strings.HasSuffix(word, "abli") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "abli") + "able"
	} else if strings.HasSuffix(word, "alli") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "li")
	} else if strings.HasSuffix(word, "entli") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "li")
	} else if strings.HasSuffix(word, "eli") && getMeasure(word[:len(word)-3]) > 0 {
		word = strings.TrimSuffix(word, "li")
	} else if strings.HasSuffix(word, "ousli") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "li")
	} else if strings.HasSuffix(word, "ization") && getMeasure(word[:len(word)-7]) > 0 {
		word = strings.TrimSuffix(word, "ization") + "ize"
	} else if strings.HasSuffix(word, "ation") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "ation") + "ate"
	} else if strings.HasSuffix(word, "ator") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "ator") + "ate"
	} else if strings.HasSuffix(word, "alism") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "ism")
	} else if strings.HasSuffix(word, "iveness") && getMeasure(word[:len(word)-7]) > 0 {
		word = strings.TrimSuffix(word, "ness")
	} else if strings.HasSuffix(word, "fulness") && getMeasure(word[:len(word)-7]) > 0 {
		word = strings.TrimSuffix(word, "ness")
	} else if strings.HasSuffix(word, "ousness") && getMeasure(word[:len(word)-7]) > 0 {
		word = strings.TrimSuffix(word, "ness")
	} else if strings.HasSuffix(word, "aliti") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "iti")
	} else if strings.HasSuffix(word, "iviti") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "iti") + "e"
	} else if strings.HasSuffix(word, "biliti") && getMeasure(word[:len(word)-6]) > 0 {
		word = strings.TrimSuffix(word, "iti") + "e"
	}
	return step3(word)
}

func step3(word string) string {
	if strings.HasSuffix(word, "icate") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "ate")
	} else if strings.HasSuffix(word, "ative") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "ative")
	} else if strings.HasSuffix(word, "alize") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "ize")
	} else if strings.HasSuffix(word, "iciti") && getMeasure(word[:len(word)-5]) > 0 {
		word = strings.TrimSuffix(word, "iti")
	} else if strings.HasSuffix(word, "ical") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "al")
	} else if strings.HasSuffix(word, "ful") && getMeasure(word[:len(word)-3]) > 0 {
		word = strings.TrimSuffix(word, "ful")
	} else if strings.HasSuffix(word, "ness") && getMeasure(word[:len(word)-4]) > 0 {
		word = strings.TrimSuffix(word, "ness")
	}
	return step4(word)
}

func step4(word string) string {
	if strings.HasSuffix(word, "al") && getMeasure(word[:len(word)-2]) > 1 {
		word = strings.TrimSuffix(word, "al")
	} else if strings.HasSuffix(word, "ance") && getMeasure(word[:len(word)-4]) > 1 {
		word = strings.TrimSuffix(word, "ance")
	} else if strings.HasSuffix(word, "ence") && getMeasure(word[:len(word)-4]) > 1 {
		word = strings.TrimSuffix(word, "ence")
	} else if strings.HasSuffix(word, "er") && getMeasure(word[:len(word)-2]) > 1 {
		word = strings.TrimSuffix(word, "er")
	} else if strings.HasSuffix(word, "ic") && getMeasure(word[:len(word)-2]) > 1 {
		word = strings.TrimSuffix(word, "ic")
	} else if strings.HasSuffix(word, "able") && getMeasure(word[:len(word)-4]) > 1 {
		word = strings.TrimSuffix(word, "able")
	} else if strings.HasSuffix(word, "ible") && getMeasure(word[:len(word)-4]) > 1 {
		word = strings.TrimSuffix(word, "ible")
	} else if strings.HasSuffix(word, "ant") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "ant")
	} else if strings.HasSuffix(word, "ement") && getMeasure(word[:len(word)-5]) > 1 {
		word = strings.TrimSuffix(word, "ement")
	} else if strings.HasSuffix(word, "ment") && getMeasure(word[:len(word)-4]) > 1 {
		word = strings.TrimSuffix(word, "ment")
	} else if strings.HasSuffix(word, "ent") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "ent")
	} else if strings.HasSuffix(word, "ion") && getMeasure(word[:len(word)-3]) > 1 &&
		(word[len(word)-4] == 's' || word[len(word)-4] == 't') {
		word = strings.TrimSuffix(word, "ion")
	} else if strings.HasSuffix(word, "ou") && getMeasure(word[:len(word)-2]) > 1 {
		word = strings.TrimSuffix(word, "ou")
	} else if strings.HasSuffix(word, "ism") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "ism")
	} else if strings.HasSuffix(word, "ate") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "ate")
	} else if strings.HasSuffix(word, "iti") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "iti")
	} else if strings.HasSuffix(word, "ous") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "ous")
	} else if strings.HasSuffix(word, "ive") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "ive")
	} else if strings.HasSuffix(word, "ize") && getMeasure(word[:len(word)-3]) > 1 {
		word = strings.TrimSuffix(word, "ize")
	}
	return step5a(word)
}

func step5a(word string) string {
	if strings.HasSuffix(word, "e") && getMeasure(word[:len(word)-1]) > 1 {
		word = strings.TrimSuffix(word, "e")
	} else if strings.HasSuffix(word, "e") && getMeasure(word[:len(word)-1]) == 1 && !endsWithCVC(word) {
		word = strings.TrimSuffix(word, "e")
	}
	return step5b(word)
}

func step5b(word string) string {
	if getMeasure(word) > 1 && endsWithDoubleConsonant(word) && word[len(word)-1] == 'l' {
		return word[:len(word)-1]
	}
	return word
}
