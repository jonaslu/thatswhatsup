package singleLinkedList;

public interface IList<E> {
	public void add(E value);

	public void remove(int index);

	public E get(int index);

	public int getLength();

	public int getIndexOf(E value);

	public void replace(int index, E value);
}
