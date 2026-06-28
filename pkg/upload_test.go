package pkg

import "testing"

func TestUploadConfigCheckPolicy(t *testing.T) {
	for _, policy := range []string{"", "skip", "overwrite"} {
		cfg := UploadConfig{Num: 1, Policy: policy}
		if err := cfg.Check(); err != nil {
			t.Fatalf("policy %q failed: %v", policy, err)
		}
	}
	cfg := UploadConfig{Num: 1, Policy: "replace"}
	if err := cfg.Check(); err == nil {
		t.Fatal("expected invalid policy to fail")
	}
}
