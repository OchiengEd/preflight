package shell

import (
	"os/exec"
	"strings"

	"github.com/komish/preflight/certification"
)

type BasedOnUbiPolicy struct{}

func (p *BasedOnUbiPolicy) Validate(image string) (bool, []byte, error) {
	stdouterr, err := exec.Command("podman", "run", "--rm", "-it", image, "cat", "/etc/os-release").CombinedOutput()
	if err != nil {
		return false, stdouterr, err
	}

	lines := strings.Split(string(stdouterr), "\n")

	var hasRHELID, hasRHELName bool
	for _, value := range lines {
		if strings.HasPrefix(value, `ID="rhel"`) {
			hasRHELID = true
		} else if strings.HasPrefix(value, `NAME="Red Hat Enterprise Linux"`) {
			hasRHELName = true
		}
	}
	if hasRHELID && hasRHELName {
		return true, []byte{}, nil
	}

	return false, []byte{}, nil
}

func (p *BasedOnUbiPolicy) Name() string {
	return "BasedOnUbi"
}

func (p *BasedOnUbiPolicy) Metadata() certification.Metadata {
	return certification.Metadata{
		Description:      "Checking if the container's base image is based on UBI",
		Level:            "best",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide", // Placeholder
		PolicyURL:        "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p *BasedOnUbiPolicy) Help() certification.HelpText {
	return certification.HelpText{
		Message:    "It is recommened that your image be based upon the Red Hat Universal Base Image (UBI)",
		Suggestion: "Change the FROM directive in your Dockerfile or Containerfile to FROM registry.access.redhat.com/ubi8/ubi",
	}
}
