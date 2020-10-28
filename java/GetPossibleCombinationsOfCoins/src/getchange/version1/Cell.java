package getchange.version1;

import java.util.Iterator;

public class Cell implements Iterator<Integer> {

	private int currentPart;
	private int denominator;
	private int total;

	public Cell(int total, int denominator) {
		this.total = total;
		currentPart = 0;
		this.denominator = denominator;
	}

	public int getCurrentValue() {
		return currentPart;
	}

	public void reset() {
		currentPart = 0;
	}

	@Override
	public boolean hasNext() {
		return currentPart <= total;
	}

	@Override
	public Integer next() {
		return currentPart += denominator;
	}
}
