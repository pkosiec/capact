package manifest

import (
	"capact.io/capact/pkg/sdk/apis/0.0.1/types"
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/pkg/errors"
	"github.com/xeipuuv/gojsonschema"
	"net/http"
	"os"
	"sigs.k8s.io/yaml"
	"sort"
)

type loadedOCFSchema struct {
	common *gojsonschema.SchemaLoader
	kind   map[types.ManifestKind]*gojsonschema.Schema
}

type JSONSchemaValidator struct {
	fs             http.FileSystem

	schemaRootPath string
	cachedSchemas  map[types.OCFVersion]*loadedOCFSchema
}

func NewJSONSchemaValidator(fs http.FileSystem, schemaRootPath string) *JSONSchemaValidator {
	return &JSONSchemaValidator{
		schemaRootPath: schemaRootPath,
		fs:             fs,
		cachedSchemas:  map[types.OCFVersion]*loadedOCFSchema{},
	}
}

func (v *JSONSchemaValidator) Do(metadata types.ManifestMetadata, yamlBytes []byte) (ValidationResult, error) {
	schema, err := v.getManifestSchema(metadata)
	if err != nil {
		return newSimpleValidationResult(), errors.Wrap(err, "failed to get JSON schema")
	}

	jsonBytes, err := yaml.YAMLToJSON(yamlBytes)
	if err != nil {
		return newSimpleValidationResult(errors.Wrap(err, "cannot convert YAML manifest to JSON")), err
	}

	manifestLoader := gojsonschema.NewBytesLoader(jsonBytes)

	jsonschemaResult, err := schema.Validate(manifestLoader)
	if err != nil {
		return newSimpleValidationResult(errors.Wrap(err, "error occurred during JSON schema validation")), err
	}

	result := newSimpleValidationResult()

	for _, err := range jsonschemaResult.Errors() {
		result.Errors = append(result.Errors, fmt.Errorf("%v", err.String()))
	}

	return result, err
}


func (v *JSONSchemaValidator) getManifestSchema(metadata types.ManifestMetadata) (*gojsonschema.Schema, error) {
	var ok bool
	var cachedSchema *loadedOCFSchema

	if cachedSchema, ok = v.cachedSchemas[metadata.OCFVersion]; !ok {
		cachedSchema = &loadedOCFSchema{
			common: nil,
			kind:   map[types.ManifestKind]*gojsonschema.Schema{},
		}
		v.cachedSchemas[metadata.OCFVersion] = cachedSchema
	}

	if schema, ok := cachedSchema.kind[metadata.Kind]; ok {
		return schema, nil
	}

	rootLoader := v.getRootSchemaJSONLoader(metadata)

	if cachedSchema.common == nil {
		sl, err := v.getCommonSchemaLoader(metadata.OCFVersion)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get common schema loader")
		}
		cachedSchema.common = sl
	}

	schema, err := cachedSchema.common.Compile(rootLoader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compile schema for %s/%s", metadata.OCFVersion, metadata.Kind)
	}

	cachedSchema.kind[metadata.Kind] = schema

	return schema, nil
}

func (v *JSONSchemaValidator) getRootSchemaJSONLoader(metadata types.ManifestMetadata) gojsonschema.JSONLoader {
	filename := strcase.ToKebab(string(metadata.Kind))
	path := fmt.Sprintf("file://%s/%s/schema/%s.json", v.schemaRootPath, metadata.OCFVersion, filename)
	return gojsonschema.NewReferenceLoaderFileSystem(path, v.fs)
}

func (v *JSONSchemaValidator) getCommonSchemaLoader(ocfVersion types.OCFVersion) (*gojsonschema.SchemaLoader, error) {
	commonDir := fmt.Sprintf("%s/%s/schema/common", v.schemaRootPath, ocfVersion)

	sl := gojsonschema.NewSchemaLoader()

	files, err := v.ReadDir(commonDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list common schemas directory")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		path := fmt.Sprintf("file://%s/%s", commonDir, file.Name())
		if err := sl.AddSchemas(gojsonschema.NewReferenceLoaderFileSystem(path, v.fs)); err != nil {
			return nil, errors.Wrapf(err, "cannot load common schema %s", path)
		}
	}

	return sl, nil
}

// ReadDir reads the directory named by dirname and returns
// a list of directory entries sorted by filename.
func (v *JSONSchemaValidator) ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := v.fs.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	if err != nil {
		return nil, err
	}
	if err := f.Close(); err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}
