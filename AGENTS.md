# AGENTS.md

## Build/CI Commands
- `npm install`: Install dependencies
- `npm run build`: Build project
- `npm run lint`: Lint code
- `npm run test`: Run all tests
- `npm run test:watch`: Watch for test changes
- `npm run test:unit`: Run unit tests
- `npm run test:integration`: Run integration tests
- `npm run test:coverage`: Run tests with coverage

## Lint/Format Commands
- `npm run lint`: Run ESLint
- `npm run format`: Format code with Prettier
- `npm run lint:fix`: Fix lint issues
- `npm run format:write`: Write formatted code

## Test Commands
- `npm run test:unit`: Run unit tests
- `npm run test:integration`: Run integration tests
- `npm run test:coverage`: Run tests with coverage
- `npm run test:watch`: Watch for test changes

## Code Style Guidelines
### Formatting
- Use 2 spaces for indentation
- Line length limit: 120 characters
- Semi-colons are required
- Use single quotes unless escaping is needed

### Naming Conventions
- Variables: snake_case
- Functions: snake_case
- Classes: PascalCase
- Constants: SCREAMING_SNAKE_CASE

### Types
- Use TypeScript for type annotations
- Prefer interfaces over types for custom types
- Use type aliases for complex types

### Error Handling
- Use try/catch blocks for error handling
- Prefer throwing errors over returning null/undefined
- Use specific error messages

### Imports
- Use relative paths for local imports
- Use absolute paths for shared imports
- Sort imports alphabetically
- Avoid importing unused modules

### Comments
- Use JSDoc for public APIs
- Comments should explain "why", not "what"
- Keep comments concise and relevant

## Rules
- No Cursor rules found
- No Copilot rules found

## Notes
- This file is used by agentic coding agents to understand the codebase
- Update this file when adding new rules or conventions
- Keep this file in the root of the repository