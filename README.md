# Solution Explained
In this project, I have created a simple layout that follows the design of one of your projects on Github. It also utilizes Cobra library for the CLI implementation.

## Important Note
PLEASE run `make unzip` in order to unzip and run with the fixtures sample provided by you. Or make sure to have a sample file with recipes called files/fixture.json.

Because the JSON file size is unknown, it could be in GB, I had to use JSON streaming read its content. This helps in avoiding memory issues, but runs a little bit slow.
While I used the standard JSON library, I wanted to mention that there are other third party libraries that could improve performance, i.e. runs better
than standard library. 

Current implementation calculate the stats within 10-15 seconds, given the fixtures.json with around **5000_0000** lines, on my Macbook pro. However, it runs slower in docker container and takes around 2m for the same file. 

## Time and Space Complexities
Maps were used where possible. Maps have time complexity of O(1). Since the task explicitly indicate that distinct recipes names is lower than 2K, I have declared maps with predefined size `recipeCounts := make(map[string]int, 2000)`, which improves performance, since golang don't have to grow the map on every new key added. The same also applies to distinct postcodes, lower than 1M, I also declared with predefined size `postCodeCounts := make(map[string]int, 1000_000)`. 

Maps used wherever possible also to increase lookup time. For instance, when searching for recipes names that contains some words, I used a map to convert them into a map, and then later used to look for matching words from each recipe name. 
```go
wordsMap := make(map[string]bool, len(s.cfg.Words))
for _, word := range s.cfg.Words {
    wordsMap[strings.ToLower(word)] = true
}
```

To show counts for each distinct recipe name, I used the same map `recipeCounts := make(map[string]int, 2000)`, took all keys, sorted them using standard library `slices.Sort(keys)`, which is based on **quicksort** algorithm, the most efficient sorting algorithm.

## Design Patterns
For the configuration I used the builder pattern. Configs runs with default, however, if you provide arguments with your command, e.g. file name, these will be applied by using the builder pattern. See cmd/cobra.go. Other patterns used is Strategy patters, e.g. in `JsonStats` struct has a dependency on the `parser.Parser`, allowing for different parser implementations to be injected. 

Additionally, making sure to make the code more modular and testable, I always try to apply SOLID principles.

## Data Sanitization

Proper data sanitization applied as per requirements. For instance, delivery formats check, postcode length checks, recipes length checks, etc.

## Unit Tests
Unit tests were applied to the most critical parts, however, not fully covering everything due to time limitations. You can run the tests via `make test`.

## How to Test/Run
First, the tool runs with default configurations, please see `pkg/config/config.go`. These configs can be overwritten via environment variables. Simply export your variables, this will allow you not to provide arguments for convenience. 

To build the binary in Docker run: `make dockerBuildRunWithDefaultArgs` this will run the tool with default args. Default args are:
- file: /app/files/fixtures.json
- words: Potato,Mushroom,Veggie
- postcode: 10120
- from time: 10AM
- to time: 3PM

Please note running with default args means loading a fixtures.json file with **5000_000** lines. It will take more than 1m to finish processing and displaying the result.

If you want to test with your own file, words, or time range, see the example `make dockerBuildRunWithCustomArgs`to run yours. 

If you want to open a shell to docker container and run the tool, then simply run `parser stats` (already in $path) and add any of your desired arguments. Otherwise it will run with default ones.

## Future Improvements
- More in depth unit tests
- Integration tests

----
END OF MY INSTRUCTIONS
----

Recipe Stats Calculator
====

In the given assignment we suggest you to process an automatically generated JSON file with recipe data and calculated some stats.

Instructions
-----

1. Clone this repository.
2. Create a new branch called `dev`.
3. Create a pull request from your `dev` branch to the master branch.
4. Reply to the thread you're having with your recruiter telling them we can start reviewing your code

Given
-----

Json fixtures file with recipe data. Download [Link](https://raw.githubusercontent.com/hellofreshdevtests/json-fixtures/content/hf_test_calculation_fixtures.tar.gz)

_Important notes_

1. Property value `"delivery"` always has the following format: "{weekday} {h}AM - {h}PM", i.e. "Monday 9AM - 5PM"
2. The number of distinct postcodes is lower than `1M`, one postcode is not longer than `10` chars.
3. The number of distinct recipe names is lower than `2K`, one recipe name is not longer than `100` chars.

Functional Requirements
------

1. Count the number of unique recipe names.
2. Count the number of occurences for each unique recipe name (alphabetically ordered by recipe name).
3. Find the postcode with most delivered recipes.
4. Count the number of deliveries to postcode `10120` that lie within the delivery time between `10AM` and `3PM`, examples _(`12AM` denotes midnight)_:
    - `NO` - `9AM - 2PM`
    - `YES` - `10AM - 2PM`
5. List the recipe names (alphabetically ordered) that contain in their name one of the following words:
    - Potato
    - Veggie
    - Mushroom

Non-functional Requirements
--------

1. The application is packaged with [Docker](https://www.docker.com/).
2. Setup scripts are provided.
3. The submission is provided as a `CLI` application.
4. The expected output is rendered to `stdout`. Make sure to render only the final `json`. If you need to print additional info or debug, pipe it to `stderr`.
5. It should be possible to (implementation is up to you):  
    a. provide a custom fixtures file as input  
    b. provide custom recipe names to search by (functional reqs. 5)  
    c. provide custom postcode and time window for search (functional reqs. 4)  

Expected output
---------------

Generate a JSON file of the following format:

```json5
{
    "unique_recipe_count": 15,
    "count_per_recipe": [
        {
            "recipe": "Mediterranean Baked Veggies",
            "count": 1
        },
        {
            "recipe": "Speedy Steak Fajitas",
            "count": 1
        },
        {
            "recipe": "Tex-Mex Tilapia",
            "count": 3
        }
    ],
    "busiest_postcode": {
        "postcode": "10120",
        "delivery_count": 1000
    },
    "count_per_postcode_and_time": {
        "postcode": "10120",
        "from": "11AM",
        "to": "3PM",
        "delivery_count": 500
    },
    "match_by_name": [
        "Mediterranean Baked Veggies", "Speedy Steak Fajitas", "Tex-Mex Tilapia"
    ]
}
```

Review Criteria
---

We expect that the assignment will not take more than 3 - 4 hours of work. In our judgement we rely on common sense
and do not expect production ready code. We are rather instrested in your problem solving skills and command of the programming language that you chose.

It worth mentioning that we will be testing your submission against different input data sets.

__General criteria from most important to less important__:

1. Functional and non-functional requirements are met.
2. Prefer application efficiency over code organisation complexity.
3. Code is readable and comprehensible. Setup instructions and run instructions are provided.
4. Tests are showcased (_no need to cover everything_).
5. Supporting notes on taken decisions and further clarifications are welcome.

