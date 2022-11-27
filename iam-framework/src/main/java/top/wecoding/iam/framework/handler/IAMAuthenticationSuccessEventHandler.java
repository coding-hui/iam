// package top.wecoding.iam.server.authentication.handler;
//
// import cn.hutool.core.map.MapUtil;
// import java.io.IOException;
// import java.time.temporal.ChronoUnit;
// import java.util.Map;
// import javax.annotation.PostConstruct;
// import javax.servlet.ServletException;
// import javax.servlet.http.HttpServletRequest;
// import javax.servlet.http.HttpServletResponse;
// import lombok.extern.slf4j.Slf4j;
// import org.springframework.http.converter.HttpMessageConverter;
// import org.springframework.http.server.ServletServerHttpResponse;
// import org.springframework.security.core.Authentication;
// import org.springframework.security.core.context.SecurityContextHolder;
// import org.springframework.security.oauth2.core.OAuth2AccessToken;
// import org.springframework.security.oauth2.core.OAuth2RefreshToken;
// import org.springframework.security.oauth2.core.endpoint.OAuth2AccessTokenResponse;
// import
// org.springframework.security.oauth2.core.http.converter.OAuth2AccessTokenResponseHttpMessageConverter;
// import
// org.springframework.security.oauth2.server.authorization.authentication.OAuth2AccessTokenAuthenticationToken;
// import org.springframework.security.web.authentication.AuthenticationSuccessHandler;
// import org.springframework.util.CollectionUtils;
// import top.wecoding.iam.common.constant.SecurityConstants;
// import top.wecoding.iam.common.userdetails.LoginUser;
// import top.wecoding.iam.server.authentication.token.AccessTokenResponseParametersConverter;
// import top.wecoding.iam.server.util.LogUtil;
//
// /**
//  * @author liuyuhui
//  * @date 2022/10/3
//  *
//  */
// @Slf4j
// public class IAMAuthenticationSuccessEventHandler implements AuthenticationSuccessHandler {
//
//   private HttpMessageConverter<OAuth2AccessTokenResponse> accessTokenHttpResponseConverter;
//
//   public IAMAuthenticationSuccessEventHandler() {
//     OAuth2AccessTokenResponseHttpMessageConverter tokenResponseHttpMessageConverter =
//         new OAuth2AccessTokenResponseHttpMessageConverter();
//     tokenResponseHttpMessageConverter.setAccessTokenResponseParametersConverter(
//         new AccessTokenResponseParametersConverter());
//     this.accessTokenHttpResponseConverter = tokenResponseHttpMessageConverter;
//   }
//
//   @PostConstruct
//   public void init() {
//     OAuth2AccessTokenResponseHttpMessageConverter oAuth2AccessTokenResponseHttpMessageConverter =
//         new OAuth2AccessTokenResponseHttpMessageConverter();
//     oAuth2AccessTokenResponseHttpMessageConverter.setAccessTokenResponseParametersConverter(
//         new AccessTokenResponseParametersConverter());
//     this.accessTokenHttpResponseConverter = oAuth2AccessTokenResponseHttpMessageConverter;
//   }
//
//   @Override
//   public void onAuthenticationSuccess(
//       HttpServletRequest request, HttpServletResponse response, Authentication authentication)
//       throws IOException, ServletException {
//
//     // TODO 发送成功日志
//     OAuth2AccessTokenAuthenticationToken accessTokenAuthentication =
//         (OAuth2AccessTokenAuthenticationToken) authentication;
//     Map<String, Object> map = accessTokenAuthentication.getAdditionalParameters();
//     if (MapUtil.isNotEmpty(map)) {
//       // 发送异步日志事件
//       LoginUser userInfo = (LoginUser) map.get(SecurityConstants.DETAILS_USER);
//       log.info("用户：{} 登录成功", userInfo.getName());
//       SecurityContextHolder.getContext().setAuthentication(accessTokenAuthentication);
//       LogUtil.successLogin(userInfo.userInfo().getUserId());
//     }
//
//     sendAccessTokenResponse(request, response, authentication);
//   }
//
//   private void sendAccessTokenResponse(
//       HttpServletRequest request, HttpServletResponse response, Authentication authentication)
//       throws IOException {
//     OAuth2AccessTokenAuthenticationToken accessTokenAuthentication =
//         (OAuth2AccessTokenAuthenticationToken) authentication;
//
//     OAuth2AccessToken accessToken = accessTokenAuthentication.getAccessToken();
//     OAuth2RefreshToken refreshToken = accessTokenAuthentication.getRefreshToken();
//     Map<String, Object> additionalParameters =
// accessTokenAuthentication.getAdditionalParameters();
//
//     OAuth2AccessTokenResponse.Builder builder =
//         OAuth2AccessTokenResponse.withToken(accessToken.getTokenValue())
//             .tokenType(accessToken.getTokenType())
//             .scopes(accessToken.getScopes());
//     if (accessToken.getIssuedAt() != null && accessToken.getExpiresAt() != null) {
//       builder.expiresIn(
//           ChronoUnit.SECONDS.between(accessToken.getIssuedAt(), accessToken.getExpiresAt()));
//     }
//     if (refreshToken != null) {
//       builder.refreshToken(refreshToken.getTokenValue());
//     }
//     if (!CollectionUtils.isEmpty(additionalParameters)) {
//       builder.additionalParameters(additionalParameters);
//     }
//     OAuth2AccessTokenResponse accessTokenResponse = builder.build();
//     ServletServerHttpResponse httpResponse = new ServletServerHttpResponse(response);
//
//     // stateless delete the auth context information
//     SecurityContextHolder.clearContext();
//     this.accessTokenHttpResponseConverter.write(accessTokenResponse, null, httpResponse);
//   }
// }
