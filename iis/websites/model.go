package websites

type AuthenticationMode string

var (
	None      AuthenticationMode = "None"
	Federated AuthenticationMode = "Federated"
	Forms     AuthenticationMode = "Forms"
	Passport  AuthenticationMode = "Passport"
	Windows   AuthenticationMode = "Windows"
)

type Website struct {
	Name                         string
	ApplicationPool              string
	PhysicalPath                 string
	State                        string
	StartsOnBoot                 bool
	MaxBandwidthPerSecondInBytes int64
	// TODO: can we return the auth mode too?
}
