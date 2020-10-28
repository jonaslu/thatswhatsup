# GetPossibleCombinationsOfCoins
Project Euler assigment #31 - https://projecteuler.net/problem=31


## Version 1
Returns the actual combinations in a list (e g ((10,0,3),(0,5,6), (0,10,3), (0,0,3)) for combinations of (10,5,1) and the sum 13 by generating a matrix and iterating over that matrix selecting the correct values. The matrix is generated from the number of times the different denominators of the coins can be divided in the sum (e g 13/10 = 1, 13/5=2, 13/1=13). Hence a matrix of ((0,10),(0,5,10),(0,1,2...13)) is generated. Resulting
((0,0,1)
 (0,0,2)
 ...
 (0,0,13)
 (0,5,1)
 (0,5,2)
 ...
 (10,10,13))
 
And the rows with sums of 13 are then selected.

## Version 2
Does what the assignment actually says :) via recursion counts the number of ways coins can be changed.

Like this:
If there are no more coins left to try - it's a dead end
If the sum less then zero - stop here, it's a dead end
If the sum is zero, we've hit one possible combination. Count it.

Keep the sum intact and try with only smaller coins
Deduct one of my current denominator and re-run the step above

