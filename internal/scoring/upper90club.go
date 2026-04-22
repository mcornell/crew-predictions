package scoring

func Upper90Club(result Result, prediction Prediction, targetTeamIsHome bool) int {
	pts := 0
	if outcome(result.Home, result.Away) == outcome(prediction.Home, prediction.Away) {
		pts++
	}
	actualTarget, predictedTarget := result.Away, prediction.Away
	actualOpponent, predictedOpponent := result.Home, prediction.Home
	if targetTeamIsHome {
		actualTarget, predictedTarget = result.Home, prediction.Home
		actualOpponent, predictedOpponent = result.Away, prediction.Away
	}
	if predictedTarget == actualTarget {
		pts++
	}
	if predictedOpponent == actualOpponent {
		pts++
	}
	return pts
}
