package dynamicArray;

import java.util.stream.IntStream;

public class DynamicArray<E> {

	private int size;
	private Object[] array;

	public DynamicArray() {
		this(2);
	}

	public DynamicArray(int initialSize) {
		this.size = initialSize;
		this.array = new Object[initialSize];
	}

	public void set(int index, E object) {
		expandIfNecessary(index);
		array[index] = object;
	}

	private void expandIfNecessary(int index) {
		if (index >= size) {
			int newSize = index * 2;
			Object[] newArray = new Object[newSize];

			IntStream.range(0, array.length).forEach(idx -> newArray[idx] = array[idx]);
			array = newArray;
			size = newSize;
		}
	}

	@SuppressWarnings("unchecked")
	public E get(int index) {
		return (E) array[index];
	}
}
