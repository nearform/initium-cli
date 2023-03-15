package helper



func ReplaceIfEmpty(value *string, replacement string) {
    if (*value == "") {
        *value = replacement
    }
}