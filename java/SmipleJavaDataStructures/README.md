# SmipleJavaDataStructures
Test of some datastructures in Java built in eclipse

Just enter Eclipse - hit Import from the File menu and import this finely cloned directory as General > Existing projects into workspace

Set to Java 1.8 FTW

## Single linked list with tests
SingleLinkedList.java implements a single linked list conforming to the the IList interface that contains the most basic list operations
```java
add(stuff)
get(index)
remove(index)
getLength()
getIndexOf(value)
replace(index, value)
```

## HashTable (or Map if you prefer it like that)
> Gimme a real hashtable... no wait, make it a double. Its been a rough day.

Ah, the good 'ole hash table map. Supports
```java
value put(key, value)
value get(key)
```

Put returns any old value associated with the key, and get.. well, you get it?

Uses the singly linked list done above as its internal storage for when **hashes collides**

## DynamicArray
An array that increases if the index is out of bounds

Supports
```java
set(index, value)
value get(index)
```

Acts up and throws ArrayIndexOutOfBounds if you don't watch its personal boundaries