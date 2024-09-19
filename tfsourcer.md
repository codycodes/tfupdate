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
