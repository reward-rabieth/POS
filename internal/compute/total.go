package compute

// OrderTotal returns price * quantity.
func OrderTotal(price float64, qty int) float64 {
	return price * float64(qty)
}
