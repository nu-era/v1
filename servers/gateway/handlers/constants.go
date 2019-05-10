package handlers

const headerContentType = "Content-Type"
const headerAccessControlAllowOrigin = "Access-Control-Allow-Origin"
const headerAccessControlAllowMethods = "Access-Control-Allow-Methods"
const headerAccessControlAllowHeader = "Access-Control-Allow-Headers"
const headerAccessControlExposeHeader = "Access-Control-Expose-Headers"
const headerXFrameOption = "X-Frame-Options"
const headerXForwarded = "X-Forwarded-For"

const contentTypeJSON = "application/json"
const contentTypeText = "text/plain"
const contentTypeHTML = "text/html"

// Twilio token and account SID
// TODO: Add to a .gitignorefile
const trialNum = "+14252121598"
const twilAuthString = "https://api.authy.com/protected/json/phones/verification/start"
const twilURLString = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
const dcMsg = `You are receiving this message because New-Era has lost contact with your device. If this was not planned, please contact us immediatly.`
