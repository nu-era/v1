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
const accountSid = "AC2da72a4a5bd7e5bc55fc7c694db69d1a"
const authToken = "f753aa5c96f96568079b6b52fbb32cb9"
const trialNum = "+14252121598"
const twilURLString = "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
const dcMsg = `You are receiving this message because New-Era has lost contact with your device. If this was not planned, please contact us immediatly.`
