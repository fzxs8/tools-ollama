# Contributing to Ollama Tools

Thank you for your interest in contributing to Ollama Tools! This document provides guidelines and information for contributors.

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- Node.js 18+
- Git
- Wails CLI

### Development Setup

1. **Fork the repository**
   ```bash
   git clone https://github.com/your-username/tools-ollama.git
   cd tools-ollama
   ```

2. **Install dependencies**
   ```bash
   cd frontend
   npm install
   cd ..
   go mod tidy
   ```

3. **Run development server**
   ```bash
   wails dev
   ```

## 📝 Code Style

### Go Code
- Follow standard Go conventions
- Use `gofmt` for formatting
- Add comments for exported functions
- Write tests for new functionality

### Frontend Code
- Use TypeScript
- Follow Vue 3 Composition API patterns
- Use ESLint and Prettier for formatting
- Follow component naming conventions

### Commit Messages
We use [Conventional Commits](https://conventionalcommits.org/):

```
type(scope): description

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes
- `refactor`: Code refactoring
- `test`: Adding tests
- `chore`: Maintenance tasks

Examples:
```
feat(chat): add message history persistence
fix(model): resolve model loading timeout issue
docs: update installation instructions
```

## 🐛 Bug Reports

When reporting bugs, please include:

1. **Environment information**
   - OS and version
   - Go version
   - Node.js version
   - Application version

2. **Steps to reproduce**
   - Clear, numbered steps
   - Expected vs actual behavior
   - Screenshots if applicable

3. **Additional context**
   - Error messages
   - Log files
   - Configuration details

## 💡 Feature Requests

For feature requests, please provide:

1. **Problem description**
   - What problem does this solve?
   - Who would benefit from this feature?

2. **Proposed solution**
   - Detailed description of the feature
   - How should it work?
   - Any alternatives considered?

3. **Additional context**
   - Mockups or examples
   - Related issues or discussions

## 🔄 Pull Request Process

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**
   - Write clean, documented code
   - Add tests if applicable
   - Update documentation

3. **Test your changes**
   ```bash
   # Run tests
   go test ./...
   
   # Test frontend
   cd frontend
   npm run test
   
   # Build application
   wails build
   ```

4. **Submit pull request**
   - Use a clear, descriptive title
   - Reference related issues
   - Provide detailed description of changes

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Tests added/updated as needed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)
- [ ] Commit messages follow conventional format

## 🧪 Testing

### Backend Tests
```bash
go test ./...
```

### Frontend Tests
```bash
cd frontend
npm run test
```

### Integration Tests
```bash
wails build
# Test the built application
```

## 📚 Documentation

When contributing documentation:

1. **Use clear, concise language**
2. **Include code examples**
3. **Update table of contents**
4. **Test all instructions**

Documentation locations:
- `README.md` - Main project documentation
- `docs/` - Detailed documentation
- Code comments - Inline documentation

## 🏗️ Architecture Guidelines

### Backend Structure
```
├── main.go              # Application entry
├── app.go               # Main app logic
├── model_manager.go     # Model management
├── chat_manager.go      # Chat functionality
├── ollama_config.go     # Configuration
└── types/               # Type definitions
```

### Frontend Structure
```
frontend/src/
├── views/               # Page components
├── components/          # Reusable components
├── stores/              # State management
├── router/              # Routing configuration
└── assets/              # Static assets
```

### State Management
- Use Pinia for Vue.js state management
- Keep state minimal and focused
- Use composables for shared logic

## 🔒 Security

- Never commit sensitive information
- Use environment variables for configuration
- Follow security best practices
- Report security issues privately

## 📞 Getting Help

- **GitHub Discussions**: General questions and ideas
- **GitHub Issues**: Bug reports and feature requests
- **Code Review**: Comments on pull requests

## 🎉 Recognition

Contributors will be recognized in:
- README.md contributors section
- Release notes
- Project documentation

Thank you for contributing to Ollama Tools! 🚀