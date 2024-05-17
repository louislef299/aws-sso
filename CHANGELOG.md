# Changelog

## [1.2.3](https://github.com/louislef299/aws-sso/compare/v1.2.2...v1.2.3) (2024-05-17)


### Bug Fixes

* Bump github.com/aws/aws-sdk-go-v2/config from 1.27.11 to 1.27.13 ([#109](https://github.com/louislef299/aws-sso/issues/109)) ([b7d0c9f](https://github.com/louislef299/aws-sso/commit/b7d0c9f7fd32ff8fab64ef6434adf7b39417be98))
* Bump github.com/aws/aws-sdk-go-v2/service/eks from 1.42.1 to 1.42.2 ([#107](https://github.com/louislef299/aws-sso/issues/107)) ([419a060](https://github.com/louislef299/aws-sso/commit/419a06053596450ee69e2ef5f82129f6f78cc780))
* Bump github.com/docker/cli ([#108](https://github.com/louislef299/aws-sso/issues/108)) ([e2455b5](https://github.com/louislef299/aws-sso/commit/e2455b578d01515980aadd5af875eb21e04a9597))
* Bump github.com/onsi/ginkgo/v2 from 2.17.2 to 2.17.3 ([#102](https://github.com/louislef299/aws-sso/issues/102)) ([ae77a4b](https://github.com/louislef299/aws-sso/commit/ae77a4b1b1054b4fd170f11a8f48e326b82d6205))
* Update token hash ([#110](https://github.com/louislef299/aws-sso/issues/110)) ([bbe7373](https://github.com/louislef299/aws-sso/commit/bbe737396d2c963fedab474a1afdf3517859e3ba))

## [1.2.2](https://github.com/louislef299/aws-sso/compare/v1.2.1...v1.2.2) (2024-05-07)


### Bug Fixes

* SEGFAULT on empty ns ([#97](https://github.com/louislef299/aws-sso/issues/97)) ([5c8f820](https://github.com/louislef299/aws-sso/commit/5c8f8209af5cd5ac6aece5206633aeda55d32918))

## [1.2.1](https://github.com/louislef299/aws-sso/compare/v1.2.0...v1.2.1) (2024-04-29)


### Bug Fixes

* Enable role regex searching ([#92](https://github.com/louislef299/aws-sso/issues/92)) ([cd7841c](https://github.com/louislef299/aws-sso/commit/cd7841c3f98b5a0d71825f1e8521df7c869cbc3e))
* Only print user-desired configs with plain list ([#89](https://github.com/louislef299/aws-sso/issues/89)) ([6d50481](https://github.com/louislef299/aws-sso/commit/6d504810c3e80a566776991517c299aca991546b))

## [1.2.0](https://github.com/louislef299/aws-sso/compare/v1.1.2...v1.2.0) (2024-04-29)


### Features

* Associate persistent region to account alias ([#86](https://github.com/louislef299/aws-sso/issues/86)) ([3db01ab](https://github.com/louislef299/aws-sso/commit/3db01abc06f7a92578c1d688474b11f701999cfd))


### Bug Fixes

* Allow for static SSO region and dynamic --region flag with session ([#80](https://github.com/louislef299/aws-sso/issues/80)) ([41f60e1](https://github.com/louislef299/aws-sso/commit/41f60e1060634079ecbc9d62702f27e63dbbce76))
* Persist namespace when switching contexts ([#85](https://github.com/louislef299/aws-sso/issues/85)) ([6d441f9](https://github.com/louislef299/aws-sso/commit/6d441f96beba94d68643cd261f23bee05e1a9a0b))

## [1.1.2](https://github.com/louislef299/aws-sso/compare/v1.1.1...v1.1.2) (2024-04-10)


### Bug Fixes

* Documentation wording and Browser docs ([#54](https://github.com/louislef299/aws-sso/issues/54)) ([cbfb8d5](https://github.com/louislef299/aws-sso/commit/cbfb8d54a3a38cdcc97f6db49ac9dab8a4dfed95))

## [1.1.1](https://github.com/louislef299/aws-sso/compare/v1.1.0...v1.1.1) (2024-03-18)


### Bug Fixes

* Add Browser support for linux and windows operating systems ([#51](https://github.com/louislef299/aws-sso/issues/51)) ([dd7395e](https://github.com/louislef299/aws-sso/commit/dd7395e39f6dee3f069c95ed5f5aeeaa28711e3e))
* Shift active token on creation ([#44](https://github.com/louislef299/aws-sso/issues/44)) ([ff33ba4](https://github.com/louislef299/aws-sso/commit/ff33ba4583543ed1c76cbc1c99bd9ade28eba6b1))
* unknown acct flow ([#52](https://github.com/louislef299/aws-sso/issues/52)) ([8c1663a](https://github.com/louislef299/aws-sso/commit/8c1663a26e3e3f8ea3c17d9ba83f2cdaf36f74da))

## [1.1.0](https://github.com/louislef299/aws-sso/compare/v1.0.3...v1.1.0) (2024-03-04)


### Features

* Add browser support on macOS ([#24](https://github.com/louislef299/aws-sso/issues/24)) ([91f3a2c](https://github.com/louislef299/aws-sso/commit/91f3a2cf851400ffc484d6edfb2cebea1ec92f00))

## [1.0.3](https://github.com/louislef299/aws-sso/compare/v1.0.2...v1.0.3) (2024-02-09)


### Bug Fixes

* Config file not found ([#12](https://github.com/louislef299/aws-sso/issues/12)) ([8a93f5f](https://github.com/louislef299/aws-sso/commit/8a93f5f54129fb2185b37586749e920db943f5e3))

## [1.0.2](https://github.com/louislef299/aws-sso/compare/v1.0.1...v1.0.2) (2024-01-19)


### Bug Fixes

* Remove known issue for EU ([#9](https://github.com/louislef299/aws-sso/issues/9)) ([635d81f](https://github.com/louislef299/aws-sso/commit/635d81f8dd45c489fce58e697462ab33408a9b71))

## [1.0.1](https://github.com/louislef299/aws-sso/compare/v1.0.0...v1.0.1) (2024-01-17)


### Bug Fixes

* Allow for auth with diff region than cluster auth ([#6](https://github.com/louislef299/aws-sso/issues/6)) ([4807df8](https://github.com/louislef299/aws-sso/commit/4807df8a9b6a1013774062bd1539feac32d39336))

## 1.0.0 (2023-12-19)


### Features

* Initialize aws-sso tool ([a981213](https://github.com/louislef299/aws-sso/commit/a981213c540bf7ffc4b928b158a9fb65625593fb))
