## aman

[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://golang.org)
[![Go Report Card](https://goreportcard.com/badge/github.com/naruhiyo/aman)](https://goreportcard.com/report/github.com/naruhiyo/aman)
[![Workflow](https://github.com/naruhiyo/aman/workflows/reviewdog/badge.svg)](https://github.com/naruhiyo/aman/actions?query=workflow%3Areviewdog)
[![GitHub issues](https://img.shields.io/github/issues/naruhiyo/aman.svg)](https://github.com/naruhiyo/aman/issues)
[![License: MIT](https://img.shields.io/badge/license-MIT-teal.svg)](https://github.com/naruhiyo/aman/blob/master/LICENSE)

***A***ssistant tool to get options easily from ***man***.

<!-- START doctoc -->
<!-- END doctoc -->

### Features

- Display options with `man` what you input.

![Display](https://user-images.githubusercontent.com/28133383/103387478-67226280-4b47-11eb-902f-34267901fd0a.gif)

- Filter options by `AND` in case-insensitive.

![Filter](https://user-images.githubusercontent.com/28133383/103393025-ffc6db80-4b63-11eb-8bc0-4512d46e0967.gif)

- Select option.
`Enter` to  select and `Esc` to finish.

![Select-option](https://user-images.githubusercontent.com/28133383/103387477-64c00880-4b47-11eb-8611-9ed28e8b3dd7.gif)

- Select multiple options.
`Enter` to select and `Esc` to finish.

![Select-multiple-options](https://user-images.githubusercontent.com/28133383/103387459-48bc6700-4b47-11eb-965a-1917ab7e1d37.gif)

### Support OS

![Ubuntu](https://img.shields.io/badge/-Ubuntu-6c272d.svg?logo=Ubuntu&style=flat-square)
![Mac or macOS](https://img.shields.io/badge/-Mac-000000.svg?logo=apple&style=flat-square)

### Installation

#### macOS

```console
$ brew tap naruhiyo/aman
$ brew install naruhiyo/aman/aman
```

#### Ubuntu

1. Download tarball file from [GitHub releases](https://github.com/naruhiyo/aman/releases/).
2. Unpack the file.

```console
$ tar -zxvf xxxx.tar.gz
```

3. Locate `aman` to directory to your `$PATH`.

For example, locate `/usr/local/bin`.

```console
$ mv aman_linux_amd64/aman /usr/local/bin/
```

### Contributors

- [narugit](https://github.com/narugit)

- [hiyoko3](https://github.com/hiyoko3)

### License

[MIT](https://github.com/naruhiyo/aman/blob/master/LICENSE)
