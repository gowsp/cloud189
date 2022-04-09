package web

import (
	"github.com/gowsp/cloud189/pkg/invoker"
	"github.com/gowsp/cloud189/pkg/util"
)

// userInfoBase: "".concat(r.apiBaseUrl, "/open/user/getUserInfoForPortal.action"),
// userInfoExt: "".concat(r.apiBaseUrl, "/open/user/getUserInfoExt.action"),
// userPrivileges: "".concat(r.apiBaseUrl, "/open/user/getUserPrivileges.action"),
// listCloudFiles: "".concat(r.apiBaseUrl, "/open/file/listFiles.action"),
// listRecycleBinFiles: "".concat(r.apiBaseUrl, "/open/file/listRecycleBinFiles.action"),
// listLatestUploadFiles: "".concat(r.apiBaseUrl, "/portal/listLatestUploadFiles.action"),
// searchFiles: "".concat(r.apiBaseUrl, "/open/file/searchFiles.action"),
// getFileInfo: "".concat(r.apiBaseUrl, "/portal/getFileInfo.action"),
// getFolderInfo: "".concat(r.apiBaseUrl, "/open/file/getFolderInfo.action"),
// getCloudFileUploadUrl: "".concat(r.apiBaseUrl, "/open/file/getUploadWebFileUrl.action"),
// renameFolder: "".concat(r.apiBaseUrl, "/open/file/renameFolder.action"),
// renameFile: "".concat(r.apiBaseUrl, "/open/file/renameFile.action"),
// createFolder: "".concat(r.apiBaseUrl, "/open/file/createFolder.action"),
// createFolderStru: "".concat(r.apiBaseUrl, "/portal/createFolders.action"),
// listShareDir: "".concat(r.apiBaseUrl, "/open/share/listShareDir.action"),
// listShares: "".concat(r.apiBaseUrl, "/portal/listShares.action"),
// getShareFileDetails: "".concat(r.apiBaseUrl, "/open/file/getShareFileDetails.action"),
// createShareLink: "".concat(r.apiBaseUrl, "/open/share/createShareLink.action"),
// checkAccessCode: "".concat(r.apiBaseUrl, "/open/share/checkAccessCode.action"),
// getShareInfoByCode: "".concat(r.apiBaseUrl, "/open/share/getShareInfoByCodeV2.action"),
// cancelShareFile: "".concat(r.apiBaseUrl, "/portal/cancelShare.action"),
// shareReport: "".concat(r.apiBaseUrl, "/portal/reportShare.action"),
// generateRsaKey: "".concat(r.apiBaseUrl, "/security/generateRsaKey.action"),
// updateUserInfoExt: "".concat(r.apiBaseUrl, "/updateUserInfoExt.action"),
// getPrivateSafeMobile: "".concat(r.apiBaseUrl, "/portal/getPrivateSafeMobile.action"),
// isSafeMobileNull: "".concat(r.apiBaseUrl, "/portal/isSafeMobileNull.action"),
// getImageCode: "".concat(r.apiBaseUrl, "/portal/image.action"),
// VerifyImageCode: "".concat(r.apiBaseUrl, "/portal/verifyCode.action"),
// bondSafeMobile: "".concat(r.apiBaseUrl, "/portal/bondSafeMobile.action"),
// getAccessKey: "".concat(r.apiBaseUrl, "/portal/getAccessKey.action"),
// validateBondMobilePass: "".concat(r.apiBaseUrl, "/portal/validateBondMobilePass.action"),
// validateSafePass: "".concat(r.apiBaseUrl, "/portal/validateSafePass.action"),
// isSetQuestions: "".concat(r.apiBaseUrl, "/portal/isSetQuestions.action"),
// listAllQuestions: "".concat(r.apiBaseUrl, "/portal/listAllQuestions.action"),
// listMyQuestions: "".concat(r.apiBaseUrl, "/portal/listMyQuestions.action"),
// saveMyQuestions: "".concat(r.apiBaseUrl, "/portal/saveMyQuestions.action"),
// verifyMyQuestions: "".concat(r.apiBaseUrl, "/portal/verifyMyQuestions.action"),
// listGrow: "".concat(r.apiBaseUrl, "/portal/listGrow.action"),
// listRookieTask: "".concat(r.apiBaseUrl, "/portal/listRookieTask.action"),
// getFileDownloadUrl: "".concat(r.apiBaseUrl, "/open/file/getFileDownloadUrl.action"),
// getClientByType: "".concat(r.apiBaseUrl, "/portal/getClientByType.action"),
// getUserBriefInfo: "".concat(r.apiBaseUrl, "/portal/v2/getUserBriefInfo.action"),
// logout: "".concat(r.apiBaseUrl, "/portal/logout.action"),
// getWebImUrl: "".concat(r.apiBaseUrl, "/portal/getWebImUrl.action"),
// listContacts: "".concat(r.apiBaseUrl, "/portal/listContacts.action"),
// createPrivateShare: "".concat(r.apiBaseUrl, "/portal/createPrivateShare.action"),
// portalListFiles: "".concat(r.apiBaseUrl, "/portal/listFiles.action"),
// getObjectFolderNodes: "".concat(r.apiBaseUrl, "/portal/getObjectFolderNodes.action"),
// getShareInfo: "".concat(r.apiBaseUrl, "/portal/getShareInfo.action"),
// getUserReadUserGuide: "".concat(r.apiBaseUrl, "/portal/getUserReadUserGuide.action"),
// doSpecifiedTask: "".concat(r.apiBaseUrl, "/portal/doSpecifiedTask.action"),
// updateUserReadUserGuide: "".concat(r.apiBaseUrl, "/portal/updateUserReadUserGuide.action"),
// isDoneTask: "".concat(r.apiBaseUrl, "/portal/isDoneTask.action"),
// increaseShareFileAccessCount: "".concat(r.apiBaseUrl, "/portal//share/increaseShareFileAccessCount.action"),
// logPreviewFile: "".concat(r.apiBaseUrl, "/portal/log/logPreviewFile.action"),
// timeStructure: "".concat(r.apiBaseUrl, "/elastic/getTimeStructure.action"),
// photoList: "".concat(r.apiBaseUrl, "/elastic/listPhotoFile.action"),
// getPhotoOpenLog: "".concat(r.apiBaseUrl, "/photo/getPhotoOpenLog.action"),
// getNewVlcVideoPlayUrl: "".concat(r.apiBaseUrl, "/portal/getNewVlcVideoPlayUrl.action")
type api struct {
	invoker    *invoker.Invoker
	sessionKey string
	conf       *invoker.Config
}

