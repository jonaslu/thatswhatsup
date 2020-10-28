package hashTable;

import singleLinkedList.SingleLinkedList;

public class HashTable<Key, Value> implements IHashTable<Key, Value> {
	private static final double loadFactor = 2d / 3;

	private int tableSize;

	private SingleLinkedList<Element<Key, Value>>[] table;
	private int items = 0;

	private static class Element<InnerKey, InnerValue> {
		private InnerKey key;
		private InnerValue value;

		public Element(InnerKey key) {
			this.key = key;
			this.value = null;
		}

		public Element(InnerKey key, InnerValue value) {
			this.key = key;
			this.value = value;
		}

		public InnerValue getValue() {
			return value;
		}

		public InnerKey getKey() {
			return key;
		}

		@Override
		public int hashCode() {
			return key.hashCode();
		}

		@Override
		public boolean equals(Object obj) {
			if (obj != null && obj instanceof Element) {
				Element<InnerKey, InnerValue> object = (Element<InnerKey, InnerValue>) obj;
				return this.key.equals(object.key);
			}

			return false;
		}
	}

	public HashTable(int tableSize) {
		this.tableSize = tableSize;
		table = new SingleLinkedList[tableSize];
	}

	@Override
	public Value put(Key key, Value value) {
		if ((items + 1 / tableSize) > loadFactor) {
			resizeTable();
		}

		Value oldValue = null;
		SingleLinkedList<Element<Key, Value>> currentBucket = table[getTableIndex(key)];

		if (currentBucket == null) {
			currentBucket = new SingleLinkedList<>();
			table[getTableIndex(key)] = currentBucket;
		} else {
			int indexOfOldKey = currentBucket.getIndexOf(new Element<Key, Value>(key));

			if (indexOfOldKey > -1) {
				Element<Key, Value> oldElement = currentBucket.get(indexOfOldKey);
				oldValue = oldElement.getValue();

				currentBucket.remove(indexOfOldKey);
			}
		}

		currentBucket.add(new Element<Key, Value>(key, value));

		items++;
		return oldValue;
	}

	private int getTableIndex(Key key) {
		return getTableIndex(key, this.tableSize);
	}

	private int getTableIndex(Key key, int hashBucketSize) {
		int hashCode = key.hashCode();
		hashCode = hashCode < 0 ? -hashCode : hashCode;
		int tableIndex = hashCode % hashBucketSize;
		return tableIndex;
	}

	private void resizeTable() {
		SingleLinkedList<Element<Key, Value>>[] newTable = new SingleLinkedList[tableSize * 2];

		for (int i = 0; i < table.length; i++) {
			SingleLinkedList<Element<Key, Value>> currentTableEntryList = table[i];

			if (currentTableEntryList == null) {
				continue;
			}

			for (int j = 0; j < currentTableEntryList.getLength(); j++) {
				Element<Key, Value> element = currentTableEntryList.get(j);

				int newTableIndex = getTableIndex(element.getKey(), (tableSize * 2));
				SingleLinkedList<Element<Key, Value>> newTableEntryList = newTable[newTableIndex];

				if (newTableEntryList == null) {
					newTableEntryList = new SingleLinkedList<>();
					newTable[newTableIndex] = newTableEntryList;
				}

				newTableEntryList.add(element);
			}
		}

		tableSize = tableSize * 2;
		table = newTable;
	}

	// @Override
	// public String toString() {
	// return "Size: " + table.length + "\n" + Arrays.stream(table).filter(e ->
	// e != null)
	// .map(e -> e != null ? "Key: " + e.key + " Value:" + e.value :
	// "").collect(Collectors.toList());
	// }

	@Override
	public Value get(Key key) {
		Element<Key, Value> elementForKey = getElementForKey(key);
		return (Value) elementForKey != null ? elementForKey.getValue() : null;
	}

	private Element<Key, Value> getElementForKey(Key key) {
		SingleLinkedList<Element<Key, Value>> elementList = table[getTableIndex(key)];

		if (elementList != null) {
			int indexOfValue = elementList.getIndexOf(new Element<Key, Value>(key));

			if (indexOfValue > -1) {
				return elementList.get(indexOfValue);
			}
		}

		return null;
	}
}
