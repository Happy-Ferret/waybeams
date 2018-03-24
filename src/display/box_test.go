package display

import (
	"testing"
)

func TestBox(t *testing.T) {
}

func TestVBox(t *testing.T) {

	t.Run("Simple Children", func(t *testing.T) {
		root, _ := VBox(NewBuilder(), Height(100), Children(func(b Builder) {
			Box(b, FlexHeight(1))
			Box(b, FlexHeight(1))
		}))
		one := root.GetChildAt(0)
		//two := root.GetChildAt(1)
		if one.GetHeight() != 50 {
			t.Errorf("Expected 50, but was %v", one.GetHeight())
		}
	})
}
