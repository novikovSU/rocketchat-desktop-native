# Rocket.Chat Desktop native

Try to implement dead simple Rocket.Chat desktop client using Go and GTK3.

First phase of active development in master. VERY buggy and dirty. **Not for any use for now!**

Application has event-oriented architecture. Full Event list:

| event                     | args  | decription                                                                         | 
|---------------------------|-------------|------------------------------------------------------------------------------|
| messages.new              | api.Message | Fires then application received new chat message                             |
| messages.read             | api.Message | Fires then user read the chat message **(Not implemented yet)**              |
| contacts.users.added      | api.User    | Fires then application detects new user has been added to server             |
| contacts.users.removed    | api.User    | Fires then application detects existing user has been removed from server    |
| contacts.channels.added   | api.Channel | Fires then application detects new channel has been added to server          |
| contacts.channels.removed | api.Channel | Fires then application detects existing channel has been removed from server |
| contacts.groups.added     | api.Group   | Fires then application detects new group has been added to server            |
| contacts.groups.removed   | api.Group   | Fires then application detects existing group has been removed from server   |
| ui.mainwindow.closed      |             | Fires then user click on main window close button **(Not implemented yet)**  |