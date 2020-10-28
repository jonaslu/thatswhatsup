package getchange.version1;

import static org.junit.Assert.assertTrue;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import org.junit.Test;

import getchange.version1.GetChange1;

public class TestGetChange1 {

	@Test
	public void testEmptyList() {
		GetChange1 getChange2 = new GetChange1(new ArrayList<>(), 13);
		assertTrue(getChange2.getPossibleChangeCombinations().isEmpty());
	}

	@Test
	public void testOneInListTotalTooSmall() {
		GetChange1 getChange2 = new GetChange1(Arrays.asList(3), 5);
		assertTrue(getChange2.getPossibleChangeCombinations().isEmpty());
	}

	@Test
	public void testOneInList() {
		GetChange1 getChange2 = new GetChange1(Arrays.asList(2), 4);
		List<List<Integer>> possibleChangeCombinations = getChange2.getPossibleChangeCombinations();
		assertTrue(possibleChangeCombinations.get(0).equals(Arrays.asList(4)));
		assertTrue(possibleChangeCombinations.size() == 1);
	}

	@Test
	public void testTwoInList() {
		GetChange1 getChange2 = new GetChange1(Arrays.asList(2, 8), 10);

		List<List<Integer>> expected = new ArrayList<>();
		expected.add(Arrays.asList(10, 0));
		expected.add(Arrays.asList(2, 8));

		List<List<Integer>> possibleChangeCombinations = getChange2.getPossibleChangeCombinations();

		testListsSameSizeAndConentContainsSame(expected, possibleChangeCombinations);
	}

	@Test
	public void testThreeInList() {
		GetChange1 getChange2 = new GetChange1(Arrays.asList(1, 5, 10), 13);

		List<List<Integer>> expected = new ArrayList<>();
		expected.add(Arrays.asList(13, 0, 0));
		expected.add(Arrays.asList(8, 5, 0));
		expected.add(Arrays.asList(3, 10, 0));
		expected.add(Arrays.asList(3, 0, 10));

		List<List<Integer>> possibleChangeCombinations = getChange2.getPossibleChangeCombinations();

		testListsSameSizeAndConentContainsSame(expected, possibleChangeCombinations);
	}

	@Test
	public void testTooLargeDenominatorsDiscarded() {
		GetChange1 getChange2 = new GetChange1(Arrays.asList(1, 5, 10, 20), 5);

		List<List<Integer>> expected = new ArrayList<>();
		expected.add(Arrays.asList(5, 0));
		expected.add(Arrays.asList(0, 5));

		List<List<Integer>> possibleChangeCombinations = getChange2.getPossibleChangeCombinations();
		testListsSameSizeAndConentContainsSame(expected, possibleChangeCombinations);
	}

	@Test
	public void testUnpossibleCombination() {
		GetChange1 getChange2 = new GetChange1(Arrays.asList(3), 5);
		assertTrue(getChange2.getPossibleChangeCombinations().isEmpty());
	}

	@Test
	public void testUnpossibleCombinationThreeDenomiators() {
		GetChange1 getChange2 = new GetChange1(Arrays.asList(3, 7, 13), 5);
		assertTrue(getChange2.getPossibleChangeCombinations().isEmpty());
	}

	private void testListsSameSizeAndConentContainsSame(List<List<Integer>> expected, List<List<Integer>> result) {
		for (int i = 0; i < result.size(); i++) {
			assertTrue(expected.contains(result.get(i)));
		}
		assertTrue(result.size() == expected.size());
	}
}
