package tfupdate

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func TestNewTfsModuleUpdater(t *testing.T) {
	cases := []struct {
		name            string
		source          string
		newSource       string
		sourceMatchType string
		version         string
		want            Updater
		ok              bool
	}{
		// TDD - first instance
		{
			name:            "terraform-aws-modules/vpc/aws",
			source: "localterraform.com/my-org/terraform-aws-modules/vpc/aws",
			newSource: "app.terraform.io/my-org/terraform-aws-modules/vpc/aws",
			sourceMatchType: "full",
			version:         "2.17.0",
			want: &TfsModuleUpdater{
				name: "terraform-aws-modules/vpc/aws",
				source: "app.terraform.io/my-org/terraform-aws-modules/vpc/aws",
				newSource: "app.terraform.io/my-org/terraform-aws-modules/vpc/aws",
				nameRegex: nil,
				version:   "2.17.0",
			},
			ok: true,
		},
		{
			name: "terraform-aws-modules/vpc/aws",
			source: "localterraform.com/my-org/terraform-aws-modules/vpc/aws",
			newSource: "./modules/vpc",
			sourceMatchType: "full",
			version:         "", // TODO: remove this as a required input
			want: &TfsModuleUpdater{
				name: "terraform-aws-modules/vpc/aws",
				source: "./modules/vpc",
				newSource: "./modules/vpc",
				nameRegex: nil,
				version: "", // TODO: remove this as a required inpu
			},
			ok: true,
		},
	}


	for _, tc := range cases {
		got, err := NewTfsModuleUpdater(tc.name, tc.source, tc.newSource, tc.version, nil)
		if tc.ok && err != nil {
			t.Errorf("NewTfsModuleUpdater() with name = %s, version = %s returns unexpected err: %+v", tc.name, tc.version, err)
		}

		if !tc.ok && err == nil {
			t.Errorf("NewTfsModuleUpdater() with name = %s, version = %s expects to return an error, but no error", tc.name, tc.version)
		}

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("NewTfsModuleUpdater() with name = %s, version = %s returns %#v, but want = %#v", tc.name, tc.version, got, tc.want)
		}
	}
}

