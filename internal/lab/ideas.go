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
		Name:        "Example Project",
		Description: "A placeholder idea. Replace this with your first real concept.",
		Status:      StatusExploring,
		Links: []Link{
			{Label: "Figma Mock", URL: "#"},
			{Label: "Tech Spec", URL: "#"},
		},
	},
}
