package handler

import (
	"boysitorus/Progamify-Restful-API/internal/config"
	"boysitorus/Progamify-Restful-API/internal/model"
	"boysitorus/Progamify-Restful-API/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// Get the Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization header is required"})
		return
	}

	// Extract the token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid token format"})
		return
	}

	// Blacklist the token
	config.BlacklistToken(tokenParts[1])

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req model.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Respond with user and token
	c.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

func (h *AuthHandler) CustomRender(c *gin.Context) {
	// Parse the input HTML
	var input = `<p><figure data-trix-attachment="{&quot;contentType&quot;:&quot;image/png&quot;,&quot;filename&quot;:&quot;TestingImage4.png&quot;,&quot;filesize&quot;:33867,&quot;height&quot;:514,&quot;href&quot;:&quot;http://localhost:8000/storage/lessons/images/93lMciqmSlgqG0im0umLdbQDv5TxpQuwNGm3JiPS.png&quot;,&quot;url&quot;:&quot;http://localhost:8000/storage/lessons/images/93lMciqmSlgqG0im0umLdbQDv5TxpQuwNGm3JiPS.png&quot;,&quot;width&quot;:514}" data-trix-content-type="image/png" data-trix-attributes="{&quot;presentation&quot;:&quot;gallery&quot;}" class="attachment attachment--preview attachment--png"><a href="http://localhost:8000/storage/lessons/images/93lMciqmSlgqG0im0umLdbQDv5TxpQuwNGm3JiPS.png"><img src="http://localhost:8000/storage/lessons/images/93lMciqmSlgqG0im0umLdbQDv5TxpQuwNGm3JiPS.png" width="514" height="514"><figcaption class="attachment__caption"><span class="attachment__name">TestingImage4.png</span> <span class="attachment__size">33.07 KB</span></figcaption></a></figure></p><h2><strong>What is Lorem Ipsum?</strong></h2><p><strong>Lorem Ipsum</strong> is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</p><h2>Why do we use it?</h2><p>It is a long established fact that a reader will be distracted by the readable content of a page when looking at its layout. The point of using Lorem Ipsum is that it has a more-or-less normal distribution of letters, as opposed to using 'Content here, content here', making it look like readable English. Many desktop publishing packages and web page editors now use Lorem Ipsum as their default model text, and a search for 'lorem ipsum' will uncover many web sites still in their infancy. Various versions have evolved over the years, sometimes by accident, sometimes on purpose (injected humour and the like).</p><p><br></p><h2>Where does it come from?</h2><p>Contrary to popular belief, Lorem Ipsum is not simply random text. It has roots in a piece of classical Latin literature from 45 BC, making it over 2000 years old. Richard McClintock, a Latin professor at Hampden-Sydney College in Virginia, looked up one of the more obscure Latin words, consectetur, from a Lorem Ipsum passage, and going through the cites of the word in classical literature, discovered the undoubtable source. Lorem Ipsum comes from sections 1.10.32 and 1.10.33 of "de Finibus Bonorum et Malorum" (The Extremes of Good and Evil) by Cicero, written in 45 BC. This book is a treatise on the theory of ethics, very popular during the Renaissance. The first line of Lorem Ipsum, "Lorem ipsum dolor sit amet..", comes from a line in section 1.10.32.</p><p><br></p>`

	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		fmt.Println("Error parsing HTML:", err)
		return
	}

	// Traverse and process the HTML nodes
	var processNode func(*html.Node)
	processNode = func(n *html.Node) {
		// If the node is a <figure> tag, check for image link and replace it with <img>
		if n.Type == html.ElementNode && n.Data == "figure" {
			for _, attr := range n.Attr {
				if attr.Key == "data-trix-attachment" {
					// Extract image URL from data-trix-attachment JSON string
					// Simple way, this could be improved with a proper JSON parser if needed
					if strings.Contains(attr.Val, `"href"`) {
						startIdx := strings.Index(attr.Val, `"href":"`) + 8
						endIdx := strings.Index(attr.Val[startIdx:], `"`) + startIdx
						if startIdx > 0 && endIdx > startIdx {
							imageURL := attr.Val[startIdx:endIdx]
							// Replace figure tag with <img> tag
							newNode := &html.Node{
								Type: html.ElementNode,
								Data: "img",
								Attr: []html.Attribute{
									{Key: "src", Val: imageURL},
									{Key: "width", Val: "514"},
									{Key: "height", Val: "514"},
								},
							}
							// Replace the <figure> tag with the new <img> node
							n.Parent.InsertBefore(newNode, n)
							n.Parent.RemoveChild(n)
						}
					}
				}
			}
		}

		// Recursively process child nodes
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			processNode(c)
		}
	}

	// Start processing from the root node
	processNode(doc)

	// Render the modified HTML back to a string
	var builder strings.Builder
	err = html.Render(&builder, doc)
	if err != nil {
		fmt.Println("Error rendering HTML:", err)
		return
	}

	result := builder.String()

	c.JSON(http.StatusOK, gin.H{
		"result": result,
	})
}
