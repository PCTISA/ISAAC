# Build-A-Bot
Boiler-plate code for developing a "custom" Discord bot. Pre-equipped with a powerful multiplexer, config processor, docker build structure, and more.

## Getting Started
1. Smash that **Use this template** button to create a new repository on your account, or in your organization.
2. Familiarize yourself with the structure of the bot (most things are fairly well commented)... Feel free to create an issue on the [template repo](https://github.com/PulseDevelopmentGroup/Build-A-Bot) for any changes you'd reccomend!
3. Change all the import and module paths in each of the Go files to the name of your project (e.g.: `github.com/PulseDevelopmentGroup/Build-A-Bot` > `github.com/[Your Username/Org]/My-Awesome-Bot`)
4. Start building!
5. _(If you need some inspiration, feel free to check out [0x626f74](https://github.com/PulseDevelopmentGroup/0x626f74), a bot made for [Carson](https://github.com/cs-5) and [Josiah's](https://github.com/JosNun) Discord server)_

### Command Structure

A command in the bot is relativly simple, with 5 required parts:

1. The command type definition:
  
   ```go
    // Example is a command
    type Example struct {
      Command  string
      HelpText string
    }
   ```

   In it's most basic form, this definition must contain the `Command` or the name of the command when used in chat and a simple 1-liner of help text (should you choose to implement a `help` command). Other properties can be added here as-needed by the specific command (See: [`example.go`](https://github.com/PulseDevelopmentGroup/Build-A-Bot/blob/master/command/example.go)).

2. The Init function:
   
   ```go
    // Init is called by the multiplexer before the bot starts to initialize any
    // variables the command needs.
    func (c Example) Init(m *multiplexer.Mux) {
      // Nothing to init
    }
   ```

   Simply put, the init function is automaticlly called before the bot offically starts to setup any important bits the command might need later. An example of this being used would be the `help` command (See the 0x626f74 project's [help command](https://github.com/PulseDevelopmentGroup/0x626f74/blob/master/command/help.go#L24)), which needs to build an array of the help messages of the other commands before being used.

3. The Handle function:

   ```go
     // Handle is called by the multiplexer whenever a user triggers the command.
     func (c Example) Handle(ctx *multiplexer.Context) {
       ctx.ChannelSend("Congradulations! You've run your first command")
     }
   ```

   The handle function is called by the multiplexer whenever a user triggers a command... it's that simple. Provided is the `ctx` struct, which contains pretty much any property  you'll need when handling a command. These properties include session info, arguments, multiplexer info, and a whole lot more. Additionally, two helper functions (`ChannelSend` and `ChannelSendf`) are made available to make sending messages to the channel where the command was called more concise.

4. The HandleHelp function:
   
   ```go
     // HandleHelp is not called by the multiplexer. It is used by the
     // `!help` command (if included) to provide a bigger description of the
     // command's functionality.
     func (c Example) HandleHelp(ctx *multiplexer.Context) {
       ctx.ChannelSend("Much bigger/more detailed command description")
     }
   ```

   Although not used by the multiplexer, this function is still required/expected/_highly_ reccomended because it can be leveraged by the help command to provide more context/help info for a specific command. Think of `HelpText` as the one-liners that show up when you type `help` into a terminal, and `HandleHelp` as the detailed response to `help [command name here]`. HandleHelp works in exactly the same way as `Handle` as far as how the command is to respond.

5. The Settings function:
   
   ```go
    // Settings is called by the multiplexer on startup to process any settings
    // associated with that command.
    func (c Example) Settings() *multiplexer.CommandSettings {
      return &multiplexer.CommandSettings{
        Command:  c.Command,
        HelpText: c.HelpText,
      }
    }
   ```

   This function is "functionally" useless, but is a way for the multiplexer to determine the kind of command it's dealing with. It returns any properties important to the multiplexer used for handling commands ([properties](https://github.com/PulseDevelopmentGroup/Build-A-Bot/blob/master/multiplexer/mux.go#L47)).

### Github Actions (Auto Build)
This repository is setup with Github Actions support to automaticlly build a docker container with the bot's code, and to subsequently publish that container on the registry associated with your repo.

In order to start a build, simply publish a release (Go to "Releases" > "Draft a new release" > Give it a version, name, and branch.) and the actions script will kick off. Note: The tag version specified in the release will be the same version tag assigned to the container. 

To pull your newly created container, simply do:

```
docker pull docker.pkg.github.com/[Your Username/Org]/My-Awesome-Bot/bot:[latest|specific version]
```

_If you're using a container such as Watchtower to auto-update your containers. Pulling the `latest` tag will allow watchtower to automaticlly pull the latest version your release._

### Other Notes
- When adding new (non-code) files that don't _need_ to be in Docker, it's probably a good idea to add them to `.dockerignore`.
- Placing a `.env` file with all your enviornment variables defined in the project root directory will automaticlly get picked up by and used by the bot. This makes development easier.
- A config file can either be loaded by a file path or a URL (both specified in `.env` or in your regular enviorment variables, or in the Docker enviorment variables passed to the container). Whatever makes life easier.
- By default, simple commands are loaded from the config file. A simple command is just a 1-liner string reply when the command is called.
- Specifying permissions is as simple as adding the name of the command (under the `permissions` object in the config file) with an array of role ID's supplied (See 0x626f74's config [here](https://github.com/PulseDevelopmentGroup/0x626f74/blob/master/config.json)). Currently, role ID's are the only supported permission type, but the goal is to change that to also support channel and user ID's.