# foundry

Simple and customizable scaffold generator as if molten metal pour into a mold in foundry.

## feature

* easy
  * can generate molds from existed sources
  * can use in any micro usecases quickly
  * need know "go template" only
  * can customize the molds same as source
  * can alias to minimize key stroke
* simple
  * need only few meta files

## usage

### To create scaffold by interactive mode

```bash
> foundry metal pour_into molds/cleanArchitecture/aggregationRoot destDir
Project: Full name identifing the project in the world. e.g. github.com/ArthurKC/foundry
Project?: github.com/ArthurKC/spiral
Name: The aggregation root name. It must be lower camel case.
Name?: user
created destDir/adapters/user/on_memory_repository.go
created destDir/domain/user/id.go
created destDir/domain/user/repository.go
created destDir/domain/user/user.go
```
