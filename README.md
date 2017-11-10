# cryptarithm

A quick and dirty attempt at brute forcing cryptarithmetic puzzles in golang. **Be warned I am very new to go 😜** If you happen to be looking at this code, and see some terrible terrible things, please do point me towards them...!

These puzzles are pretty simple, so the algorithm is relatively fast.

Following the pattern mentioned [here](https://talks.golang.org/2012/concurrency.slide#48) I was able to construct a race, such that I could spin up multiple go routines running the solver, where the fastest one wins. I added this feature due to the non-deterministic nature of the algorithm. It relies on randomness, therefore the worst case, well is pretty bad.

To get some insight into runtimes etc, I have added some basic support for running the algorithm multiple times and getting some statistics (mean/median).