# DocProcessor

- **GitHub URL**: <https://github.com/HelixDevelopment/DocProcessor>
- **Description**: Documentation processing and feature map extraction engine for QA automation. Parses project documentation, extracts structured feature inventories, and generates feature maps that drive automated test planning and coverage analysis. Bridges the gap between written specifications and executable test strategies.
- **Category**: HelixDevelopment
- **Status**: Active

## Capabilities

- Multi-format document parsing (Markdown, HTML, PDF, DOCX, plain text)
- Feature map extraction -- identifies testable features, acceptance criteria, and edge cases from prose
- Structured feature inventory generation with dependency graphs
- Coverage gap analysis -- compares extracted features against existing test suites
- Auto-generation of test case skeletons from extracted acceptance criteria
- Document diff tracking -- detects specification changes between versions
- Integration with QA pipelines for specification-driven test planning
- Batch processing of documentation trees with incremental update support

## Technology

- **Language**: Go
- **Frameworks**: Go standard library, document parsing libraries
- **Architecture**: Pipeline-based processing with pluggable format parsers and extractors
- **Key patterns**: AST-based document analysis, structured output schemas

## Integration

- Consumed by helixqa for specification-driven QA orchestration and test planning
- Feeds extracted feature maps into HelixSpecifier for spec-driven development workflows
- Integrates with HelixAgent for documentation-aware reasoning and context injection
- Used by HelixBuilder to extract requirements from project documentation during build planning
- Outputs feed the workable-items system for tracking specification coverage

## Status

Active development. Core document parsing and feature extraction pipelines are operational. Supports Markdown and HTML natively with PDF/DOCX via conversion layers. Feature map extraction accuracy is continuously improved through corpus-based training.
