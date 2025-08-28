package handlers

import (
	"context"
	"net/http"
	"strings"
	"wrench/app/auth"
	contexts "wrench/app/contexts"
	"wrench/app/manifest/api_settings"
)

type AuthValidatorHandler struct {
	Next             Handler
	EndpointSettings *api_settings.EndpointSettings
	ApiSettings      *api_settings.ApiSettings
}

func (handler *AuthValidatorHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	authorizationSettings := handler.ApiSettings.Authorization
	endpointSettings := handler.EndpointSettings

	if !endpointSettings.EnableAnonymous {

		if authorizationSettings.Type == api_settings.JWKSAuthorizationType {
			tokenString := wrenchContext.Request.Header.Get("Authorization")
			if len(tokenString) == 0 {
				handler.setHasError("Unauthorized", http.StatusUnauthorized, wrenchContext, bodyContext)
			} else {
				tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

				tokenIsValid := auth.JwksValidationAuthentication(ctx, tokenString, authorizationSettings)
				if tokenIsValid {
					tokenIsAuthorized := auth.JwksValidationAuthorization(tokenString, endpointSettings.Roles, endpointSettings.Scopes, endpointSettings.Claims)
					if !tokenIsAuthorized {
						handler.setHasError("Forbidden", http.StatusForbidden, wrenchContext, bodyContext)
					}
				} else {
					handler.setHasError("Unauthorized", http.StatusUnauthorized, wrenchContext, bodyContext)
				}
			}
		}

		if authorizationSettings.Type == api_settings.HMACAuthorizationType {
			isHMACValid := auth.HMACValidate(wrenchContext, bodyContext, authorizationSettings)
			if !isHMACValid {
				handler.setHasError("Unauthorized", http.StatusUnauthorized, wrenchContext, bodyContext)
			}
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *AuthValidatorHandler) SetNext(next Handler) {
	handler.Next = next
}

func (handler *AuthValidatorHandler) setHasError(msg string, httpStatusCode int, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	wrenchContext.SetHasError2()
	bodyContext.ContentType = "text/plain"
	bodyContext.HttpStatusCode = httpStatusCode
	bodyContext.SetBody([]byte(msg))
}
