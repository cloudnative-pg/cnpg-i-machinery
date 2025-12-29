# Changelog

## [0.4.2](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.4.1...v0.4.2) (2025-12-24)


### Bug Fixes

* **deps:** update all non-major go dependencies ([#229](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/229)) ([219b005](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/219b0052b1f2675db044e39ebb915f46954ae0a1))

## [0.4.1](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.4.0...v0.4.1) (2025-09-24)


### Bug Fixes

* **deps:** update all non-major go dependencies ([#203](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/203)) ([c5feb94](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/c5feb94e90fef9717113705526a15733af66863a))
* **deps:** update all non-major go dependencies ([#212](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/212)) ([c7ddf28](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/c7ddf28d72acea9ec5523d6cd1a8000bd604d2e1))
* **deps:** update all non-major go dependencies ([#215](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/215)) ([adc1c66](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/adc1c663c1c4acb07a30c3d26f7fedee8d44caf7))
* **deps:** update all non-major go dependencies ([#221](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/221)) ([ecbb3c4](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/ecbb3c47bf5d0a1d117d049079348ceec32fc004))
* reduce noisy logs for expected gRPC errors ([#227](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/227)) ([b88a426](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/b88a42687dd20c19d57ea72afa3a4b86b5604c38))

## [0.4.0](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.3.0...v0.4.0) (2025-07-03)


### Features

* support for Kubernetes native sidecar injection ([#196](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/196)) ([0ea207c](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/0ea207c829e5deff7fc51edfbc1c57719133ad20))
* **tls:** add client CA reload ([#197](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/197)) ([9f0abd6](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/9f0abd656267e101c99c2fc16de23fcfd9826989)), closes [#174](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/174)


### Bug Fixes

* **deps:** update all non-major go dependencies ([#181](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/181)) ([5b6e4ae](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/5b6e4aea6f18b211943fd9ec1ee6d3e68c3ad266))
* **deps:** update all non-major go dependencies ([#191](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/191)) ([c7e2dcf](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/c7e2dcf290d3012ccc53d81b54ecd8aa82239fef))
* **deps:** update all non-major go dependencies ([#195](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/195)) ([cf6f8ee](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/cf6f8ee1fe72526f53e4bdad007a8a525d404796))
* **deps:** update module github.com/cloudnative-pg/machinery to v0.3.0 ([#199](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/199)) ([34c5d73](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/34c5d7325ff15c63341f3c125e2a4c1fc6f93ad1))

## [0.3.0](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.2.1...v0.3.0) (2025-04-17)


### Features

* reload ssl cert to allow zero downtime rotation ([#174](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/174)) ([063997e](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/063997e162673e44439671aa4653715f424fa30e))


### Bug Fixes

* **deps:** downgrade github.com/grpc-ecosystem/go-grpc-middleware/v2 to v2.2.0 ([#179](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/179)) ([94ce92a](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/94ce92ab204379750d0a055f6b20df742e124520)), closes [#178](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/178)
* **deps:** update all non-major go dependencies ([#170](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/170)) ([3c6e63f](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/3c6e63f539a0c0fbd3f3c394d8fdb97f0c422af3))
* **deps:** update all non-major go dependencies ([#180](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/180)) ([cf9800a](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/cf9800a72f3b734ba2079be964ad479f40106bd6))
* logging setup in http.CreateMainCmd ([#169](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/169)) ([251db26](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/251db268655d5f00433c68480a7c2fce40c32479))

## [0.2.1](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.2.0...v0.2.1) (2025-03-28)


### Bug Fixes

* **deps:** update all non-major go dependencies ([#140](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/140)) ([1626c16](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/1626c16d54b6b311475d0d44cbfd1a1604d884d3))
* **deps:** update all non-major go dependencies ([#158](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/158)) ([6a3bad8](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/6a3bad87d20506e8249fcd7e2e5113a46e29ba46))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.23.1 ([#157](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/157)) ([4431fd1](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/4431fd1a5701c9b88eadb4de4b57a318cff0a514))

## [0.2.0](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.1.2...v0.2.0) (2025-03-13)


### Features

* configurable decoder of CNPG resources, lenient and strict ([#149](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/149)) ([6e06857](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/6e0685758c8d564facbc0d0386e151509f5910bb))

## [0.1.2](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.1.1...v0.1.2) (2025-02-27)


### Bug Fixes

* **deps:** set cloudnative-pg API version to latest ([#144](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/144)) ([2de477f](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/2de477fe4caa1bc95ea91b89e14b513cf879f099))

## [0.1.1](https://github.com/cloudnative-pg/cnpg-i-machinery/compare/v0.1.0...v0.1.1) (2025-02-27)


### Bug Fixes

* **deps:** downgrade go-grpc-middleware/v2 to v2.2.0 ([#139](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/139)) ([f02364a](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/f02364aea613e0cd8b32bf2bc237ef17d0acdcc1))

## 0.1.0 (2025-02-26)


### ⚠ BREAKING CHANGES

* use cnpg-machinery logging
* reorganize `pluginhelper` pkg ([#56](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/56))

### Features

* add `BuildSetStatusResponse` to plugin helper ([#47](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/47)) ([a3bd4ce](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/a3bd4ce2d72de59b1259df7b70ba7937d9c3abc0))
* add `GetKind` helper ([#46](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/46)) ([6509c5a](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/6509c5ad7f9a0dcbdef54b89eca7a2d32ac11b10))
* add backup helper functions ([0a49f6d](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/0a49f6de1ae86aabb0fb76f4f2404164acb87610))
* add pod spec methods for env injection ([#95](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/95)) ([2f1f208](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/2f1f20869d5cddbd5e1c9778e4a6d26f216fd644))
* initial import ([f9cade4](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/f9cade4b50973c72b2049d80202a96b1d23c420f))
* inject PG env ([#90](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/90)) ([e5495e9](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/e5495e9c5ed6fd1ee14a700d74fc3a395ffe866f))
* inject sidecar ([#65](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/65)) ([95a7e6c](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/95a7e6cb16f921e34f4188c6fed2f96a55f664e9))
* log server errors ([#54](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/54)) ([2a955f3](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/2a955f3116cd5faf8a83565b4e88df5a2c8441b1))
* propagate the logr.Logger into the gRPC server handlers ([#31](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/31)) ([26aafa5](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/26aafa55c7bf37e9f70e3db098e3fa9f52c463c1))
* refine logging context management and enrichers error management ([#2](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/2)) ([76ce219](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/76ce219b15a6f81494d9c374cfe3ad3db586f65f))
* support TLS connection ([#13](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/13)) ([32ad400](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/32ad400d28865d2659683b55a8be059be25e154a))


### Bug Fixes

* avoid explicit signal handling ([#111](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/111)) ([bd94f16](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/bd94f16685d31ee692b28f4b74603d80c515e864))
* **deps:** update all non-major go dependencies ([1886321](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/1886321540447e2f0fcaf19dab4011f067c59702))
* **deps:** update all non-major go dependencies ([c09e2e2](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/c09e2e24c34ef00ab950db84cad71d2224324356))
* **deps:** update all non-major go dependencies ([#122](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/122)) ([51a318b](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/51a318b7132ada7c8aba7bd6716c30ed6eca0976))
* **deps:** update all non-major go dependencies ([#124](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/124)) ([02474a2](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/02474a2d4040efd3d536a8ba5f5572b5a8fdc3bc))
* **deps:** update all non-major go dependencies ([#130](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/130)) ([991bf0e](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/991bf0e266c7516cd0994eec2418d7d901ee6369))
* **deps:** update all non-major go dependencies ([#133](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/133)) ([d1af2d2](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/d1af2d2478b4332a503c6422ec32f17848ae8501))
* **deps:** update all non-major go dependencies ([#39](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/39)) ([cd686e0](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/cd686e019766731bb79d494cf7bcbfb8979d4c6d))
* **deps:** update all non-major go dependencies ([#69](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/69)) ([16817f1](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/16817f104f8bcd7e07c96e5a4b5642bc743c12b2))
* **deps:** update all non-major go dependencies ([#89](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/89)) ([a545ade](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/a545adeb31e1c95505daec42e916ff7afc2ee877))
* **deps:** update github.com/cloudnative-pg/cnpg-i digest to bc221c3 ([#120](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/120)) ([df88260](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/df88260dbb6083dbc33b8f2cb6a4f593bdb7db6f))
* **deps:** update github.com/cloudnative-pg/machinery digest to 66cd032 ([#117](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/117)) ([285637c](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/285637c6574456fc53c39c9105cff33b0db916c9))
* **deps:** update module github.com/cloudnative-pg/api to v1 ([#118](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/118)) ([312c6c8](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/312c6c8764eecf9c4507762b5a2dae20987870e2))
* **deps:** update module github.com/onsi/ginkgo/v2 to v2.22.2 ([#123](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/123)) ([6ee50a1](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/6ee50a1f950fadb68bc8e43c015eaf418c52e471))
* **deps:** update module google.golang.org/grpc to v1.65.0 ([8baee90](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/8baee90500a40094f55348ccb25686c44bcebe0e))
* use logr-compatible log context access ([#12](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/12)) ([8576847](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/8576847b3449cf636fb1f85065fc052a10b767a7))


### Code Refactoring

* reorganize `pluginhelper` pkg ([#56](https://github.com/cloudnative-pg/cnpg-i-machinery/issues/56)) ([e9a72db](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/e9a72db1bfef00db871a6626bc06f54173e0c18d))
* use cnpg-machinery logging ([b2622f8](https://github.com/cloudnative-pg/cnpg-i-machinery/commit/b2622f81a69dcdb47a425399ee0c0128b03df15c))
