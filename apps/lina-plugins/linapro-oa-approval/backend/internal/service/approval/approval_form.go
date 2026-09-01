// This file implements the per-flow configurable form field definitions:
// validation, snapshot serialization, and value validation against a frozen
// snapshot. Flows without configured fields fall back to the built-in
// standard form so existing flows keep working unchanged.

package approval

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"

	"lina-core/pkg/bizerr"
)

// Form field type values supported by the dynamic form.
const (
	FieldTypeText       = "text"       // Single-line text
	FieldTypeTextarea   = "textarea"   // Multi-line text
	FieldTypeNumber     = "number"     // Number, optional as-amount flag
	FieldTypeDate       = "date"       // Date, YYYY-MM-DD
	FieldTypeSelect     = "select"     // Select with configured options
	FieldTypeAttachment = "attachment" // Attachment URL address list
	FieldTypeDetail     = "detail"     // Detail sub-table with configured columns
)

// Detail column types supported inside one detail field.
const (
	DetailColumnText   = "text"
	DetailColumnNumber = "number"
	DetailColumnDate   = "date"
)

// MaxFormFields bounds the configured field count of one flow.
const MaxFormFields = 30

// MaxDetailRows bounds one submitted detail sub-table.
const MaxDetailRows = 50

// fieldKeyPattern constrains storage keys to identifier-safe characters.
var fieldKeyPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_]{0,63}$`)

// FieldConfig defines one configurable form field on a flow.
type FieldConfig struct {
	Key      string         `json:"key"`      // Storage key, identifier-safe and unique
	Label    string         `json:"label"`    // Display name
	Type     string         `json:"type"`     // text/textarea/number/date/select/attachment/detail
	Required bool           `json:"required"` // Whether the field must be filled
	AsAmount bool           `json:"asAmount"` // Whether a number field projects to the request amount
	Options  []string       `json:"options"`  // Select options
	Columns  []ColumnConfig `json:"columns"`  // Detail sub-table columns
}

// ColumnConfig defines one column of a detail sub-table.
type ColumnConfig struct {
	Key   string `json:"key"`   // Storage key
	Label string `json:"label"` // Display name
	Type  string `json:"type"`  // text/number/date
}

// defaultFormFields returns the built-in standard form used when a flow has
// no configured fields: amount, statement content, and attachment URLs.
func defaultFormFields() []FieldConfig {
	return []FieldConfig{
		{Key: "amount", Label: "Amount", Type: FieldTypeNumber, AsAmount: true},
		{Key: "content", Label: "Content", Type: FieldTypeTextarea},
		{Key: "attachments", Label: "Attachments", Type: FieldTypeAttachment},
	}
}

// ValidateFormFields validates the caller-submitted field definitions and
// returns the normalized list, or a structured business error.
func ValidateFormFields(fields []FieldConfig) ([]FieldConfig, error) {
	if len(fields) == 0 {
		return nil, nil
	}
	if len(fields) > MaxFormFields {
		return nil, bizerr.NewCode(CodeFormFieldsInvalid)
	}
	seenKeys := make(map[string]struct{}, len(fields))
	normalized := make([]FieldConfig, 0, len(fields))
	for _, field := range fields {
		field.Key = strings.TrimSpace(field.Key)
		field.Label = strings.TrimSpace(field.Label)
		field.Type = strings.TrimSpace(field.Type)
		if !fieldKeyPattern.MatchString(field.Key) {
			return nil, bizerr.NewCode(CodeFormFieldKeyInvalid)
		}
		if _, dup := seenKeys[field.Key]; dup {
			return nil, bizerr.NewCode(CodeFormFieldKeyInvalid)
		}
		seenKeys[field.Key] = struct{}{}
		if field.Label == "" {
			return nil, bizerr.NewCode(CodeFormFieldLabelRequired)
		}
		switch field.Type {
		case FieldTypeText, FieldTypeTextarea, FieldTypeNumber, FieldTypeDate, FieldTypeAttachment:
		case FieldTypeSelect:
			if len(field.Options) == 0 {
				return nil, bizerr.NewCode(CodeFormFieldOptionsRequired)
			}
		case FieldTypeDetail:
			if len(field.Columns) == 0 || len(field.Columns) > MaxFormFields {
				return nil, bizerr.NewCode(CodeFormFieldColumnsRequired)
			}
			columnKeys := make(map[string]struct{}, len(field.Columns))
			for _, column := range field.Columns {
				column.Key = strings.TrimSpace(column.Key)
				column.Label = strings.TrimSpace(column.Label)
				column.Type = strings.TrimSpace(column.Type)
				if !fieldKeyPattern.MatchString(column.Key) {
					return nil, bizerr.NewCode(CodeFormFieldColumnsRequired)
				}
				if _, dup := columnKeys[column.Key]; dup {
					return nil, bizerr.NewCode(CodeFormFieldColumnsRequired)
				}
				columnKeys[column.Key] = struct{}{}
				if column.Label == "" {
					return nil, bizerr.NewCode(CodeFormFieldColumnsRequired)
				}
				switch column.Type {
				case DetailColumnText, DetailColumnNumber, DetailColumnDate:
				default:
					return nil, bizerr.NewCode(CodeFormFieldColumnsRequired)
				}
			}
		default:
			return nil, bizerr.NewCode(CodeFormFieldTypeInvalid)
		}
		normalized = append(normalized, field)
	}
	return normalized, nil
}

// marshalFormSnapshot serializes the frozen field definitions.
func marshalFormSnapshot(fields []FieldConfig) (string, error) {
	if len(fields) == 0 {
		return "", nil
	}
	content, err := json.Marshal(fields)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// parseFormSnapshot deserializes the frozen field definitions and falls back
// to the built-in standard form for empty snapshots so legacy requests stay
// readable.
func parseFormSnapshot(snapshot string) ([]FieldConfig, error) {
	trimmed := strings.TrimSpace(snapshot)
	if trimmed == "" {
		return defaultFormFields(), nil
	}
	fields := make([]FieldConfig, 0)
	if err := json.Unmarshal([]byte(trimmed), &fields); err != nil {
		return nil, errors.New("parse approval form snapshot failed: " + err.Error())
	}
	return fields, nil
}

// validateFormValues checks the submitted values against the frozen field
// snapshot and returns the canonical stored JSON text. Detail tables and
// attachments are stored as JSON; everything else keeps its scalar value.
func validateFormValues(fields []FieldConfig, values map[string]any) (string, error) {
	stored := make(map[string]any, len(values))
	for _, field := range fields {
		value, present := values[field.Key]
		switch field.Type {
		case FieldTypeDetail:
			rows, err := normalizeDetailRows(field, value, present)
			if err != nil {
				return "", err
			}
			if len(rows) > 0 {
				stored[field.Key] = rows
			}
		case FieldTypeAttachment:
			urls, err := normalizeFormFieldAttachments(field, value, present)
			if err != nil {
				return "", err
			}
			if len(urls) > 0 {
				stored[field.Key] = urls
			}
		default:
			normalized, hasValue, err := normalizeScalarField(field, value, present)
			if err != nil {
				return "", err
			}
			if hasValue {
				stored[field.Key] = normalized
			}
		}
	}
	if len(stored) == 0 {
		return "", nil
	}
	content, err := json.Marshal(stored)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// normalizeScalarField validates one scalar field value against its type and
// required flag. It returns the normalized value and whether a value exists.
func normalizeScalarField(field FieldConfig, value any, present bool) (any, bool, error) {
	if !present || value == nil {
		if field.Required {
			return nil, false, bizerr.NewCode(CodeFormValueRequired)
		}
		return nil, false, nil
	}
	switch field.Type {
	case FieldTypeNumber:
		number, ok := toFloat64(value)
		if !ok {
			return nil, false, bizerr.NewCode(CodeFormValueInvalid)
		}
		return number, true, nil
	case FieldTypeDate:
		text, ok := value.(string)
		if !ok || len(text) < 8 {
			return nil, false, bizerr.NewCode(CodeFormValueInvalid)
		}
		return strings.TrimSpace(text), true, nil
	case FieldTypeSelect:
		text, ok := value.(string)
		if !ok {
			return nil, false, bizerr.NewCode(CodeFormValueInvalid)
		}
		text = strings.TrimSpace(text)
		found := false
		for _, option := range field.Options {
			if option == text {
				found = true
				break
			}
		}
		if !found {
			return nil, false, bizerr.NewCode(CodeFormValueInvalid)
		}
		return text, true, nil
	default:
		text, ok := value.(string)
		if !ok {
			return nil, false, bizerr.NewCode(CodeFormValueInvalid)
		}
		text = strings.TrimSpace(text)
		if text == "" {
			if field.Required {
				return nil, false, bizerr.NewCode(CodeFormValueRequired)
			}
			return nil, false, nil
		}
		return text, true, nil
	}
}

// normalizeDetailRows validates one submitted detail sub-table against the
// configured columns and returns the canonical row list.
func normalizeDetailRows(field FieldConfig, value any, present bool) ([]map[string]any, error) {
	if !present || value == nil {
		if field.Required {
			return nil, bizerr.NewCode(CodeFormValueRequired)
		}
		return nil, nil
	}
	rawRows, ok := value.([]any)
	if !ok || len(rawRows) == 0 {
		if field.Required {
			return nil, bizerr.NewCode(CodeFormValueRequired)
		}
		return nil, nil
	}
	if len(rawRows) > MaxDetailRows {
		return nil, bizerr.NewCode(CodeFormValueInvalid)
	}
	rows := make([]map[string]any, 0, len(rawRows))
	for _, rawRow := range rawRows {
		rowMap, ok := rawRow.(map[string]any)
		if !ok {
			return nil, bizerr.NewCode(CodeFormValueInvalid)
		}
		row := make(map[string]any, len(field.Columns))
		hasValue := false
		for _, column := range field.Columns {
			cellValue, cellPresent := rowMap[column.Key]
			if column.Type == DetailColumnNumber {
				if !cellPresent || cellValue == nil {
					continue
				}
				number, ok := toFloat64(cellValue)
				if !ok {
					return nil, bizerr.NewCode(CodeFormValueInvalid)
				}
				row[column.Key] = number
				hasValue = true
				continue
			}
			text, _ := cellValue.(string)
			text = strings.TrimSpace(text)
			if column.Type == DetailColumnDate && text != "" && len(text) < 8 {
				return nil, bizerr.NewCode(CodeFormValueInvalid)
			}
			if text != "" {
				row[column.Key] = text
				hasValue = true
			}
		}
		if hasValue {
			rows = append(rows, row)
		}
	}
	return rows, nil
}

// normalizeFormFieldAttachments validates one attachment field value list.
func normalizeFormFieldAttachments(field FieldConfig, value any, present bool) ([]string, error) {
	if !present || value == nil {
		if field.Required {
			return nil, bizerr.NewCode(CodeFormValueRequired)
		}
		return nil, nil
	}
	rawList, ok := value.([]any)
	if !ok {
		return nil, bizerr.NewCode(CodeFormValueInvalid)
	}
	urls := make([]string, 0, len(rawList))
	for _, item := range rawList {
		text, ok := item.(string)
		if !ok {
			return nil, bizerr.NewCode(CodeFormValueInvalid)
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}
		if !strings.HasPrefix(text, "http://") && !strings.HasPrefix(text, "https://") {
			return nil, bizerr.NewCode(CodeAttachmentsInvalid)
		}
		urls = append(urls, text)
	}
	if field.Required && len(urls) == 0 {
		return nil, bizerr.NewCode(CodeFormValueRequired)
	}
	return urls, nil
}

// extractAmountValue reads the first as-amount number field value from the
// submitted values for the request amount column projection.
func extractAmountValue(fields []FieldConfig, values map[string]any) float64 {
	for _, field := range fields {
		if !field.AsAmount || field.Type != FieldTypeNumber {
			continue
		}
		if number, ok := toFloat64(values[field.Key]); ok {
			return number
		}
	}
	return 0
}

// toFloat64 converts JSON number values (float64 or json.Number-like strings)
// into float64.
func toFloat64(value any) (float64, bool) {
	switch typed := value.(type) {
	case float64:
		return typed, true
	case int:
		return float64(typed), true
	case int64:
		return float64(typed), true
	case string:
		parsed, err := strconv.ParseFloat(strings.TrimSpace(typed), 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}
