package assets

import (
    "embed"
    "k8s.io/api/apps/v1"
    "k8s.io/apimachinery/pkg/runtime"
    "k8s.io/apimachinery/pkg/runtime/serializer"
)

// Embed all files in the manifests directory.
//go:embed manifests/*
var manifests embed.FS

var (
    appsScheme = runtime.NewScheme()
    appsCodecs = serializer.NewCodecFactory(appsScheme)
)

func init() {
    if err := v1.AddToScheme(appsScheme); err != nil {
        panic(err)
    }
}

func GetDeploymentFromFile(name string) *v1.Deployment {
    deploymentBytes, err := manifests.ReadFile(name)
    if err != nil {
        panic(err)
    }

    deploymentObject, err := runtime.Decode(
        appsCodecs.UniversalDecoder(v1.SchemeGroupVersion),
        deploymentBytes,
    )
    if err != nil {
        panic(err)
    }

    return deploymentObject.(*v1.Deployment)
}

