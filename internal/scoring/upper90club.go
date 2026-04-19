package scoring

func Upper90Club(result Result, prediction Prediction, targetTeamIsHome bool) int {
	pts := 0
	if outcome(result.Home, result.Away) == outcome(prediction.Home, prediction.Away) {
		pts++
	}
	actualTargetGoals, predictedTargetGoals := result.Away, prediction.Away
	if targetTeamIsHome {
		actualTargetGoals, predictedTargetGoals = result.Home, prediction.Home
	}
	if predictedTargetGoals == actualTargetGoals {
		pts++
	}
	return pts
}
