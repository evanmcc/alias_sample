# An Weighted Sampling library for go, using the Alias Method

This is all explained better than I could ever manage at the (long) link below:
http://www.keithschwarz.com/darts-dice-coins/

But the long and the short of it is that if you have a list of probabilities (as float64s, no need to normalize them), it will quite efficiently give you a matching distribution of samples, which you can use for weighted dice rolls, etc.

Most of the code here was written by the author of the above blog post, I ported it to go, added the normalization routine (the original code assumed but didn't check prior normalization), and added a property test that attempts to make sure that the distribution isn't entirely out of whack.
