// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

const (
	BasicAuthScopes  = "BasicAuth.Scopes"
	BearerAuthScopes = "BearerAuth.Scopes"
)

// Defines values for HarborSbomReportMediaType.
const (
	HarborSbomReportMediaTypeApplicationspdxJson         HarborSbomReportMediaType = "application/spdx+json"
	HarborSbomReportMediaTypeApplicationvndCyclonedxJson HarborSbomReportMediaType = "application/vnd.cyclonedx+json"
)

// Defines values for SbomParametersSbomMediaTypes.
const (
	SbomParametersSbomMediaTypesApplicationspdxJson         SbomParametersSbomMediaTypes = "application/spdx+json"
	SbomParametersSbomMediaTypesApplicationvndCyclonedxJson SbomParametersSbomMediaTypes = "application/vnd.cyclonedx+json"
)

// Defines values for ScanRequestEnabledCapabilitiesType.
const (
	ScanRequestEnabledCapabilitiesTypeSbom          ScanRequestEnabledCapabilitiesType = "sbom"
	ScanRequestEnabledCapabilitiesTypeVulnerability ScanRequestEnabledCapabilitiesType = "vulnerability"
)

// Defines values for ScannerCapabilityType.
const (
	ScannerCapabilityTypeSbom          ScannerCapabilityType = "sbom"
	ScannerCapabilityTypeVulnerability ScannerCapabilityType = "vulnerability"
)

// Defines values for Severity.
const (
	Critical   Severity = "Critical"
	High       Severity = "High"
	Low        Severity = "Low"
	Medium     Severity = "Medium"
	Negligible Severity = "Negligible"
	Unknown    Severity = "Unknown"
)

// Artifact defines model for Artifact.
type Artifact struct {
	// Digest The artifact's digest, consisting of an algorithm and hex portion.
	Digest *string `json:"digest,omitempty"`

	// MimeType The MIME type of the artifact.
	MimeType *string `json:"mime_type,omitempty"`

	// Repository The name of the Docker Registry repository containing the artifact.
	Repository *string `json:"repository,omitempty"`

	// Tag The artifact's tag
	Tag *string `json:"tag,omitempty"`
}

// CVSSDetails defines model for CVSSDetails.
type CVSSDetails struct {
	// ScoreV2 The CVSS 2.0 score for the vulnerability.
	ScoreV2 *float32 `json:"score_v2,omitempty"`

	// ScoreV3 The CVSS 3.0 score for the vulnerability.
	ScoreV3 *float32 `json:"score_v3,omitempty"`

	// VectorV2 The CVSS 2.0 vector for the vulnerability. The string is of the form AV:L/AC:M/Au:N/C:P/I:N/A:N
	VectorV2 *string `json:"vector_v2,omitempty"`

	// VectorV3 The CVSS 3.0 vector for the vulnerability.
	VectorV3 *string `json:"vector_v3,omitempty"`
}

