# x-git-sync

Backup your repository at configured intervals

https://user-images.githubusercontent.com/47723417/213995030-a72ab64f-3e64-403e-bda7-57279b37780d.mp4

## Why use xgs

- highly configurable
- 0 external dependencies (except git)
- backup with the help of git, every interval
- inspired by [obsidian-git](https://github.com/denolehov/obsidian-git) and its automatic backup system
- alternative to obsidian git that does not require obsidian and is faster
- sane defaults
- JSON based configuration

### Why use this project and not something else

XGS is a lot more minimal and configurable than [git-auto-sync](https://github.com/GitJournal/git-auto-sync),
doesn't require obsidian or VScode to work ([Obsidian Git](https://github.com/denolehov/obsidian-git), [VS Code GitDoc](https://marketplace.visualstudio.com/items?itemName=vsls-contrib.gitdoc)) and
isn't as complicated or unintelligible as [Git Annex](https://git-annex.branchable.com/) or [Git Sync](https://github.com/simonthum/git-sync).

Alternatives:

- [git-auto-sync](https://github.com/GitJournal/git-auto-sync)
- [Obsidian Git](https://github.com/denolehov/obsidian-git)
- [VS Code GitDoc](https://marketplace.visualstudio.com/items?itemName=vsls-contrib.gitdoc)
- [Git Annex](https://git-annex.branchable.com/)
- [Git Sync](https://github.com/simonthum/git-sync)

## How to use xgs

### Installing xgs

#### From Source

> Requires:
>
> - go

```bash
git clone https://github.com/xnacly/x-git-sync xgs
cd xgs
go build
```

```
./xgs # unix
xgs.exe # windows
```

#### From release (unix)

- download executable from latest release
- move the `xgs`-executable to a directory in the path, for linux: `mv ./xgs /usr/bin` (this might require elevated privileges)

### Running xgs

Prerequisites:

- git needs to be installed, xgs will panic if it isn't

  1.  projects needs to be a git repository with a remote set up
  2.  git user needs be authenticated to the remote
  3.  you should be able to run the following commands in your project without issues before using xgs in it:

  - `git add -A`
  - `git commit -m "test"`
  - `git push`

  4.  you can now use xgs in your project

1. Navigate to the git project you want to backup
2. run `xgs` in your terminal

> If you have no `xgs.json`, xgs will use its default configuration.

### Config path

- On Unix systems, `$XDG_CONFIG_HOME/xgs.json` or `$HOME/.config/xgs.json`
- On Darwin, `$HOME/Library/Application Support/xgs.json`
- On Windows, `%AppData%/xgs.json`
- On Plan 9, `$home/lib/xgs.json`

### Config options and defaults

If xgs can't find its config file (`xgs.json`) it will fallback to its default config:

```json
{
  // will be inserted into the %title% placeholder in the commit_format string
  "commit_title": "backup",

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
  "commit_date": "2006-01-02 15:04:05",

  // specifies the format of the commit, currently supports:
  // - commit_title: %title%
  // - commit_date: %date%
  "commit_format": "%title% %date%",

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

  // commit command, which xgs runs after running `git add -A`
  "commit_cmd": "git commit -m",

  // enables debug mode (verbose logging, extra infos, etc.), default: false
  "debug": false,

  // enables pulling the latest changes from remote on start, default: true
  "pull_on_start": true
}
```
