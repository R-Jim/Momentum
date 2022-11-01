package carrier

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_JetIDs_removeJetID(t *testing.T) {
	jetIDs := JetIDs{"jet1", "jet2", "jet3", "jet4"}

	jetIDs = jetIDs.removeJetID("jet2")

	require.Equal(t, len(jetIDs), 3)
	require.Equal(t, jetIDs, JetIDs{"jet1", "jet4", "jet3"}) // due to reverse
}
