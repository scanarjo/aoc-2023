package seedmapping

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MappingRange struct {
	destinationStart int
	sourceStart      int
	rangeLength      int
}

type Mapping struct {
	source      string
	destination string

	ranges []MappingRange
}

type Measurement struct {
	value int
	unit  string
}

func (mapping *Mapping) Map(measurement Measurement) (*Measurement, error) {
	if measurement.unit != mapping.source {
		return nil, errors.New("invalid unit")
	}

	destinationMeasurement := Measurement{value: measurement.value, unit: mapping.destination}

	for _, mappingRange := range mapping.ranges {
		if measurement.value >= mappingRange.sourceStart && measurement.value < mappingRange.sourceStart+mappingRange.rangeLength {
			destinationMeasurement.value = mappingRange.destinationStart + (measurement.value - mappingRange.sourceStart)
			break
		}
	}

	return &destinationMeasurement, nil
}

func (mapping *Mapping) ReverseMap(destination Measurement) (*Measurement, error) {
	if destination.unit != mapping.destination {
		return nil, errors.New("invalid unit")
	}

	source := Measurement{value: destination.value, unit: mapping.source}

	for _, mappingRange := range mapping.ranges {
		if destination.value >= mappingRange.destinationStart && destination.value < mappingRange.destinationStart+mappingRange.rangeLength {
			source.value = mappingRange.sourceStart + (destination.value - mappingRange.destinationStart)
			break
		}
	}

	return &source, nil
}

func parseHeader(header string) (string, string, error) {
	values := strings.Fields(header)

	if len(values) != 2 {
		return "", "", fmt.Errorf("invalid header: %s", header)
	}

	if values[1] != "map:" {
		return "", "", fmt.Errorf("invalid header: %s", header)
	}

	parts := strings.Split(values[0], "-to-")

	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid header: %s", header)
	}

	return parts[0], parts[1], nil
}

func DeserializeMapping(scanner *bufio.Scanner) (*Mapping, error) {
	if !scanner.Scan() {
		return nil, errors.New("scanner exhausted")
	}

	header := scanner.Text()

	source, destination, err := parseHeader(header)

	if err != nil {
		return nil, err
	}

	mapping := Mapping{
		source:      source,
		destination: destination,
	}

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		values := strings.Fields(line)
		if len(values) != 3 {
			return nil, fmt.Errorf("invalid mapping range: %s", line)
		}

		destinationStart, err := strconv.Atoi(values[0])
		if err != nil {
			return nil, fmt.Errorf("failed to parse destination start: %w", err)
		}

		sourceStart, err := strconv.Atoi(values[1])
		if err != nil {
			return nil, fmt.Errorf("failed to parse source start: %w", err)
		}

		rangeLength, err := strconv.Atoi(values[2])
		if err != nil {
			return nil, fmt.Errorf("failed to parse range length: %w", err)
		}

		mapping.ranges = append(mapping.ranges, MappingRange{
			destinationStart: destinationStart,
			sourceStart:      sourceStart,
			rangeLength:      rangeLength,
		})
	}

	return &mapping, nil
}

func DeserializeSeedList(serializedSeedList string) ([]Measurement, error) {
	stringSeeds := strings.Fields(serializedSeedList)

	if stringSeeds[0] != "seeds:" {
		return nil, fmt.Errorf("invalid seed list: %s", serializedSeedList)
	}

	seeds := []Measurement{}
	for i := 1; i < len(stringSeeds); i++ {
		seed := Measurement{
			unit: "seed",
		}

		seed.value, _ = strconv.Atoi(stringSeeds[i])

		seeds = append(seeds, seed)
	}

	return seeds, nil
}

type MeasurementPair struct {
	min Measurement
	max Measurement
}

func DeserializeSeedListV2(serializedSeedList string) ([]MeasurementPair, error) {
	stringSeeds := strings.Fields(serializedSeedList)

	if stringSeeds[0] != "seeds:" {
		return nil, fmt.Errorf("invalid seed list: %s", serializedSeedList)
	}

	seeds := []MeasurementPair{}
	for i := 1; i+1 < len(stringSeeds); i += 2 {
		start, _ := strconv.Atoi(stringSeeds[i])
		rangeLength, _ := strconv.Atoi(stringSeeds[i+1])

		min := Measurement{
			value: start,
			unit:  "seed",
		}

		max := Measurement{
			value: start + rangeLength - 1,
			unit:  "seed",
		}

		seeds = append(seeds, MeasurementPair{min: min, max: max})
	}

	fmt.Fprintln(os.Stdout, []any{"seeds %v", seeds}...)

	return seeds, nil
}
