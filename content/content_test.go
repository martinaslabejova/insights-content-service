/*
Copyright © 2020 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package content_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/RedHatInsights/insights-operator-utils/tests/helpers"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"

	"github.com/RedHatInsights/insights-content-service/content"
)

const errYAMLBadToken = "yaml: line 14: found character that cannot start any token"

func init() {
	zerolog.SetGlobalLevel(zerolog.WarnLevel)
}

// TestContentParseOK checks whether reading from directory of correct content works as expected
func TestContentParseOK(t *testing.T) {
	con, err := content.ParseRuleContentDir("../tests/content/ok/")
	helpers.FailOnError(t, err)

	rule1Content, exists := con.Rules["rule1"]
	assert.True(t, exists, "'rule1' content is missing")

	_, exists = rule1Content.ErrorKeys["err_key"]
	assert.True(t, exists, "'err_key' error content is missing")
}

// TestContentParseOKNoContent checks that parsing content when there is no rule
// content available, but the file structure is otherwise okay, succeeds.
func TestContentParseOKNoContent(t *testing.T) {
	con, err := content.ParseRuleContentDir("../tests/content/ok_no_content/")
	helpers.FailOnError(t, err)

	assert.Empty(t, con.Rules)
}

// TestContentParseContentWithoutSummaryMD checks that parsing content when
// there is NOT summary.md file available, but the file structure is otherwise
// okay, succeeds.
func TestContentParseContentWithoutSummaryMD(t *testing.T) {
	con, err := content.ParseRuleContentDir("../tests/content/ok_missing_summary_md/")

	assert.Nil(t, err)
	assert.NotEmpty(t, con.Rules)
}

// TestContentParseInvalidDir checks how incorrect (non-existing) directory is handled
func TestContentParseInvalidDir(t *testing.T) {
	const invalidDirPath = "../tests/content/not-a-real-dir"
	_, err := content.ParseRuleContentDir(invalidDirPath)
	assert.EqualError(t, err, fmt.Sprintf("open %s/config.yaml: no such file or directory", invalidDirPath))
}

// TestContentParseNotDirectory1 checks how incorrect (non-existing) directory is handled
func TestContentParseNotDirectory1(t *testing.T) {
	// this is not a proper directory
	const notADirPath = "../tests/tests.toml"
	_, err := content.ParseRuleContentDir(notADirPath)
	assert.EqualError(t, err, fmt.Sprintf("open %s/config.yaml: not a directory", notADirPath))
}

// TestContentParseNotDirectory2 checks how incorrect (non-existing) directory is handled
func TestContentParseInvalidDir2(t *testing.T) {
	// this is not a proper directory
	const notADirPath = "/dev/null"
	_, err := content.ParseRuleContentDir(notADirPath)
	assert.EqualError(t, err, fmt.Sprintf("open %s/config.yaml: not a directory", notADirPath))
}

// TestContentParseMissingFile checks how missing mandatory file(s) in content directory are handled
func TestContentParseMissingFile(t *testing.T) {
	buf := new(bytes.Buffer)
	log.Logger = zerolog.New(buf)

	// has mandatory plugin.yaml missing
	_, err := content.ParseRuleContentDir("../tests/content/missing/")

	assert.Nil(t, err)
	assert.Contains(t, buf.String(), "Error trying to parse rule in dir")
}

// TestContentParseBadPluginYAML tests handling bad/incorrect plugin.yaml file
func TestContentParseBadPluginYAML(t *testing.T) {
	buf := new(bytes.Buffer)
	log.Logger = zerolog.New(buf)

	_, err := content.ParseRuleContentDir("../tests/content/bad_plugin/")

	assert.Nil(t, err)
	assert.Contains(t, buf.String(), errYAMLBadToken)
}

// TestContentParseBadMetadataYAML tests handling bad/incorrect metadata.yaml file
func TestContentParseBadMetadataYAML(t *testing.T) {
	buf := new(bytes.Buffer)
	log.Logger = zerolog.New(buf)

	_, err := content.ParseRuleContentDir("../tests/content/bad_metadata/")

	assert.Nil(t, err)
	assert.Contains(t, buf.String(), errYAMLBadToken)
}

// TestContentParseBadMetadataYAML tests handling bad/incorrect metadata.yaml file
func TestContentParseNoExternal(t *testing.T) {
	noExternalPath := "../tests/content/no_external"
	_, err := content.ParseRuleContentDir(noExternalPath)
	assert.EqualError(t, err, fmt.Sprintf("open %s/external: no such file or directory", noExternalPath))
}

// TestContentParseNoInternal tests case where there is no folder for internal rules
func TestContentParseNoInternal(t *testing.T) {
	noInternalPath := "../tests/content/no_internal"
	_, err := content.ParseRuleContentDir(noInternalPath)
	assert.EqualError(t, err, fmt.Sprintf("open %s/internal: no such file or directory", noInternalPath))
}
