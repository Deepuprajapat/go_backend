package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/VI-IM/im_backend_go/internal/utils"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) GetProjectSEOContent(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	ctx := r.Context()
	vars := mux.Vars(r)
	url := vars["url"]

	// Get project by canonical URL
	project, err := h.app.GetProjectByCanonicalURL(ctx, url)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to fetch project", err.Error())
	}

	if project == nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Project not found", "")
	}

	// Generate canonical URL
	canonicalURL := fmt.Sprintf("https://investmango.com%s", r.RequestURI)

	// Get meta info
	metaInfo := project.MetaInfo

	// Generate HTML response
	htmlResponse := fmt.Sprintf(`<!DOCTYPE html>
<html>
	<head>
		<title>%s</title>
		<meta name="description" content="%s">
		<meta name="keywords" content="%s">
		<link rel="canonical" href="%s">
	</head>
	<body>
		<h1>%s</h1><br><br>
		<p>%s</p><br>
		<p><strong>Keywords:</strong> %s</p><br>
		<p><strong>Canonical URL:</strong> %s</p>
	</body>
</html>`,
		metaInfo.Title,
		metaInfo.Description,
		metaInfo.Keywords,
		canonicalURL,
		metaInfo.Title,
		strings.ReplaceAll(metaInfo.Description, "\n", "<br>"),
		metaInfo.Keywords,
		canonicalURL,
	)

	return &imhttp.Response{
		Data:       htmlResponse,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetPropertySEOContent(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	ctx := r.Context()
	url := r.URL.Query().Get("url")

	property, err := h.app.GetPropertyByName(ctx, url)
	if err != nil {
		return nil, err
	}

	if property == nil {
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Property not found with this url: "+url, "")
	}

	// Generate HTML response
	htmlResponse := fmt.Sprintf(`<!DOCTYPE html>
<html>
    <head>
        <title>%s</title>
        <meta name="description" content="%s">
        <meta name="keywords" content="%s">
        <link rel="canonical" href="%s">
    </head>
    <body>
        <h1>%s</h1><br><br>
        <p>%s</p><br>
        <p><strong>Keywords:</strong> %s</p><br>
        <p><strong>Canonical URL:</strong> %s</p>
    </body>
</html>`,
		property.MetaInfo.Title,
		property.MetaInfo.Description,
		property.MetaInfo.Keywords,
		url,
		property.MetaInfo.Title,
		strings.ReplaceAll(property.MetaInfo.Description, "\n", "<br>"),
		property.MetaInfo.Keywords,
		url)

	return &imhttp.Response{
		Data:       htmlResponse,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *Handler) GetHTMLContent(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	url := r.URL.Query().Get("url")
	ctx := r.Context()
	// ctx := r.Context()

	if strings.HasPrefix(url, "https://www.investmango.com/") {
		url = strings.TrimPrefix(url, "https://www.investmango.com/")
	} else if strings.HasPrefix(url, "https://investmango.com/") {
		url = strings.TrimPrefix(url, "https://investmango.com/")
	}

	cleanUrl := strings.TrimPrefix(url, "www.")
	cleanUrl = strings.TrimSuffix(cleanUrl, "/")
	cleanUrl = strings.TrimPrefix(cleanUrl, "/")
	logger.Get().Info().Msg("cleanUrl from gethtmlcontent: " + cleanUrl)

	var htmlResponse string

	if strings.Contains(cleanUrl, "propertyforsale/") {
		parts := strings.Split(cleanUrl, "/")
		if len(parts) >= 2 {
			propertyURL := parts[1]
			logger.Get().Info().Msg("propertyURL from gethtmlcontent: " + propertyURL)
			property, err := h.app.GetPropertyByName(ctx, propertyURL)
			logger.Get().Info().Msg("property from gethtmlcontent: " + property.MetaInfo.Title)

			if err != nil {
				return nil, err
			}
			if property == nil {
				return nil, imhttp.NewCustomErr(http.StatusNotFound, "Property not found with this url: "+url, "")
			}

			// Get OG image from product schema if available
			ogImage := utils.GetOgImageFromSchema(property.ProductSchema)
			htmlResponse = fmt.Sprintf(`<!DOCTYPE html>
	<html>
	    <head>
	        <meta name="google-site-verification" content="ItwxGLnb2pNSeyJn0kZsRa3DZxRZO3MSCQs5G3kTLgA">
	        <title>%s</title>
	        <meta name="description" content="%s">
	        <meta name="keywords" content="%s">
	        <meta property="og:title" content="%s">
	        <meta property="og:description" content="%s">
	        <meta property="og:image" content="%s">
	        <meta property="og:url" content="%s">
	        <meta property="og:type" content="website">
	        <meta name="robots" content="index, follow">
	        <link rel="canonical" href="%s">
	        %s
	    </head>
	    <body>
	        <h1>%s</h1>
	        <p>%s</p>
	    </body>
	</html>`,
				property.MetaInfo.Title,
				property.MetaInfo.Description,
				property.MetaInfo.Keywords,
				property.MetaInfo.Title,
				property.MetaInfo.Description,
				ogImage,
				url,
				url,
				property.ProductSchema,
				property.MetaInfo.Title,
				strings.ReplaceAll(property.MetaInfo.Description, "\n", "<br>"))

			return &imhttp.Response{
				Data:       htmlResponse,
				StatusCode: http.StatusOK,
			}, nil
		}
	}

	switch cleanUrl {
	case "", "/":
		htmlResponse = fmt.Sprintf(`<!DOCTYPE html><html><head>
            <meta name="google-site-verification" content="ItwxGLnb2pNSeyJn0kZsRa3DZxRZO3MSCQs5G3kTLgA">
            <title>Real Estate Portfolio &amp; Strategic Management Company</title>
            <meta name="description" content="Invest Mango: Real estate portfolio and strategic management services. Elevate your financial future through informed decisions and prime opportunities.">
            <meta property="og:title" content="Real Estate Portfolio &amp; Strategic Management Company">
            <meta property="og:description" content="Invest Mango: Real estate portfolio and strategic management services. Elevate your financial future through informed decisions and prime opportunities.">
            <meta name="robots" content="index, follow">
            <meta property="og:url" content="%s">
            <link rel="canonical" href="https://www.investmango.com/">
            </head><body><h1>Welcome to Invest Mango</h1></body></html>`, url)

	case "contact":
		htmlResponse = fmt.Sprintf(`<!DOCTYPE html><html><head>
            <meta name="google-site-verification" content="ItwxGLnb2pNSeyJn0kZsRa3DZxRZO3MSCQs5G3kTLgA">
            <title>Contact Us</title>
            <meta name="description" content="For more details on On-Site Visits, Locations, Developers' Information, Property Age, Documentation, Bank Assistance, and many more feel free to connect.">
            <meta property="og:title" content="Contact Us">
            <meta property="og:description" content="For more details on On-Site Visits, Locations, Developers' Information, Property Age, Documentation, Bank Assistance, and many more feel free to connect.">
            <meta name="robots" content="index, follow">
            <meta property="og:url" content="%s">
            <link rel="canonical" href="https://www.investmango.com/contact">
            </head></html>`, url)

	case "career", "career/":
		htmlResponse = fmt.Sprintf(`<!DOCTYPE html><html><head>
            <meta name="google-site-verification" content="ItwxGLnb2pNSeyJn0kZsRa3DZxRZO3MSCQs5G3kTLgA">
            <title>Careers | Invest Mango</title>
            <meta name="description" content="Invest Mango offers a dynamic culture that is all about winning and innovating oneself. Unleash your true potential much more than you have ever imagined.">
            <meta property="og:title" content="Careers | Invest Mango">
            <meta property="og:description" content="Invest Mango offers a dynamic culture that is all about winning and innovating oneself. Unleash your true potential much more than you have ever imagined.">
            <meta name="robots" content="index, follow">
            <meta property="og:url" content="%s">
            <link rel="canonical" href="https://www.investmango.com/career">
            </head></html>`, url)

	case "about", "about/":
		htmlResponse = fmt.Sprintf(`<!DOCTYPE html><html><head>
            <meta name="google-site-verification" content="ItwxGLnb2pNSeyJn0kZsRa3DZxRZO3MSCQs5G3kTLgA">
            <title>About Us – Know About Invest Mango</title>
            <meta name="description" content="Invest Mango | Reputed Investment and Real Estate Portfolio Management Organisation. We provide Real Estate Consulting services in Delhi NCR.">
            <meta name="keywords" content="real estate investment consultant, investment and portfolio management, top real estate agents in noida, property investment consultant, commercial property management, real estate consultants services, real estate advisory services, top real estate companies in noida.">
            <meta property="og:title" content="About Us – Know About Invest Mango">
            <meta property="og:description" content="Invest Mango | Reputed Investment and Real Estate Portfolio Management Organisation. We provide Real Estate Consulting services in Delhi NCR.">
            <meta property="og:url" content="%s">
            <meta name="robots" content="index, follow">
            <link rel="canonical" href="https://www.investmango.com/about">
            </head></html>`, url)

	case "sitemap":
		htmlResponse = fmt.Sprintf(`<!DOCTYPE html><html><head>
            <title>About Us – Know About Invest Mango</title>
            <meta property="og:title" content="About Us – Know About Invest Mango">
            <meta property="og:description" content="Invest Mango | Reputed Investment and Real Estate Portfolio Management Organisation. We provide Real Estate Consulting services in Delhi NCR.">
            <meta property="og:url" content="%s">
            <link rel="canonical" href="https://www.investmango.com/sitemap">
            </head></html>`, url)

	case "privacy-policy", "privacy-policy/":
		htmlResponse = "<!DOCTYPE html><html><head>" +
			"<meta name=\"google-site-verification\" content=\"ItwxGLnb2pNSeyJn0kZsRa3DZxRZO3MSCQs5G3kTLgA\">" +
			"<title>Privacy Policy | Invest Mango</title>" +
			"<meta name=\"description\" content=\"Read the privacy policy of Invest Mango to learn how we collect, use, and protect your data.\">" +
			"<meta property=\"og:title\" content=\"Privacy Policy | Invest Mango\">" +
			"<meta property=\"og:description\" content=\"Read the privacy policy of Invest Mango to learn how we collect, use, and protect your data.\">" +
			"<meta name=\"robots\" content=\"index, follow\">" +
			"<meta property=\"og:url\" content=\"" + url + "\">" +
			"<link rel=\"canonical\" href=\"https://www.investmango.com/privacy-policy\">" +
			"</head><body>" +
			"<div style='max-width: 800px; margin: auto; padding: 40px; font-family: sans-serif; color: #333'>" +
			"<h1 style='text-align:center;'>Privacy Policy</h1>" +
			"<p><strong>Effective Date:</strong> 10-05-2025</p>" +
			"<p style='margin-top: 20px;'>At Invest Mango, we are committed to protecting your privacy and ensuring that your personal information is handled safely and responsibly. This Privacy Policy outlines how we collect, use, disclose, and protect the information you provide to us when you visit our website, use our services, or interact with us in any way.</p>" +
			"<h2 style='margin-top: 30px;'>1. Information Collected From Interactive Forms</h2>" +
			"<p>When you voluntarily send us electronic mail / fillup the form, we will keep a record of this information so that we can respond to you. We only collect information from you when you register on our site or fill out a form. Also, when filling out a form on our site, you may be asked to enter your: name, e-mail address or phone number. You may, however, visit our site anonymously. In case you have submitted your personal information and contact details, we reserve the rights to Call, SMS, Email or WhatsApp about our products and offers, even if your number has DND activated on it.</p>" +
			"<h2 style='margin-top: 30px;'>2. Information We Collect</h2>" +
			"<p>We may collect personal information such as your name, email address, phone number, location or address, investment preferences, and any other information you voluntarily provide via forms, chats, or phone calls. In addition, we automatically collect data such as your IP address, browser type and version, device information, pages visited, time spent on the site, and data collected through cookies and tracking technologies.</p>" +
			"<h2 style='margin-top: 30px;'>3. How We Use Your Information</h2>" +
			"<p>Your data is used to provide property recommendations and investment options, respond to your inquiries or service requests, send newsletters, updates, and promotional offers (only if you have opted in), improve website performance and user experience, conduct market research and analysis, and comply with legal obligations.</p>" +
			"<h2 style='margin-top: 30px;'>4. Sharing Your Information</h2>" +
			"<p>We do not sell your personal information. However, we may share it with trusted developers and real estate partners to facilitate property inquiries, service providers who support our marketing, data analysis, or IT functions, and legal or regulatory authorities when required by law or for legal processes.</p>" +
			"<h2 style='margin-top: 30px;'>5. Data Security</h2>" +
			"<p>We implement industry-standard security measures to protect your information from unauthorized access, alteration, disclosure, or destruction. Despite our efforts, no method of online transmission or storage is completely secure.</p>" +
			"<h2 style='margin-top: 30px;'>6. Cookies & Tracking Technologies</h2>" +
			"<p>We use cookies to understand user behavior, personalize content, and enable essential website functionalities. You can manage or disable cookies through your browser settings at any time.</p>" +
			"<h2 style='margin-top: 30px;'>7. Your Rights</h2>" +
			"<p>You have the right to access, correct, or delete your personal data, opt out of marketing communications, and request a copy of the data we hold about you. To exercise these rights, please contact us at: <a href='mailto:info@investmango.com'>info@investmango.com</a></p>" +
			"<h2 style='margin-top: 30px;'>8. Third Party Links</h2>" +
			"<p>Our website may include links to third-party websites. We are not responsible for the privacy practices or content of those sites.</p>" +
			"<h2 style='margin-top: 30px;'>9. Children's Privacy</h2>" +
			"<p>Our services are not intended for children under the age of 18, and we do not knowingly collect personal information from minors.</p>" +
			"<h2 style='margin-top: 30px;'>10. Changes to this Policy</h2>" +
			"<p>We may update this Privacy Policy from time to time. Any changes will be posted on this page along with the revised effective date.</p>" +
			"<h2 style='margin-top: 30px;'>11. Contact Us</h2>" +
			"<p><strong>Invest Mango</strong><br/>" +
			"Email: <a href='mailto:info@investmango.com'>info@investmango.com</a><br/>" +
			"Website: <a href='https://www.investmango.com' target='_blank'>https://www.investmango.com</a></p>" +
			"</div></body></html>"

	case "terms-and-conditions", "terms-and-conditions/":
		htmlResponse = "<!DOCTYPE html><html><head>" +
			"<meta name=\"google-site-verification\" content=\"ItwxGLnb2pNSeyJn0kZsRa3DZxRZO3MSCQs5G3kTLgA\">" +
			"<title>Terms & Conditions | Invest Mango</title>" +
			"<meta name=\"description\" content=\"Read the Terms & Conditions for using Invest Mango's website and services.\">" +
			"<meta property=\"og:title\" content=\"Terms & Conditions | Invest Mango\">" +
			"<meta property=\"og:description\" content=\"Read the Terms & Conditions for using Invest Mango's website and services.\">" +
			"<meta name=\"robots\" content=\"index, follow\">" +
			"<meta property=\"og:url\" content=\"" + url + "\">" +
			"<link rel=\"canonical\" href=\"https://www.investmango.com/terms-and-conditions\">" +
			"</head><body>" +
			"<div style='max-width: 800px; margin: auto; padding: 40px; font-family: sans-serif; color: #333'>" +
			"<h1 style='text-align:center;'>Terms of Service</h1>" +
			"<p><strong>Effective Date:</strong> 10-05-2025</p>" +
			"<p style='margin-top: 20px;'>Welcome to Invest Mango. By accessing or using our website or services, you agree to comply with and be bound by the following terms and conditions. Please read them carefully.</p>" +
			"<h2 style='margin-top: 30px;'>1. Acceptance of Terms</h2>" +
			"<p>By using this website and our services, you agree to these Terms & Conditions and our Privacy Policy. If you do not agree, please refrain from using the site.</p>" +
			"<h2 style='margin-top: 30px;'>2. Services Offered</h2>" +
			"<p>Invest Mango is a real estate advisory platform offering:</p>" +
			"<ul>" +
			"<li>Property listings (residential & commercial)</li>" +
			"<li>Investment consultancy</li>" +
			"<li>Developer tie-ups and project promotions</li>" +
			"<li>Site visit coordination</li>" +
			"<li>Pre-launch and resale opportunities</li>" +
			"</ul>" +
			"<p>Note: We act as facilitators and not as developers or builders.</p>" +
			"<h2 style='margin-top: 30px;'>3. User Responsibilities</h2>" +
			"<p>You agree to:</p>" +
			"<ul>" +
			"<li>Provide accurate and complete information during inquiries</li>" +
			"<li>Use the website for lawful purposes only</li>" +
			"<li>Refrain from submitting false or misleading information</li>" +
			"<li>Not copy, reproduce, or misuse any content from the website</li>" +
			"</ul>" +
			"<h2 style='margin-top: 30px;'>4. Disclaimer</h2>" +
			"<ul>" +
			"<li>We make reasonable efforts to ensure information on our website is accurate, but we do not guarantee its completeness or suitability for any purpose.</li>" +
			"<li>Real estate pricing, availability, and offers are subject to change without notice.</li>" +
			"<li>Invest Mango is not liable for decisions made based on our advice or third-party listings.</li>" +
			"</ul>" +
			"<h2 style='margin-top: 30px;'>5. Intellectual Property</h2>" +
			"<p>All website content, including logos, text, images, graphics, and layout, is the property of Invest Mango or its content providers. Unauthorized use is prohibited.</p>" +
			"<h2 style='margin-top: 30px;'>6. Limitation of Liability</h2>" +
			"<p>Invest Mango will not be liable for:</p>" +
			"<ul>" +
			"<li>Any direct, indirect, incidental, or consequential damage</li>" +
			"<li>Inaccuracies or omissions in property data</li>" +
			"<li>Third-party service interruptions or errors</li>" +
			"<li>Investment losses resulting from reliance on our content</li>" +
			"</ul>" +
			"<h2 style='margin-top: 30px;'>7. Third Party Links</h2>" +
			"<p>Our platform may contain links to third-party websites. We are not responsible for the content, policies, or reliability of those sites.</p>" +
			"<h2 style='margin-top: 30px;'>8. Termination</h2>" +
			"<p>We reserve the right to:</p>" +
			"<ul>" +
			"<li>Terminate access to our website/services at any time without notice</li>" +
			"<li>Remove any content that violates our policies or legal obligations</li>" +
			"</ul>" +
			"<h2 style='margin-top: 30px;'>9. Changes to Terms</h2>" +
			"<p>We may update these Terms & Conditions at any time. Updates will be posted on this page and your continued use constitutes acceptance of the new terms.</p>" +
			"<h2 style='margin-top: 30px;'>10. Governing Law</h2>" +
			"<p>These Terms shall be governed by and construed in accordance with the laws of India, and any disputes will be subject to the exclusive jurisdiction of courts in Noida, Pune, or Gurugram.</p>" +
			"<h2 style='margin-top: 30px;'>Contact Us</h2>" +
			"<p>Email: <a href='mailto:info@investmango.com'>info@investmango.com</a>, <a href='mailto:hr@investmango.com'>hr@investmango.com</a></p>" +
			"</div></body></html>"

	default:
		return nil, imhttp.NewCustomErr(http.StatusNotFound, "Page not found", "")
	}

	return &imhttp.Response{
		Data:       htmlResponse,
		StatusCode: http.StatusOK,
	}, nil
}
