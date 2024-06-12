package redact

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type Pattern struct {
	Regex *regexp.Regexp
	num   int
}

func RedactFile(file string) string {
	file = CSWF(file)
	file = Punctuation(file)
	file = Atoan(file)
	file = Atoan(file)
	file = IsFunc(file)
	file = Antoa(file)
	file = Atoan(file)
	file = Punctuation(file)
	return file
}

func Punctuation(file string) string {
	re := regexp.MustCompile(`([\w]+\b|['"])[\s]+([,.:;!?]+)`)
	file = re.ReplaceAllString(file, "$1$2")
	re = regexp.MustCompile(`([,.:;!?'"]+)([\w]+)`)
	file = re.ReplaceAllString(file, "$1 $2")
	re = regexp.MustCompile(`([,.:;!?])\s+([,.:;!?])`)
	file = re.ReplaceAllString(file, "$1$2")
	re = regexp.MustCompile(`(don|can|wan|couldn|woudn|mastn)\s*(')\s*(t)`)
	file = re.ReplaceAllString(file, "$1$2$3")
	re = regexp.MustCompile(`(i|you|she|he|)\s*(')\s*(ll)`)
	file = re.ReplaceAllString(file, "$1$2$3")
	re = regexp.MustCompile(`(i)\s*(')\s*(m)`)
	file = re.ReplaceAllString(file, "$1$2$3")
	re = regexp.MustCompile(`(')(\s*|\s+)([^']*)([.,?!:;]|\b|)(\s+|\s*)(')`)
	file = re.ReplaceAllString(file, "$1$3$4$6")
	re = regexp.MustCompile(`(")(\s*|\s+)([^"]*)([.,?!:;]|\b|)(\s+|\s*)(")`)
	file = re.ReplaceAllString(file, "$1$3$4$6")
	re = regexp.MustCompile(`(')(\s*|\s+)([^']*)([.,?!:;]|\b)(\s+|\s*)(')`)
	file = re.ReplaceAllString(file, "$1$3$4$6")
	re = regexp.MustCompile(`(")(\s*|\s+)([^"]*)([.,?!:;']|\b)(\s+|\s*)(")`)
	file = re.ReplaceAllString(file, "$1$3$4$6")
	re = regexp.MustCompile(`([\w]|[,.!?:;])(')([^']+)('[\s])`)
	file = re.ReplaceAllString(file, "$1 $2$3$4")
	re = regexp.MustCompile(`([\w]|[,.!?:;])(")([^"]+)("[\s])`)
	file = re.ReplaceAllString(file, "$1 $2$3$4")
	re = regexp.MustCompile(`([\w]|[,.!?:;])(')([^']+)('[\s])`)
	file = re.ReplaceAllString(file, "$1$3$5")
	re = regexp.MustCompile(`([\w]|[,.!?:;])(")([^"]+)("[\s])`)
	file = re.ReplaceAllString(file, "$1$3$5")
	re = regexp.MustCompile(`(")(\s*|\s+)([^"]*)\b(\s+|\s*)(")`)
	file = re.ReplaceAllString(file, "$1$3$5")
	re = regexp.MustCompile(`(')(\s*|\s+)([^']*)\b(\s+|\s*)(')`)
	file = re.ReplaceAllString(file, "$1$3$5")

	return file
}

func Atoan(file string) string {
	pattern := regexp.MustCompile(`A([\s]+[aeoiuhAEOIUH])`)
	file = pattern.ReplaceAllString(file, "An$1")
	pattern = regexp.MustCompile(`a([\s]+[aeoiuhAEOIUH])`)
	file = pattern.ReplaceAllString(file, "an$1")
	return file
}

func Antoa(file string) string {
	pattern := regexp.MustCompile(`(A|a)(n|N)([\s]+\b[^a^e^o^i^u^h^A^E^O^I^U^H])`)
	file = pattern.ReplaceAllString(file, "$1$3")
	return file
}

