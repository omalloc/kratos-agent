package agent

// HasHang 该组件是否挂起
func (x *Microservice) HasHang() bool {
	if v, ok := x.Metadata["hang"]; ok && v == "true" {
		return true
	}
	return false
}
