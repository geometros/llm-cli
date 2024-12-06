const readline = require('readline');
const Anthropic = require('@anthropic-ai/sdk');

const arg = process.argv.slice(2).join(' ');

if (arg) {
  const Anthropic = require('@anthropic-ai/sdk');
  main(arg);
} else {
  const chatHistory = [];

  const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout,
    prompt: '> '
  });
  
  inputLoop()

  function inputLoop() {
    input = rl.question('> ', async (input) => {if (input === 'quit' || input == 'exit'){
      rl.close();
      return;
    } else {

      chatHistory.push({
        role: "user",
        content: input
      })
      
      const assistantContent = await main(input,chatHistory);
      
      chatHistory.push({
        role: "assistant",
        content: assistantContent
      })
      
      rl.prompt();
      
      inputLoop();
    }});
  }
}

async function main(userInput, chatHistory = []) {
  const client = new Anthropic({
  apiKey: process.env['ANTHROPIC_API_KEY'], 
});

const stream = await client.messages.create({
  max_tokens: 1024,
  messages: [
    ...chatHistory,
    { role: 'user', content: userInput}
  ],
  model: 'claude-3-5-sonnet-latest',
  stream: true,
});

let replyStream = '';

for await (const event of stream) {
  if (event.type === 'content_block_delta') {
    replyStream += event.delta.text
    process.stdout.write(event.delta.text)
  }
}
process.stdout.write('\n') //newline to prevent shell weirdness

return replyStream;
}

