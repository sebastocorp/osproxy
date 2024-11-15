package modifiers

import (
	"osproxy/api/v1alpha5"
	"strings"

	"github.com/aws/smithy-go/transport/http"
)

var (
	modifiers = map[string]ModifierT{
		"path":   PathModifier,
		"header": HeaderModifier,
	}
)

type ModifierT func(*http.Request, v1alpha5.ProxyModifierConfigT)

func PathModifier(r *http.Request, mod v1alpha5.ProxyModifierConfigT) {
	r.URL.Path = mod.Path.AddPrefix + strings.TrimPrefix(r.URL.Path, mod.Path.RemovePrefix)
}

func HeaderModifier(r *http.Request, mod v1alpha5.ProxyModifierConfigT) {
	r.Header.Set(mod.Header.Name, mod.Header.Value)
	if mod.Header.Remove {
		r.Header.Del(mod.Header.Name)
	}
}
