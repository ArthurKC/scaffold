# scaffold

Simple and customizable scaffold generator with go template.

## feature

* easy
  * can generate templates from existed sources
  * can use in any micro usecases quickly
  * need know "go template" only
  * can be contained in any projects
  * can customize the templates same as source
  * can alias to minimize key stroke
* simple
  * need only few meta files
  * do only go template


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
