package scoring

func Upper90Club(result Result, prediction Prediction, columbusIsHome bool) int {
	pts := 0
	if outcome(result.Home, result.Away) == outcome(prediction.Home, prediction.Away) {
		pts++
	}
	actualColumbusGoals, predictedColumbusGoals := result.Away, prediction.Away
	if columbusIsHome {
		actualColumbusGoals, predictedColumbusGoals = result.Home, prediction.Home
	}
	if predictedColumbusGoals == actualColumbusGoals {
		pts++
	}
	return pts
}
