# Repository Ranking Tool
This tool ranks Git repositories based on commit activity using a time-decay algorithm. Recent commits are weighted more heavily than older ones, providing a way to identify the most actively developed repositories.

## Algorithm Explanation

The Repository Ranking Tool uses a time-decay algorithm to weight repository activity, giving more importance to recent commits while gradually reducing the impact of older contributions.

### Time-Decay Scoring Algorithm

1. **Initial Setup**:
   - The algorithm defines a decay rate based on configurable parameters:
     - `defaultWeight` (0.3): The target weight for commits from `defaultDaysPeriod` ago
     - `defaultDaysPeriod` (30): The number of days used as a reference period

2. **Decay Rate Calculation**:
   - The decay rate is calculated using an exponential decay formula:
     ```
     decayRate = -ln(targetWeight) / (daysPeriod * 86400)
     ```
   - This ensures that commits from 30 days ago will have 30% of the weight of current commits

3. **Commit Processing**:
   - For each commit, the algorithm:
     - Records the repository, user, files changed, and lines added/deleted
     - Tracks the time range (oldest and newest commits)
     - Aggregates metrics per repository

4. **Score Calculation**:
   - For each commit, a weight is calculated based on its age:
     ```
     weight = e^(-decayRate * age)
     ```
   - Where `age` is the time difference (in seconds) between the commit and the most recent commit
   - The repository score is incremented by: `totalLinesChanged × weight`
   - This means:
     - The most recent commits have a weight of approximately 1.0
     - Older commits have progressively smaller weights
     - Very old commits contribute minimally to the final score

5. **Repository Ranking**:
   - Repositories are sorted by their final score in descending order
   - Only the top N repositories (specified by the `-tail` parameter) are included in the output

This approach ensures that actively maintained repositories rank higher than those with historical activity but little recent development.


## Building and Running
### Prerequisites
- Go 1.22 or higher

### Build Options
#### Option 1: Run directly
```bash
go run . -input="/path/to/commits.csv" -output="/path/to/output.json" -tail=20
```

#### Option 2: Build and run executable
```bash
# Build the executable
go build -o rank
# Run the executable
./rank -input="/path/to/commits.csv" -output="/path/to/output.json" -tail=20
```

## Command Line Flags
| Flag      | Description                                  | Default Value |
|-----------|----------------------------------------------|---------------|
| `-input`  | Path to the input CSV file with commit data  | `""`          |
| `-output` | Path to the output file (CSV or JSON)        | `""`          |
| `-tail`   | Number of top repositories to include        | `"20"`        |

## Input Format
The input CSV file should contain commit records with the following columns:
1. timestamp
2. username
3. repository
4. files
5. additions
6. deletions

## Output
The program outputs a JSON file containing the ranked repositories, sorted by score in descending order.
```json
[
  {
    "Name": "repo260",
    "Additions": 611871,
    "Deletions": 534997,
    "TotalLinesChanged": 1146868,
    "Files": 1497,
    "Score": 540982.6874018862
  },
  {
    "Name": "repo920",
    "Additions": 653495,
    "Deletions": 3004,
    "TotalLinesChanged": 656499,
    "Files": 3193,
    "Score": 346024.89946474356
  }
]
```

## Example Usage
```bash
go run . -input="./data/commits.csv" -output="./results/ranked_repos.json" -tail=50
```

This will:
1. Read commit data from `./data/commits.csv`
2. Rank repositories using the time-decay algorithm
3. Output the top 50 repositories to `./results/ranked_repos.json`

## Future Improvements

### Configuration File Support
- Implement support for a YAML/JSON configuration file that would allow:
  - Customizing time-decay parameters
  - Adjusting weighting factors for additions, deletions, and file counts
  - Setting different time windows for analysis
- This would make the tool more flexible without requiring code changes

### Alternative Ranking Algorithms
- Implement multiple scoring algorithms to address different business needs
- Add a `-algorithm` flag to select the desired scoring method

### Testing Infrastructure
- Develop a comprehensive testing suite:
  - Unit tests for core algorithm components
  - Integration tests using various sample datasets
  - Benchmark tests to ensure performance with large repositories
- Add test coverage reporting to identify untested code paths
