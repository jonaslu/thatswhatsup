package getchange.version2;

import static org.junit.Assert.assertEquals;

import java.util.ArrayList;
import java.util.Arrays;

import org.junit.Test;

import getchange.version2.GetChange2;

public class TestGetChange2 {

	@Test
	public void testEmptyList() {
		assertEquals(GetChange2.getNumberOfWays(new ArrayList<>(), 0), 0);
		assertEquals(GetChange2.getNumberOfWays(new ArrayList<>(), 15), 0);
	}

	@Test
	public void testOneInList() {
		assertEquals(GetChange2.getNumberOfWays(Arrays.asList(1), 1), 1);
		assertEquals(GetChange2.getNumberOfWays(Arrays.asList(1), 15), 1);
	}

	@Test
	public void testOneInComboNotPossible() {
		assertEquals(GetChange2.getNumberOfWays(Arrays.asList(3), 5), 0);
		assertEquals(GetChange2.getNumberOfWays(Arrays.asList(5), 3), 0);
	}

	@Test
	public void testVariousCases() {
		assertEquals(GetChange2.getNumberOfWays(Arrays.asList(5, 1), 13), 3);
		assertEquals(GetChange2.getNumberOfWays(Arrays.asList(20, 10, 5, 1), 13), 4);
		assertEquals(GetChange2.getNumberOfWays(Arrays.asList(20, 10, 5, 1), 11), 4);
	}

}
