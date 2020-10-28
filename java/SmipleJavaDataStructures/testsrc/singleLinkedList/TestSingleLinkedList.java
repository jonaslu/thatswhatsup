package singleLinkedList;

import static org.junit.Assert.assertTrue;

import org.junit.Before;
import org.junit.Test;

public class TestSingleLinkedList {

	private IList<Integer> integerList;

	@Before
	public void setUp() {
		integerList = new SingleLinkedList<>();
	}

	@Test
	public void testEmptyListSizeZero() {
		assertTrue("[]".equals(integerList.toString()));
		assertTrue(integerList.getLength() == 0);
	}

	@Test(expected = IndexOutOfBoundsException.class)
	public void testEmptyListThrowsOutOfBoundsException() {
		integerList.get(4);
	}

	@Test
	public void testAddElements() {
		integerList.add(3);
		assertTrue("[3]".equals(integerList.toString()));
		assertTrue(integerList.get(0) == 3);
		assertTrue(integerList.getLength() == 1);

		integerList.add(4);
		assertTrue("[3,4]".equals(integerList.toString()));
		assertTrue(integerList.get(0) == 3);
		assertTrue(integerList.get(1) == 4);
		assertTrue(integerList.getLength() == 2);
	}

	@Test
	public void testRemoveElements() {
		integerList.add(5);
		integerList.add(6);

		integerList.remove(1);
		assertTrue(integerList.get(0) == 5);
		assertTrue("[5]".equals(integerList.toString()));
		assertTrue(integerList.getLength() == 1);

		integerList.remove(0);
		assertTrue("[]".equals(integerList.toString()));
		assertTrue(integerList.getLength() == 0);
	}

	@Test(expected = IndexOutOfBoundsException.class)
	public void testRemoveElementsThrowsOutOfBoundsException() {
		integerList.remove(0);
		integerList.remove(10);
	}

	@Test(expected = IndexOutOfBoundsException.class)
	public void testReplaceThrowsOutOfBounds() {
		integerList.replace(0, 1);

		integerList.add(2);
		integerList.replace(1, 1);
	}

	@Test
	public void testGetIndexOfNotPresentValue() {
		assertTrue(integerList.getIndexOf(3) == -1);
		integerList.add(1);
		assertTrue(integerList.getIndexOf(3) == -1);
	}

	@Test
	public void testGetIndexOf() {
		integerList.add(2);
		assertTrue(integerList.getIndexOf(2) == 0);

		integerList.add(3);

		assertTrue(integerList.getIndexOf(2) == 0);
		assertTrue(integerList.getIndexOf(3) == 1);

		integerList.remove(0);
		assertTrue(integerList.getIndexOf(2) == -1);
		assertTrue(integerList.getIndexOf(3) == 0);

		integerList.remove(0);
		assertTrue(integerList.getIndexOf(3) == -1);

		integerList.add(3);
		assertTrue(integerList.getIndexOf(3) == 0);
	}

	@Test
	public void testReplace() {
		integerList.add(2);
		assertTrue(integerList.getIndexOf(2) == 0);
		assertTrue(integerList.get(0) == 2);
		integerList.replace(0, 3);

		assertTrue(integerList.get(0) == 3);
		assertTrue(integerList.getIndexOf(3) == 0);

		integerList.add(3);
		integerList.replace(1, 4);

		assertTrue(integerList.get(1) == 4);
		assertTrue(integerList.getIndexOf(4) == 1);

		integerList.add(5);
		integerList.replace(1, 6);

		assertTrue(integerList.get(1) == 6);
		assertTrue(integerList.getIndexOf(6) == 1);

		integerList.remove(1);

		assertTrue(integerList.get(1) == 5);
		assertTrue(integerList.getIndexOf(5) == 1);
	}
}
