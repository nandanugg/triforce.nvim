package utility

// SlimMap generic function to convert from slice of one type to slice of another type
func SlimMap[SourceType any, DestinationType any](
	input []SourceType,
	mapperFn func(value SourceType) DestinationType,
) []DestinationType {
	result := make([]DestinationType, len(input))
	for index, value := range input {
		result[index] = mapperFn(value)
	}
	return result
}