func IsFunc(file string) string {
	cap := regexp.MustCompile(`\(\s*cap\s*\)`)
	if cap.MatchString(file) {
		file = cap.ReplaceAllString(file, "(cap, 1)")
	}
	cap = regexp.MustCompile(`\(\s*cap\s*,\s*([-*\d]+)\s*\)`)
	low := regexp.MustCompile(`\(\s*low\s*\)`)
	if low.MatchString(file) {
		file = low.ReplaceAllString(file, "(low, 1)")
	}
	low = regexp.MustCompile(`\(\s*low\s*,\s*([-*\d]+)\s*\)`)
	up := regexp.MustCompile(`\(\s*up\s*\)`)
	if up.MatchString(file) {
		file = up.ReplaceAllString(file, "(up, 1)")
	}
	up = regexp.MustCompile(`\(\s*up\s*,\s*([-*\d]+)\s*\)`)
	bin := regexp.MustCompile(`(\b[\w]+\b\s*[,.;:!?]*\s*)(\(\s*bin\s*\))`)
	hex := regexp.MustCompile(`(\b[\w]+\b\s*[,.;:!?]*\s*)(\(\s*hex\s*\))`)
	funcArr := []func(string) string{Cap, Low, Up, Bin, Hex}
	patterns := []Pattern{
		{Regex: cap, num: 0},
		{Regex: low, num: 1},
		{Regex: up, num: 2},
		{Regex: bin, num: 3},
		{Regex: hex, num: 4},
	}
	re := regexp.MustCompile(`\(\s*cap\s*,\s*([-*\d]+)\s*\)|\(\s*low\s*,\s*([-*\d]+)\s*\)|\(\s*up\s*,\s*([-*\d]+)\s*\)|(\b[\w]+\b\s*[,.;:!?]*\s*)(\(\s*bin\s*\))|(\b[\w]+\b\s*[,.;:!?]*\s*)(\(\s*hex\s*\))`)
	index := re.FindAllStringIndex(file, -1)
	if len(index) > 0 {
		start := file[:index[0][1]]
		end := file[index[0][1]:]
		if cap.MatchString(start) || low.MatchString(start) || up.MatchString(start) || bin.MatchString(start) || hex.MatchString(start) {
			for _, pattern := range patterns {
				if pattern.Regex.MatchString(start) {
					start = funcArr[pattern.num](start)
					file = start + end
				}
			}
		}
		return (IsFunc(file))
	} else {
		return file
	}
}

func Hex(input string) string {
	re := regexp.MustCompile(`(\b[\w]+\b)(\s*[,.;:!?]*\s*)(\(\s*hex\s*\))`)
	matches := re.FindAllStringSubmatch(input, 1)
	index := re.FindAllStringIndex(input, -1)
	start := input[:index[0][0]]
	end := input[index[0][0]:]
	end = re.ReplaceAllString(end, "$1$2")
	hex, err := strconv.ParseInt(matches[0][1], 16, 64)
	if err != nil {
		log.Fatal(err, " функция HEX")
	}
	heximalstr := strconv.Itoa(int(hex))
	end = strings.Replace(end, matches[0][1], heximalstr, 1)
	input = start + end
	return input
}

func Bin(input string) string {
	re := regexp.MustCompile(`(\b[\w]+\b)(\s*[,.;:!?]*\s*)(\(\s*bin\s*\))`)
	matches := re.FindAllStringSubmatch(input, 1)
	index := re.FindAllStringIndex(input, -1)
	start := input[:index[0][0]]
	end := input[index[0][0]:]
	end = re.ReplaceAllString(end, "$1$2")
	binary := matches[0][1]
	bin, err := strconv.ParseInt(binary, 2, 64)
	if err != nil {
		log.Fatal(err, " функция BIN")
	}
	decimalStr := strconv.FormatInt(bin, 10)
	end = strings.Replace(end, matches[0][1], decimalStr, 1)
	input = start + end
	return input
}

func Cap(input string) string {
	re := regexp.MustCompile(`\(\s*cap\s*,\s*([-*\d]+)\s*\)`)
	matches1 := re.FindAllStringSubmatch(input, 1)
	n := 0
	for _, sub := range matches1 {
		n, _ = strconv.Atoi(sub[1])
	}
	if n < 0 {
		log.Fatal("cap, введите положительное число")
	}
	input = re.ReplaceAllString(input, "")
	re = regexp.MustCompile(`\b[a-zA-Z]+\b`)
	matches := re.FindAllStringSubmatch(input, -1)
	index := re.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return input
	}
	if n > 0 && len(matches) >= n {
		matches = matches[len(matches)-n:]
		index = index[len(index)-n:]
		for _, sub := range matches {
			sub[0] = strings.Title(strings.ToLower(sub[0]))
		}
		input = Replace(input, matches, index, n)
	} else if n > len(matches) {
		fmt.Println(input, matches1[0][0])
		log.Fatalln("Не хватает слов для (cap)")
	}
	return input
}

