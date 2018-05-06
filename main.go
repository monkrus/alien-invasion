package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// QUESTIONS
// 1. Is numAliens < numCities? Any restriction on these?
// 2. Is a road always a 2-way connection? Meaning, if NY says it is connected
// to Boston, is does Boston say it is connected to NY?

var cities = make(map[string][]string)
var aliens = make(map[string][]int)
var uniqueCities = make([]string, 0)

func main() {
	populateCities()

	// Get number of aliens from program arg
	numAliens, err := strconv.Atoi(os.Args[1:2][0])
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	populateAliens(numAliens)

	runSimulation()

	fmt.Println(cities)
	fmt.Println(aliens)
}

// populateCities reads data from a file called "cities"
// and populates the cities map
func populateCities() {
	data, err := os.Open("cities")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer data.Close()

	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanLines)

	// Iterate over lines of cities file
	for scanner.Scan() {
		// Split city line by spaces
		c := strings.Split(scanner.Text(), " ")

		// If this city does not already exist
		if len(cities[c[0]]) == 0 {
			// Create new city with name of current line
			cities[c[0]] = make([]string, 0)

			// Add to uniqueCities for use in alien assignments
			uniqueCities = append(uniqueCities, c[0])
		}

		// Regex to filter out direction and equals sign
		reg, err := regexp.Compile(`^(.*?)=`)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		// Add neighbors of current city
		// TODO Probably good to make sure this isn't the same as current city
		for _, neighbor := range c[1:] {
			filtered := reg.ReplaceAllString(neighbor, "")

			cities[c[0]] = append(cities[c[0]], filtered)
		}
	}
}

func populateAliens(numAliens int) {
	rand.Seed(time.Now().Unix())

	for i := 0; i < numAliens; i++ {
		city := uniqueCities[rand.Intn(len(uniqueCities))]

		// Ensure no cities have more than 2 aliens
		for len(aliens[city]) == 2 {
			city = uniqueCities[rand.Intn(len(uniqueCities))]
		}

		aliens[city] = append(aliens[city], i)
	}
}

func runSimulation() {
	for city := range aliens {
		if len(aliens[city]) == 2 {
			fmt.Printf("%v has been destroyed by alien %v and alien %v!\n", city, aliens[city][0], aliens[city][1])

			// Go to neighbors of deleted and delete itself from their lists
			for _, neighbor := range cities[city] {
				for i, n := range cities[neighbor] {
					if n == city {
						cities[neighbor][i] = cities[neighbor][len(cities[neighbor])-1]
						cities[neighbor] = cities[neighbor][:len(cities[neighbor])-1]
					}
				}
			}
			delete(aliens, city)
			delete(cities, city)
		}
	}
}
