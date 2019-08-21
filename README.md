# Rocket.Chat Desktop native

Try to implement dead simple Rocket.Chat desktop client using Go and GTK3.

First phase of active development in master. VERY buggy and dirty. **Not for any use for now!**

Application has event-oriented architecture. First part in event name is package name, when event is fired. Full Event list:

| event                         | args                                              | decription                                                                   | 
|-------------------------------|---------------------------------------------------|------------------------------------------------------------------------------|
| rocket.messages.new           | api.Message                                       | Fires then new chat message is received by rocket package                    |
| model.messages.received       | model.ChatModel, string, api.Message              | Fires then new chat message is received by model package                     |
| model.unread_counters.updated | model.ChatModel, string                           | Fires then unread counter for model updated (set or cleared)                 |
| rocket.messages.new           | api.Message                                       | Fires then application received new chat message **(Not implemented yet)**   |
| rocket.users.load             | []api.User                                        | Fires then users is loaded by rocket package                                 |
| rocket.channels.load          | []api.Channel                                     | Fires then channels is loaded by rocket package                              |
| rocket.groups.load            | []api.Group                                       | Fires then groups is loaded by rocket package                                |
| messages.read                 | api.Message                                       | Fires then user read the chat message **(Not implemented yet)**              |
| contacts.update.started       |                                                   | Fires then application starts to load/update contact list                    |
| contacts.update.finished      |                                                   | Fires then application finish to load/update contact list                    |
| model.user.added              | model.ChatModel, model.UserModel                  | Fires then application detects new user has been added to server             |
| model.user.removed            | model.ChatModel, model.UserModel                  | Fires then application detects existing user has been removed from server    |
| model.channel.added           | model.ChatModel, model.ChannelModel               | Fires then application detects new channel has been added to server          |
| model.channel.removed         | model.ChatModel, model.ChannelModel               | Fires then application detects existing channel has been removed from server |
| model.group.added             | model.ChatModel, model.GroupModel                 | Fires then application detects new group has been added to server            |
| model.group.removed           | model.ChatModel, model.GroupModel                 | Fires then application detects existing group has been removed from server   |
| ui.mainwindow.closed          |                                                   | Fires then user click on main window close button **(Not implemented yet)**  |