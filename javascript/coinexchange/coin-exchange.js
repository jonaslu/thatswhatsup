/*
Word explanation of what it does:

We fill an array of the length of the sum + 1 (first element is the sum 0) with infinity.
This is how many coins we initially need to obtain that sum (except for the sum of 0 which requires 0 coins)

We then check the (partial) sums all the way from 0 to sum with the lowest coin denominator:
Can we exchange our current count of coins for this denominator?

Since the amount is infinity, if the sum is divisible by the coin denominator it means we can.
But we have already calculated how many coins we need for the sum minus one of this coins valor
(E g 6 = 1 + 5 where 5 is the number of coins needed to build the sum 5).

So we only need to calculate if we can exchange one coin for whatever amount of coins
we needed for the sum - 1 coin. This is already stored in the array at position [sum - coinValor].

If it's possible and the number of coins are less than what we already needed, exchange
our current valor for the amount of coins in the lower valor and move on.

We then check the position at the sum we're trying to obtain. If it's infinity - it
means we couldn't exchange any valor of coins for that sum and it's not obtainable.

If it's less than infinity it's the minimum number of coins we needed to obtain that sum.

Improvements:
Prune the coins array if the sum is less than some valors (eg the sum is 3 and the coin valors are [1,2,5]
the coin with valor 5 never needs checking since it's bigger than 3).
Boundary checks also o f c.
*/

const sum = 11;
const coins = [1,2,5];

const arr = Array(sum + 1).fill(Infinity);
arr[0] = 0;

for (let coinIndex=0; coinIndex < coins.length; coinIndex++) {
  const coinValor = coins[coinIndex];
  for (let currentSum=1;currentSum < arr.length; currentSum++) {
    if ((currentSum - coinValor) >= 0) {
      const exchangedCoin = arr[currentSum - coinValor] + 1;
      arr[currentSum] = Math.min(exchangedCoin, arr[currentSum]);
    }
  }
}

console.log("Min number of coins:", arr[sum] < Infinity ? arr[sum] : -1);
