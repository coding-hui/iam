package top.wecoding.iam.server.service.impl;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;
import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.server.authorization.OAuth2Authorization;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.security.oauth2.server.authorization.OAuth2TokenType;
import org.springframework.stereotype.Service;
import top.wecoding.commons.core.cache.CacheKey;
import top.wecoding.commons.core.model.PageInfo;
import top.wecoding.commons.lang.Strings;
import top.wecoding.iam.common.constant.RedisConstant;
import top.wecoding.iam.common.model.request.TokenInfoPageRequest;
import top.wecoding.iam.common.model.response.TokenInfoResponse;
import top.wecoding.iam.framework.cache.UserDetailsCacheKeyBuilder;
import top.wecoding.iam.framework.cache.UserTokenCacheKeyBuilder;
import top.wecoding.iam.server.service.TokenService;
import top.wecoding.redis.util.RedisUtils;

/**
 * @author liuyuhui
 * @date 2022/10/6
 */
@Service
@RequiredArgsConstructor
public class TokenServiceImpl implements TokenService {

  private final RedisTemplate<String, Object> redisTemplate;

  private final OAuth2AuthorizationService authorizationService;

  @Override
  public boolean delete(String tokenValue) {
    OAuth2Authorization authorization =
        authorizationService.findByToken(tokenValue, OAuth2TokenType.ACCESS_TOKEN);
    if (authorization == null) {
      return true;
    }

    OAuth2Authorization.Token<OAuth2AccessToken> accessToken = authorization.getAccessToken();
    if (accessToken == null || Strings.isBlank(accessToken.getToken().getTokenValue())) {
      return true;
    }
    // 清空用户信息
    RedisUtils.del(new UserDetailsCacheKeyBuilder().build(authorization.getPrincipalName()));
    // 清空access token
    authorizationService.remove(authorization);

    return true;
  }

  @Override
  public PageInfo<TokenInfoResponse> infoPage(TokenInfoPageRequest tokenInfoPageRequest) {
    CacheKey key = new UserTokenCacheKeyBuilder().build(RedisConstant.OAUTH_ACCESS_PREFIX);
    Set<String> keys = redisTemplate.keys(key.getKey());
    if (keys == null || keys.isEmpty()) {
      return PageInfo.empty();
    }

    int total = keys.size();
    int pageSize = tokenInfoPageRequest.getSize();
    long offset = tokenInfoPageRequest.getOffset();

    List<String> selectKeys =
        keys.stream().skip(offset).limit(pageSize).collect(Collectors.toList());

    List<TokenInfoResponse> tokens =
        Optional.ofNullable(redisTemplate.opsForValue().multiGet(selectKeys))
            .map(
                list ->
                    list.stream()
                        .map(
                            token -> {
                              OAuth2Authorization authorization = (OAuth2Authorization) token;
                              OAuth2Authorization.Token<OAuth2AccessToken> accessToken =
                                  authorization.getAccessToken();

                              TokenInfoResponse.TokenInfoResponseBuilder builder =
                                  TokenInfoResponse.builder()
                                      .id(authorization.getId())
                                      .clientId(authorization.getRegisteredClientId())
                                      .username(authorization.getPrincipalName())
                                      .accessToken(accessToken.getToken().getTokenValue())
                                      .expiresAt(accessToken.getToken().getExpiresAt())
                                      .expiresAt(accessToken.getToken().getExpiresAt())
                                      .issuedAt(accessToken.getToken().getIssuedAt());

                              return builder.build();
                            })
                        .collect(Collectors.toList()))
            .orElse(Collections.emptyList());

    return PageInfo.of(tokens, tokenInfoPageRequest, total);
  }
}
