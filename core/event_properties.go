package core

import "strconv"

// BoolProperty implements EventProperty.BoolProperty.
func (e *EventSource) BoolProperty(key string, defaultValue bool) bool {
	return extractBoolProperty(e.Properties, key, defaultValue)
}

// BoolProperty implements EventProperty.BoolProperty.
func (e *EventResult) BoolProperty(key string, defaultValue bool) bool {
	return extractBoolProperty(e.Properties, key, defaultValue)
}

func extractBoolProperty(props EventProperties, key string, defaultValue bool) bool {
	if props == nil {
		return defaultValue
	}
	raw, exists := props[key]
	if !exists {
		return defaultValue
	}
	if boolValue, isBool := raw.(bool); isBool {
		return boolValue
	}
	if intValue, isNumber := raw.(int); isNumber {
		if intValue == 0 {
			return false
		}
		if intValue == 1 {
			return true
		}
		return defaultValue
	}
	if strValue, isString := raw.(string); isString {
		if value, err := strconv.ParseBool(strValue); err == nil {
			return value
		}
		return defaultValue
	}
	return defaultValue
}

// StringProperty implements EventProperty.StringProperty.
func (e *EventSource) StringProperty(key string, defaultValue string) string {
	return extractStringProperty(e.Properties, key, defaultValue)
}

// StringProperty implements EventProperty.StringProperty.
func (e *EventResult) StringProperty(key string, defaultValue string) string {
	return extractStringProperty(e.Properties, key, defaultValue)
}

func extractStringProperty(props EventProperties, key, defaultValue string) string {
	if props == nil {
		return defaultValue
	}
	raw, exists := props[key]
	if !exists {
		return defaultValue
	}
	if strValue, isString := raw.(string); isString {
		return strValue
	}
	return defaultValue
}
