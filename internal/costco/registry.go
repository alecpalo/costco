package costco

import (
	"context"
	_ "costco/docs/swag"
	"costco/internal/storage"
	"costco/internal/storage/filesystem"
	"fmt"

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
	g *gin.Engine
	s storage.Storage
	a Authenticator
	k *kubernetes.Clientset
}

// ParseConfigs takes in the data from the costco configs ConfigMap and returns
// an storage implementation, authenticator and an error. Upon success returning nil.
// If more than one storage or authenticator is specified ParseConfigs will return an
// error. If there is an error initializing the interfaces, ParseConfigs will return an error.
func ParseConfigs(configs map[string]string) (storage.Storage, Authenticator, error) {
	var auth Authenticator
	var store storage.Storage

	storageType, ok := configs["storage"]
	if !ok {
		return nil, nil, errNoStorageSpecified
	}

	switch storageType {
	case S3:
		fmt.Println("s3 like storage")
	case FILESYSTEM:
		store = filesystem.Init()
	default:
		return nil, nil, errInvalidStorageType
	}

	authType, ok := configs["auth"]
	if !ok {
		return nil, nil, errNoAuthSpecified
	}

	switch authType {
	case BASIC:
		basicAuth := InitBasicClient()
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

	storage, auth, err := ParseConfigs(configs.Data)

	g := gin.Default()

	r := Registry{
		g: g,
		s: storage,
		a: auth,
		k: kClient,
	}

	return r
}

// RegisterEndpoints registers all endpoints for the container registry
func (r *Registry) RegisterEndpoints() {

	v2 := r.g.Group("/v2")
	{
		v2.GET("/", r.a.Authenticate(), r.CheckV2)

		v2.GET("/:name/tag/list", r.a.Authenticate(), r.ListTags)

		v2.GET("/:name/manifests/:reference", r.a.Authenticate(), r.PullManifest)
		v2.DELETE("/:name/manifests/:reference", r.a.Authenticate(), r.DeleteManifest)
		v2.HEAD("/:name/manifests/:reference", r.a.Authenticate(), r.CheckManifest)

		v2.GET("/:name/blobs/:digest", r.a.Authenticate(), r.GetBlob)
		v2.DELETE("/:name/blobs/:digest", r.a.Authenticate(), r.DeleteBlob)
		v2.HEAD("/:name/blobs/:digest", r.a.Authenticate(), r.CheckLayer)

		v2.POST("/:name/blobs/uploads", r.a.Authenticate(), r.StartUpload)

		v2.GET("/:name/blobs/uploads/:uuid", r.a.Authenticate(), r.CheckUpload)
		v2.PATCH("/:name/blobs/uploads/:uuid", r.a.Authenticate(), r.UploadLayer)
		v2.PUT("/:name/blobs/uploads/:uuid", r.a.Authenticate(), r.CompleteUpload)
		v2.DELETE("/:name/blobs/uploads/:uuid", r.a.Authenticate(), r.CancelUpload)

		v2.GET("_catalogs", r.a.Authenticate(), r.GetRepos)

		v2.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}

// Start runs the API server.
func (r *Registry) Start() {
	r.g.Run()
}
