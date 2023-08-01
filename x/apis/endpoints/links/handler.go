package links

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/forbole/juno/v5/types/config"
	"google.golang.org/api/firebasedynamiclinks/v1"
	"google.golang.org/api/option"

	"github.com/desmos-labs/djuno/v2/x/apis/utils"
)

type Handler struct {
	cfg           *Config
	firebaseLinks *firebasedynamiclinks.Service
}

func NewHandler(junoCfg config.Config) *Handler {
	cfgBz, err := junoCfg.GetBytes()
	if err != nil {
		panic(err)
	}

	cfg, err := ParseConfig(cfgBz)
	if err != nil {
		panic(err)
	}

	var firebaseLinksService *firebasedynamiclinks.Service
	if cfg.FirebaseCredentialsFilePath != "" {
		dynamicLinksService, err := firebasedynamiclinks.NewService(context.Background(), option.WithCredentialsFile(cfg.FirebaseCredentialsFilePath))
		if err != nil {
			panic(err)
		}
		firebaseLinksService = dynamicLinksService
	}

	return &Handler{
		cfg:           cfg,
		firebaseLinks: firebaseLinksService,
	}
}

// ParseGenerateLinkRequest parses the given body into a GenerateLinkRequest
func (h *Handler) ParseGenerateLinkRequest(body []byte) (*GenerateLinkRequest, error) {
	var req GenerateLinkRequest
	err := json.Unmarshal(body, &req)
	return &req, err
}

// ValidateLinkRequest validates the given GenerateLinkRequest and returns an error if something is wrong
func (h *Handler) ValidateLinkRequest(req *GenerateLinkRequest) error {
	// Parse the URL and get the values we need
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return utils.WrapErr(http.StatusBadRequest, fmt.Sprintf("error while parsing the URL: %s", err))
	}

	urlPath := strings.TrimPrefix(parsedURL.Path, "/")
	params := parsedURL.Query()

	// Find the path inside the registered ones
	var path *Path
	for i, registeredPath := range RegisteredPaths {
		if strings.EqualFold(registeredPath.Path, urlPath) {
			path = &RegisteredPaths[i]
		}
	}

	if path == nil {
		return utils.WrapErr(http.StatusBadRequest, "invalid path")
	}

	// Find the action inside the registered ones
	var action *Action
	for i, registeredActions := range path.Actions {
		if strings.EqualFold(registeredActions.Name, params.Get(ParamAction)) {
			action = &path.Actions[i]
		}
	}

	if action == nil {
		return utils.WrapErr(http.StatusBadRequest, "invalid action")
	}

	// Make sure the action required params are all present
	for _, requiredParam := range action.RequiredParams {
		if !params.Has(requiredParam) {
			return utils.WrapErr(http.StatusBadRequest, fmt.Sprintf("missing required param: %s", requiredParam))
		}
	}

	return nil
}

func (h *Handler) GetLinkDesktopInfo() *firebasedynamiclinks.DesktopInfo {
	return &firebasedynamiclinks.DesktopInfo{
		DesktopFallbackLink: h.cfg.Desktop.FallbackLink,
	}
}

func (h *Handler) GetLinkAndroidInfo() *firebasedynamiclinks.AndroidInfo {
	return &firebasedynamiclinks.AndroidInfo{
		AndroidPackageName:           h.cfg.Android.PackageName,
		AndroidMinPackageVersionCode: h.cfg.Android.MinPackageVersionCode,
	}
}

func (h *Handler) GetLinkIOSInfo() *firebasedynamiclinks.IosInfo {
	return &firebasedynamiclinks.IosInfo{
		IosBundleId:       h.cfg.Ios.BundleID,
		IosMinimumVersion: h.cfg.Ios.MinimumVersion,
		IosAppStoreId:     h.cfg.Ios.AppStoreID,
	}
}

// HandleGenerateLinkRequest handles the given GenerateLinkRequest and returns a GenerateLinkResponse,
// or an error if something goes wrong
func (h *Handler) HandleGenerateLinkRequest(req *GenerateLinkRequest) (*GenerateLinkResponse, error) {
	if h.firebaseLinks == nil {
		return nil, nil
	}

	// Generate the link
	res, err := h.firebaseLinks.ShortLinks.Create(&firebasedynamiclinks.CreateShortDynamicLinkRequest{
		DynamicLinkInfo: &firebasedynamiclinks.DynamicLinkInfo{
			DynamicLinkDomain: h.cfg.Domain,

			DesktopInfo: h.GetLinkDesktopInfo(),
			AndroidInfo: h.GetLinkAndroidInfo(),
			IosInfo:     h.GetLinkIOSInfo(),

			Link: fmt.Sprintf("%s/%s", h.cfg.Domain, req.URL),

			// TODO: Integrate this properly
			// SocialMetaTagInfo: &firebasedynamiclinks.SocialMetaTagInfo{
			//	 SocialTitle:       event.Name,
			//	 SocialDescription: event.Description,
			//	 SocialImageLink:   event.CoverPictureUrl,
			// },
		},
		Suffix: &firebasedynamiclinks.Suffix{
			Option: "SHORT",
		},
	}).Do()
	if err != nil {
		return nil, err
	}

	// Return the short link
	return NewGenerateLinkResponse(res.ShortLink), nil
}