// Error defines model for Error.
type Error struct {
	Message *string `json:"message,omitempty"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error *Error `json:"error,omitempty"`
}

// HarborSbomReport defines model for HarborSbomReport.
type HarborSbomReport struct {
	Artifact *Artifact `json:"artifact,omitempty"`

	// GeneratedAt The time of the report generated.
	GeneratedAt *time.Time `json:"generated_at,omitempty"`

	// MediaType The format of the sbom data.
	MediaType *HarborSbomReportMediaType `json:"media_type,omitempty"`

	// Sbom The raw data of the sbom generated by the scanner.
	Sbom *map[string]interface{} `json:"sbom,omitempty"`

	// Scanner Basic scanner properties such as name, vendor, and version.
	Scanner *Scanner `json:"scanner,omitempty"`

	// VendorAttributes The additional attributes of the vendor.
	VendorAttributes *map[string]interface{} `json:"vendor_attributes,omitempty"`
}

// HarborSbomReportMediaType The format of the sbom data.
type HarborSbomReportMediaType string

// HarborVulnerabilityReport defines model for HarborVulnerabilityReport.
type HarborVulnerabilityReport struct {
	Artifact    *Artifact  `json:"artifact,omitempty"`
	GeneratedAt *time.Time `json:"generated_at,omitempty"`

	// Scanner Basic scanner properties such as name, vendor, and version.
	Scanner *Scanner `json:"scanner,omitempty"`

	// Severity A standard scale for measuring the severity of a vulnerability.
	//
	// * `Unknown` - either a security problem that has not been assigned to a priority yet or a priority that the
	//   scanner did not recognize.
	// * `Negligible` - technically a security problem, but is only theoretical in nature, requires a very special
	//   situation, has almost no install base, or does no real damage.
	// * `Low` - a security problem, but is hard to exploit due to environment, requires a user-assisted attack,
	//   a small install base, or does very little damage.
	// * `Medium` - a real security problem, and is exploitable for many people. Includes network daemon denial of
	//   service attacks, cross-site scripting, and gaining user privileges.
	// * `High` - a real problem, exploitable for many people in a default installation. Includes serious remote denial
	//   of service, local root privilege escalations, or data loss.
	// * `Critical` - a world-burning problem, exploitable for nearly all people in a default installation. Includes
	//   remote root privilege escalations, or massive data loss.
	Severity        *Severity            `json:"severity,omitempty"`
	Vulnerabilities *[]VulnerabilityItem `json:"vulnerabilities,omitempty"`
}

// Registry defines model for Registry.
type Registry struct {
	// Authorization An optional value of the HTTP Authorization header sent with each request to the Docker Registry v2 API.
	// It's used to exchange Base64 encoded robot account credentials to a short lived JWT access token which
	// allows the underlying scanner to pull the artifact from the Docker Registry.
	Authorization *string `json:"authorization,omitempty"`

	// Url A base URL or the Docker Registry v2 API.
	Url *string `json:"url,omitempty"`
}

// SbomParameters defines model for SbomParameters.
type SbomParameters struct {
	SbomMediaTypes *[]SbomParametersSbomMediaTypes `json:"sbom_media_types,omitempty"`
}

// SbomParametersSbomMediaTypes defines model for SbomParameters.SbomMediaTypes.
type SbomParametersSbomMediaTypes string

// ScanRequest defines model for ScanRequest.
type ScanRequest struct {
	Artifact Artifact `json:"artifact"`

	// EnabledCapabilities Enable which capabilities supported by scanner, for backward compatibility, without this field scanner can be considered to enable all capabilities by default.
	EnabledCapabilities *[]struct {
		Parameters *SbomParameters `json:"parameters,omitempty"`

		// ProducesMimeTypes The set of MIME types of reports generated by the scanner for the consumes_mime_types of the same capability record, it is a subset or fullset of the
		// produces_mime_types of the capability returned by the metadata API, used for client to fine grained control of the expected report type. It's a optional
		// field, only applied when client needs to customize it, otherwise the scanner can think it's a fullset as before behavior if without this field.
		ProducesMimeTypes *[]string `json:"produces_mime_types,omitempty"`

		// Type The type of the scan capability.
		Type ScanRequestEnabledCapabilitiesType `json:"type"`
	} `json:"enabled_capabilities,omitempty"`
	Registry Registry `json:"registry"`
}

// ScanRequestEnabledCapabilitiesType The type of the scan capability.
type ScanRequestEnabledCapabilitiesType string

// ScanRequestId A unique identifier returned by the [/scan](#/operation/AcceptScanRequest] operations. The format of the
// identifier is not imposed but it should be unique enough to prevent collisons when polling for scan reports.
type ScanRequestId = string

// ScanResponse defines model for ScanResponse.
type ScanResponse struct {
	// Id A unique identifier returned by the [/scan](#/operation/AcceptScanRequest] operations. The format of the
	// identifier is not imposed but it should be unique enough to prevent collisons when polling for scan reports.
	Id ScanRequestId `json:"id"`
}

// Scanner Basic scanner properties such as name, vendor, and version.
type Scanner struct {
	// Name The name of the scanner.
	Name *string `json:"name,omitempty"`

	// Vendor The name of the scanner's provider.
	Vendor *string `json:"vendor,omitempty"`

	// Version The version of the scanner.
	Version *string `json:"version,omitempty"`
}

// ScannerAdapterMetadata Represents metadata of a Scanner Adapter which allows Harbor to lookup a scanner capabilities
// of scanning a given Artifact stored in its registry and making sure that it
// can interpret a returned result.
type ScannerAdapterMetadata struct {
	Capabilities []ScannerCapability `json:"capabilities"`

	// Properties A set of custom properties that can further describe capabilities of a given scanner.
	Properties *ScannerProperties `json:"properties,omitempty"`

	// Scanner Basic scanner properties such as name, vendor, and version.
	Scanner Scanner `json:"scanner"`
}

// ScannerCapability Capability consists of the set of recognized artifact MIME types and the set of scanner report MIME types.
//
// For example, a scanner capable of analyzing Docker images and producing a vulnerabilities report recognizable
// by Harbor web console might be represented with the following capability:
// - consumes MIME types:
//   - `application/vnd.oci.image.manifest.v1+json`
//   - `application/vnd.docker.distribution.manifest.v2+json`
//
// - produces MIME types:
//   - `application/vnd.scanner.adapter.vuln.report.harbor+json; version=1.0`
//
// For example, a scanner capable of analyzing artifacts and producing a sbom report recognizable
// by Harbor might be represented with the following capability:
// - type: sbom
// - consumes MIME types:
//   - `application/vnd.oci.image.manifest.v1+json`
//   - `application/vnd.docker.distribution.manifest.v2+json`
//
// - produces MIME types:
//   - `application/vnd.security.sbom.report+json; version=1.0`
type ScannerCapability struct {
	AdditionalAttributes *map[string]interface{} `json:"additional_attributes,omitempty"`

	// ConsumesMimeTypes The set of MIME types of the artifacts supported by the scanner to produce the reports specified in the "produces_mime_types". A given
	// mime type should only be present in one capability item.
	ConsumesMimeTypes []string `json:"consumes_mime_types"`

	// ProducesMimeTypes The set of MIME types of reports generated by the scanner for the consumes_mime_types of the same capability record.
	ProducesMimeTypes []string `json:"produces_mime_types"`

	// Type The type of the capability, for example, 'vulnerability' represents analyzing the artifact then producing the vulnerabilities report,
	// 'sbom' represents generating the corresponding sbom for the artifact which be scanned. In order to the backward and forward compatible,
	// the field is optional, we think it's a original 'vulnerability' scan if no such field.
	Type *ScannerCapabilityType `json:"type,omitempty"`
}

// ScannerCapabilityType The type of the capability, for example, 'vulnerability' represents analyzing the artifact then producing the vulnerabilities report,
// 'sbom' represents generating the corresponding sbom for the artifact which be scanned. In order to the backward and forward compatible,
// the field is optional, we think it's a original 'vulnerability' scan if no such field.
type ScannerCapabilityType string

// ScannerProperties A set of custom properties that can further describe capabilities of a given scanner.
type ScannerProperties map[string]string

// Severity A standard scale for measuring the severity of a vulnerability.
//
//   - `Unknown` - either a security problem that has not been assigned to a priority yet or a priority that the
//     scanner did not recognize.
//   - `Negligible` - technically a security problem, but is only theoretical in nature, requires a very special
//     situation, has almost no install base, or does no real damage.
//   - `Low` - a security problem, but is hard to exploit due to environment, requires a user-assisted attack,
//     a small install base, or does very little damage.
//   - `Medium` - a real security problem, and is exploitable for many people. Includes network daemon denial of
//     service attacks, cross-site scripting, and gaining user privileges.
//   - `High` - a real problem, exploitable for many people in a default installation. Includes serious remote denial
//     of service, local root privilege escalations, or data loss.
//   - `Critical` - a world-burning problem, exploitable for nearly all people in a default installation. Includes
//     remote root privilege escalations, or massive data loss.
type Severity string

// VulnerabilityItem defines model for VulnerabilityItem.
type VulnerabilityItem struct {
	// CweIds The Common Weakness Enumeration Identifiers associated with this vulnerability.
	CweIds *[]string `json:"cwe_ids,omitempty"`

	// Description The detailed description of the vulnerability.
	Description *string `json:"description,omitempty"`

	// FixVersion The version of the package containing the fix if available.
	FixVersion *string `json:"fix_version,omitempty"`

	// Id The unique identifier of the vulnerability.
	Id *string `json:"id,omitempty"`

	// Links The list of links to the upstream databases with the full description of the vulnerability.
	Links *[]string `json:"links,omitempty"`

	// Package An operating system package containing the vulnerability.
	Package       *string      `json:"package,omitempty"`
	PreferredCvss *CVSSDetails `json:"preferred_cvss,omitempty"`

	// Severity A standard scale for measuring the severity of a vulnerability.
	//
	// * `Unknown` - either a security problem that has not been assigned to a priority yet or a priority that the
	//   scanner did not recognize.
	// * `Negligible` - technically a security problem, but is only theoretical in nature, requires a very special
	//   situation, has almost no install base, or does no real damage.
	// * `Low` - a security problem, but is hard to exploit due to environment, requires a user-assisted attack,
	//   a small install base, or does very little damage.
	// * `Medium` - a real security problem, and is exploitable for many people. Includes network daemon denial of
	//   service attacks, cross-site scripting, and gaining user privileges.
	// * `High` - a real problem, exploitable for many people in a default installation. Includes serious remote denial
	//   of service, local root privilege escalations, or data loss.
	// * `Critical` - a world-burning problem, exploitable for nearly all people in a default installation. Includes
	//   remote root privilege escalations, or massive data loss.
	Severity         *Severity               `json:"severity,omitempty"`
	VendorAttributes *map[string]interface{} `json:"vendor_attributes,omitempty"`

	// Version The version of the package containing the vulnerability.
	Version *string `json:"version,omitempty"`
}

// GetScanReportParams defines parameters for GetScanReport.
type GetScanReportParams struct {
	// SbomMediaType media_type specifies the format of SBOM to be retrieved from the scanner adapter, it should either SPDX SBOM or CycloneDX
	SbomMediaType *string `form:"sbom_media_type,omitempty" json:"sbom_media_type,omitempty"`
	Accept        *string `json:"Accept,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get scanner metadata
	// (GET /metadata)
	GetMetadata(ctx echo.Context) error
	// Accept artifact scanning request
	// (POST /scan)
	AcceptScanRequest(ctx echo.Context) error
	// Get scan report
	// (GET /scan/{scan_request_id}/report)
	GetScanReport(ctx echo.Context, scanRequestId ScanRequestId, params GetScanReportParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetMetadata converts echo context to params.
func (w *ServerInterfaceWrapper) GetMetadata(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetMetadata(ctx)
	return err
}

// AcceptScanRequest converts echo context to params.
func (w *ServerInterfaceWrapper) AcceptScanRequest(ctx echo.Context) error {
	var err error

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(BearerAuthScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AcceptScanRequest(ctx)
	return err
}

// GetScanReport converts echo context to params.
func (w *ServerInterfaceWrapper) GetScanReport(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scan_request_id" -------------
	var scanRequestId ScanRequestId

	err = runtime.BindStyledParameterWithLocation("simple", false, "scan_request_id", runtime.ParamLocationPath, ctx.Param("scan_request_id"), &scanRequestId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scan_request_id: %s", err))
	}

	ctx.Set(BasicAuthScopes, []string{})

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetScanReportParams
	// ------------- Optional query parameter "sbom_media_type" -------------

	err = runtime.BindQueryParameter("form", true, false, "sbom_media_type", ctx.QueryParams(), &params.SbomMediaType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sbom_media_type: %s", err))
	}

	headers := ctx.Request().Header
	// ------------- Optional header parameter "Accept" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("Accept")]; found {
		var Accept string
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for Accept, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "Accept", runtime.ParamLocationHeader, valueList[0], &Accept)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter Accept: %s", err))
		}

		params.Accept = &Accept
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetScanReport(ctx, scanRequestId, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/metadata", wrapper.GetMetadata)
	router.POST(baseURL+"/scan", wrapper.AcceptScanRequest)
	router.GET(baseURL+"/scan/:scan_request_id/report", wrapper.GetScanReport)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"DkwAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
