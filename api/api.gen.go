// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.1.0 DO NOT EDIT.
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

// AcceptScanRequestApplicationVndScannerAdapterScanRequestPlusJSONVersion10RequestBody defines body for AcceptScanRequest for application/vnd.scanner.adapter.scan.request+json; version=1.0 ContentType.
type AcceptScanRequestApplicationVndScannerAdapterScanRequestPlusJSONVersion10RequestBody = ScanRequest

// AcceptScanRequestApplicationVndScannerAdapterScanRequestPlusJSONVersion11RequestBody defines body for AcceptScanRequest for application/vnd.scanner.adapter.scan.request+json; version=1.1 ContentType.
type AcceptScanRequestApplicationVndScannerAdapterScanRequestPlusJSONVersion11RequestBody = ScanRequest

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

	err = runtime.BindStyledParameterWithOptions("simple", "scan_request_id", ctx.Param("scan_request_id"), &scanRequestId, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
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

		err = runtime.BindStyledParameterWithOptions("simple", "Accept", valueList[0], &Accept, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationHeader, Explode: false, Required: false})
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
