# llm-cli

### Command line tool for LLMs

#### reasons
* tired of signing back in to web
* you have too many browser tabs rn
* text is right there in stdout
* convenient one off questions & system prompt
* I want to learn Go
___
Currently does nothing but take a quoted string as an argument and return a response from Claude via Anthropic API.

#### To Do: 

~~* Add streaming so you don't have to wait on longer replies~~
* Add an interactive mode to allow conversational threading
* Merge streaming branch back into main with both modes available via flags (assuming non-streaming mode will be better for piping, redirects?\)
* Allow saving threads
* Flags for setting max tokens, system prompt
* Support for other models
