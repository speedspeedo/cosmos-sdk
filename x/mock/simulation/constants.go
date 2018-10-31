package simulation

const (
	// Fraction of double-signing evidence from a past height
	pastEvidenceFraction float64 = 0.5

	// Minimum time per block
	minTimePerBlock int64 = 1000 / 2

	// Maximum time per block
	maxTimePerBlock int64 = 1000

	// Number of keys
	numKeys int = 250

	// Chance that double-signing evidence is found on a given block
	evidenceFraction float64 = 0.5

	// TODO Remove in favor of binary search for invariant violation
	onOperation bool = false
)

var (
	// Currently there are 3 different liveness types, fully online, spotty connection, offline.
	initialLivenessWeightings   = []int{40, 5, 5}
	livenessTransitionMatrix, _ = CreateTransitionMatrix([][]int{
		{90, 20, 1},
		{10, 50, 5},
		{0, 10, 1000},
	})
	// 3 states: rand in range [0, 4*provided blocksize], rand in range [0, 2 * provided blocksize], 0
	blockSizeTransitionMatrix, _ = CreateTransitionMatrix([][]int{
		{85, 5, 0},
		{15, 92, 1},
		{0, 3, 99},
	})
)
