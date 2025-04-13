package costco

import (
	"context"
	_ "costco/docs/swag"
	"costco/internal/storage"
	"costco/internal/storage/s3"
	"fmt"
	"sigs.k8s.io/yaml"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
)

// Registry is an object containing all required fields for the container registry
type Registry struct {
	web   *gin.Engine
	store storage.Store
	auth  Authenticator
	k8s   *kubernetes.Clientset
}

// ParseConfigs takes in the data from the costco configs ConfigMap and returns
// a storage implementation, authenticator and an error. Upon success returning nil.
// If more than one storage or authenticator is specified ParseConfigs will return an
// error. If there is an error initializing the interfaces, ParseConfigs will return an error.
func ParseConfigs(configs map[string]string) (storage.Store, Authenticator, error) {
	var auth Authenticator
	var store storage.Store

	yamlConfigs, ok := configs["configs.yaml"]
	if !ok {
		return nil, nil, fmt.Errorf("configs.yaml not found in configs.yaml")
	}

	var costcoConfigs Configs
	err := yaml.Unmarshal([]byte(yamlConfigs), &costcoConfigs)
	if err != nil {
		return nil, nil, err
	}

	switch costcoConfigs.Storage.kind {
	case S3:
		store = s3.Init()
	case FILESYSTEM:
		if costcoConfigs.Storage.FileSystemConfig == nil {
			return nil, nil, errNoStorageSpecified
		}
		//store = filesystem.Init(costcoConfigs.Storage.FileSystemConfig.filePath)
	default:
		return nil, nil, errInvalidStorageType
	}

	switch costcoConfigs.Auth.kind {
	case BASIC:
		if costcoConfigs.Auth.BasicAuthConfig == nil {
			return nil, nil, errNoAuthSpecified
		}
		basicAuth := InitBasicClient(costcoConfigs.Auth.BasicAuthConfig.Username, costcoConfigs.Auth.BasicAuthConfig.Password)
		auth = &basicAuth
	default:
		return nil, nil, errInvalidAuthType
	}

	return store, auth, nil
}

// ParseCommands parses the cobra commands returning the namespace specified
// in the command. Errors occurring during parsing will be logged then ignored.
func ParseCommands(cmd *cobra.Command) string {
	namespace, err := cmd.Flags().GetString("namespace")
	if err != nil {
		klog.Error(err)
		namespace = "costco"
	}

	// debug turns on debugging logs
	debug, err := cmd.Flags().GetBool("debug")
	if err != nil {
		klog.Error(err)
		debug = false
	}

	// verbose turns on informational logs
	verbose, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		klog.Error(err)
		verbose = false
	}

	if debug {
		klog.V(1)
	}

	if verbose {
		klog.V(2)
	}

	return namespace
}

// Init returns an initialized Registry object by parsing the configs specified in the
// costco-configs ConfigMap.
func Init(cmd *cobra.Command) Registry {
	namespace := ParseCommands(cmd)

	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatal(err)
	}

	kClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatal(err)
	}

	configs, err := kClient.CoreV1().ConfigMaps(namespace).Get(context.Background(), "costco-configs", v1.GetOptions{})
	if err != nil {
		klog.Fatal(err)
	}

	store, auth, err := ParseConfigs(configs.Data)
	if err != nil {
		klog.Fatal(err)
	}

	g := gin.Default()

	r := Registry{
		web:   g,
		store: store,
		auth:  auth,
		k8s:   kClient,
	}

	return r
}

// RegisterEndpoints registers all endpoints for the container registry
func (r *Registry) RegisterEndpoints() {

	v2 := r.web.Group("/v2")
	{
		v2.GET("/", r.auth.Authenticate(), r.CheckV2)

		v2.GET("/:name/tag/list", r.auth.Authenticate(), r.ListTags)

		v2.GET("/:name/manifests/:reference", r.auth.Authenticate(), r.PullManifest)
		v2.DELETE("/:name/manifests/:reference", r.auth.Authenticate(), r.DeleteManifest)
		v2.HEAD("/:name/manifests/:reference", r.auth.Authenticate(), r.CheckManifest)

		v2.GET("/:name/blobs/:digest", r.auth.Authenticate(), r.GetBlob)
		v2.DELETE("/:name/blobs/:digest", r.auth.Authenticate(), r.DeleteBlob)
		v2.HEAD("/:name/blobs/:digest", r.auth.Authenticate(), r.CheckLayer)

		v2.POST("/:name/blobs/uploads", r.auth.Authenticate(), r.StartUpload)

		v2.GET("/:name/blobs/uploads/:uuid", r.auth.Authenticate(), r.CheckUpload)
		v2.PATCH("/:name/blobs/uploads/:uuid", r.auth.Authenticate(), r.UploadLayer)
		v2.PUT("/:name/blobs/uploads/:uuid", r.auth.Authenticate(), r.CompleteUpload)
		v2.DELETE("/:name/blobs/uploads/:uuid", r.auth.Authenticate(), r.CancelUpload)

		v2.GET("_catalogs", r.auth.Authenticate(), r.GetRepos)

		v2.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

// Start runs the API server.
func (r *Registry) Start() {
	err := r.web.Run()
	if err != nil {
		klog.Fatal(err)
	}
}
