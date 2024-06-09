# pgopher

[![Go Report Card](https://goreportcard.com/badge/github.com/mycreepy/pgopher)](https://goreportcard.com/report/github.com/mycreepy/pgopher)
[![Known Vulnerabilities](https://snyk.io/test/github/mycrEEpy/pgopher/badge.svg)](https://snyk.io/test/github/mycrEEpy/pgopher)

`pgopher` schedules the collection of CPU pprof profiles and allows you to download them via an API. These profiles can then be used for PGO during your next compilation.

## What is PGO?

[Profile-guided optimization](https://go.dev/doc/pgo) (PGO), also known as feedback-directed optimization (FDO), is a compiler optimization technique that feeds information (a profile) from representative runs of the application back into to the compiler for the next build of the application, which uses that information to make more informed optimization decisions. For example, the compiler may decide to more aggressively inline functions which the profile indicates are called frequently.

## Getting Started

1. Change `deploy/pgopher.yml` to your needs.

2. Deploy to Kubernetes:

```sh
kubectl apply -k deploy/
```

3. Download profiles:

```sh
curl http://[url]/api/v1/profile/[name] -o default.pgo
```

## Contributing

tba

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
