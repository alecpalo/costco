# Store

Store is the interface for the object store of the container registry. A basic container registry needs to be able to handle uploading layers, both as monoliths and as chunks, getting layers, deleting layers, checking if layers exist, completing layer uploads and uploading images.

## Images

Images are essentially just metadata about how the image is made. The file shows what the layers are for the image, which is essentially the instructions for how to create it.

When an image is pulled from a container registry, the manifest is pulled to determine which layers need to be pulled. Each layer is pulled which then allows the local docker daemon to build run the container.

When uploading an image, the essential steps are upload each layer, then upload a manifest that links each layer to a given step.

## Layers

Layers are tar files that contain the output from that step of building the container. The final layer in the step is used as the container, so to speak, after the addition of a writable layer. Layers can be uploaded either as chunks or as monoliths.

## Endpoints

Based on the above information, these are the endpoints that a store must implement.

1. Put Chunk 
2. Complete Layer 
3. Put Layer 
4. Get Layer
5. List Layers
6. Find Layer
7. Delete Layer
8. Create Manifest 
9. Delete Manifest 
