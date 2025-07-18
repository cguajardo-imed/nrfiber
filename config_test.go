package nrfiber

import (
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestConfigNoticeErrorEnabled(t *testing.T) {
	cfg := ConfigNoticeErrorEnabled(true)
	if cfg.key != configKeyNoticeErrorEnabled {
		t.Errorf("expected key %s, got %s", configKeyNoticeErrorEnabled, cfg.key)
	}
	if cfg.value != true {
		t.Errorf("expected value true, got %v", cfg.value)
	}
}

func TestConfigStatusCodeIgnored(t *testing.T) {
	codes := []int{400, 404, 500}
	cfg := ConfigStatusCodeIgnored(codes)
	if cfg.key != configKeyStatusCodeIgnored {
		t.Errorf("expected key %s, got %s", configKeyStatusCodeIgnored, cfg.key)
	}
	if !reflect.DeepEqual(cfg.value, codes) {
		t.Errorf("expected value %v, got %v", codes, cfg.value)
	}
}

func TestConfigCustomTransactionNameFunc(t *testing.T) {
	customFunc := func(c *fiber.Ctx) string { return "test" }
	cfg := ConfigCustomTransactionNameFunc(customFunc)
	if cfg.key != configCustomTransactionNameFunc {
		t.Errorf("expected key %s, got %s", configCustomTransactionNameFunc, cfg.key)
	}
}

func TestCreateConfigMap(t *testing.T) {
	cfg1 := ConfigNoticeErrorEnabled(true)
	cfg2 := ConfigStatusCodeIgnored([]int{404})

	configMap := createConfigMap(cfg1, cfg2)

	if len(configMap) != 2 {
		t.Errorf("expected map length 2, got %d", len(configMap))
	}

	if v, ok := configMap[configKeyNoticeErrorEnabled]; !ok || v != true {
		t.Errorf("expected NoticeErrorEnabled true, got %v", v)
	}
}

func TestNoticeErrorEnabled(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]any
		expected bool
	}{
		{"enabled", map[string]any{configKeyNoticeErrorEnabled: true}, true},
		{"disabled", map[string]any{configKeyNoticeErrorEnabled: false}, false},
		{"missing", map[string]any{}, false},
		{"wrong type", map[string]any{configKeyNoticeErrorEnabled: "true"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := noticeErrorEnabled(tt.config); got != tt.expected {
				t.Errorf("noticeErrorEnabled() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStatusCodeIgnored(t *testing.T) {
	tests := []struct {
		name     string
		config   map[string]any
		expected []int
	}{
		{"valid codes", map[string]any{configKeyStatusCodeIgnored: []int{404, 500}}, []int{404, 500}},
		{"empty codes", map[string]any{configKeyStatusCodeIgnored: []int{}}, []int{}},
		{"missing", map[string]any{}, []int{}},
		{"wrong type", map[string]any{configKeyStatusCodeIgnored: "404"}, []int{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statusCodeIgnored(tt.config); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("statusCodeIgnored() = %v, want %v", got, tt.expected)
			}
		})
	}
}
