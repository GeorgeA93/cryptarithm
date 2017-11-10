package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Sample struct {
	result  Result
	runTime time.Duration
}

type Result struct {
	mapping map[string]int
	seed    int64
}

type Try func() Result

func try(i int64, words []string, answer string) Try {
	seed := time.Now().UnixNano() * i
	rand.Seed(seed)
	return func() Result {
		firstLetters := make(map[string]bool)

		answerChars := strings.Split(answer, "")
		firstLetters[answerChars[0]] = true

		var allChars []string
		for _, word := range words {
			wordChars := strings.Split(word, "")
			firstLetters[wordChars[0]] = true
			allChars = append(allChars, wordChars...)
		}
		allChars = append(allChars, answerChars...)
		allChars = uniqueStrings(allChars)

		for {
			// Construct the map of char to number e.g. a: 1, b: 8, c: 2 ...
			mapping := make(map[string]int)
			shuffled := shuffleStrings(allChars)
			length := 9
			shuffledLen := len(shuffled)
			if shuffledLen > length {
				length = shuffledLen
			}
			nums := rand.Perm(length)
			for i, val := range shuffled {
				_, ok := firstLetters[val]
				// Numbers cannot start with 0
				if ok && nums[i] == 0 {
					addr := 1
					if i+addr > len(nums)-1 {
						addr = -1
					}
					// Swap the the 0 with another number
					nums[i], nums[i+addr] = nums[i+addr], nums[i]
				}
				mapping[val] = nums[i]
			}
			if len(mapping) != len(allChars) {
				continue
			}

			// Build the sum
			var sum []int
			for _, word := range words {
				wordChars := strings.Split(word, "")
				var numBuffer bytes.Buffer
				for _, char := range wordChars {
					str := strconv.FormatInt(int64(mapping[char]), 10)
					numBuffer.WriteString(str)
				}
				numStr := numBuffer.String()
				if strings.Split(numStr, "")[0] == "0" {
					continue
				}
				num, _ := strconv.Atoi(numStr)
				if num != 0 {
					sum = append(sum, num)
				}
			}

			// Build the expected answer
			var answerBuffer bytes.Buffer
			for _, char := range answerChars {
				str := strconv.FormatInt(int64(mapping[char]), 10)
				answerBuffer.WriteString(str)
			}
			answerStr := answerBuffer.String()
			expectedAnswer, err := strconv.Atoi(answerStr)
			if err != nil {
				continue
			}

			// Add the sum and compare to the answer
			guess := 0
			for _, num := range sum {
				guess += num
			}
			if guess == expectedAnswer {
				return Result{
					mapping: mapping,
					seed:    seed,
				}
			}
		}
	}
}

func shuffleStrings(input []string) []string {
	var shuffled []string
	shuffled = append(shuffled, input...)
	for i := range shuffled {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled
}

func uniqueStrings(input []string) []string {
	output := make([]string, 0, len(input))
	m := make(map[string]bool)

	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			output = append(output, val)
		}
	}

	return output
}

func race(words []string, answer string, contestants int64) Result {
	c := make(chan Result)
	contestant := func(i int64) { c <- try(i, words, answer)() }
	var i int64
	for i = 0; i < contestants; i++ {
		go contestant(i + 1)
	}
	return <-c
}

func meanDuration(durations []time.Duration) time.Duration {
	var totalDuration time.Duration
	for _, duration := range durations {
		totalDuration += duration
	}
	return time.Duration(int64(totalDuration) / int64(len(durations)))
}

func totalDuration(durations []time.Duration) time.Duration {
	var totalDuration time.Duration
	for _, duration := range durations {
		totalDuration += duration
	}
	return totalDuration
}

func sortDurations(durations []time.Duration, asc bool) []time.Duration {
	var sorted []time.Duration
	sorted = append(sorted, durations...)
	less := func(i, j int) bool {
		if asc {
			return int64(sorted[i]) < int64(sorted[j])
		}
		return int64(sorted[i]) > int64(sorted[j])
	}
	sort.Slice(sorted, less)
	return sorted
}

func fastestDuration(durations []time.Duration) time.Duration {
	return sortDurations(durations, true)[0]
}

func slowestDuration(durations []time.Duration) time.Duration {
	return sortDurations(durations, false)[0]
}

func medianDuration(durations []time.Duration) time.Duration {
	sorted := sortDurations(durations, true)
	size := len(sorted)
	if size%2 == 0 {
		return (sorted[size/2] + sorted[(size/2)+1]) / 2
	}
	return sorted[size/2]
}

func sortSamples(samples []Sample, asc bool) []Sample {
	var sorted []Sample
	sorted = append(sorted, samples...)
	less := func(i, j int) bool {
		if asc {
			return int64(sorted[i].runTime) < int64(sorted[j].runTime)
		}
		return int64(sorted[i].runTime) > int64(sorted[j].runTime)
	}
	sort.Slice(sorted, less)
	return sorted
}

func fastestSample(samples []Sample) Sample {
	return sortSamples(samples, true)[0]
}

func slowestSample(samples []Sample) Sample {
	return sortSamples(samples, false)[0]
}

func (s Sample) String() string {
	return fmt.Sprintf("\n\tMapping: %v\n\tSeed: %v\n\tRun Time: %v", s.result.mapping, s.result.seed, s.runTime)
}

func main() {
	words := []string{"alas", "lass", "no", "more"}
	answer := "cash"

	var samples []Sample
	samplesToTake := 100
	numContestants := int64(2)

	fmt.Println("Configuration")
	fmt.Println("=========================")
	fmt.Println("Number of Samples:", samplesToTake)
	fmt.Println("Parallelism:", numContestants)
	fmt.Println("")

	for i := 0; i < samplesToTake; i++ {
		start := time.Now()
		result := race(words, answer, numContestants)
		elapsed := time.Since(start)
		sample := Sample{result, elapsed}
		samples = append(samples, sample)
	}

	var runTimes []time.Duration
	for _, sample := range samples {
		runTimes = append(runTimes, sample.runTime)
	}
	fmt.Println("Results")
	fmt.Println("=========================")
	fmt.Println("Total Run Time:", totalDuration(runTimes))
	fmt.Println("Mean Run Time:", meanDuration(runTimes))
	fmt.Println("Median Run Time:", medianDuration(runTimes))
	fmt.Println("Fastest Sample:", fastestSample(samples))
	fmt.Println("Slowest Sample:", slowestSample(samples))
}
