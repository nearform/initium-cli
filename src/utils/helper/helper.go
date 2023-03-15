package helper

// Replace with replacement string if
func ReplaceIfEmpty(value *string, replacement string) {
    if (*value == "") {
        *value = replacement
    }
}