func TestTfsModuleUpdaterHashiCorpModuleSources(t *testing.T) {
	cases := []struct {
		filename        string
		src             string
		name            string
		source					string
		newSource				string
		sourceMatchType string
		version         string
		want            string
		ok              bool
	}{
		// Follow the cases shown on https://developer.hashicorp.com/terraform/language/modules/sources
		// TODO: create a test case such that the newSource is not a valid source, which should error out

		{
			// TODO: implement the logic to make this test case pass
			filename: "main.tf",
			src: `
module "consul" {
  source = "./consul"
}
`,
			name: "Local Path",
			source: "./consul",
			newSource: "./consul-new",
			sourceMatchType: "full",
			want: `
module "consul" {
  source = "./consul-new"
}
`,
			ok: true,
		},
		{
			// TODO: implement the logic to make this test case pass
			filename: "main.tf",
			src: `
module "vpc" {
  source = "terraform-aws-modules/vpc/aws"
}
`,
			name: "Terraform Registry (public registry) to Terraform Registry (public registry)",
			source: "terraform-aws-modules/vpc/aws",
			newSource: "terraform-aws-modules/vpc2/aws",
			sourceMatchType: "full",
			want: `
module "vpc" {
  source = "terraform-aws-modules/vpc2/aws"
}
`,
			ok: true,
		},

		{
			// TODO: implement the logic to make this test case pass
			filename: "main.tf",
			src: `
module "vpc" {
  source  = "app.terraform.io/my-org/terraform-aws-modules/vpc/aws"
  version = "2.17.0"
}
`,
			name: "Terraform Cloud / Terraform Enterprise Registry to Terraform Cloud / Terraform Enterprise Registry",
			source: "localterraform.com/my-org/terraform-aws-modules/vpc/aws",
			newSource: "app.terraform.io/my-new-org/terraform-aws-modules/vpc/aws",
			sourceMatchType: "full",
			want: `
module "vpc" {
  source  = "app.terraform.io/my-new-org/terraform-aws-modules/vpc/aws"
  version = "2.17.0"
}
`,
			ok: true,
		},

		{
			// TODO: implement the logic to make this test case pass
			filename: "main.tf",
			src: `
module "example" {
  source = "github.com/hashicorp/example"
}
`,
			name: "GitHub (git)to GitHub (git) via HTTPS",
			source: "github.com/hashicorp/example",
			newSource: "github.com/hashicorp-new/example",
			sourceMatchType: "full",
			want: `
module "example" {
  source = "github.com/hashicorp-new/example"
}
`,
			ok: true,
		},

		{
			// TODO: implement the logic to make this test case pass
			filename: "main.tf",
			src: `
module "vpc" {
  source  = "localterraform.com/my-org/terraform-aws-modules/vpc/aws"
  version = "2.17.0"
}
`,
			name: "single module update - localterraform.com to some registry",
			source: "localterraform.com/my-org/terraform-aws-modules/vpc/aws",
			newSource: "app.terraform.io/my-org/terraform-aws-modules/vpc/aws",
			sourceMatchType: "full",
			want: `
module "vpc" {
  source  = "app.terraform.io/my-org/terraform-aws-modules/vpc/aws"
  version = "2.17.0"
}
`,
			ok: true,
		},
		{
			// TODO: implement the logic to make this test case pass
			filename: "main.tf",
			src: `
module "vpc1" {
  source  = "localterraform.com/my-org/terraform-aws-modules/vpc/aws"
  version = "2.18.0"
}
module "vpc2" {
  source  = "localterraform.com/my-org/terraform-aws-modules/vpc/aws"
  version = "2.18.0"
}
`,
			name: "multiple module update",
			source: "localterraform.com/my-org/terraform-aws-modules/vpc/aws",
			newSource: "app.terraform.io/my-org/terraform-aws-modules/vpc/aws",
			sourceMatchType: "full",
			want: `
module "vpc1" {
  source  = "app.terraform.io/my-org/terraform-aws-modules/vpc/aws"
  version = "2.18.0"
}
module "vpc2" {
  source  = "app.terraform.io/my-org/terraform-aws-modules/vpc/aws"
  version = "2.18.0"
}
`,
			ok: true,
		},
		{
			filename: "main.tf",
			src: `
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "2.17.0"
}
`,
			name: "do not update non-matching module name from source / newSource",
			source: "terraform-aws-modules/hoge/aws",
			newSource: "terraform-aws-modules2/hoge/aws",
			sourceMatchType: "full",
			version:         "2.17.0",
			want: `
module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "2.17.0"
}
`,
			ok: true,
		},
		{
			// TODO: implement the logic to make this test case pass
			filename: "main.tf",
			src: `
module "vpc" {
  source = "git::https://example.com/vpc.git"
}
`,
			name: "git source",
			source: "git::https://example.com/vpc.git",
			newSource: "git::https://example2.com/vpc.git",
			sourceMatchType: "full",
			want: `
module "vpc" {
  source = "git::https://example2.com/vpc.git"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		u := &TfsModuleUpdater{
			name: tc.name,
			source: tc.source,
			newSource: tc.newSource,
			nameRegex: func() *regexp.Regexp {
				if tc.sourceMatchType == "regex" {
					return regexp.MustCompile(tc.name)
				}
				return nil
			}(),
			version: tc.version,
		}
		f, diags := hclwrite.ParseConfig([]byte(tc.src), tc.filename, hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			t.Fatalf("unexpected diagnostics: %s", diags)
		}

		err := u.Update(context.Background(), nil, tc.filename, f)
		if tc.ok && err != nil {
			t.Errorf("Update() with src = %s, newSrc = %s, name = %s, version = %s returns unexpected err: %+v", tc.src, tc.newSource, tc.name, tc.version, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("Update() with src = %s, newSrc = %s, name = %s, version = %s expects to return an error, but no error", tc.src, tc.newSource, tc.name, tc.version)
		}

		got := string(hclwrite.Format(f.BuildTokens(nil).Bytes()))
		if got != tc.want {
			t.Errorf("Update() with src = %s, newSrc = %s, name = %s, version = %s returns %s, but want = %s", tc.src, tc.newSource, tc.name, tc.version, got, tc.want)
		}
	}
}

func TestTfsModuleUpdaterRegex(t *testing.T) {
	cases := []struct {
		filename        string
		src             string
		name            string
		source					string
		newSource				string
		sourceMatchType string
		version         string
		want            string
		ok              bool
	}{
// TODO: add new test cases for regex support & keep them separated for now
		{
			filename: "main.tf",
			src: `
module "vpc1" {
  source  = "terraform-aws-modules.git/vpc/aws1"
  version = "2.17.0"
}
module "vpc2" {
  source  = "terraform-aws-modules.git/vpc/aws2"
  version = "2.17.0"
}
`,
			name: "terraform-aws-modules.git/",
			source: "terraform-aws-modules.git/",
			newSource: "terraform-aws-modules.git/", // TODO:
			// version:         "2.18.0",
			sourceMatchType: "regex",
			want: `
module "vpc1" {
  source  = "terraform-aws-modules.git/vpc/aws1"
  version = "2.17.0"
}
module "vpc2" {
  source  = "terraform-aws-modules.git/vpc/aws2"
  version = "2.17.0"
}
`,
			ok: true,
		},
		{
			filename: "main.tf",
			src: `
module "vpc1" {
  source  = "terraform-aws-modules.git/vpc/aws1"
  version = "2.17.0"
}
module "vpc2" {
  source  = "terraform-aws-modules.git/vpc/aws2"
  version = "2.17.0"
}
`,
			name:            "terraform-aws-modules\\.git/.+",
			version:         "2.18.0",
			sourceMatchType: "regex",
			want: `
module "vpc1" {
  source  = "terraform-aws-modules.git/vpc/aws1"
  version = "2.18.0"
}
module "vpc2" {
  source  = "terraform-aws-modules.git/vpc/aws2"
  version = "2.18.0"
}
`,
			ok: true,
		},
	}

	for _, tc := range cases {
		u := &TfsModuleUpdater{
			name: tc.name,
			source: tc.source,
			newSource: tc.newSource,
			nameRegex: func() *regexp.Regexp {
				if tc.sourceMatchType == "regex" {
					return regexp.MustCompile(tc.name)
				}
				return nil
			}(),
			version: tc.version,
		}
		f, diags := hclwrite.ParseConfig([]byte(tc.src), tc.filename, hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			t.Fatalf("unexpected diagnostics: %s", diags)
		}

		err := u.Update(context.Background(), nil, tc.filename, f)
		if tc.ok && err != nil {
			t.Errorf("Update() with src = %s, newSrc = %s, name = %s, version = %s returns unexpected err: %+v", tc.src, tc.newSource, tc.name, tc.version, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("Update() with src = %s, newSrc = %s, name = %s, version = %s expects to return an error, but no error", tc.src, tc.newSource, tc.name, tc.version)
		}

		got := string(hclwrite.Format(f.BuildTokens(nil).Bytes()))
		if got != tc.want {
			t.Errorf("Update() with src = %s, newSrc = %s, name = %s, version = %s returns %s, but want = %s", tc.src, tc.newSource, tc.name, tc.version, got, tc.want)
		}
	}
}


// TODO: update the name of the original test function to be for registry to registry
func TestParseTfsModuleSource(t *testing.T) {
	cases := []struct {
		src     string
		name    string
		version string
	}{
		{
			src: `
module "vpc" {
  source = "git::https://example.com/vpc.git"
}
`,
			name:    "git::https://example.com/vpc.git",
			version: "",
		},
		{
			src: `
module "vpc" {
  source = "git::https://example.com/vpc.git?ref=v1"
}
`,
			name:    "git::https://example.com/vpc.git",
			version: "1",
		},
		{
			src: `
module "vpc" {
  source = "git::https://example.com/vpc.git?ref=v1.2"
}
`,
			name:    "git::https://example.com/vpc.git",
			version: "1.2",
		},
		{
			src: `
module "vpc" {
  source = "git::https://example.com/vpc.git?ref=v1.2.0"
}
`,
			name:    "git::https://example.com/vpc.git",
			version: "1.2.0",
		},
		{
			src: `
module "vpc" {
  source = "git::https://example.com/vpc.git?ref=v1.2.0-rc1"
}
`,
			name:    "git::https://example.com/vpc.git",
			version: "1.2.0-rc1",
		},
		{
			src: `
module "vpc" {
  source = "git::https://example.com/vpc.git?ref=vhoge"
}
`,
			name:    "git::https://example.com/vpc.git?ref=vhoge",
			version: "",
		},
	}

	for _, tc := range cases {
		f, diags := hclwrite.ParseConfig([]byte(tc.src), "", hcl.Pos{Line: 1, Column: 1})
		if diags.HasErrors() {
			t.Fatalf("unexpected diagnostics: %s", diags)
		}

		m := allMatchingBlocksByType(f.Body(), "module")
		if len(m) != 1 {
			t.Fatalf("failed to get module block: %s", tc.src)
		}
		s := m[0].Body().GetAttribute("source")
		if s == nil {
			t.Fatalf("failed to get module source attribute: %s", tc.src)
		}
		name, version := parseModuleSource(s)

		if !(name == tc.name && version == tc.version) {
			t.Errorf("parseModuleSource() with src = %s returns (%s, %s), but want = (%s, %s)", tc.src, name, version, tc.name, tc.version)
		}
	}
}
