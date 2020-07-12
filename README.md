# foundry

Simple and customizable foundry metal with go mold.

## feature

* easy
  * can generate molds from existed sources
  * can use in any micro usecases quickly
  * need know "go mold" only
  * can be contained in any projects
  * can customize the molds same as source
  * can alias to minimize key stroke
* simple
  * need only few meta files
  * do only go mold


## usage

### To create foundry by interactive mode

```bash
> foundry create example/cleanArch/aggregationRoot destDir
Project: Full name identifing the project in the world. e.g. github.com/ArthurKC/foundry
Project?: github.com/ArthurKC/spiral
Name: The aggregation root name. It must be lower camel case.
Name?: user
created destDir/adapters/user/on_memory_repository.go
created destDir/domain/user/id.go
created destDir/domain/user/repository.go
created destDir/domain/user/user.go
```
