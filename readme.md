# Pack calculator

## Project structure

The project has the following structure:
- endpoints - Defines logic for the API endpoints (the project could be improved by moving requests and responses logic outside of this package)
- helpers - Multiple helper functions (including the algorithm for solving the optimal pack combination challenge)
- logger - A single instance logger using logrus used across the project
- storage - Persistent storage options (for this project I used redis for simplicity)
- main.go - Loads .env, setup logger, setup persistent storage, setup api endpoint and start the API server

## Run

To run the project is enough to do:

```
docker-compose up --build
```

To view the user interface just open index.html in a browser. 

**Important** - The API uses port 8080 and the user interface does the request directly to localhost:8080. Changes regarding this must be updated in the user interface as well.

## Tests

To run tests just use the following command:

```go test ./tests```

## Algorithm

The optimal pack combination algorithm is implemented in `helpers/optimal_pack_combination.go`.

### How It Works

- The algorithm uses a **priority queue** (min-heap) of partial solutions (`Solution` objects).
- Each solution tracks:
  - Total number of items
  - Total number of packs used
  - A map of pack sizes to their quantities (the desired result)
- The algorithm it's like Djikstra's Algorithm for finding the shortest-path in a graph.

### Initial Setup

- The queue starts with a single initial solution: `0 items, 0 packs`.

### Main Loop

The algorithm processes the queue until it's empty:

1. **Pop the best solution** (lowest total items, fewest packs) from the queue.
2. **Check if it meets or exceeds the required amount**.  
   - If yes, this is our optimal solution (and the loop ends).
3. **Otherwise**, try adding each pack size to the current solution:
   - Create a new solution with updated item total and pack count.
   - **Only add it to the queue if**:
     - This total number of items hasn't been seen before, or
     - It uses fewer packs than the previously seen solution for this item total.

### Why It Works

Because the priority queue always explores **the most promising solution first**, the first solution that reaches the target amount is guaranteed to:
- Use the fewest possible total items,
- And, if multiple ways exist with the same item count, use the fewest packs.

This is effectively a **best-first search**, similar to Dijkstra’s algorithm but without a fixed graph.
