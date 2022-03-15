package azure

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	azureTagNameToPrometheusNameRegExp = regexp.MustCompile("[^_a-zA-Z0-9]")
)

func AddResourceTagsToPrometheusLabels(labels prometheus.Labels, resourceTags interface{}, tags []string) prometheus.Labels {
	resourceTagMap := map[string]string{}

	switch v := resourceTags.(type) {
	case map[string]*string:
		resourceTagMap = to.StringMap(v)
	case map[string]string:
		resourceTagMap = v
	}

	// normalize
	resourceTagMap = normalizeTags(resourceTagMap)

	for _, tag := range tags {
		tagParts := strings.SplitN(tag, "?", 2)
		tag = tagParts[0]

		tagSettings := ""
		if len(tagParts) == 2 {
			tagSettings = tagParts[1]
		}

		tag = strings.ToLower(tag)
		tagLabel := fmt.Sprintf("tag_" + azureTagNameToPrometheusNameRegExp.ReplaceAllLiteralString(tag, "_"))
		labels[tagLabel] = ""

		if val, exists := resourceTagMap[tag]; exists {
			if tagSettings != "" {
				val = applyTagValueSettings(val, tagSettings)
			}
			labels[tagLabel] = val
		}
	}

	return labels
}

func applyTagValueSettings(val string, settings string) string {
	ret := val
	settingQuery, _ := url.ParseQuery(settings)

	if settingQuery.Has("toLower") {
		ret = strings.ToLower(ret)
	}

	if settingQuery.Has("toUpper") {
		ret = strings.ToUpper(ret)
	}

	return ret
}

func normalizeTags(tags map[string]string) map[string]string {
	ret := map[string]string{}

	for tagName, tagValue := range tags {
		tagName = strings.ToLower(tagName)
		ret[tagName] = strings.TrimSpace(tagValue)
	}

	return ret
}