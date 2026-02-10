package drive

import (
	"testing"
)

func TestParseShareURL(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantCode   string
		wantAccess string
		wantErr    bool
	}{
		{
			name:       "完整URL无访问码",
			input:      "https://cloud.189.cn/t/reiQruJJnyEv",
			wantCode:   "reiQruJJnyEv",
			wantAccess: "",
		},
		{
			name:       "完整URL有访问码",
			input:      "https://cloud.189.cn/t/reiQruJJnyEv（访问码：v5oy）",
			wantCode:   "reiQruJJnyEv（访问码：v5oy）",
			wantAccess: "v5oy",
		},
		{
			name:       "纯shareCode",
			input:      "reiQruJJnyEv",
			wantCode:   "reiQruJJnyEv",
			wantAccess: "",
		},
		{
			name:       "纯shareCode有访问码",
			input:      "reiQruJJnyEv（访问码：v5oy）",
			wantCode:   "reiQruJJnyEv（访问码：v5oy）",
			wantAccess: "v5oy",
		},
		{
			name:    "空字符串",
			input:   "",
			wantErr: true,
		},
		{
			name:       "URL带尾部斜杠",
			input:      "https://cloud.189.cn/t/reiQruJJnyEv/",
			wantCode:   "reiQruJJnyEv",
			wantAccess: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, access, err := ParseShareURL(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if code != tt.wantCode {
				t.Errorf("shareCode = %q, want %q", code, tt.wantCode)
			}
			if access != tt.wantAccess {
				t.Errorf("accessCode = %q, want %q", access, tt.wantAccess)
			}
		})
	}
}
