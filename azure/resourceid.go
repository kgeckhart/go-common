package azure

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	resourceIdRegExp = regexp.MustCompile(`(?i)/subscriptions/(?P<subscription>[^/]+)(/resourceGroups/(?P<resourceGroup>[^/]+))?(/providers/(?P<resourceProvider>[^/]*)/(?P<resourceProviderNamespace>[^/]*)/(?P<resourceName>[^/]+)(/(?P<resourceSubPath>.+))?)?`)
)

type (
	AzureResourceDetails struct {
		OriginalResourceId        string
		Subscription              string
		ResourceGroup             string
		ResourceProviderName      string
		ResourceProviderNamespace string
		ResourceType              string
		ResourceName              string
		ResourceSubPath           string
	}
)

func (resource *AzureResourceDetails) ResourceId() (resourceId string) {
	if resource.Subscription != "" {
		resourceId += fmt.Sprintf(
			"/subscriptions/%s",
			resource.Subscription,
		)
	} else {
		return
	}

	if resource.ResourceGroup != "" {
		resourceId += fmt.Sprintf(
			"/resourceGroups/%s",
			resource.ResourceGroup,
		)
	}

	if resource.ResourceProviderName != "" && resource.ResourceProviderNamespace != "" && resource.ResourceName != "" {
		resourceId += fmt.Sprintf(
			"/providers/%s/%s/%s",
			resource.ResourceProviderName,
			resource.ResourceProviderNamespace,
			resource.ResourceName,
		)

		if resource.ResourceSubPath != "" {
			resourceId += fmt.Sprintf(
				"/%s",
				resource.ResourceSubPath,
			)
		}
	}

	return
}

func ParseResourceId(resourceId string) (resource *AzureResourceDetails, err error) {
	resource = &AzureResourceDetails{}

	if matches := resourceIdRegExp.FindStringSubmatch(resourceId); len(matches) >= 1 {
		resource.OriginalResourceId = resourceId
		for i, name := range resourceIdRegExp.SubexpNames() {
			if i != 0 && name != "" {
				switch name {
				case "subscription":
					resource.Subscription = strings.ToLower(matches[i])
				case "resourceGroup":
					resource.ResourceGroup = strings.ToLower(matches[i])
				case "resourceProvider":
					resource.ResourceProviderName = strings.ToLower(matches[i])
				case "resourceProviderNamespace":
					resource.ResourceProviderNamespace = strings.ToLower(matches[i])
				case "resourceName":
					resource.ResourceName = strings.ToLower(matches[i])
				case "resourceSubPath":
					resource.ResourceSubPath = strings.Trim(matches[i], "/")
				}
			}
		}

		// build resourcetype
		if resource.ResourceProviderName != "" && resource.ResourceProviderNamespace != "" {
			resource.ResourceType = fmt.Sprintf(
				"%s/%s",
				resource.ResourceProviderName,
				resource.ResourceProviderNamespace,
			)
		}

	} else {
		err = fmt.Errorf("unable to parse Azure resourceID \"%v\"", resourceId)
	}

	return
}
