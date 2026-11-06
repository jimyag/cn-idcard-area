package main

import "testing"

func TestLookup(t *testing.T) {
	for prefix, want := range map[string]string{
		"440305": "广东省 / 深圳市 / 南山区",
		"110103": "北京市 / 崇文区",
		"130103": "河北省 / 石家庄市 / 桥东区",
		"120110": "天津市 / 东丽区",
		"419001": "河南省 / 济源市",
		"441900": "广东省 / 东莞市",
		"440000": "广东省",
		"999999": "未知",
	} {
		t.Run(prefix, func(t *testing.T) {
			got, err := lookup(prefix)
			if err != nil {
				t.Fatal(err)
			}
			if got != want {
				t.Errorf("lookup(%q) = %q, want %q", prefix, got, want)
			}
		})
	}
}

func TestValidPrefix(t *testing.T) {
	for prefix, want := range map[string]bool{
		"440305": true,
		"000000": true,
		"44030":  false,
		"44030x": false,
		"４４":     false,
	} {
		if got := validPrefix(prefix); got != want {
			t.Errorf("validPrefix(%q) = %t, want %t", prefix, got, want)
		}
	}
}
