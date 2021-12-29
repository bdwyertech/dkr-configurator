package main

import (
	"context"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Jeffail/gabs/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/ghodss/yaml"
)

func awsConfig() aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func GetParameter(key string) (value string, err error) {
	prm := strings.Split(key, ":parameter")[1]
	region := strings.Split(key, ":")[3]

	ssmClient := ssm.NewFromConfig(awsConfig(), func(o *ssm.Options) {
		o.Region = region
	})

	resp, err := ssmClient.GetParameter(context.Background(), &ssm.GetParameterInput{
		Name:           aws.String(prm),
		WithDecryption: true,
	})
	if err != nil {
		return
	}

	value = *resp.Parameter.Value

	return
}

func GetParametersByPath(path string) (parameters []ssmTypes.Parameter) {
	// SSM Client
	ssmclient := ssm.NewFromConfig(awsConfig())

	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		WithDecryption: true,
		Recursive:      true,
	}

	p := ssm.NewGetParametersByPathPaginator(ssmclient, input)
	for p.HasMorePages() {
		resp, err := p.NextPage(context.Background())
		if err != nil {
			log.Fatalf("ERROR: ssm.GetParametersByPath:: %s\n%s", path, err)
		}
		parameters = append(parameters, resp.Parameters...)
	}

	return
}

func GetParametersByPathJSON(path string) (jsonBytes []byte) {
	parameters := GetParametersByPath(path)
	cfg := &gabs.Container{}
	for _, parameter := range parameters {
		// JSON Pointer https://tools.ietf.org/html/rfc6901
		cfg.SetJSONPointer(*parameter.Value, *parameter.Name)
	}

	// Pretty JSON
	jsonBytes = []byte(cfg.StringIndent("", "  "))

	return
}

func GetParametersByPathYAML(path string) (ymlBytes []byte) {
	jsonBytes := GetParametersByPathJSON(path)
	ymlBytes, err := yaml.JSONToYAML(jsonBytes)
	if err != nil {
		log.Fatalf("ERROR: JSONtoYAML\n%s", err)
	}

	return
}
