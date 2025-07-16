package compute

import "testing"

func TestOrderTotal(t *testing.T) {
	cases := []struct {
		name string
		p    float64
		q    int
		want float64
	}{
		{"normal", 1000, 3, 3000},
		{"zero qty", 500, 0, 0},
		{"zero price", 0, 10, 0},
		{"large", 1_999.99, 4, 7_999.96},
	}
	for _, c := range cases {
		got := OrderTotal(c.p, c.q)
		if got != c.want {
			t.Fatalf("%s: got %.2f want %.2f", c.name, got, c.want)
		}
	}
}
