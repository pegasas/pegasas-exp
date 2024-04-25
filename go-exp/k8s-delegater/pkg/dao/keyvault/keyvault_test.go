package keyvault

import (
	"context"
	"testing"
)

func TestAdd(t *testing.T) {
	secretRaw := "secret"
	kv, _ := NewLocalKeyVaultClient(secretRaw)
	if secret, err := kv.GetSecret(context.Background(), "asd"); secret != secretRaw || err != nil {
		t.Errorf("kv.GetSecret ut error %v %v", secretRaw, secret)
	}
}
