package kubernetes

import (
    "flag"
    "path/filepath"
    "io/ioutil"
    "context"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/client-go/util/homedir"
    "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
    "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"


    "k8s-kurated-addons.cli/src/utils/logger"

)

type KubernetesService struct {
    Client kubernetes.Clientset
}

func New() KubernetesService {
    var kubeconfig *string
    if home := homedir.HomeDir(); home != "" {
        kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
    } else {
        kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
    }

    return KubernetesService{
        Client: getClient(*kubeconfig),
    }
}

func getClient(kubeconfig string) kubernetes.Clientset {
    config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
    if err != nil {
        logger.PrintError("Unable to read config", err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        logger.PrintError("Unable to load k8s client", err)
    }

    return *clientset
}

func (ks KubernetesService) DeployManifest(manifestPath string) {

    logger.PrintInfo("Deploying Manifest...")
    yamlData, err := ioutil.ReadFile(manifestPath)

    if (err != nil) {
        logger.PrintError("Unable to read manifest", err)
    }

    decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
    obj     := &unstructured.Unstructured{}
    if _, _, err := decoder.Decode(yamlData, nil, obj); err != nil {
        logger.PrintError("Cannot decode YAML", err)
    }

    coreV1Client := ks.Client.CoreV1()
    req := coreV1Client.
        RESTClient().
        Post().
        Resource(obj.GetKind() + "s").
        Namespace("default").
        Body(obj)

    ctx := context.TODO()
    _, err = req.DoRaw(ctx)

    if (err != nil) {
        logger.PrintError("Failed to deploy YAML", err)
    }

    logger.PrintInfo("Deployed manifest successfully.")
}
