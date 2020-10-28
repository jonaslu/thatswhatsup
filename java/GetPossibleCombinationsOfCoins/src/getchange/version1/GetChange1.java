package getchange.version1;

import java.util.ArrayList;
import java.util.List;

public class GetChange1 {

	private CellList cellList;
	private int total;

	public GetChange1(List<Integer> denominators, int total) {
		this.total = total;
		cellList = new CellList(denominators, total);
	}

	public List<List<Integer>> getPossibleChangeCombinations() {
		List<List<Integer>> possibleCombinationsList = new ArrayList<>();

		while (cellList.hasNext()) {
			List<Integer> currentIteration = cellList.next();
			if (currentIteration.stream().mapToInt(i -> i).sum() == total) {
				possibleCombinationsList.add(currentIteration);
			}
		}

		return possibleCombinationsList;
	}
}
