# Uploading an Image

The image upload process starts by getting a UUID from the container registry, then uploading each layer (either in chunks or as a monolith) and lastly completing the upload by uploading the manifest. Images are uploaded in layers, meaning that the config at the end is simply a metadata store to enable recreating images when pulling.

Below are the steps to upload a layer (and thus an image)

1. ```POST /v2/<name>/blobs/uploads/``` -> generate a UUID
    a. returns a UUID that should be used in the next steps
2. ```PATCH /v2/<name>/blobs/uplaods/<uuid>``` -> uploads a layer
    a. layers can be chunked or monoliths, a monolithic upload is an upload with one chunk
3. ```PUT /v2/<name>/blobs/uploads/<uuid>?digest=<digest>``` -> complete a layer upload
4. ```PUT /v2/<name>/manifests/<reference>``` -> push an image manifest

## Useful Links

1. [Pushing Images](https://docker-docs.uclv.cu/registry/spec/api/#pushing-an-image)