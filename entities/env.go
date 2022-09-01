package entities

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gosimple/slug"
)

const (
	EnvRootUser = "yolo"
)

type EnvSSHHostKey struct {
	Algorithm   string `json:"algorithm"`
	Fingerprint string `json:"fingerprint"`
}

type EnvStatus string

const (
	EnvStatusCreating EnvStatus = "creating"
	EnvStatusCreated  EnvStatus = "created"
	EnvStatusRemoving EnvStatus = "removing"
)

type Env struct {
	ID                       string                `json:"id"`
	Name                     string                `json:"name"`
	InfrastructureJSON       string                `json:"infrastructure_json"`
	InstanceType             string                `json:"instance_type"`
	InstancePublicIPAddress  string                `json:"instance_public_ip_address"`
	SSHHostKeys              []EnvSSHHostKey       `json:"ssh_host_keys"`
	SSHKeyPairPEMContent     string                `json:"ssh_key_pair_pem_content"`
	ResolvedRepository       ResolvedEnvRepository `json:"resolved_repository"`
	OpenedPorts              map[string]bool       `json:"opened_ports"`
	Status                   EnvStatus             `json:"status"`
	AdditionalPropertiesJSON string                `json:"additional_properties_json"`
	CreatedAtTimestamp       int64                 `json:"created_at_timestamp"`
}

func NewEnv(
	envName string,
	instanceType string,
	resolvedRepository ResolvedEnvRepository,
) *Env {

	return &Env{
		ID:                 uuid.NewString(),
		Name:               envName,
		InstanceType:       instanceType,
		SSHHostKeys:        []EnvSSHHostKey{},
		ResolvedRepository: resolvedRepository,
		OpenedPorts:        map[string]bool{},
		Status:             EnvStatusCreating,
		CreatedAtTimestamp: time.Now().Unix(),
	}
}

func (e *Env) GetNameSlug() string {
	return BuildEnvNameSlug(e.Name)
}

func (e *Env) GetSSHKeyPairName() string {
	return "yolo-" + e.GetNameSlug() + "-key-pair"
}

func (e *Env) SetInfrastructureJSON(infrastructure interface{}) error {
	infrastructureJSON, err := json.Marshal(infrastructure)

	if err != nil {
		return err
	}

	e.InfrastructureJSON = string(infrastructureJSON)

	return nil
}

func (e *Env) SetAdditionalPropertiesJSON(additionalProperties interface{}) error {
	additionalPropsJSON, err := json.Marshal(additionalProperties)

	if err != nil {
		return err
	}

	e.AdditionalPropertiesJSON = string(additionalPropsJSON)

	return nil
}

func BuildEnvNameSlug(name string) string {
	return slug.Make(name)
}

func BuildEnvNameFromResolvedRepo(
	resolvedRepo ResolvedEnvRepository,
) string {

	return resolvedRepo.Owner + "/" + resolvedRepo.Name
}

func ParseSSHHostKeys(hostKeysContent string) ([]EnvSSHHostKey, error) {
	sanitizedHostKeysContent := strings.TrimSpace(hostKeysContent)
	hostKeys := strings.Split(sanitizedHostKeysContent, "\n")

	parsedHostKeys := []EnvSSHHostKey{}

	for _, hostKey := range hostKeys {
		sanitizedHostKey := strings.TrimSpace(hostKey)
		hostKeyComponents := strings.Split(sanitizedHostKey, " ")

		// eg: (ssh-rsa) (AAAAB3NzaC1yc===) (root@ip-10-0-0-200)
		if len(hostKeyComponents) != 3 {
			return nil, fmt.Errorf("invalid host key (\"%s\")", hostKey)
		}

		parsedHostKeys = append(parsedHostKeys, EnvSSHHostKey{
			Algorithm:   hostKeyComponents[0],
			Fingerprint: hostKeyComponents[1],
		})
	}

	return parsedHostKeys, nil
}

func CheckPortValidity(port string, reservedPorts []string) error {
	portAsInt, err := strconv.Atoi(port)

	valid := err == nil && portAsInt >= 1 && portAsInt <= 65535

	if !valid {
		return ErrInvalidPort{
			InvalidPort: port,
		}
	}

	for _, reservedPort := range reservedPorts {
		if reservedPort == port {
			return ErrReservedPort{
				ReservedPort: port,
			}
		}
	}

	return nil
}
