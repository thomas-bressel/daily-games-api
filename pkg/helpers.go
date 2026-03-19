package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"unicode/utf8"
)

// WriteJSON writes a JSON-encoded ApiResponse to the ResponseWriter with the given HTTP status code.
func WriteJSON(w http.ResponseWriter, status int, payload ApiResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// WriteSuccess writes a 200 OK JSON response with the given data payload.
func WriteSuccess(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, ApiResponse{
		Status: "ok",
		Data:   data,
	})
}

// WriteError writes a JSON error response with the given HTTP status code and message.
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, ApiResponse{
		Status: "error",
		Error:  message,
	})
}

// ParseJSON decodes the request body into dst.
// It enforces a 1MB body size limit, rejects unknown fields,
// and ensures the body contains exactly one JSON object.
func ParseJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB max

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		if errors.As(err, &syntaxErr) {
			return fmt.Errorf("JSON invalide (erreur de syntaxe)")
		}
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("corps JSON manquant")
		}
		return fmt.Errorf("JSON invalide: %w", err)
	}

	// Ensure no trailing data after the JSON object
	if dec.Decode(&struct{}{}) != io.EOF {
		return fmt.Errorf("JSON contient plusieurs objets ou des données en trop")
	}

	return nil
}

// CleanDescription strips HTML tags, decodes common HTML entities,
// and truncates the result to maxLength characters.
func CleanDescription(raw string, maxLength int) string {
	// Remove HTML tags
	re := regexp.MustCompile(`<[^>]+>`)
	clean := re.ReplaceAllString(raw, "")

	// Decode common HTML entities
	replacer := strings.NewReplacer(
		"&amp;", "&",
		"&lt;", "<",
		"&gt;", ">",
		"&quot;", `"`,
		"&#39;", "'",
		"&nbsp;", " ",
	)
	clean = replacer.Replace(clean)
	clean = strings.TrimSpace(clean)

	// Truncate to maxLength runes (not bytes)
	if utf8.RuneCountInString(clean) > maxLength {
		runes := []rune(clean)
		clean = string(runes[:maxLength]) + "..."
	}

	return clean
}

// ExtractTags scans a title and description for known retro gaming keywords
// and returns a deduplicated list of matching tags.
func ExtractTags(title, description string) []string {
	keywords := []string{
		"amstrad", "atari", "amiga", "commodore", "spectrum",
		"demo", "homebrew", "emulator", "hardware", "review",
		"game", "indie", "retro", "coding", "news", "sales",
	}

	combined := strings.ToLower(title + " " + description)
	seen := make(map[string]bool)
	tags := make([]string, 0) // never nil — serialises as [] not null in JSON

	for _, kw := range keywords {
		if strings.Contains(combined, kw) && !seen[kw] {
			seen[kw] = true
			tags = append(tags, kw)
		}
	}

	return tags
}
