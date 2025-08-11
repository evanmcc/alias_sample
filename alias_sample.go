/******************************************************************************
 * File: AliasMethod.go
 * Author: Keith Schwarz (htiek@cs.stanford.edu)
 * Porter: Evan Vigil-McClanahan (mcclanahan@gmail.com)
 *
 * An implementation of the alias method implemented using Vose's algorithm.
 * The alias method allows for efficient sampling of random values from a
 * discrete probability distribution (i.e. rolling a loaded die) in O(1) time
 * each after O(n) preprocessing time.
 *
 * For a complete writeup on the alias method, including the intuition and
 * important proofs, please see the article "Darts, Dice, and Coins: Smpling
 * from a Discrete Distribution" at
 *
 *                 http://www.keithschwarz.com/darts-dice-coins/
 *
 * In his Archive, Keith says that anyone can use the code, and as such I have
 * done a very basic port from Java to Go. I have left the comments in place where
 * it made sense to me to do so.  Since no tests were provided, all code in
 * alias_sampler_test.go is my own.  - PEVM
 */

package alias_sample

import (
	r "math/rand"

	"github.com/gammazero/deque"
)

type AliasSampler struct {
	seed int64 // I save the initial seed since I want to use it in a different project
	rand *r.Rand

	probability []float64
	alias       []int
}

type SampleError struct {
	message string
}

func (e *SampleError) Error() string {
	return e.message
}

func Init(probs []float64) (*AliasSampler, error) {
	// grab a random seed
	seed := r.Int63()
	return InitWithSeed(probs, seed)
}

func InitWithSeed(probs []float64, seed int64) (*AliasSampler, error) {
	source := r.NewSource(seed)
	rand := r.New(source)

	if len(probs) == 0 {
		return nil, &SampleError{"no probabilities provided"}
	}

	/* Make a copy of the probabilities list, since we will be making
	 * changes to it.
	 */
	probs2 := make([]float64, len(probs))
	copy(probs2, probs)

	var tot float64
	for _, p := range probs2 {
		tot += p
	}

	for i := range probs2 {
		probs2[i] /= tot
	}

	probability := make([]float64, len(probs))
	alias := make([]int, len(probs))

	/* Compute the average probability and cache it for later use. */
	average := 1.0 / float64(len(probs))

	var small deque.Deque[int]
	var large deque.Deque[int]

	/* Populate the stacks with the input probabilities. */
	for i := range len(probs) {
		/* If the probability is below the average probability, then we add
		 * it to the small list; otherwise we add it to the large list.
		 */
		if probs2[i] >= average {
			large.PushBack(i)
		} else {
			small.PushBack(i)
		}
	}

	/* As a note: in the mathematical specification of the algorithm, we
	 * will always exhaust the small list before the big list.  However,
	 * due to floating point inaccuracies, this is not necessarily true.
	 * Consequently, this inner loop (which tries to pair small and large
	 * elements) will have to check that both lists aren't empty.
	 */
	for !(small.Len() == 0) && !(large.Len() == 0) {
		/* Get the index of the small and the large probabilities. */
		less := small.PopBack()
		more := large.PopBack()

		/* These probabilities have not yet been scaled up to be such that
		 * 1/n is given weight 1.0.  We do this here instead.
		 */
		probability[less] = probs2[less] * float64(len(probs))
		alias[less] = more

		/* Decrease the probability of the larger one by the appropriate
		 * amount.
		 */
		probs2[more] = (probs2[more] + probs2[less]) - average

		/* If the new probability is less than the average, add it into the
		 * small list; otherwise add it to the large list.
		 */
		if probs2[more] >= 1.0/float64(len(probs2)) {
			large.PushBack(more)
		} else {
			small.PushBack(more)
		}
	}

	/* At this point, everything is in one list, which means that the
	 * remaining probabilities should all be 1/n.  Based on this, set them
	 * appropriately.  Due to numerical issues, we can't be sure which
	 * stack will hold the entries, so we empty both.
	 */
	for small.Len() != 0 {
		probability[small.PopBack()] = 1.0
	}

	for large.Len() != 0 {
		probability[large.PopBack()] = 1.0
	}

	return &AliasSampler{
		seed:        seed,
		rand:        rand,
		probability: probability,
		alias:       alias,
	}, nil
}

func (s *AliasSampler) Next() int {
	/* Generate a fair die roll to determine which column to inspect. */
	column := s.rand.Intn(len(s.probability))

	/* Generate a biased coin toss to determine which option to pick. */
	coinToss := s.rand.Float64() < s.probability[column]

	/* Based on the outcome, return either the column or its alias. */
	if coinToss {
		return column
	} else {
		return s.alias[column]
	}
}
