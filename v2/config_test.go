package nrfiber

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestConfigNoticeErrorEnabled(t *testing.T) {
	config := ConfigNoticeErrorEnabled(true)
	assert.NotNil(t, config)
	assert.Equal(t, configKeyNoticeErrorEnabled, config.key)
	assert.Equal(t, true, config.value)
}

func TestConfigStatusCodeIgnored(t *testing.T) {
	statusCodes := []int{404, 500, 503}
	config := ConfigStatusCodeIgnored(statusCodes)
	assert.NotNil(t, config)
	assert.Equal(t, configKeyStatusCodeIgnored, config.key)
	assert.Equal(t, statusCodes, config.value)
}

func TestConfigCustomTransactionNameFunc(t *testing.T) {
	customFunc := func(c *fiber.Ctx) string {
		return "custom-name"
	}
	config := ConfigCustomTransactionNameFunc(customFunc)
	assert.NotNil(t, config)
	assert.Equal(t, configCustomTransactionNameFunc, config.key)
	assert.NotNil(t, config.value)
}

func TestCreateConfigMap(t *testing.T) {
	config1 := ConfigNoticeErrorEnabled(true)
	config2 := ConfigStatusCodeIgnored([]int{404})

	configMap := createConfigMap(config1, config2)

	assert.Len(t, configMap, 2)
	assert.Equal(t, true, configMap[configKeyNoticeErrorEnabled])
	assert.Equal(t, []int{404}, configMap[configKeyStatusCodeIgnored])
}

func TestCreateConfigMapEmpty(t *testing.T) {
	configMap := createConfigMap()
	assert.NotNil(t, configMap)
	assert.Len(t, configMap, 0)
}

func TestNoticeErrorEnabledDefault(t *testing.T) {
	configMap := createConfigMap()
	enabled := noticeErrorEnabled(configMap)
	assert.False(t, enabled)
}

func TestNoticeErrorEnabledTrue(t *testing.T) {
	configMap := createConfigMap(ConfigNoticeErrorEnabled(true))
	enabled := noticeErrorEnabled(configMap)
	assert.True(t, enabled)
}

func TestNoticeErrorEnabledFalse(t *testing.T) {
	configMap := createConfigMap(ConfigNoticeErrorEnabled(false))
	enabled := noticeErrorEnabled(configMap)
	assert.False(t, enabled)
}

func TestStatusCodeIgnoredDefault(t *testing.T) {
	configMap := createConfigMap()
	codes := statusCodeIgnored(configMap)
	assert.NotNil(t, codes)
	assert.Len(t, codes, 0)
}

func TestStatusCodeIgnoredWithCodes(t *testing.T) {
	expectedCodes := []int{400, 404, 500}
	configMap := createConfigMap(ConfigStatusCodeIgnored(expectedCodes))
	codes := statusCodeIgnored(configMap)
	assert.Equal(t, expectedCodes, codes)
}

func TestCustomTransactionNameFuncDefault(t *testing.T) {
	defaultFunc := func(c *fiber.Ctx) string {
		return "default"
	}
	configMap := createConfigMap()
	resultFunc := customTransactionNameFunc(configMap, defaultFunc)

	// Create a dummy context to test
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	result := resultFunc(ctx)
	assert.Equal(t, "default", result)
}

func TestCustomTransactionNameFuncCustom(t *testing.T) {
	customFunc := func(c *fiber.Ctx) string {
		return "custom"
	}
	defaultFunc := func(c *fiber.Ctx) string {
		return "default"
	}
	configMap := createConfigMap(ConfigCustomTransactionNameFunc(customFunc))
	resultFunc := customTransactionNameFunc(configMap, defaultFunc)

	// Create a dummy context to test
	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	result := resultFunc(ctx)
	assert.Equal(t, "custom", result)
}

func TestMultipleConfigs(t *testing.T) {
	customFunc := func(c *fiber.Ctx) string {
		return "custom-transaction"
	}

	config1 := ConfigNoticeErrorEnabled(true)
	config2 := ConfigStatusCodeIgnored([]int{404, 500})
	config3 := ConfigCustomTransactionNameFunc(customFunc)

	configMap := createConfigMap(config1, config2, config3)

	assert.True(t, noticeErrorEnabled(configMap))
	assert.Equal(t, []int{404, 500}, statusCodeIgnored(configMap))

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
	defer app.ReleaseCtx(ctx)

	nameFunc := customTransactionNameFunc(configMap, nil)
	assert.Equal(t, "custom-transaction", nameFunc(ctx))
}
