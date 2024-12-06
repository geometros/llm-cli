const readline = require('readline');
const Anthropic = require('@anthropic-ai/sdk');

const arg = process.argv.slice(2).join(' ');

if (arg) {
  const Anthropic = require('@anthropic-ai/sdk');
  main(arg);
} else {
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
      await main(input);
      rl.prompt();
      inputLoop();
    }});
  }
}

async function main(userInput) {
  const client = new Anthropic({
  apiKey: process.env['ANTHROPIC_API_KEY'], 
});

const stream = await client.messages.create({
  max_tokens: 1024,
  messages: [{ role: 'user', content: userInput}],
  model: 'claude-3-opus-20240229',
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
}

