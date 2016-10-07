package type_ref

import "testing"

func TestPerson_GetName(t *testing.T) {
	p := &Person{Name: "oinume"}
	println(p.GetName())
}
