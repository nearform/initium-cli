# Changelog

## [0.8.4](https://github.com/nearform/initium-cli/compare/v0.8.3...v0.8.4) (2023-12-01)


### Bug Fixes

* 0.8.4 ([49b8fbc](https://github.com/nearform/initium-cli/commit/49b8fbc0a00df5955f2dd4da13990afd33cd989d))

## [0.8.3](https://github.com/nearform/initium-cli/compare/v0.8.2...v0.8.3) (2023-12-01)


### Bug Fixes

* test npm publish action ([2d000df](https://github.com/nearform/initium-cli/commit/2d000df740afa45027629066228a74a0ba4ea655))

## [0.8.2](https://github.com/nearform/initium-cli/compare/v0.8.1...v0.8.2) (2023-11-30)


### Bug Fixes

* release please action needs access to packages ([#160](https://github.com/nearform/initium-cli/issues/160)) ([caff703](https://github.com/nearform/initium-cli/commit/caff703a62558f4da24439d6992eb1f7b5a3dcee))

## [0.8.1](https://github.com/nearform/initium-cli/compare/v0.8.0...v0.8.1) (2023-11-30)


### Bug Fixes

* release-please and pipeline templates ([#158](https://github.com/nearform/initium-cli/issues/158)) ([d53d290](https://github.com/nearform/initium-cli/commit/d53d29050c40bf0d25d2ef710795629e2bf6f26f))

## [0.8.0](https://github.com/nearform/initium-cli/compare/v0.7.0...v0.8.0) (2023-11-30)


### Features

* add secrets management ([#145](https://github.com/nearform/initium-cli/issues/145)) ([9623e63](https://github.com/nearform/initium-cli/commit/9623e638883fd87fa4f46e1e7b006d4eea3758b3))
* enforce linux/amd64 images ([#154](https://github.com/nearform/initium-cli/issues/154)) ([90aa6d8](https://github.com/nearform/initium-cli/commit/90aa6d8210fdb78514062659275a407b2abbf8f3))
* post application url and commit hash as a PR comment after every deploy ([#148](https://github.com/nearform/initium-cli/issues/148)) ([79d5a2e](https://github.com/nearform/initium-cli/commit/79d5a2e4e6423a409aea36643be4dcc1e9405674))


### Bug Fixes

* cluster name for INITIUM_CLUSTER_ENDPOINT envvar ([#151](https://github.com/nearform/initium-cli/issues/151)) ([05c14db](https://github.com/nearform/initium-cli/commit/05c14db4836afb518c91427d5a39375cbab52e1c))
* code refactoring and cleanup ([#156](https://github.com/nearform/initium-cli/issues/156)) ([b08bdd7](https://github.com/nearform/initium-cli/commit/b08bdd7f54c174c7db0247c100b8cd7ac756c6b7))

## [0.7.0](https://github.com/nearform/initium-cli/compare/v0.6.0...v0.7.0) (2023-10-31)


### Features

* initium project type flag ([#101](https://github.com/nearform/initium-cli/issues/101)) ([321141e](https://github.com/nearform/initium-cli/commit/321141ea3fd5f6133708e5e4568ef141b73a197f))
* pass secrets to app ([#142](https://github.com/nearform/initium-cli/issues/142)) ([3edf2ea](https://github.com/nearform/initium-cli/commit/3edf2ea350d3e5997c95313d6b56ab412460b711))
* private registry with multiple secrets ([#130](https://github.com/nearform/initium-cli/issues/130)) ([a4a20f1](https://github.com/nearform/initium-cli/commit/a4a20f142dc57fb84d37549a621f4fd26656cd56))
* provide an option to deploy private services ([#132](https://github.com/nearform/initium-cli/issues/132)) ([9528449](https://github.com/nearform/initium-cli/commit/9528449c88eff8473bc75ad6df2def176b7defb2))


### Bug Fixes

* add initium-cli workflow permissions to write packages ([#141](https://github.com/nearform/initium-cli/issues/141)) ([ba6ec24](https://github.com/nearform/initium-cli/commit/ba6ec246be76ef4da735cc006242d4d01611600e))
* npm readme ([#136](https://github.com/nearform/initium-cli/issues/136)) ([7cfac95](https://github.com/nearform/initium-cli/commit/7cfac95bb8f93d28cdf10fefc1b165aebc0949e2))

## [0.6.0](https://github.com/nearform/initium-cli/compare/v0.5.0...v0.6.0) (2023-10-03)


### Features

* Come up with a standard to pass configuration to the deployed application ([#115](https://github.com/nearform/initium-cli/issues/115)) ([f1b3c2b](https://github.com/nearform/initium-cli/commit/f1b3c2baf4f47e19d59a7089b790c02e9c50b25c))
* unify all release steps in a single pipeline ([#113](https://github.com/nearform/initium-cli/issues/113)) ([93f4be9](https://github.com/nearform/initium-cli/commit/93f4be9305056e7c67f6cadeca1f7809a17efb88))


### Bug Fixes

* add annotations to ensure new docker image is downloaded ([#117](https://github.com/nearform/initium-cli/issues/117)) ([2627252](https://github.com/nearform/initium-cli/commit/262725205322c44870eccfbe351ff2ff448f1d94))

## [0.5.0](https://github.com/nearform/initium-cli/compare/v0.4.0...v0.5.0) (2023-09-21)


### Features

* allow release please to publish to npm ([5129767](https://github.com/nearform/initium-cli/commit/51297674339ff204afc71d6f6ee2ed38027fa9fa))

## [0.4.0](https://github.com/nearform/initium-cli/compare/v0.3.0...v0.4.0) (2023-09-21)


### Features

* change released binary name to initium ([#107](https://github.com/nearform/initium-cli/issues/107)) ([798a813](https://github.com/nearform/initium-cli/commit/798a813687a4c5356016b02dbfc292a65a5f772d))


### Bug Fixes

* execute publish steps only on release ([#110](https://github.com/nearform/initium-cli/issues/110)) ([abb4439](https://github.com/nearform/initium-cli/commit/abb4439d6638fca609316fa6b0fd621135cf75e6))
* release-please action failing for missing `)` ([#109](https://github.com/nearform/initium-cli/issues/109)) ([b12c919](https://github.com/nearform/initium-cli/commit/b12c919c241fcedb2594979cff666d63176a080d))

## [0.3.0](https://github.com/nearform/initium-cli/compare/v0.2.0...v0.3.0) (2023-09-04)


### Features

* get smarter at detecting app name ([#87](https://github.com/nearform/initium-cli/issues/87)) ([4fbe947](https://github.com/nearform/initium-cli/commit/4fbe947e9478d9452eac29495339cbee8ef5ea67))


### Bug Fixes

* avoid running closed PR action immediately after reopen ([#98](https://github.com/nearform/initium-cli/issues/98)) ([107b505](https://github.com/nearform/initium-cli/commit/107b505c943bf3c487d04c720f0ff0ac8e6576da))
* the init github command does not require the shared flags ([#93](https://github.com/nearform/initium-cli/issues/93)) ([6c9a32c](https://github.com/nearform/initium-cli/commit/6c9a32cfc682b5a081a84eb93c2d3550730720e8))

## [0.2.0](https://github.com/nearform/initium-cli/compare/v0.1.0...v0.2.0) (2023-08-18)


### Features

* add option to persist configuration ([#80](https://github.com/nearform/initium-cli/issues/80)) ([6c9f3ed](https://github.com/nearform/initium-cli/commit/6c9f3ed5ae7f9cd05a3f5a75610bab12f5bf57bf))


### Bug Fixes

* move from repo-name to container-registry ([#79](https://github.com/nearform/initium-cli/issues/79)) ([cebfa95](https://github.com/nearform/initium-cli/commit/cebfa954d362d9651596ed415abe53f6a428fc17))
* remove outdated quick-start ([#82](https://github.com/nearform/initium-cli/issues/82)) ([d3180a8](https://github.com/nearform/initium-cli/commit/d3180a833ac340223b33816c0ba42b5c1711ac89))

## [0.1.0](https://github.com/nearform/initium-cli/compare/v0.0.1...v0.1.0) (2023-08-11)


### Features

* use app token to generate relases from release-please ([#75](https://github.com/nearform/initium-cli/issues/75)) ([bc82eda](https://github.com/nearform/initium-cli/commit/bc82eda1b3767f2244b58d2982e3cf8da2059166))


### Bug Fixes

* update release-please.yaml ([#68](https://github.com/nearform/initium-cli/issues/68)) ([b06c556](https://github.com/nearform/initium-cli/commit/b06c556b9b393172d49945130569dd749c3af672))
* update to latest github-app-token action ([#76](https://github.com/nearform/initium-cli/issues/76)) ([539722f](https://github.com/nearform/initium-cli/commit/539722f5b82240ceae69aa307ace5ac9f40183df))

## 0.0.1 (2023-08-10)


### Features

* add action ([c3e489e](https://github.com/nearform/initium-cli/commit/c3e489e4f949959479be6c1e133e1b4b4be3fe0c))
* Add github init command ([#20](https://github.com/nearform/initium-cli/issues/20)) ([5631ef3](https://github.com/nearform/initium-cli/commit/5631ef392757dd39dc4ccda78fb8ca868d4fd576))
* Add Installation Command ([#46](https://github.com/nearform/initium-cli/issues/46)) ([233f4bd](https://github.com/nearform/initium-cli/commit/233f4bd593a1390730fc3954cb413383f2855143))
* Add support for configuration file ([#57](https://github.com/nearform/initium-cli/issues/57)) ([20a057d](https://github.com/nearform/initium-cli/commit/20a057d9e773b8b0a31d1d9c4800357bf37ee54a))
* adding checkout action ([e8ab9f8](https://github.com/nearform/initium-cli/commit/e8ab9f8cda3f5e32b16410c569eb0c1652d6d834))
* adding correct folder ([8a5f4b0](https://github.com/nearform/initium-cli/commit/8a5f4b05cb22bb439a0cf3331b59b25a51993588))
* adding initial tests ([eb1f803](https://github.com/nearform/initium-cli/commit/eb1f803c5d1a18291a1c2a283b299a78c3915112))
* adding job for commit and push nodejs project ([c6449ef](https://github.com/nearform/initium-cli/commit/c6449ef4a5878e5086c410a2e352b2f595e72349))
* adding make project_build to test the workflow ([45fa49f](https://github.com/nearform/initium-cli/commit/45fa49fb36012c65e9f8e47af6334ad23c2a9e10))
* adding more tests ([d998753](https://github.com/nearform/initium-cli/commit/d998753908df0b1cd9a2343bee861bb39e24c3eb))
* adding more tests ([0d21958](https://github.com/nearform/initium-cli/commit/0d2195864dc4b5925172e910bfdf3897dc722d91))
* adding small and simple test ([5e231b4](https://github.com/nearform/initium-cli/commit/5e231b41143e5f3e85f228ba4e431c7c7e54636f))
* adding test pipeline and make build only on master ([f9fff0d](https://github.com/nearform/initium-cli/commit/f9fff0d8cec7bcb296b4062d3a9527bec4680997))
* changes in the code ([d985443](https://github.com/nearform/initium-cli/commit/d985443ae863cb74f67855c5a333761290f40642))
* Implement the on branch workflow ([#58](https://github.com/nearform/initium-cli/issues/58)) ([2238866](https://github.com/nearform/initium-cli/commit/2238866bbfcd122757429a3d7e5a86e798ac1d2d))
* improve workflow ([#59](https://github.com/nearform/initium-cli/issues/59)) ([315926b](https://github.com/nearform/initium-cli/commit/315926bb016f6659bfff9f0c4a929b29cc9971c2))
* **issue-24:** Refactor Docker Image attributes and add files for knative templating ([#32](https://github.com/nearform/initium-cli/issues/32)) ([04e43fe](https://github.com/nearform/initium-cli/commit/04e43feb76c9893adc00a6d67b89e2277806161c))
* making things simpler ([f74d5da](https://github.com/nearform/initium-cli/commit/f74d5da28f89caa2788b18560b82f420acf20ae7))
* moving checkout action ([1d08b9f](https://github.com/nearform/initium-cli/commit/1d08b9fff6aa8daabe656cb1cedd46bb7577ea2c))
* moving the env for the whole job ([e6750b7](https://github.com/nearform/initium-cli/commit/e6750b716d26815b9fb589ece7d22ebb111d5a8d))
* remove action ([a45d5e0](https://github.com/nearform/initium-cli/commit/a45d5e01e62bc942d02b9a56b1652061dc2ea22a))
* remove duplicated tests ([01e31ab](https://github.com/nearform/initium-cli/commit/01e31ab42dc463377202ca4ad7dbac55d002912d))
* rename module to initium-cli ([#62](https://github.com/nearform/initium-cli/issues/62)) ([677885d](https://github.com/nearform/initium-cli/commit/677885df1d969ea2a5275a1b95ca37ed0314c173))
* updating token value ([4694813](https://github.com/nearform/initium-cli/commit/4694813c2044182a02a1fbae3464892765633e17))
* using local path ([09ea0d6](https://github.com/nearform/initium-cli/commit/09ea0d6a8003ee0233cdc4c49b6142582de29555))
* using reusable actions to allow add more jobs easily ([bb0bd27](https://github.com/nearform/initium-cli/commit/bb0bd270ca5003246da49776cbff74ece3330fc7))


### Bug Fixes

* correct shell, correct env ([68497e3](https://github.com/nearform/initium-cli/commit/68497e3220f6516d30b234549e54f7ecf45fe7a9))
* fixing reference ([411612f](https://github.com/nearform/initium-cli/commit/411612fc26c22b2345b67c536c6c8b5ab5140ddb))
* fixing reference ([0cabf57](https://github.com/nearform/initium-cli/commit/0cabf576237f2a0d4b8c16b2d0f20c82fc6d762f))
* fixing reference ([5613ff0](https://github.com/nearform/initium-cli/commit/5613ff01c301ce4fc366b933369889b3f2370e92))
* move secret input ([e53cbc2](https://github.com/nearform/initium-cli/commit/e53cbc2b21cbe665942d0d0fdc2d16335910d6cd))
* Update release-please.yaml ([#67](https://github.com/nearform/initium-cli/issues/67)) ([1876eb6](https://github.com/nearform/initium-cli/commit/1876eb6d3e98898bd0aaec6b0fba76fece118d99))
