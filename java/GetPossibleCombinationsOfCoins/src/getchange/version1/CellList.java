package getchange.version1;

import java.util.ArrayList;
import java.util.Iterator;
import java.util.List;
import java.util.NoSuchElementException;
import java.util.stream.Collectors;

public class CellList implements Iterator<List<Integer>> {

	private int total;
	private List<Cell> cells = new ArrayList<>();

	// denominations in ascending order e g [1,5,10]
	public CellList(List<Integer> denominations, int total) {
		this.total = total;
		constructCells(denominations);
	}

	private void constructCells(List<Integer> denominations) {
		for (Integer denomination : denominations) {
			if (total >= denomination) {
				cells.add(new Cell(total, denomination));
			}
		}
	}

	@Override
	public boolean hasNext() {
		boolean returnValue = cells.stream().anyMatch(cell -> cell.hasNext());
		return returnValue;
	}

	@Override
	public List<Integer> next() {
		if (!hasNext()) {
			throw new NoSuchElementException();
		}

		Iterator<Cell> cellIterator = cells.iterator();
		Cell currentCell = cellIterator.next();

		while (!currentCell.hasNext()) {
			currentCell.reset();
			currentCell = cellIterator.next();
		}

		currentCell.next();
		return cells.stream().map(Cell::getCurrentValue).collect(Collectors.toList());
	}
}
