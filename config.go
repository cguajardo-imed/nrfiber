package nrfiber

import "github.com/gofiber/fiber/v2"

const (
	configKeyNoticeErrorEnabled     = "NoticeErrorEnabled"
	configKeyStatusCodeIgnored      = "StatusCodeIgnored"
	configCustomTransactionNameFunc = "CustomTransactionNameFunc"
)

type config struct {
	key   string
	value any
}

func ConfigNoticeErrorEnabled(enabled bool) *config {
	return &config{
		key:   configKeyNoticeErrorEnabled,
		value: enabled,
	}
}

func ConfigStatusCodeIgnored(statusCode []int) *config {
	return &config{
		key:   configKeyStatusCodeIgnored,
		value: statusCode,
	}
}

func ConfigCustomTransactionNameFunc(customTransactionNameFunc func(c *fiber.Ctx) string) *config {
	return &config{
		key:   configCustomTransactionNameFunc,
		value: customTransactionNameFunc,
	}
}

func createConfigMap(configs ...*config) map[string]any {
	configMap := make(map[string]any, len(configs))
	for _, c := range configs {
		configMap[c.key] = c.value
	}
	return configMap
}

func noticeErrorEnabled(configMap map[string]any) bool {
	if val, ok := configMap[configKeyNoticeErrorEnabled]; ok {
		if boolVal, ok := val.(bool); ok {
			return boolVal
		}
	}
	return false
}

func statusCodeIgnored(configMap map[string]any) []int {
	if val, ok := configMap[configKeyStatusCodeIgnored]; ok {
		if v, ok := val.([]int); ok {
			return v
		}
	}
	return []int{}
}

func customTransactionNameFunc(configMap map[string]any, defaultFunc func(c *fiber.Ctx) string) func(c *fiber.Ctx) string {
	if val, ok := configMap[configCustomTransactionNameFunc]; ok {
		if v, ok := val.(func(c *fiber.Ctx) string); ok {
			return v
		}
	}
	return defaultFunc
}
