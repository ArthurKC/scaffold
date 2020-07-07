# scaffold

Simple and customizable scaffold generator with go template.

## usage

### To create scaffold by interactive mode

```bash
> scaffold create example/cleanArch/aggregationRoot destDir
Project: Full name identifing the project in the world. e.g. github.com/ArthurKC/scaffold
Project?: github.com/ArthurKC/spiral
Name: The aggregation root name. It must be lower camel case.
Name?: user
created destDir/adapters/user/on_memory_repository.go
created destDir/domain/user/id.go
created destDir/domain/user/repository.go
created destDir/domain/user/user.go
```
