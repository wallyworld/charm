// Copyright 2015 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package resource_test

import (
	"github.com/juju/errors"
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"gopkg.in/juju/charm.v6-unstable/resource"
)

var _ = gc.Suite(&MetaSuite{})

type MetaSuite struct{}

func (s *MetaSuite) TestParseMetaOkay(c *gc.C) {
	name := "my-resource"
	data := map[string]interface{}{
		"type":     "file",
		"filename": "filename.tgz",
		"comment":  "One line that is useful when operators need to push it.",
	}
	res := resource.ParseMeta(name, data)

	c.Check(res, jc.DeepEquals, resource.Meta{
		Name:    "my-resource",
		Type:    resource.TypeFile,
		Path:    "filename.tgz",
		Comment: "One line that is useful when operators need to push it.",
	})
}

func (s *MetaSuite) TestParseMetaMissingName(c *gc.C) {
	name := ""
	data := map[string]interface{}{
		"type":     "file",
		"filename": "filename.tgz",
		"comment":  "One line that is useful when operators need to push it.",
	}
	res := resource.ParseMeta(name, data)

	c.Check(res, jc.DeepEquals, resource.Meta{
		Name:    "",
		Type:    resource.TypeFile,
		Path:    "filename.tgz",
		Comment: "One line that is useful when operators need to push it.",
	})
}

func (s *MetaSuite) TestParseMetaMissingType(c *gc.C) {
	name := "my-resource"
	data := map[string]interface{}{
		"filename": "filename.tgz",
		"comment":  "One line that is useful when operators need to push it.",
	}
	res := resource.ParseMeta(name, data)

	c.Check(res, jc.DeepEquals, resource.Meta{
		Name:    "my-resource",
		Type:    resource.TypeUnknown,
		Path:    "filename.tgz",
		Comment: "One line that is useful when operators need to push it.",
	})
}

func (s *MetaSuite) TestParseMetaMissingPath(c *gc.C) {
	name := "my-resource"
	data := map[string]interface{}{
		"type":    "file",
		"comment": "One line that is useful when operators need to push it.",
	}
	res := resource.ParseMeta(name, data)

	c.Check(res, jc.DeepEquals, resource.Meta{
		Name:    "my-resource",
		Type:    resource.TypeFile,
		Path:    "",
		Comment: "One line that is useful when operators need to push it.",
	})
}

func (s *MetaSuite) TestParseMetaMissingComment(c *gc.C) {
	name := "my-resource"
	data := map[string]interface{}{
		"type":     "file",
		"filename": "filename.tgz",
	}
	res := resource.ParseMeta(name, data)

	c.Check(res, jc.DeepEquals, resource.Meta{
		Name:    "my-resource",
		Type:    resource.TypeFile,
		Path:    "filename.tgz",
		Comment: "",
	})
}

func (s *MetaSuite) TestParseMetaEmpty(c *gc.C) {
	name := "my-resource"
	data := make(map[string]interface{})
	res := resource.ParseMeta(name, data)

	c.Check(res, jc.DeepEquals, resource.Meta{
		Name: "my-resource",
	})
}

func (s *MetaSuite) TestParseMetaNil(c *gc.C) {
	name := "my-resource"
	var data map[string]interface{}
	res := resource.ParseMeta(name, data)

	c.Check(res, jc.DeepEquals, resource.Meta{
		Name: "my-resource",
	})
}

func (s *MetaSuite) TestValidateFull(c *gc.C) {
	res := resource.Meta{
		Name:    "my-resource",
		Type:    resource.TypeFile,
		Path:    "filename.tgz",
		Comment: "One line that is useful when operators need to push it.",
	}
	err := res.Validate()

	c.Check(err, jc.ErrorIsNil)
}

func (s *MetaSuite) TestValidateZeroValue(c *gc.C) {
	var res resource.Meta
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
}

func (s *MetaSuite) TestValidateMissingName(c *gc.C) {
	res := resource.Meta{
		Type:    resource.TypeFile,
		Path:    "filename.tgz",
		Comment: "One line that is useful when operators need to push it.",
	}
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `resource missing name`)
}

func (s *MetaSuite) TestValidateMissingType(c *gc.C) {
	res := resource.Meta{
		Name:    "my-resource",
		Path:    "filename.tgz",
		Comment: "One line that is useful when operators need to push it.",
	}
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `resource missing type`)
}

func (s *MetaSuite) TestValidateUnknownType(c *gc.C) {
	res := resource.Meta{
		Name:    "my-resource",
		Type:    "repo",
		Path:    "repo-root",
		Comment: "One line that is useful when operators need to push it.",
	}
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `.*unsupported resource type .*`)
}

func (s *MetaSuite) TestValidateMissingPath(c *gc.C) {
	res := resource.Meta{
		Name:    "my-resource",
		Type:    resource.TypeFile,
		Comment: "One line that is useful when operators need to push it.",
	}
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `resource missing filename`)
}

func (s *MetaSuite) TestValidateNestedPath(c *gc.C) {
	res := resource.Meta{
		Name: "my-resource",
		Type: resource.TypeFile,
		Path: "spam/eggs",
	}
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `.*filename cannot contain "/" .*`)
}

func (s *MetaSuite) TestValidateAbsolutePath(c *gc.C) {
	res := resource.Meta{
		Name: "my-resource",
		Type: resource.TypeFile,
		Path: "/spam/eggs",
	}
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `.*filename cannot contain "/" .*`)
}

func (s *MetaSuite) TestValidateSuspectPath(c *gc.C) {
	res := resource.Meta{
		Name: "my-resource",
		Type: resource.TypeFile,
		Path: "git@github.com:juju/juju.git",
	}
	err := res.Validate()

	c.Check(err, jc.Satisfies, errors.IsNotValid)
	c.Check(err, gc.ErrorMatches, `.*filename cannot contain "/" .*`)
}

func (s *MetaSuite) TestValidateMissingComment(c *gc.C) {
	res := resource.Meta{
		Name: "my-resource",
		Type: resource.TypeFile,
		Path: "filename.tgz",
	}
	err := res.Validate()

	c.Check(err, jc.ErrorIsNil)
}
