package addresscode

// Prefix extracts the first six digits from a code or an ID number with
// valid syntax. It returns an empty string for invalid input.
func Prefix(input string) string {
	digits := len(input)
	switch digits {
	case 6, 15:
	case 18:
		if last := input[17]; last == 'X' || last == 'x' {
			digits--
		}
	default:
		return ""
	}

	for i := 0; i < digits; i++ {
		if input[i] < '0' || input[i] > '9' {
			return ""
		}
	}
	return input[:6]
}
