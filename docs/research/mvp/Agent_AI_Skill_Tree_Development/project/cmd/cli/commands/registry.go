package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/helixdevelopment/skill-system/internal/models"
	"github.com/spf13/cobra"
)

// NewRegistryCommand creates the registry management command group
func NewRegistryCommand() *cobra.Command {
	var domain string

	cmd := &cobra.Command{
		Use:   "registry",
		Short: "Skill registry management",
		Long:  `Monitor registry health, review skill completeness, and check coverage.`,
	}

	// registry status
	statusCmd := &cobra.Command{
		Use:     "status",
		Short:   "Show registry health status",
		Long:    `Display overall registry health including skill counts, coverage, and issues.`,
		Example: `  skill-system registry status`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegistryStatus(cmd)
		},
	}

	// registry missing
	missingCmd := &cobra.Command{
		Use:     "missing",
		Short:   "List missing dependencies",
		Long:    `Show all skills that have unresolved or missing dependencies.`,
		Example: `  skill-system registry missing`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegistryMissing(cmd)
		},
	}

	// registry stale
	staleCmd := &cobra.Command{
		Use:     "stale",
		Short:   "List stale skills",
		Long:    `Show skills that haven't been reviewed recently and may need attention.`,
		Example: `  skill-system registry stale`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegistryStale(cmd)
		},
	}

	// registry review <name>
	reviewCmd := &cobra.Command{
		Use:     "review <name>",
		Short:   "Review a skill's registry entry",
		Long:    `Show detailed registry information for a specific skill including coverage and issues.`,
		Example: `  skill-system registry review go-concurrency`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegistryReview(cmd, args[0])
		},
	}

	// registry coverage [--domain]
	coverageCmd := &cobra.Command{
		Use:   "coverage",
		Short: "Show coverage statistics",
		Long:  `Display coverage metrics for skills, optionally filtered by domain.`,
		Example: `  skill-system registry coverage
  skill-system registry coverage --domain backend`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRegistryCoverage(cmd, domain)
		},
	}
	coverageCmd.Flags().StringVar(&domain, "domain", "", "Filter by domain")

	cmd.AddCommand(statusCmd, missingCmd, staleCmd, reviewCmd, coverageCmd)
	return cmd
}

