package ksqldb

import "strings"

// SanitizeQuery sanitizes the given content
// eventually we can use the KSqlParser to rewrite the query, so its automatically sanitized
// whitespaces will be eaten by the KSqlParser
func SanitizeQuery(content string) string {
	content = strings.ReplaceAll(content, "\t", "")
	content = strings.ReplaceAll(content, "\n", "")
	return content
}