func NewApi(path string) *api {
	conf, _ := invoker.OpenConfig(path)
	api := &api{conf: conf}
	api.invoker = invoker.NewInvoker("https://cloud.189.cn/api", api.refresh, conf)
	return api
}

func NewMemApi(username, password string) *api {
	conf := &invoker.Config{User: &invoker.User{Name: username, Password: password}}
	api := &api{conf: conf}
	api.invoker = invoker.NewInvoker("https://cloud.189.cn/api", api.refresh, conf)
	return api
}

func (i *api) login(user *invoker.User) error {
	result, err := i.invoker.PwdLogin("https://cloud.189.cn/api/portal/loginUrl.action", nil, user)
	if err != nil {
		return err
	}
	resp, err := i.invoker.Fetch(result.ToUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	i.conf.User = user
	i.conf.SSON = result.SSON
	i.conf.Auth = i.invoker.Cookie("https://cloud.189.cn", "COOKIE_LOGIN_USER")
	return i.conf.Save()
}
func (i *api) refresh() error {
	resp, err := i.invoker.Fetch("https://cloud.189.cn/api/portal/loginUrl.action")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	cookies := i.invoker.Cookies(resp.Request.URL)
	user := util.FindCookie(cookies, "COOKIE_LOGIN_USER")
	if user != nil {
		i.conf.Auth = user.Value
		i.conf.Save()
		return nil
	}
	return i.Login(i.conf.User.Name, i.conf.User.Password)
}
