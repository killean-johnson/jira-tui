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

* TODO
    - Learn Go
    - Learn Bubble Tea and Bubbles OR _gocui_, using lazygit and other examples to grow from
    - Use go-jira for jira access and management

* Keybindings
    - overall keybinds
        - Ctrl-c exit

    - issuelist keybinds
        - j/k -- up and down
        - enter -- select issue (load issue into issueinfo, then move view to issueinfo)
        - s -- open sortmenu
        - a -- add issue
        - m -- open movemenu
        - esc -- exit

    - sortmenu keybinds
        - j/k -- up and down
        - enter -- confirm the sorting method, and sort the issues accordingly. transfer view back to issuelist
        - esc -- return to issuelist without doing anything

    - movemenu keybinds
        - j/k -- up and down
        - enter -- moves the issue to the proper stage. transfer view back to issuelist
        - esc -- return to issuelist without doing anything

    - issueinfo keybinds
        - j/k -- up and down
        - esc -- return to issuelist without doing anything

    - addissue keybinds
        - arrowkeys -- navigate the setup file
        - Ctrl-a -- add the new issue, returns to the updated issuelist
        - esc -- go back to the issuelist without making any changes

