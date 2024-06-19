package keyring

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/zalando/go-keyring"
)

const (
	windowsMaxLength = 2000
	darwinMaxLength  = 2000
	serviceName      = "evermarkable"
)

func GetSecretFromStore(secretName string) (string, error) {
	secret, err := keyring.Get(serviceName, secretName)
	if err != nil {
		return "", err
	}

	if secretNames, err := decodeSecretNames(secretName, secret); err == nil {
		var res strings.Builder
		for _, secretName := range secretNames {
			if tmpGet, err := keyring.Get(serviceName, secretName); err != nil {
				return "", fmt.Errorf("error retrieving partial key: %v", err)
			} else {
				res.WriteString(tmpGet)
			}
		}
		return res.String(), nil
	}
	return secret, err
}

func SaveSecretInStore(secretName string, secret string) error {
	var maxLen int
	switch os := runtime.GOOS; os {
	case "darwin":
		maxLen = darwinMaxLength
	case "windows":
		maxLen = windowsMaxLength
	default:
		maxLen = len(secret)
	}

	if len(secret) > maxLen {
		totalLen := len(secret)
		count, charsSent := 0, 0
		for charsSent < totalLen {
			sendingChars := minOf(maxLen, totalLen-charsSent)
			secTmp := secret[charsSent : charsSent+sendingChars]
			if err := keyring.Set(serviceName, getDividedSecretName(secretName, count), secTmp); err != nil {
				return fmt.Errorf("error saving out partial key: %v", err)
			}
			count++
			charsSent += sendingChars
		}
		return keyring.Set(serviceName, secretName, encodesecretNames(secretName, count))
	}
	return keyring.Set(serviceName, secretName, secret)
}

func DeleteSecretFromStore(secretName string) error {
	if secretNames, err := decodeSecretNames(secretName, ""); err == nil {
		for _, secretName := range secretNames {
			if err := keyring.Delete(serviceName, secretName); err != nil {
				return fmt.Errorf("error deleting partial key: %v", err)
			}
		}
		return nil
	}

	return keyring.Delete(serviceName, secretName)
}

func ErrorIsNotFound(err error) bool {
	return errors.Is(err, keyring.ErrNotFound)
}

func minOf(vars ...int) int {
	min := vars[0]
	for _, i := range vars {
		if min > i {
			min = i
		}
	}
	return min
}

func decodeSecretNames(name string, value string) ([]string, error) {
	if !strings.Contains(value, fmt.Sprintf("%s-0:", name)) {
		return nil, fmt.Errorf("name is not contained in value")
	}
	return strings.Split(value, ":"), nil
}

func encodesecretNames(name string, count int) string {
	outputBuilder := make([]string, count)
	for i := 0; i < count; i++ {
		outputBuilder[i] = fmt.Sprintf("%s-%d", name, i)
	}
	return strings.Join(outputBuilder, ":")
}

func getDividedSecretName(name string, count int) string {
	return fmt.Sprintf("%s-%d", name, count)
}