func Low(input string) string {
	re := regexp.MustCompile(`\(\s*low\s*,\s*([-*\d]+)\s*\)`)
	matches1 := re.FindAllStringSubmatch(input, 1)
	n := 0
	for _, sub := range matches1 {
		n, _ = strconv.Atoi(sub[1])
	}
	if n < 0 {
		log.Fatal("low, введите положительное число")
	}
	input = re.ReplaceAllString(input, "")
	re = regexp.MustCompile(`\b[a-zA-Z]+\b`)
	matches := re.FindAllStringSubmatch(input, -1)
	index := re.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return input
	}
	if n > 0 && len(matches) >= n {
		matches = matches[len(matches)-n:]
		index = index[len(index)-n:]
		for _, sub := range matches {
			sub[0] = strings.ToLower(sub[0])
		}
		input = Replace(input, matches, index, n)
	} else if n > len(matches) {
		fmt.Println(input, matches1[0][0])
		log.Fatalln("Не хватает слов для (low)")
	}
	return input
}

func Up(input string) string {
	re := regexp.MustCompile(`\(\s*up\s*,\s*([-*\d]+)\s*\)`)
	matches1 := re.FindAllStringSubmatch(input, 1)
	n := 0
	for _, sub := range matches1 {
		n, _ = strconv.Atoi(sub[1])
	}
	if n < 0 {
		log.Fatal("low, введите положительное число")
	}
	input = re.ReplaceAllString(input, "")
	re = regexp.MustCompile(`\b[a-zA-Z]+\b`)
	matches := re.FindAllStringSubmatch(input, -1)
	index := re.FindAllStringIndex(input, -1)
	if len(matches) == 0 {
		return input
	}
	if n > 0 && len(matches) >= n {
		matches = matches[len(matches)-n:]
		index = index[len(index)-n:]
		for _, sub := range matches {
			sub[0] = strings.ToTitle(sub[0])
		}
		input = Replace(input, matches, index, n)
	} else if n > len(matches) {
		fmt.Println(input, matches1[0][0])
		log.Fatalln("Не хватает слов для (Up)")
	}
	return input
}

func Replace(str string, matches [][]string, index [][]int, n int) string {
	words := []string{}
	for _, array := range matches {
		for _, arr := range array {
			words = append(words, arr)
		}
	}
	for i := 0; i < len(index); i++ {
		str = str[:index[i][0]] + words[i] + str[index[i][1]:]
	}
	return str
}

func CSWF(file string) string {
	capnum := regexp.MustCompile(`(\b)\s(\(\s*cap\s*,\s*\d\s*\))(\s)`)
	file = capnum.ReplaceAllString(file, "$1$2$3")
	cap := regexp.MustCompile(`(\b)\s(\(\s*cap\s*\))(\s)`)
	file = cap.ReplaceAllString(file, "$1$2$3")
	lownum := regexp.MustCompile(`(\b)\s(\(\s*low\s*,\s*\d\s*\))(\s)`)
	file = lownum.ReplaceAllString(file, "$1$2$3")
	low := regexp.MustCompile(`(\b)\s(\(\s*low\s*\))(\s)`)
	file = low.ReplaceAllString(file, "$1$2$3")
	upnum := regexp.MustCompile(`(\b)\s(\(\s*up\s*,\s*\d\s*\))(\s)`)
	file = upnum.ReplaceAllString(file, "$1$2$3")
	up := regexp.MustCompile(`(\b)\s(\(\s*up\s*\))(\s)`)
	file = up.ReplaceAllString(file, "$1$2$3")
	hex := regexp.MustCompile(`(\b)\s(\(\s*hex\s*\))(\s)`)
	file = hex.ReplaceAllString(file, "$1$2$3")
	bin := regexp.MustCompile(`(\b)\s(\(\s*bin\s*\))(\s)`)
	file = bin.ReplaceAllString(file, "$1$2$3")
	return file
}
