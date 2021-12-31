# Jira-Tui
## Description
Terminal Client for Jira. Can be used for basic manipulation of issues on the active sprint.

## Setup
When jira-tui is ran for the first time, it will create a config.json file un %HOME%/.config/jira-tui/
You'll need to add your email/username and a [Jira API Token](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/), as well as the project URL.
Once that's done, it's good to go.

## Features
### Select Project
### View Issue Info
### Change Issue Assignee
### Change Issue Description
### Change Issue Status
### Create New Issue On Active Sprint

## Config
All of the keybindings can be modified for any of the views. The config is under %HOME%/.config/jira-tui/config.json

## TODO
	- Helpbar is not actively updating when doing any operation until the TUI Layout func runs again.
		This is due to how the layout functions get updated in order. Find a GOOD way around this?

## ROADMAP
    - Backlog view
    - Sorting options for issue list
    - Caching board, as well as issues for faster startup
    - Ability to change board
    - Another keybind for changing the issue summary
    - Add primary/secondary boolean to the keymaps, where only primary are shown in help bar, and a help button can be hit to show all commands
    - ~~Rewrite with each view section being a layout with it's own handlers~~
