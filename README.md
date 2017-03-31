# dynamokv

Use AWS dynamodb as a simple Key Value storage.
Dynamokv is specially designed to store configuration in a dynamodb table and load them as environment variables.

## Usage

dynamokv store TABLENAME data.yml

dynamokv fetch TABLENAME

dynamokv set TABLENAME KEY VALUE

dynamokv get TABLENAME KEY

## Key Value File Format

```yaml
SIMPLE_KEY: VALUE GOES HERE
SERIALIZED_KEY:
  serialization: 'base64'
  value: |
    SOME LONG STRING
    WITH MULTIPLE LINES
ENCRYPTED_KEY:
  serialization:
    type: kms
    options:
      key: 'alias/key'
  value: YOUR SECRET VALUE
READ_FILE_KEY:
  value:
    file: 'config'
```

Supported Serialization types: base64 and kms. For KMS you need to provide key as option.
