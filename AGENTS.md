# AI Agent Working Guidelines

## Code Style & Refactoring
- **Commit Granularity**: Commits should be granular and logical. Avoid massive commits that include DB changes, logic, and tests all at once. Split them into:
  - Phase 1: DB Schema & Models
  - Phase 2: Repository / Infrastructure
  - Phase 3: Domain / UseCase Logic
  - Phase 4: Interface / Handler / Routes
  - Phase 5: Tests
- **Variable Cleanliness**: Remove unused variables and imports immediately. Linting tools must report zero errors.
- **Refactoring**: When refactoring, always update the corresponding tests. Ensure 100% pass rate before finishing.

## Testing
- **Coverage**: Aim for >80% coverage for UseCase and Handler layers.
- **Mocking**: Use `go.uber.org/mock/gomock` for Unit Tests.
- **Table Driven**: Use Table Driven Tests for readability and maintainability.

## Architecture
- **Clean Architecture**: Strictly follow the Hexagonal/Clean Architecture layers:
  - `domain/model`: Pure business logic structs (No external deps)
  - `domain/repository`: Interfaces (No implementation)
  - `usecase`: Application business rules (Depends on domain)
  - `infrastructure`: Database, External APIs (Depends on domain repository interfaces)
  - `server/handler`: HTTP Transport (Depends on usecase)

## Workflow / Standard Operating Procedure
Strictly follow these 5 phases for every task:

1. **Plan**: Create a work plan and get approval.
2. **Worktree**: Create a git worktree for the task.
3. **Execute**: Perform the work.
4. **Review**: Confirm the work with the user.
5. **Merge**: Merge the work into the base branch.

### Execution Guidelines
- **Task Boundaries**: Use `task_boundary` tool frequently to update progress.
- **Verification**: Run tests (`go test ./...`) after every logical change set.
