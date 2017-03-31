# dynamokv

Use AWS dynamodb as a simple Key Value storage.
Dynamokv is specially designed to store configuration in a dynamodb table and load them as environment variables.

## Usage

dynamokv store TABLENAME data.yml

dynamokv fetch TABLENAME

dynamokv set TABLENAME KEY VALUE

dynamokv get TABLENAME KEY

dynamokv template TABLENAME TEMPLATEFILE

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

## Template File Format

```
The value for KEY_NAME is {{KEY_NAME}}
The value without deserializing for KEY_NAME is {{RAW:KEY_NAME}}
```


Supported Serialization types: base64 and kms. For KMS you need to provide key as option.
