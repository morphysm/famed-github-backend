package regex_test

import (
	"testing"

	"github.com/morphysm/famed-github-backend/pkg/regex"
	"github.com/stretchr/testify/assert"
)

const testIssue = "UID: CL-2020-06\n\n" +
	"Severity: high\n\n" +
	"Type: DoS\n\n" +
	"Affected Clients: All clients\n\n" +
	"Summary: A DoS attack that exploits an RLP ecoding error (and lack of packet size validation) that eventually causes client crash and reply with a flood of WHOAREYOU messages that are larger than the attackers message.\n\n" +
	"Links:\n\n" +
	"Reported: 2020-08-25\n\n" +
	"Fixed: 2020-10-07\n\n" +
	"Published: 2021-12-01\n\n" +
	"Bounty Hunter: Test Hunter\n\n" +
	"Bounty Points: 5000\n\n" +
	"**Test Bold**"

func TestFindRightOfKey_Valid(t *testing.T) {
	t.Parallel()

	value, err := regex.FindRightOfKey(testIssue, "Bounty Hunter:")
	assert.NoError(t, err)
	assert.Equal(t, "Test Hunter", value)
}

func TestFindRightOfKey_ValueNotFound(t *testing.T) {
	t.Parallel()

	value, err := regex.FindRightOfKey(testIssue, "Links:")
	assert.Error(t, err)
	assert.Equal(t, "", value)
}

func TestFindRightOfKey_KeyNotFound(t *testing.T) {
	t.Parallel()

	value, err := regex.FindRightOfKey(testIssue, "Unknown:")
	assert.Error(t, err)
	assert.Equal(t, "", value)
}
