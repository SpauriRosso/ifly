package presence

// DefaultApplicationID identifies IFLY's public Discord application when no override is configured.
const DefaultApplicationID = "707372989921886250"

const (
	drawingDetails  = "Drawing a piece of art"
	admiringDetails = "Admiring a piece of art"
)

// activityDetails takes whether an artwork is open and returns the short English line shown on Discord.
func activityDetails(artworkOpen bool) string {
	if artworkOpen {
		return admiringDetails
	}
	return drawingDetails
}
