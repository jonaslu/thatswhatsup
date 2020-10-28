package dynamicArray;

import static org.junit.Assert.assertEquals;

import java.util.stream.IntStream;

import org.junit.Test;

public class DynamicArrayTest {

	@Test
	public void testIncrease() {
		DynamicArray<Integer> array = new DynamicArray<>(1);

		IntStream.range(0, 5).forEach(idx -> array.set(idx, idx));
		IntStream.range(0, 5).forEach(idx -> assertEquals(array.get(idx).intValue(), idx));
	}

	@Test
	public void testSetAtLargeIndex() {
		DynamicArray<Integer> array = new DynamicArray<>(1);
		array.set(100, 100);
		assertEquals(array.get(100).intValue(), 100);
	}

	@Test(expected = IndexOutOfBoundsException.class)
	public void testIndexOutOfBounds() {
		DynamicArray<Integer> array = new DynamicArray<>(1);
		array.get(1);
	}
}
