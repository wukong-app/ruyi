package ruyi_test

import (
	"fmt"
	"testing"

	"github.com/wukong-app/ruyi"
)

func TestRuyiExpandAndShrink(t *testing.T) {
	ry, err := ruyi.New()
	if err != nil {
		t.Fatalf("Failed to create Ruyi: %v", err)
	}
	
	fmt.Printf("Ruyi description is %v \n", ry.GetDescription())
	fmt.Printf("Ruyi size is %v \n", ry.GetSize())

	_, _ = ry.Expand()
	fmt.Printf("Ruyi expanded size is %v \n", ry.GetSize())

	_, _ = ry.Expand()
	fmt.Printf("Ruyi expanded size is %v \n", ry.GetSize())

	_, _ = ry.Shrink()
	fmt.Printf("Ruyi shrunk size is %v \n", ry.GetSize())
}
