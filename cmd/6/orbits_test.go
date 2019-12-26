package six

import (
	"bytes"
	"reflect"
	"testing"
)

func TestTree(t *testing.T) {
	testInput := `COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
`

	buf := bytes.NewBufferString(testInput)
	tree := createTree(buf)
	for i := 0; i < 100; i++ {
		count := GetOrbitCount(tree)
		if count != 42 {
			t.Fatalf("attempt %d expected count == 42 but was %d", i, count)
		}
	}
}

func TestTreePath(t *testing.T) {
	testInput := `COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
COM)B
B)C
C)D
D)E
E)F
B)G
G)H
D)I
E)J
J)K
K)L
K)YOU
I)SAN
`
	buf := bytes.NewBufferString(testInput)
	tree := createTree(buf)

	path := tree.FindPath("YOU", "SAN")
	if !reflect.DeepEqual(path, []string{"YOU", "K", "J", "E", "D", "I", "SAN"}) {
		t.Fatalf("path: %+v", path)
	}

	orbittransfers, path := GetOrbitalTransfers(tree, "YOU", "SAN")
	if orbittransfers != 4 || reflect.DeepEqual(path, []string{"K", "J", "E", "D", "I"}) {
		t.Fatalf("orbittransfers: %d, path: %+v", orbittransfers, path)
	}
}
