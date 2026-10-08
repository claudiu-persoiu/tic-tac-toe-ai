# Tic-Tac-Toe with AI

This is a simple implementation of a Tic-Tac-Toe game driven by an AI.

This project uses Ollama to run the models locally.

Possible models:

    nimble:9b
    tev1:4b

## Instructions

1. Make sure you have Ollama installed and set up on your machine. You can find the installation instructions on the [Ollama website](https://ollama.com/).
2. Pull the model you want to use with Ollama. For example, to pull the nimble:9b model, run the following command in your terminal:

```bash
ollama pull nimble:9b
```
3. Clone this repository to your local machine.
4. Navigate to the project directory.
5. Run the Tic-Tac-Toe game using the following command:

```bash
go run main.go
```

At runtime, you can specify a model and additional prompt instructions.