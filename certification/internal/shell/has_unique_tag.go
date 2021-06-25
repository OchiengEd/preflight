package shell

import (
	"github.com/komish/preflight/certification"
	"github.com/komish/preflight/certification/errors"
)

type HasUniqueTagPolicy struct{}

func (p *HasUniqueTagPolicy) Validate(image string) (bool, []byte, error) {
	return false, []byte{}, errors.ErrFeatureNotImplemented
}
func (p *HasUniqueTagPolicy) Name() string {
	return "HasUniqueTag"
}
func (p *HasUniqueTagPolicy) Metadata() certification.Metadata {
	return certification.Metadata{
		Description:      "Checking if container has a tag other than 'latest'.",
		Level:            "best",
		KnowledgeBaseURL: "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
		PolicyURL:        "https://connect.redhat.com/zones/containers/container-certification-policy-guide",
	}
}

func (p *HasUniqueTagPolicy) Help() certification.HelpText {
	return certification.HelpText{
		Message:    "Containers should have a tag other than latest, so that the image can be uniquely identfied.",
		Suggestion: "Add a tag to your image. Consider using Semantic Versioning. https://semver.org/",
	}
}
