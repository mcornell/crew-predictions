package scoring

type Result struct{ Home, Away int }
type Prediction struct{ Home, Away int }

func AcesRadio(result Result, prediction Prediction) int {
	if prediction == Prediction(result) {
		return 15
	}
	return 0
}
