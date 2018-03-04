package builder

import (
	"testing"
)

func TestSecretNames(t *testing.T) {
	original := "wordpress-wordpress-mysql-secret"
	shadow := ShadowSecretName(original)

	bound := BoundSecretName(shadow)

	if bound != original {
		t.Fatalf("expected %q got %q", original, bound)
	}
}
