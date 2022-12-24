package top.wecoding.iam.framework.security.web;

import jakarta.servlet.http.HttpServletRequest;
import java.util.regex.Matcher;
import java.util.regex.Pattern;
import org.springframework.http.HttpHeaders;
import org.springframework.http.MediaType;
import org.springframework.security.oauth2.core.OAuth2AuthenticationException;
import org.springframework.security.oauth2.server.resource.BearerTokenError;
import org.springframework.security.oauth2.server.resource.BearerTokenErrors;
import org.springframework.security.oauth2.server.resource.web.BearerTokenResolver;
import org.springframework.stereotype.Component;
import org.springframework.util.AntPathMatcher;
import org.springframework.util.PathMatcher;
import org.springframework.util.StringUtils;
import top.wecoding.iam.framework.props.IgnoreWhiteProperties;

/**
 * @author liuyuhui
 * @date 2022/10/3
 */
@Component
public class WeCodingBearerTokenExtractor implements BearerTokenResolver {

  private static final Pattern authorizationPattern =
      Pattern.compile("^Bearer (?<token>[a-zA-Z0-9-:._~+/]+=*)$", Pattern.CASE_INSENSITIVE);

  private final PathMatcher pathMatcher = new AntPathMatcher();

  private final IgnoreWhiteProperties urlProperties;

  private final boolean allowFormEncodedBodyParameter = false;

  private final boolean allowUriQueryParameter = false;

  private final String bearerTokenHeaderName = HttpHeaders.AUTHORIZATION;

  public WeCodingBearerTokenExtractor(IgnoreWhiteProperties urlProperties) {
    this.urlProperties = urlProperties;
  }

  private static String resolveFromRequestParameters(HttpServletRequest request) {
    String[] values = request.getParameterValues("access_token");
    if (values == null || values.length == 0) {
      return null;
    }
    if (values.length == 1) {
      return values[0];
    }
    BearerTokenError error =
        BearerTokenErrors.invalidRequest("Found multiple bearer tokens in the request");
    throw new OAuth2AuthenticationException(error);
  }

  @Override
  public String resolve(HttpServletRequest request) {
    boolean match =
        urlProperties.getWhites().stream()
            .anyMatch(url -> pathMatcher.match(url, request.getRequestURI()));

    if (match) {
      return null;
    }

    final String authorizationHeaderToken = resolveFromAuthorizationHeader(request);
    final String parameterToken =
        isParameterTokenSupportedForRequest(request) ? resolveFromRequestParameters(request) : null;
    if (authorizationHeaderToken != null) {
      if (parameterToken != null) {
        final BearerTokenError error =
            BearerTokenErrors.invalidRequest("Found multiple bearer tokens in the request");
        throw new OAuth2AuthenticationException(error);
      }
      return authorizationHeaderToken;
    }
    if (parameterToken != null && isParameterTokenEnabledForRequest(request)) {
      return parameterToken;
    }
    return null;
  }

  private String resolveFromAuthorizationHeader(HttpServletRequest request) {
    String authorization = request.getHeader(this.bearerTokenHeaderName);
    if (!StringUtils.startsWithIgnoreCase(authorization, "bearer")) {
      return null;
    }
    Matcher matcher = authorizationPattern.matcher(authorization);
    if (!matcher.matches()) {
      BearerTokenError error = BearerTokenErrors.invalidToken("Bearer token is malformed");
      throw new OAuth2AuthenticationException(error);
    }
    return matcher.group("token");
  }

  private boolean isParameterTokenSupportedForRequest(final HttpServletRequest request) {
    return (("POST".equals(request.getMethod())
            && MediaType.APPLICATION_FORM_URLENCODED_VALUE.equals(request.getContentType()))
        || "GET".equals(request.getMethod()));
  }

  private boolean isParameterTokenEnabledForRequest(final HttpServletRequest request) {
    return ((this.allowFormEncodedBodyParameter
            && "POST".equals(request.getMethod())
            && MediaType.APPLICATION_FORM_URLENCODED_VALUE.equals(request.getContentType()))
        || (this.allowUriQueryParameter && "GET".equals(request.getMethod())));
  }
}
