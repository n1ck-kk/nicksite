package lab

type Status string

const (
	StatusExploring  Status = "Exploring"
	StatusBuilding   Status = "Building"
	StatusValidating Status = "Validating"
	StatusShipped    Status = "Shipped"
	StatusPaused     Status = "Paused"
)

type Link struct {
	Label string
	URL   string
}

type Idea struct {
	Name        string
	Description string
	Status      Status
	Links       []Link
}

// Ideas is the list of projects in the lab.
// Add new entries here and push to deploy.
var Ideas = []Idea{
	{
		Name:        "Military AR Glasses · Drone HUD",
		Description: "Tactical glasses with a live drone feed projected on the bottom-right lens. Clear AR overlay — no goggles, no separate screen. Altitude, speed, and heading overlaid on real-world view via a 5.8 GHz downlink.",
		Status:      StatusExploring,
		Links:       []Link{},
	},
}
