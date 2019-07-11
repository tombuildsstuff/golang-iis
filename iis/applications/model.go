package applications

// IIS Application Configuration Data
type Application struct {
	Name            string
	Path            string
	ApplicationPool string
	PhysicalPath    string
}
