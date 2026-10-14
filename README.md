# How to run?
## Setup folder
mkdir code-editing-agent
cd code-editing-agent
go mod init agent
touch main.go

## Download the dependencies
go mod tidy
# Run it
go run main.go

# What defines an Agent?
agent = llm + harness

an LLM with access to tools, giving it the ability to modify something outside the context window.

# What defines a tool?
Each tool we’re going to add will require the following:

A name
A description to tell the model what the tool does, when to use it, when to not use it, what it returns and so on
An input schema that describes, as a JSON schema, what inputs this tool expects and in which form
A function that actually executes the tool with the input the model sends to us and returns the result

Under the hood, this will all end up as strings that are sent to the model.

# Milestones
## V1: basic harness with 3 tools (read, list, edit files) [DONE]
## V2: select models to use (ollama)
## V3: apply 6 layer harness
## V4: apply deep eval

# Demonstration
## Edit tool use case
PS D:\code-editing-agent> go run main.go                                                                                                                                                                                  
Chat with gpt-oss:120b-cloud (use ctrl + c to quit)                                                                                                                                                                       
You: hey can you create fizzbuzz.js that I can run with Nodejs and that has fizzbuzz in it and executes it                                                                                                                
tool:edit_file({"new_str":"// fizzbuzz.js\n// This script prints the FizzBuzz sequence from 1 to 100.\n// Run with: node fizzbuzz.js\n\nfor (let i = 1; i \u003c= 100; i++) {\n  let output = '';\n  if (i % 3 === 0) output += 'Fizz';\n  if (i % 5 === 0) output += 'Buzz';\n  console.log(output || i);\n}\n","old_str":"","path":"fizzbuzz.js"})                                                                                                
GPT: The **fizzbuzz.js** file has been created. You can run it with Node.js:                                                                                                                                              
                                                                                                                                                                                                                          
```bash                                                                                                                                                                                                                   
node fizzbuzz.js                                                                                                                                                                                                          
```                                                                                                                                                                                                                       
                                                                                                                                                                                                                          
It will output the FizzBuzz sequence from 1 to 100. Enjoy!                                                                                                                

## Read file tool use case demonstration
PS D:\code-editing-agent> go run main.go
Chat with gpt-oss:120b-cloud (use ctrl + c to quit)
You: help me solve the riddle in the secret.txt 
tool:read_file({"path":"secret.txt"})
GPT: **Answer:** The animal is a **horse**.  

**Why?** The riddle says the creature “always says ‘neigh’,” which is the sound a horse makes. And “disagreeable” is a playful way of saying “neigh‑saying,” i.e., it’s always refusing or disagreeing, which sounds like “neigh‑saying.” So the answer to the riddle is a horse.
You:


## List files tool use case demonstration
PS D:\code-editing-agent> go run main.go
Chat with gpt-oss:120b-cloud (use ctrl + c to quit)
You: what do you see in this directory
tool:list_files({"path":""})
GPT: Here’s a listing of everything in the repository’s top‑level directory:

```
.git/                         (Git metadata)
README.md                     – Project overview / documentation
go.mod                        – Go module definition
go.sum                        – Checksums for module dependencies
main.go                       – Main Go source file for the application
secret.txt                    – A text file (contents are hidden until opened)
llmresponse.md                – Markdown file (likely containing a prior LLM response)
```

The `.git/` folder contains all of the usual Git objects, refs, logs, hooks, etc., but the only source‑code and documentation files at the root are the ones listed above. If you’d like to view the contents of any of these files (e.g., `README.md`, `main.go`, or `secret.txt`), just let me know!
You:

## Read file + List file use case demonstration
PS D:\code-editing-agent> go run main.go
Chat with gpt-oss:120b-cloud (use ctrl + c to quit)
You: Tell me about all the Go files in here. Be brief!
tool:list_files({"path":""})
tool:list_files({"path":"."})
tool:read_file({"path":"main.go"})
GPT: **Go files in this repository**

| File | Brief description |
|------|-------------------|
| **`main.go`** | The only Go source file. It implements a simple command‑line chat client that talks to an Ollama model (`gpt-oss:120b-cloud`). Key parts: <br> • **`ToolDefinition`** – describes a tool (name, description, JSON schema, Go function). <br> • **`ReadFile`** and **`ListFiles`** – concrete tool implementations for reading a file’s contents and recursively listing directory entries. <br> • **`GenerateSchema`** – uses `invopop/jsonschema` to turn Go types into JSON‑Schema parameters for Ollama tools. <br> • **`Agent`** – holds the Ollama client, a user‑input callback, and the available tools. Its **`Run`** method manages the conversation loop, sends user messages to the model, handles tool calls, and prints responses. <br> • **`main`** – wires everything together: creates the Ollama client, defines the tools, builds the `Agent`, and starts the chat. |

No other `.go` files are present.
You: 

## Edit tool use case demonstration
You: Create a congrats.js script that rot13-decodes the following string 'Pbatenghyngvbaf ba ohvyqvat n pbqr-rqvgvat ntrag!' and prints it
tool:edit_file({"new_str":"// congrats.js\n// This script decodes a ROT13-encoded message and prints it.\n// Run with: node congrats.js\n\nfunction rot13(str) {\n  return str.replace(/[A-Za-z]/g, function (c) {\n    return String.fromCharCode(\n      (c \u003c= 'Z' ? 90 : 122) \u003e= (c = c.charCodeAt(0) + 13)\n        ? c\n        : c - 26\n    );\n  });\n}\n\nconst encoded = 'Pbatenghyngvbaf ba ohvyqvat n pbqr-rqvgvat ntrag!';\nconst decoded = rot13(encoded);\nconsole.log(decoded);\n","old_str":"","path":"congrats.js"})
GPT: A new **`congrats.js`** script has been created. It decodes the ROT13‑encoded message and prints the result.

```javascript
// congrats.js
// This script decodes a ROT13-encoded message and prints it.
// Run with: node congrats.js

function rot13(str) {
  return str.replace(/[A-Za-z]/g, function (c) {
    return String.fromCharCode(
      (c <= 'Z' ? 90 : 122) >= (c = c.charCodeAt(0) + 13)
        ? c
        : c - 26
    );
  });
}

const encoded = 'Pbatenghyngvbaf ba ohvyqvat n pbqr-rqvgvat ntrag!';
const decoded = rot13(encoded);
console.log(decoded);
```

Run it in your terminal:

```bash
node congrats.js
```

It will output the decoded message. Enjoy!