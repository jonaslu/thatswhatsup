package hashTable;

public interface IHashTable<K, V> {

	public V put(K key, V value);

	public V get(K key);

}
