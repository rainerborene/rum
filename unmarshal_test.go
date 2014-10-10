package rum

import (
	"encoding/hex"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	session_bytes, _ := hex.DecodeString("04087b0849220f73657373696f6e5f6964063a0645544922253434303233343166346464383736356132343764656631363935653639366238063b005449221977617264656e2e757365722e757365722e6b6579063b00545b075b066903311a01492219793339383232576f7753724d5046653768366974063b00544922105f637372665f746f6b656e063b00464922315a6d324354466161696c4b716b4479613161785054784f3942354a384975334e7a612f74576c4f6750616f3d063b0046")

	var params map[string]interface{}

	Unmarshal(session_bytes, params)

	assert.Contains(t, string(session_bytes), "session_id")
	assert.Equal(t, params["_csrf_token"], "Zm2CTFaailKqkDya1axPTxO9B5J8Iu3Nza/tWlOgPao=")
	assert.Equal(t, params["session_id"], "4402341f4dd8765a247def1695e696b8")
	assert.Equal(t, len(params["warden.user.user.key"].([]interface{})), 2)
}
