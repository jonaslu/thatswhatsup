package hashTable;

import static org.junit.Assert.assertTrue;

import org.junit.Before;
import org.junit.Test;

public class testHashTable {

	private class CollidingInteger {
		private int value;

		public CollidingInteger(int value) {
			this.value = value;
		}

		@Override
		public int hashCode() {
			return 10;
		}

		@Override
		public boolean equals(Object obj) {
			if (obj != null && obj instanceof CollidingInteger) {
				CollidingInteger object = (CollidingInteger) obj;
				return object.value == object.value;
			}
			return false;
		}
	}

	// Integer hashCode returns the int so we can test the objects when the
	// table resizes
	private HashTable<Integer, Integer> integerHashTable;

	@Before
	public void setUp() {
		integerHashTable = new HashTable<>(2);
	}

	@Test
	public void testPutOneValueReturnsSame() {
		integerHashTable.put(1, 3);
		assertTrue(integerHashTable.get(1) == 3);
	}

	@Test
	public void testPutValueRetrieveOtherReturnsNull() {
		integerHashTable.put(1, 56);
		assertTrue(integerHashTable.get(2) == null);
	}

	@Test
	public void testTableResizes() {
		integerHashTable.put(0, 56);
		integerHashTable.put(1, 57);
		integerHashTable.put(3, 58);

		assertTrue(integerHashTable.get(0) == 56);
		assertTrue(integerHashTable.get(1) == 57);
		assertTrue(integerHashTable.get(3) == 58);
	}

	@Test
	public void hashCollides() {
		HashTable<CollidingInteger, Integer> collidingIntegerHashTable = new HashTable<>(1);
		collidingIntegerHashTable.put(new CollidingInteger(1), 2);
		collidingIntegerHashTable.put(new CollidingInteger(2), 3);

		assertTrue(collidingIntegerHashTable.get(new CollidingInteger(2)) == 3);
	}

	@Test
	public void testWithStrings() {
		HashTable<String, Integer> stringHashTable = new HashTable<>(1);
		stringHashTable.put("blorgh", 2);
		stringHashTable.put("Blargh", 3);

		assertTrue(stringHashTable.get("blorgh") == 2);
		assertTrue(stringHashTable.get("Blargh") == 3);

		Integer prevVal = stringHashTable.put("blorgh", 4);
		assertTrue(prevVal == 2);
		assertTrue(stringHashTable.get("blorgh") == 4);
	}
}
