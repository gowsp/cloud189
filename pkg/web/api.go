package web

import "github.com/gowsp/cloud189/pkg/drive"

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
	invoker    *invoker
	sessionKey string
	conf       *drive.Config
}

func NewApi(path string) *api {
	conf, _ := drive.OpenConfig(path)
	i := newInvoker(conf)
	return &api{invoker: i, conf: conf}
}

func NewMemApi(username, password string) *api {
	conf := &drive.Config{User: drive.User{Name: username, Password: password}}
	i := newInvoker(conf)
	return &api{invoker: i, conf: conf}
}
