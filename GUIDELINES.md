# Development Guidelines for Gemini CLI Agent

This document outlines the recommended workflow for the Gemini CLI Agent when undertaking software engineering tasks. Adhering to these guidelines ensures a structured, robust, and reliable approach to development.

## Workflow Steps:

1.  **Understand and Plan (Pre-computation)**:
    *   Thoroughly analyze the user's request and the existing codebase. Utilize `read_file`, `read_many_files`, `search_file_content`, and `glob` tools to gather all necessary context, understand existing conventions, and identify relevant areas for modification.
    *   Formulate a clear, concise, and grounded plan for addressing the task. This plan should detail the intended changes, affected files, and the sequence of operations.
    *   **Crucially, identify or establish automated testing mechanisms.** If the feature lacks tests, or if the changes might impact existing functionality, prioritize writing new tests or extending existing ones to create a safety net. This includes unit tests, integration tests, or any other relevant automated checks.

2.  **Implement and Iterate (Build & Verify Loop)**:
    *   Execute the planned changes using `write_file`, `replace`, or `run_shell_command`.
    *   **Immediately after making changes, run the relevant automated tests.** This includes any newly created tests for the feature, as well as existing project tests (e.g., `go test`, `npm test`, `pytest`) to detect regressions.
    *   **Continuously check for build, linting, and type-checking errors** (e.g., `go build`, `ruff check`, `tsc`).
    *   If tests fail or errors occur, diagnose and fix the issues. Repeat the build and verify loop until all tests pass and no errors are present.
    *   This iterative process ensures that issues are caught early and that the codebase remains stable.

3.  **Final Verification and Completion**:
    *   Once all automated tests pass and the code adheres to project standards, perform a final review of the changes against the original request to ensure all requirements have been met.
    *   Confirm that no existing functionality has been inadvertently broken.
    *   Only consider the task complete when all tests pass, the code is clean, and the feature is fully implemented as requested.

## Core Principles:

*   **Test-Driven (where applicable):** Prioritize writing tests before or alongside code changes to guide development and ensure correctness.
*   **No Regressions:** Ensure that new changes do not break existing functionality. Automated tests are paramount for this.
*   **Iterative Development:** Build and test in small, manageable steps to quickly identify and resolve issues.
*   **Adherence to Conventions:** Always mimic the style, structure, and architectural patterns of existing code in the project.
*   **Safety First:** Be mindful of potential side effects and always explain critical commands before execution.
