package commands

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/chef/automate/lib/io/fileutils"
	"github.com/chef/automate/lib/platform/secrets"
)

func init() {
	RootCmd.AddCommand(initCmd)
}

const (
	// Location of secrets_key generated by compliance service in
	// previous A2 versions. This key is migrated to
	// secrets-service if it exists.
	complianceKeyPath = "/hab/svc/compliance-service/data/secrets_key"
	// Location of secrets_key generated by older versions of the
	// secrets service.
	secretsServiceKeyPath = "/hab/svc/secrets-service/data/secrets_key"
	// Number of bytes of random data to use to generate the key
	defaultKeyLength = 16
)

var initCmd = &cobra.Command{
	Use:   "init KEYPATH",
	Short: "Initializes the shared secret used for symmetric encryption. Migrated from compliance if needed",
	Run: func(_ *cobra.Command, args []string) {
		platformSecrets, err := secrets.NewDefaultSecretStore()
		if err != nil {
			fatalf("failed to initialize platform secrets store: %s", err.Error())
		}

		keyExists, err := platformSecrets.Exists(secrets.SecretsServiceKeyName)
		if err != nil {
			fatalf("failed to determine if key already exists: %s", err.Error())
		}

		if keyExists {
			return
		}

		err = generateOrMigrateKey(platformSecrets)
		if err != nil {
			fatalf(err.Error())
		}
	},
}

func generateOrMigrateKey(store secrets.SecretStore) error {
	oldSecretsServiceKeyExists, err := fileutils.PathExists(secretsServiceKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to check for old secrets-service key")
	}
	if oldSecretsServiceKeyExists {
		return keyFromPath(store, secretsServiceKeyPath)
	}

	complianceKeyExists, err := fileutils.PathExists(complianceKeyPath)
	if err != nil {
		return errors.Wrap(err, "failed to check for compliance key")
	}
	if complianceKeyExists {
		return keyFromPath(store, complianceKeyPath)
	}

	fmt.Println("Generating new secret key")
	key, err := generateKey()
	if err != nil {
		return errors.Wrap(err, "could not generate key")
	}
	return store.SetSecret(secrets.SecretsServiceKeyName, key)
}

func keyFromPath(store secrets.SecretStore, path string) error {
	fmt.Printf("Migrating secret key from %s\n", path)
	key, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "failed to read key from %s", path)
	}

	return store.SetSecret(secrets.SecretsServiceKeyName, key)
}

func generateKey() ([]byte, error) {
	randomBytes := make([]byte, defaultKeyLength)
	encodedBytes := make([]byte, defaultKeyLength*2)
	if _, err := rand.Read(randomBytes); err != nil {
		return nil, errors.Wrap(err, "failed to read random data while generating random bytes")
	}

	hex.Encode(encodedBytes, randomBytes)
	return encodedBytes, nil
}

func fatalf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg, args...)
	os.Exit(1)
}
