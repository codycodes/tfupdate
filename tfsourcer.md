# TODO

- [ ] pass in a module name and change its version (current)
- [ ] pass in a module name and recursively change version of all child modules (current)
- [ ] pass in a module name and remove its version (new) - for all child modules?
- [ ] pass in a module name and add its version (new) - for all child modules?

🤔 Determine on the last two what the best structure is for argument parsing:
Currently tfupdate only has: `tfupdate module terraform-aws-modules/vpc/aws -v 2.0.0`

See what happens when:

- A file is passed in with no version attr on the existing file
  - A version must be specified to update to
  - What happens with local modules?

- absolute paths not supported by afero.ts

## tests

should match from a prefixed source of something to something else. I want to start with localterraform.com to a terraform cloud/enterprise source and vice-versa.

Simple things to start with:

- [ ] Do not change a source if newSource is not valid
- [ ] Ability to change a single module that does not have a specified version
- [ ] Do not modify modules which do not match
