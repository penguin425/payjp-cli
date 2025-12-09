package output

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"gopkg.in/yaml.v3"
)

// Format represents the output format
type Format string

const (
	FormatJSON  Format = "json"
	FormatTable Format = "table"
	FormatYAML  Format = "yaml"
	FormatQuiet Format = "quiet"
)

// Formatter is the interface for output formatters
type Formatter interface {
	Format(data interface{}) error
}

// NewFormatter creates a new formatter based on the format type
func NewFormatter(format Format) Formatter {
	switch format {
	case FormatJSON:
		return &JSONFormatter{}
	case FormatYAML:
		return &YAMLFormatter{}
	case FormatQuiet:
		return &QuietFormatter{}
	default:
		return &TableFormatter{}
	}
}

// Output outputs the data in the specified format
func Output(format string, data interface{}) error {
	f := NewFormatter(Format(format))
	return f.Format(data)
}

// OutputQuiet outputs only the ID
func OutputQuiet(data interface{}) error {
	f := &QuietFormatter{}
	return f.Format(data)
}

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

// Format formats the data as JSON
func (f *JSONFormatter) Format(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// YAMLFormatter formats output as YAML
type YAMLFormatter struct{}

// Format formats the data as YAML
func (f *YAMLFormatter) Format(data interface{}) error {
	encoder := yaml.NewEncoder(os.Stdout)
	encoder.SetIndent(2)
	defer encoder.Close()
	return encoder.Encode(data)
}

// QuietFormatter outputs only the ID
type QuietFormatter struct{}

// Format outputs only the ID field
func (f *QuietFormatter) Format(data interface{}) error {
	id := extractID(data)
	if id != "" {
		fmt.Println(id)
	}
	return nil
}

// extractID extracts the ID from a struct
func extractID(data interface{}) string {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	// Try common ID field names
	for _, fieldName := range []string{"ID", "Id", "id"} {
		field := v.FieldByName(fieldName)
		if field.IsValid() && field.Kind() == reflect.String {
			return field.String()
		}
	}

	return ""
}

// TableFormatter formats output as a table
type TableFormatter struct{}

// Format formats the data as a table
func (f *TableFormatter) Format(data interface{}) error {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Handle slice of items
	if v.Kind() == reflect.Slice {
		return f.formatSlice(v)
	}

	// Handle single item
	return f.formatSingle(v)
}

// formatSlice formats a slice of items as a table
func (f *TableFormatter) formatSlice(v reflect.Value) error {
	if v.Len() == 0 {
		fmt.Println("No items found.")
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	// Get headers from first element
	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}

	headers, keys := getTableHeaders(first)
	table.SetHeader(headers)

	// Add rows
	for i := 0; i < v.Len(); i++ {
		item := v.Index(i)
		if item.Kind() == reflect.Ptr {
			item = item.Elem()
		}
		row := getTableRow(item, keys)
		table.Append(row)
	}

	table.Render()
	fmt.Printf("Total: %d items\n", v.Len())
	return nil
}

// formatSingle formats a single item as a table
func (f *TableFormatter) formatSingle(v reflect.Value) error {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("expected struct, got %v", v.Kind())
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(true)
	table.SetRowLine(false)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeader([]string{"FIELD", "VALUE"})

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		fieldName := getFieldName(field)
		fieldValue := formatFieldValue(value)

		table.Append([]string{fieldName, fieldValue})
	}

	table.Render()
	return nil
}

// getTableHeaders returns headers for a table
func getTableHeaders(v reflect.Value) ([]string, []string) {
	if v.Kind() != reflect.Struct {
		return nil, nil
	}

	t := v.Type()
	headers := []string{}
	keys := []string{}

	// Common fields to display in list view
	commonFields := []string{"ID", "Amount", "Currency", "Status", "Paid", "Captured", "Refunded", "Email", "Description", "Name", "Interval", "CreatedAt", "Created"}

	for _, fieldName := range commonFields {
		field, ok := t.FieldByName(fieldName)
		if ok && field.IsExported() {
			headers = append(headers, strings.ToUpper(getFieldName(field)))
			keys = append(keys, fieldName)
		}
	}

	// If no common fields found, use first few fields
	if len(headers) == 0 {
		for i := 0; i < t.NumField() && i < 6; i++ {
			field := t.Field(i)
			if field.IsExported() {
				headers = append(headers, strings.ToUpper(getFieldName(field)))
				keys = append(keys, field.Name)
			}
		}
	}

	return headers, keys
}

// getTableRow returns a row for a table
func getTableRow(v reflect.Value, keys []string) []string {
	row := []string{}

	for _, key := range keys {
		field := v.FieldByName(key)
		if field.IsValid() {
			row = append(row, formatFieldValue(field))
		} else {
			row = append(row, "")
		}
	}

	return row
}

// getFieldName returns the display name for a field
func getFieldName(field reflect.StructField) string {
	// Try JSON tag first
	if tag := field.Tag.Get("json"); tag != "" {
		parts := strings.Split(tag, ",")
		if parts[0] != "" && parts[0] != "-" {
			return parts[0]
		}
	}

	// Convert CamelCase to snake_case
	return toSnakeCase(field.Name)
}

// formatFieldValue formats a field value for display
func formatFieldValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}

	// Handle pointers
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Bool:
		return fmt.Sprintf("%v", v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		// Check if it's a timestamp (Unix time)
		if v.Int() > 1000000000 && v.Int() < 2000000000 {
			t := time.Unix(v.Int(), 0)
			return t.Format("2006-01-02 15:04:05")
		}
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%.2f", v.Float())
	case reflect.String:
		s := v.String()
		if len(s) > 50 {
			return s[:47] + "..."
		}
		return s
	case reflect.Struct:
		// Handle time.Time
		if t, ok := v.Interface().(time.Time); ok {
			return t.Format("2006-01-02 15:04:05")
		}
		return "{...}"
	case reflect.Map, reflect.Slice:
		if v.Len() == 0 {
			return ""
		}
		return fmt.Sprintf("[%d items]", v.Len())
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}

// toSnakeCase converts a CamelCase string to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// PrintError prints an error message
func PrintError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// PrintSuccess prints a success message
func PrintSuccess(message string) {
	fmt.Println(message)
}

// PrintInfo prints an info message
func PrintInfo(message string) {
	fmt.Println(message)
}
