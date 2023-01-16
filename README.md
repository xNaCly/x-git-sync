# git-auto-sync

Backup your repository at configured intervals

## Why use gas

- highly configurable
- 0 external dependencies (except git)
- backup with the help of git, every interval
- inspired by [obsidian-git](https://github.com/denolehov/obsidian-git) and its automatic backup system
- alternative to obsidian git that does not require obsidian and is faster
- sane defaults
- JSON based configuration

## How to use gas

### Installing gas

#### From Source

> Requires:
>
> - go

```bash
git clone https://github.com/xnacly/git-auto-sync gas
cd gas
go build
```

```
./gas # unix
gas.exe # windows
```

#### From release

### Running gas

Prerequisites:

- git needs to be installed, gas will panic if it isn't

  1.  projects needs to be a git repository with a remote set up
  2.  git user needs be authenticated to the remote
  3.  you should be able to run the following commands in your project without issues before using gas in it:

  - `git add -A`
  - `git commit -m "test"`
  - `git push`

  4.  you can now use gas in your project

1. Navigate to the git project you want to backup
2. run `gas` in your terminal

> If you have no `gas.json`, gas will use its default configuration.

### Config path

- On Unix systems, `$XDG_CONFIG_HOME/gas.json` or `$HOME/.config/gas.json`
- On Darwin, `$HOME/Library/Application Support/gas.json`
- On Windows, `%AppData%/gas.json`
- On Plan 9, `$home/lib/gas.json`

### Config options and defaults

If gas can't find its config file (`gas.json`) it will fallback to its default config:

```jsonc
{
  // will be inserted before the local date string in the commit title
  "auto_commit_prefix": "backup: ",

  // specifies the date format which the date will be formatted as
  //
  //  - 2006 for the year, 06 would only be the last two integer
  //  - 01 for the month
  //  - 02 for the day
  //  - 15 for the hour (24-hour format), 05 for 12-hour format
  //  - 04 for the minute
  //  - 05 for the second
  //
  // time formatting in go is weird, see docs:
  //
  // https://www.digitalocean.com/community/tutorials/how-to-use-dates-and-times-in-go
  "commit_title_date_format": "2006-01-02 15:04:05",

  // List filenames affected by the commit in the commit body
  // together with the type of change which happend to the file:
  //
  //      Affected Files:
  //      <filename> <change>
  //
  // possible change types:
  //  - modified
  //  - added
  //  - renamed
  //  - deleted
  //  - copied
  //  - updated but unmerged
  "add_affected_files": true,

  // time interval between backups (in s)
  "backup_interval": 300,

  // commit command, which gas runs after running `git add -A`
  "commit_cmd": "git commit -m",

  // enables debug mode (verbose logging, extra infos, etc.), default: false
  "debug": false,

	// enables pulling the latest changes from remote on start, default: true
  "pull_on_start": true
}
```