func runRegistryStatus(cmd *cobra.Command) error {
	client := getAPIClient(cmd)
	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	resp, err := client.Request(ctx, http.MethodGet, "/api/v1/registry/status", nil)
	if err != nil {
		return fmt.Errorf("get registry status: %w", err)
	}
	defer resp.Body.Close()

	var status struct {
		TotalSkills int            `json:"total_skills"`
		TotalDeps   int            `json:"total_dependencies"`
		MissingDeps int            `json:"missing_dependencies"`
		StaleSkills int            `json:"stale_skills"`
		Coverage    float64        `json:"average_coverage"`
		ByStatus    map[string]int `json:"by_status"`
		ByDomain    map[string]int `json:"by_domain"`
		Health      string         `json:"health"`
		LastScan    *time.Time     `json:"last_scan,omitempty"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		// Fallback: build status from entries
		return printSimpleRegistryStatus(ctx, client)
	}

	// Print formatted status
	healthColor := colorGreen
	switch status.Health {
	case "critical":
		healthColor = colorRed
	case "warning":
		healthColor = colorYellow
	}

	fmt.Println("Registry Health Status")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Overall Health: %s%s%s\n", healthColor, status.Health, colorReset)
	fmt.Printf("Total Skills:   %d\n", status.TotalSkills)
	fmt.Printf("Dependencies:   %d total\n", status.TotalDeps)
	fmt.Printf("Missing Deps:   %s%d%s\n", colorRed, status.MissingDeps, colorReset)
	fmt.Printf("Stale Skills:   %s%d%s\n", colorYellow, status.StaleSkills, colorReset)
	fmt.Printf("Avg Coverage:   %.1f%%\n", status.Coverage*100)
	if status.LastScan != nil {
		fmt.Printf("Last Scan:      %s\n", status.LastScan.Format(time.RFC3339))
	}

	if len(status.ByStatus) > 0 {
		fmt.Printf("\nSkills by Status:\n")
		for s, c := range status.ByStatus {
			fmt.Printf("  %-12s %d\n", s, c)
		}
	}

	if len(status.ByDomain) > 0 {
		fmt.Printf("\nSkills by Domain:\n")
		for d, c := range status.ByDomain {
			fmt.Printf("  %-12s %d\n", d, c)
		}
	}

	return nil
}

func printSimpleRegistryStatus(ctx context.Context, client *APIClient) error {
	// Get all registry entries and compute stats
	resp, err := client.Request(ctx, http.MethodGet, "/api/v1/registry", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var entries []models.SkillRegistryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return fmt.Errorf("decode entries: %w", err)
	}

	total := len(entries)
	missing := 0
	stale := 0
	var totalCoverage float64

	for _, e := range entries {
		if len(e.MissingDeps) > 0 {
			missing++
		}
		if e.Stale {
			stale++
		}
		totalCoverage += e.Coverage
	}

	avgCoverage := 0.0
	if total > 0 {
		avgCoverage = totalCoverage / float64(total)
	}

	health := "healthy"
	if missing > total/4 || stale > total/4 {
		health = "critical"
	} else if missing > 0 || stale > 0 {
		health = "warning"
	}

	healthColor := colorGreen
	if health == "critical" {
		healthColor = colorRed
	} else if health == "warning" {
		healthColor = colorYellow
	}

	fmt.Println("Registry Health Status")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Overall Health: %s%s%s\n", healthColor, health, colorReset)
	fmt.Printf("Total Skills:   %d\n", total)
	fmt.Printf("Missing Deps:   %s%d%s\n", colorRed, missing, colorReset)
	fmt.Printf("Stale Skills:   %s%d%s\n", colorYellow, stale, colorReset)
	fmt.Printf("Avg Coverage:   %.1f%%\n", avgCoverage*100)

	return nil
}

func runRegistryMissing(cmd *cobra.Command) error {
	client := getAPIClient(cmd)
	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	resp, err := client.Request(ctx, http.MethodGet, "/api/v1/registry?has_missing=true", nil)
	if err != nil {
		return fmt.Errorf("get missing: %w", err)
	}
	defer resp.Body.Close()

	var entries []models.SkillRegistryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No missing dependencies found. All skills are fully resolved.")
		return nil
	}

	fmt.Printf("Skills with Missing Dependencies (%d)\n", len(entries))
	fmt.Println(strings.Repeat("-", 60))
	for _, e := range entries {
		if len(e.MissingDeps) > 0 {
			fmt.Printf("\n%s%s%s\n", colorRed, e.SkillName, colorReset)
			for _, dep := range e.MissingDeps {
				fmt.Printf("  - %s\n", dep)
			}
		}
	}
	return nil
}

func runRegistryStale(cmd *cobra.Command) error {
	client := getAPIClient(cmd)
	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	resp, err := client.Request(ctx, http.MethodGet, "/api/v1/registry?stale=true", nil)
	if err != nil {
		return fmt.Errorf("get stale: %w", err)
	}
	defer resp.Body.Close()

	var entries []models.SkillRegistryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	if len(entries) == 0 {
		fmt.Println("No stale skills found. Registry is up to date.")
		return nil
	}

	fmt.Printf("Stale Skills (%d)\n", len(entries))
	fmt.Println(strings.Repeat("-", 70))
	fmt.Printf("%-30s %-20s %s\n", "NAME", "LAST REVIEW", "COVERAGE")
	fmt.Println(strings.Repeat("-", 70))
	for _, e := range entries {
		lastReview := "never"
		if e.LastReview != nil {
			lastReview = e.LastReview.Format("2006-01-02")
		}
		cov := fmt.Sprintf("%.0f%%", e.Coverage*100)
		fmt.Printf("%-30s %-20s %s\n", e.SkillName, lastReview, cov)
	}
	return nil
}

func runRegistryReview(cmd *cobra.Command, name string) error {
	client := getAPIClient(cmd)
	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	resp, err := client.Request(ctx, http.MethodGet, "/api/v1/registry/"+name, nil)
	if err != nil {
		return fmt.Errorf("get review: %w", err)
	}
	defer resp.Body.Close()

	var entry models.SkillRegistryEntry
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	fmt.Printf("Registry Review: %s\n", entry.SkillName)
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Coverage:     %.1f%%\n", entry.Coverage*100)
	fmt.Printf("Auto-Expand:  %v\n", entry.AutoExpand)
	fmt.Printf("Stale:        %v\n", entry.Stale)
	if entry.LastReview != nil {
		fmt.Printf("Last Review:  %s\n", entry.LastReview.Format(time.RFC3339))
	} else {
		fmt.Println("Last Review:  never")
	}

	if len(entry.MissingDeps) > 0 {
		fmt.Printf("\n%sMissing Dependencies (%d):%s\n", colorRed, len(entry.MissingDeps), colorReset)
		for _, dep := range entry.MissingDeps {
			fmt.Printf("  - %s\n", dep)
		}
	} else {
		fmt.Printf("\n%sAll dependencies resolved.%s\n", colorGreen, colorReset)
	}

	return nil
}

func runRegistryCoverage(cmd *cobra.Command, domain string) error {
	client := getAPIClient(cmd)
	ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
	defer cancel()

	query := "/api/v1/registry/coverage"
	if domain != "" {
		query += "?domain=" + domain
	}

	resp, err := client.Request(ctx, http.MethodGet, query, nil)
	if err != nil {
		return fmt.Errorf("get coverage: %w", err)
	}
	defer resp.Body.Close()

	var coverage struct {
		Overall  float64            `json:"overall"`
		ByDomain map[string]float64 `json:"by_domain"`
		BySkill  []struct {
			Name     string  `json:"name"`
			Coverage float64 `json:"coverage"`
		} `json:"by_skill"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&coverage); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	fmt.Printf("Registry Coverage")
	if domain != "" {
		fmt.Printf(" - Domain: %s", domain)
	}
	fmt.Println()
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Overall: %.1f%%\n", coverage.Overall*100)

	if len(coverage.ByDomain) > 0 {
		fmt.Printf("\nBy Domain:\n")
		for d, c := range coverage.ByDomain {
			bar := renderBar(c, 30)
			fmt.Printf("  %-15s [%s] %.0f%%\n", d, bar, c*100)
		}
	}

	if len(coverage.BySkill) > 0 {
		fmt.Printf("\nBy Skill:\n")
		for _, s := range coverage.BySkill {
			bar := renderBar(s.Coverage, 20)
			fmt.Printf("  %-25s [%s] %.0f%%\n", s.Name, bar, s.Coverage*100)
		}
	}

	return nil
}

// renderBar creates an ASCII progress bar
func renderBar(percentage float64, width int) string {
	filled := int(percentage * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	empty := width - filled
	return strings.Repeat("#", filled) + strings.Repeat("-", empty)
}

func getGlobalFormat(cmd *cobra.Command) string {
	f, _ := cmd.Root().PersistentFlags().GetString("format")
	return f
}
