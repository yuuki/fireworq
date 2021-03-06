package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fireworq/fireworq/config"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Usage: gendoc <type>")
		os.Exit(1)
	}

	if os.Args[1] == "config" {
		printConfigDoc()
	}
}

func printConfigDoc() {
	fmt.Println(`<!-- DO NOT EDIT: this document is automatically generated by script/gendoc/gendoc.go -->

Configuration
=============

You can configure Fireworq by providing environment variables on
starting a daemon (for both docker-composed and manually-set-up
instances) or by specifying command line arguments (for a manual
setup).  Command line arguments precede the values of environment
variables.

The following variables/arguments are available.  Some of them are
applicable only to a [manual setup][section-manual-setup].
`)

	categorized := make(map[string]configItems)
	for _, d := range config.Descriptions() {
		item := &configItem{d}
		categorized[item.Category] = append(categorized[item.Category], item)
	}

	fmt.Println(`- [Common Variables/Arguments][section-config-common]`)
	categorized["common"].printTableOfContents()

	fmt.Println(`- [Variables/Arguments only Applicable to Manual Setup][section-config-manual-setup]`)
	categorized["manual"].printTableOfContents()

	fmt.Println(`- [Variables only Applicable to a Docker-composed Instance][section-config-docker]
  - [` + "`" + `FIREWORQ_PORT` + "`" + `](#env-port)`)

	fmt.Println(`
## <a name="config-common">Common Variables/Arguments</a>`)
	categorized["common"].printDescriptions()

	fmt.Println(`
## <a name="config-manual-setup">Variables/Arguments only Applicable to Manual Setup</a>`)
	categorized["manual"].printDescriptions()

	fmt.Println(`
## <a name="config-docker">Variables only Applicable to a Docker-composed Instance</a>

### <a name="env-port">` + "`" + `FIREWORQ_PORT` + "`" + `</a>

Default: ` + "`" + `8080` + "`" + `

Specifies the port number of a daemon.`)

	fmt.Print(`
[section-config-common]: #config-common
[section-config-manual-setup]: #config-manual-setup
[section-config-docker]: #config-docker
[section-manual-setup]: ./production.md#manual-setup
[section-graceful-restart]: ./production.md#graceful-restart

[api-put-queue]: ./api.md#api-put-queue
[api-put-routing]: ./api.md#api-put-routing
`)
}

type configItems []*configItem

func (items configItems) printTableOfContents() {
	for _, item := range items {
		fmt.Println("  - " + item.Link())
	}
}

func (items configItems) printDescriptions() {
	for _, item := range items {
		item.printDescription()
	}
}

type configItem struct {
	config.Item
}

func (item *configItem) printDescription() {
	heading := fmt.Sprintf(
		"### <a name=\"%s\">%s</a>",
		item.AnchorName(),
		item.LinkLabel(),
	)
	fmt.Println(heading)
	if item.DefaultValue != "" {
		fmt.Println(fmt.Sprintf("Default: `%s`", item.DefaultValue))
	}
	fmt.Println(item.Description)
}

func (item *configItem) EnvironmentVariable() string {
	return "FIREWORQ_" + strings.ToUpper(item.Name)
}

func (item *configItem) Link() string {
	return fmt.Sprintf("[%s](#%s)", item.LinkLabel(), item.AnchorName())
}

func (item *configItem) LinkLabel() string {
	return fmt.Sprintf(
		"`%s`, `%s`",
		item.EnvironmentVariable(),
		item.Argument(),
	)
}

func (item *configItem) AnchorName() string {
	return "env-" + strings.Replace(item.Name, "_", "-", -1)
}
