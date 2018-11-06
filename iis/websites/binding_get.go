package websites

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Binding struct {
	IPAddress  string
	Port       int
	DomainName string
	Protocol   string
}

type getBindingResponse struct {
	BindingInformation string `json="bindingInformation"`
	Protocol           string `json="protocol"`
}

func (c *WebsitesClient) GetBindings(name string) (*[]Binding, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
$bindings = Get-WebBinding -name %q
if ($bindings.Count -gt 1) {
    $v = $bindings | ConvertTo-Json -Compress
	Write-Host $v
} else {
    $v = "[""{0}""]" -f $bindings[0].ToString()
    Write-Host $v
}
  `, name)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Bindings for Website %q: %+v", name, err)
	}

	var resp []getBindingResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &resp)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling Bindings for Website %q: %+v", name, err)
		}
	}

	bindings := make([]Binding, 0)
	for _, b := range resp {
		// ip:port:domain e.g. `*:80:mysite.com`
		segments := strings.Split(b.BindingInformation, ":")
		ip := segments[0]
		port := segments[1]
		domain := segments[2]

		portInt, err := strconv.Atoi(port)
		if err != nil {
			return nil, fmt.Errorf("Error converting %q to an int: %+v", port, err)
		}

		binding := Binding{
			IPAddress:  ip,
			DomainName: domain,
			Port:       portInt,
			Protocol:   b.Protocol,
		}
		bindings = append(bindings, binding)
	}

	return &bindings, nil
}
