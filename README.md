[![Build](https://github.com/nao1215/contributor/actions/workflows/build.yml/badge.svg)](https://github.com/nao1215/contributor/actions/workflows/build.yml)
[![UnitTest](https://github.com/nao1215/contributor/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/contributor/actions/workflows/unit_test.yml)
[![reviewdog](https://github.com/nao1215/contributor/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/contributor/actions/workflows/reviewdog.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/contributor)](https://goreportcard.com/report/github.com/nao1215/contributor)
![GitHub](https://img.shields.io/github/license/nao1215/contributor)  
# contributor - print contributor List (only support git)
contributor command print list of people who have modified code or documentation in the git project.

```
$ cd <PROJECT_ROOT_DIR>
$ contributor 
+-------------------------+-----------------------------------------------------------+-----------+-----------+
|          NAME           |                           EMAIL                           | +(APPEND) | -(DELETE) |
+-------------------------+-----------------------------------------------------------+-----------+-----------+
| Ichinose Shogo          | shogo82148@gmail.com                                      |     11042 |      6044 |
| Daisuke Maki            | lestrrat+github@gmail.com                                 |       866 |       223 |
| Songmu                  | y.songmu@gmail.com                                        |       237 |        65 |
| Stefan Tudose           | stefan.tudose@data4life.care                              |        14 |        12 |
| mattn                   | mattn.jp@gmail.com                                        |         9 |         9 |
| yusuke-enomoto          | yusuke.enomoto@dena.com                                   |         8 |         6 |
| pyros2097               | pyros2097@gmail.com                                       |         3 |         1 |
| catatsuy                | catatsuy@catatsuy.org                                     |         2 |         2 |
| Shoma Okamoto           | 32533860+shmokmt@users.noreply.github.com                 |         1 |         1 |
| nasa9084                | nasa.9084.bassclarinet@gmail.com                          |         1 |         1 |
| dependabot-preview[bot] | 27856297+dependabot-preview[bot]@users.noreply.github.com |         0 |         0 |
+-------------------------+-----------------------------------------------------------+-----------+-----------+
```

# How to install
### Step.1 Install golang
contributor command only supports installation with `$ go install`. If you does not have the golang development environment installed on your system, please install golang from the [golang official website](https://go.dev/doc/install).

### Step2. Install contributor
```
$ go install github.com/nao1215/contributor@latest
```
# Contributing
First off, thanks for taking the time to contribute! ❤️
See [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.  

# Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.

- [GitHub Issue](https://github.com/nao1215/contributor/issues)

# LICENSE
The contributor project is licensed under the terms of [MIT LICENSE](./LICENSE).
