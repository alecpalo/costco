# Configuration

## Storage

Two storage types are accepted by costco:

1. file-system
2. s3

```yaml
storage: file-system
storage: s3
```

### File System

The file system storage requires the following arguments:

1. base-path

```yaml
storage: file-system
storage-base-path: example/file/path
```

### S3 

The S3 storage works with any S3 compatible storage system. Below are the required and optional arguments

## Authentication

There is currently one authentication type:

1. basic

```yaml
auth: basic
```

### Basic

Basic auth accepts the following configurations:

1. username
2. password

```yaml
auth: basic
auth-username: username
auth-password: password
```