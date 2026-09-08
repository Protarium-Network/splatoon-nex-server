package globals

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"os"
	"strconv"

	"github.com/PretendoNetwork/nex-go/v2"
	"github.com/PretendoNetwork/nex-go/v2/types"
	nexprotocolsglobals "github.com/PretendoNetwork/nex-protocols-go/v2/globals"
)

// PasswordFromPID derives the NEX password directly from the PID using the
// same shared secret (PN_NEX_PASSWORD_SECRET) as the wsc-account proxy that
// issues the NEX token, instead of looking the account up through the
// account gRPC service. This lets self-hosted Splatoon authenticate PIDs
// that were verified against the real Pretendo Network account service
// (no local account registration required) while still deriving a password
// only this server and the token issuer can compute.
func PasswordFromPID(pid types.PID) (string, uint32) {
	secret, err := hex.DecodeString(os.Getenv("PN_NEX_PASSWORD_SECRET"))
	if err != nil || len(secret) < 32 {
		Logger.Error("PN_NEX_PASSWORD_SECRET must contain at least 32 bytes encoded as hexadecimal")
		return "", nex.ResultCodes.RendezVous.InvalidUsername
	}

	pidBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(pidBytes, uint64(pid))
	mac := hmac.New(sha256.New, secret)
	_, _ = mac.Write(pidBytes)

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), 0
}

// This is the same format as nex-viewer's settings.json
type jsonAccount struct {
	Platform string  `json:"platform"`
	Username string  `json:"username"`
	Pid      float64 `json:"pid"`
	Password string  `json:"password"`
}

type settingsJson struct {
	Accounts []jsonAccount `json:"accounts"`
}

// PasswordFromPIDLocal is an alternative NEX password validator that can be used offline
func PasswordFromPIDLocal(pid types.PID) (string, uint32) {
	file, err := os.ReadFile("settings.json")
	if err != nil {
		Logger.Error(err.Error())
		return "", nex.ResultCodes.RendezVous.InvalidUsername
	}

	var data *settingsJson
	err = json.Unmarshal(file, &data)
	if err != nil {
		Logger.Error(err.Error())
		return "", nex.ResultCodes.RendezVous.InvalidUsername
	}

	for _, account := range data.Accounts {
		if account.Username == strconv.FormatUint(uint64(pid), 10) {
			nexprotocolsglobals.Logger.Infof("Using local account details for %v", account.Username)
			return account.Password, 0
		}
	}

	return "", nex.ResultCodes.RendezVous.InvalidUsername
}
