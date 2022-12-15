# cli-plugin-airdrop

## Developer instruction

- Use `ignite` binary compiled from this branch https://github.com/ignite/cli/pull/3306 (unless it's merged!)
- clone this repo locally
- Run `ignite plugin add -g /absolute/path/to/cli-plugin-airdrop` to add the plugin to global config
- `ignite airdrop` command is now available.

Then repeat that loop :
- Hack plugin code
- Rerun `ignite airdrop` to recompile the plugin and test

