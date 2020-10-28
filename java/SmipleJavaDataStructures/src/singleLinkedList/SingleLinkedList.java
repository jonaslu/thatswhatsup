package singleLinkedList;

public class SingleLinkedList<E> implements IList<E> {

	public class Element<T> {
		T value;
		Element<T> next = null;

		public Element(T value) {
			this(value, null);
		}

		public Element(T value, Element<T> next) {
			this.value = value;
			this.next = next;
		}

		public void setNext(Element<T> element) {
			this.next = element;
		}

		public Element<T> getNext() {
			return next;
		}

		public T getValue() {
			return value;
		}
	}

	Element<E> head = null;
	Element<E> tail = null;
	int length;

	@Override
	public void add(E value) {
		if (tail == null) {
			head = new Element<>(value);
			tail = head;
			length = 1;
		} else {
			Element<E> newTail = new Element<>(value);
			tail.setNext(newTail);
			tail = newTail;
			length += 1;
		}
	}

	@Override
	public void remove(int index) {
		assertIndexIsWithinBounds(index);

		if (index == 0) {
			head = head.next;
			if (head == null) {
				tail = null;
			}
		} else {
			Element<E> next = head;
			for (int i = 0; i < index - 1; i++) {
				next = next.getNext();
			}

			Element<E> nextInList = next.getNext().getNext();
			next.setNext(nextInList);

			if (index == length - 1) {
				tail = nextInList;
			}
		}

		length -= 1;
	}

	@Override
	public E get(int index) {
		assertIndexIsWithinBounds(index);

		Element<E> getElement = head;

		for (int i = 0; i < index; i++) {
			getElement = head.getNext();
		}

		return getElement.getValue();
	}

	@Override
	public String toString() {
		if (head == null) {
			return "[]";
		} else {
			String retVal = "[";

			Element<E> next = head;
			while (next != null) {
				retVal += next.getValue();
				next = next.getNext();

				if (next != null) {
					retVal += ",";
				}
			}

			retVal += "]";
			return retVal;
		}
	}

	@Override
	public int getLength() {
		return length;
	}

	@Override
	public int getIndexOf(E value) {
		Element<E> next = head;

		for (int i = 0; i < length; i++) {
			if (next.getValue().equals(value)) {
				return i;
			}
			next = next.getNext();
		}

		return -1;
	}

	private void assertIndexIsWithinBounds(int index) {
		if (index < 0 || length == 0 || index >= length) {
			throw new IndexOutOfBoundsException();
		}
	}

	@Override
	public void replace(int index, E value) {
		assertIndexIsWithinBounds(index);

		Element<E> insertElement = new Element<E>(value);

		if (index == 0) {
			if (tail.equals(head)) {
				tail = insertElement;
			} else {
				insertElement.setNext(head.getNext());
			}

			head = insertElement;
		} else if (index < length - 1) {
			Element<E> prev = getElementAtIndex(index - 1);
			Element<E> next = getElementAtIndex(index + 1);

			insertElement.setNext(next);
			prev.setNext(insertElement);
		} else {
			Element<E> prev = getElementAtIndex(index - 1);
			prev.setNext(insertElement);
			tail = insertElement;
		}
	}

	private Element<E> getElementAtIndex(int index) {
		Element<E> next = head;

		for (int i = 0; i < index; i++) {
			next = next.getNext();
		}

		return next;
	}
}