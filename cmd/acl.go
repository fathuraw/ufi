package cmd

import (
	"fmt"
	"strings"

	"github.com/fathuraw/ufi/internal/output"
	"github.com/fathuraw/ufi/internal/unifi"
	"github.com/spf13/cobra"
)

var aclCmd = &cobra.Command{
	Use:   "acl",
	Short: "Manage ACL rules",
}

var aclListCmd = &cobra.Command{
	Use:   "list",
	Short: "List ACL rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		rules, err := client.ListACLRules(listParams(cmd))
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(rules)
		}
		t := output.NewTable("ID", "NAME", "ACTION", "SOURCE MAC", "ENABLED", "INDEX")
		for _, r := range rules {
			t.AddRow(r.ID, r.Name, r.Action, r.SourceMAC, r.Enabled, r.Index)
		}
		t.Flush()
		return nil
	},
}

var aclCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create an ACL rule",
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		action, _ := cmd.Flags().GetString("action")
		sourceMac, _ := cmd.Flags().GetString("source-mac")
		desc, _ := cmd.Flags().GetString("description")
		enabled, _ := cmd.Flags().GetBool("enabled")

		req := unifi.ACLRuleCreateRequest{
			Name:        name,
			Action:      action,
			SourceMAC:   sourceMac,
			Description: desc,
			Enabled:     enabled,
		}
		r, err := client.CreateACLRule(req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(r)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "ACL rule %s created (ID: %s).\n", r.Name, r.ID)
		return nil
	},
}

var aclUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an ACL rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		action, _ := cmd.Flags().GetString("action")
		sourceMac, _ := cmd.Flags().GetString("source-mac")
		desc, _ := cmd.Flags().GetString("description")
		enabled, _ := cmd.Flags().GetBool("enabled")

		req := unifi.ACLRuleCreateRequest{
			Name:        name,
			Action:      action,
			SourceMAC:   sourceMac,
			Description: desc,
			Enabled:     enabled,
		}
		r, err := client.UpdateACLRule(args[0], req)
		if err != nil {
			return err
		}
		if jsonOut {
			return output.PrintJSON(r)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "ACL rule %s updated.\n", args[0])
		return nil
	},
}

var aclDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an ACL rule",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := client.DeleteACLRule(args[0]); err != nil {
			return err
		}
		fmt.Fprintf(cmd.OutOrStdout(), "ACL rule %s deleted.\n", args[0])
		return nil
	},
}

var aclReorderCmd = &cobra.Command{
	Use:   "reorder <id1,id2,...>",
	Short: "Reorder ACL rules",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := strings.Split(args[0], ",")
		if err := client.ReorderACLRules(ids); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), "ACL rules reordered.")
		return nil
	},
}

func init() {
	addListFlags(aclListCmd)

	aclCreateCmd.Flags().String("name", "", "rule name (required)")
	aclCreateCmd.Flags().String("action", "", "action (ALLOW, DENY)")
	aclCreateCmd.Flags().String("source-mac", "", "source MAC address")
	aclCreateCmd.Flags().String("description", "", "description")
	aclCreateCmd.Flags().Bool("enabled", true, "enable rule")
	aclCreateCmd.MarkFlagRequired("name")
	aclCreateCmd.MarkFlagRequired("action")

	aclUpdateCmd.Flags().String("name", "", "rule name")
	aclUpdateCmd.Flags().String("action", "", "action")
	aclUpdateCmd.Flags().String("source-mac", "", "source MAC address")
	aclUpdateCmd.Flags().String("description", "", "description")
	aclUpdateCmd.Flags().Bool("enabled", true, "enable rule")

	aclCmd.AddCommand(aclListCmd)
	aclCmd.AddCommand(aclCreateCmd)
	aclCmd.AddCommand(aclUpdateCmd)
	aclCmd.AddCommand(aclDeleteCmd)
	aclCmd.AddCommand(aclReorderCmd)

	rootCmd.AddCommand(aclCmd)
}
