package sources

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// DotEnvSource loads configuration from a .env file (KEY=VALUE lines).
type DotEnvSource struct {
    path string
}

// NewDotEnvSource creates a new DotEnvSource for the given file path.
func NewDotEnvSource(path string) *DotEnvSource {
    return &DotEnvSource{path: path}
}

// Load reads the .env file and returns configuration as a nested map.
// Keys are normalized like EnvSource: underscores become dots and keys are lowercased.
func (d *DotEnvSource) Load() (map[string]any, error) {
    f, err := os.Open(d.path)
    if err != nil {
        return nil, fmt.Errorf("opening .env file %s: %w", d.path, err)
    }
    defer f.Close()

    out := make(map[string]any)
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        // Support optional leading `export`
        if strings.HasPrefix(line, "export ") {
            line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
        }
        // Strip inline comments that start with an unescaped '#'
        // Find first '#' not preceded by '\\'
        if idx := indexOfUnescapedHash(line); idx >= 0 {
            line = strings.TrimSpace(line[:idx])
        }
        if line == "" {
            continue
        }
        key, val, ok := strings.Cut(line, "=")
        if !ok {
            // malformed line; skip
            continue
        }
        key = strings.TrimSpace(key)
        val = strings.TrimSpace(val)

        // Unquote if quoted
        if len(val) >= 2 {
            if (val[0] == '\'' && val[len(val)-1] == '\'') || (val[0] == '"' && val[len(val)-1] == '"') {
                val = val[1 : len(val)-1]
            }
        }
        // Unescape common sequences for double/single quoted values
        val = unescapeValue(val)

        dotKey := strings.ToLower(strings.ReplaceAll(key, "_", "."))
        setNestedValue(out, strings.Split(dotKey, "."), val)
    }
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("reading .env file %s: %w", d.path, err)
    }
    return out, nil
}

func indexOfUnescapedHash(s string) int {
    for i := 0; i < len(s); i++ {
        if s[i] == '#' {
            if i == 0 || s[i-1] != '\\' {
                return i
            }
        }
    }
    return -1
}

func unescapeValue(s string) string {
    // Handle common \n, \t, \\ and unescape \#
    replacer := strings.NewReplacer(
        "\\n", "\n",
        "\\r", "\r",
        "\\t", "\t",
        "\\#", "#",
        "\\\\", "\\",
    )
    return replacer.Replace(s)
}

