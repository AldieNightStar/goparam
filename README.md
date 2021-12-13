# Go Param Generator

# Overview
* Generates `struct`'s for you. No need to write it by hands
* Adds Getters and Setters
* Adds Constructor
* Easy to use
* Adds builder (Killer feature)

# Usage
```
goparam -package main -name User -params "Age int, Id int, Type UserType" -get -set -ctor -builder

Options:
  -name    - Type/Struct name
  -params  - Parameter list, separated by comma: Name Type, Name2 Type2, ...
  -package - Package name
  -get     - Generate getters
  -set     - Generate setters
  -ctor    - Generate constructor New...()
  -builder - Adds builder feature to your struct
```
* Sample
```
//go:generate goparam -name User -package main -params "Age int, Id int, Name string" -get -set -ctor
```
