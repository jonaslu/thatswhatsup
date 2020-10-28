package getchange.version2;

import java.util.ArrayList;
import java.util.List;

public class GetChange2 {

	public static int getNumberOfWays(List<Integer> denominators, int sum) {
		if (denominators.isEmpty()) {
			return 0;
		}

		if (sum < 0) {
			return 0;
		}

		if (sum == 0) {
			return 1;
		}

		int head = denominators.get(0);
		List<Integer> tail = new ArrayList<>(denominators.subList(1, denominators.size()));

		return getNumberOfWays(tail, sum) + getNumberOfWays(new ArrayList<>(denominators), sum - head);
	}
}
