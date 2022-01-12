package main

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
    "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
    Host         string `yaml:"host" env:"HOST" env-default:"0.0.0.0"`
    Port         string `yaml:"port" env:"PORT" env-default:"8080"`
    Loglevel     string `yaml:"logLevel" env:"LOG_LEVEL" env-default:"error"`
	RedisHost    string `yaml:"redisHost" env:"REDIS_SERVICE_HOST" env-default:"redis"`
	RedisPort    string `yaml:"redisPort" env:"REDIS_SERVICE_PORT_MAIN" env-default:"6379"`
	K8sApiHost   string `yaml:"k8sApiHost" env:"KUBERNETES_SERVICE_HOST"`
	K8sApiPort   string `yaml:"k8sApiPort" env:"KUBERNETES_SERVICE_PORT"`
	K8sNamespace string `yaml:"k8sNamespace" env:"K8S_NAMESPACE"`
	K8sToken     string `yaml:"k8sToken" env:"K8S_TOKEN"`
	K8sService   string `yaml:"k8sService" env:"K8S_SERVICE"`
}

var (
	config Config
    k8sTokenPath string = "/run/secrets/kubernetes.io/serviceaccount/token"
    k8sNamespacePath string = "/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func init() {
    err := cleanenv.ReadConfig("config.yaml", &config)
    if err != nil {
        log.Error("Can`t read config.yaml")
    }
	k8sToken, err := ioutil.ReadFile(k8sTokenPath)
    if err != nil {
        log.Fatal(err)
    }
    k8sNamespace, err := ioutil.ReadFile(k8sNamespacePath)
    if err != nil { 
        log.Fatal(err)
    }
	config.K8sToken = string(k8sToken)
	config.K8sNamespace = string(k8sNamespace)
}
