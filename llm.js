const Anthropic = require('@anthropic-ai/sdk');

async function main() {
    const client = new Anthropic({
    apiKey: process.env['ANTHROPIC_API_KEY'], 
  });

  const stream = await client.messages.create({
    max_tokens: 1024,
    messages: [{ role: 'user', content: 'Hello, Claude' }],
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

main();