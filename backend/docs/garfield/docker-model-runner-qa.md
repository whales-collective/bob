# Docker Model Runner FAQ - 30 Common Questions and Answers

## 1. What is Docker Model Runner?

Docker Model Runner is a plugin that allows you to pull, run, and manage AI models directly from Docker Hub. It enables you to interact with models through the command line or Docker Desktop Dashboard using OpenAI-compatible APIs.

## 2. How do I enable Docker Model Runner?

Navigate to Settings > Features in development > Beta tab in Docker Desktop, check "Enable Docker Model Runner", apply and restart. You may need to first enable experimental features in the Experimental features tab.

## 3. What models are available with Docker Model Runner?

All available models are hosted in the public Docker Hub namespace at https://hub.docker.com/u/ai. Popular models include ai/smollm2 and other AI models optimized for local execution.

## 4. How do I check if Docker Model Runner is active?

Use the command `docker model status` to check whether the Docker Model Runner is currently active and running.

## 5. How do I pull a model from Docker Hub?

Use `docker model pull <model>` to download a model. For example: `docker model pull ai/smollm2`. The model will be cached locally for faster future access.

## 6. How do I list all locally available models?

Use `docker model list` to see all models currently pulled to your local environment, including their parameters, quantization, architecture, and size information.

## 7. How do I run a model with a one-time prompt?

Use `docker model run <model> "your prompt"`. For example: `docker model run ai/smollm2 "Hi"` will run the model once with the specified prompt.

## 8. How do I start an interactive chat session with a model?

Use `docker model run <model>` without a prompt to start interactive chat mode. Type your messages and use `/bye` to exit the chat session.

## 9. How do I remove a model from my local system?

Use `docker model rm <model>` to remove a downloaded model from your system. This will free up disk space but you'll need to pull the model again to use it.

## 10. How do I push a model to Docker Hub?

Use `docker model push <namespace>/<model>` to push your model to Docker Hub. You need appropriate permissions and authentication to push to the specified namespace.

## 11. How do I tag a model with a specific version?

Use `docker model tag` to specify a particular version or variant of the model. If no tag is provided, Docker defaults to "latest".

## 12. How do I view Docker Model Runner logs?

Use `docker model logs` to fetch logs for monitoring activity or debugging issues. Add `-f` or `--follow` for real-time streaming, or `--no-engines` to exclude inference engine logs.

## 13. What API endpoints are available?

Docker Model Runner provides OpenAI-compatible endpoints including `/v1/chat/completions`, `/v1/completions`, `/v1/embeddings`, and model management endpoints at `http://model-runner.docker.internal/`.

## 14. How do I call the API from within a container?

Use `http://model-runner.docker.internal/engines/llama.cpp/v1/chat/completions` as the endpoint URL when making API calls from within other containers.

## 15. How do I call the API from the host using Unix socket?

Use `curl --unix-socket $HOME/.docker/run/docker.sock localhost/exp/vDD4.40/engines/llama.cpp/v1/chat/completions` to make API calls from the host.

## 16. How do I enable TCP access for the API?

Use `docker desktop enable model-runner --tcp <port>` to enable host-side TCP support, then access the API at `http://localhost:<port>/engines/llama.cpp/v1/chat/completions`.

## 17. Where are models stored locally?

Models are pulled from Docker Hub and stored locally in Docker's cache. They're loaded into memory only at runtime when requests are made and unloaded when not in use to optimize resources.

## 18. How long does it take to pull a model?

The initial pull may take some time since models can be large (hundreds of MB to GB). However, after the first pull, models are cached locally for much faster access.

## 19. Can I use Docker Model Runner with Testcontainers?

Yes, Testcontainers for Java and Go now support Docker Model Runner, allowing you to integrate AI models into your testing workflows.

## 20. Can I use Docker Model Runner with Docker Compose?

Yes, Docker Compose now supports Docker Model Runner, enabling you to define AI model services alongside your application services in compose files.

## 21. How do I integrate Docker Model Runner into my development workflow?

You can clone the hello-genai repository (`git clone https://github.com/docker/hello-genai.git`) and run the sample application to see how to build GenAI applications powered by Docker Model Runner.

## 22. What should I do if `docker model` is not recognized?

Create a symlink: `ln -s /Applications/Docker.app/Contents/Resources/cli-plugins/docker-model ~/.docker/cli-plugins/docker-model` so Docker can detect the plugin.

## 23. What happens if I try to run a model that's too large for my system?

Currently, there are no safeguards to prevent running oversized models. This may result in severe slowdowns or render your system temporarily unusable, especially with insufficient GPU memory or RAM.

## 24. Can I specify models by digest instead of name?

Currently, Docker Model CLI lacks consistent support for specifying models by image digest. As a workaround, refer to models by name instead of digest.

## 25. How do I view available Docker Model Runner commands?

Use `docker model help` to display help information and a list of available subcommands including list, pull, rm, run, status, and version.

## 26. Can I interact with models through Docker Desktop GUI?

Yes, you can view and interact with your local models in the Models tab in the Docker Desktop Dashboard, including using chat mode directly in the interface.

## 27. What programming languages can I use with the API?

You can use any programming language that supports HTTP requests to interact with the OpenAI-compatible API endpoints, including Python, JavaScript, Java, Go, and others.

## 28. How do I optimize resource usage with Docker Model Runner?

Models are automatically loaded into memory only when requests are made and unloaded when not in use. This optimizes resource usage, but ensure your system has sufficient RAM for the models you want to run.

## 29. How do I provide feedback or report bugs?

Use the "Give feedback" link next to the "Enable Docker Model Runner" setting in Docker Desktop to report bugs or provide feedback about the feature.

## 30. How do I disable Docker Model Runner?

Go to Settings > Features in development > Beta tab in Docker Desktop, clear the "Enable Docker Model Runner" checkbox, and select "Apply & restart" to disable the feature.
