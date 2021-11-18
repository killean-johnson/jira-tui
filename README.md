* ROADMAP
    - Backlog view
    - Ability to add cards
    - Ability to view cards
    - Ability to edit information on a card


* Ideaboard
    - Work it like a state machine (Carry previous state and current state)
        * main view state
        * card view state
        * add ticket dialogue state
        * sort settings state (Pop-up box?)
    - Sorting is a small popup that you can use j and k to walk up and down and enter to select the sort option
    - Config file for all auth setup and user preferences
    - Cli
    - Cache card/ticket info and list

* File Structure
    - Main Folder
        - Jira API Helpers Folder
        - TUI Helpers Folder
        - State Interaction Folder (Main view, ticket view, etc)
        - Cache Folder
        * Config File
