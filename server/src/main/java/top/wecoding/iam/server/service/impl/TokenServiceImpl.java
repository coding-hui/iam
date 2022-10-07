package top.wecoding.iam.server.service.impl;

import cn.hutool.core.util.StrUtil;
import lombok.RequiredArgsConstructor;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.security.oauth2.core.OAuth2AccessToken;
import org.springframework.security.oauth2.core.OAuth2TokenType;
import org.springframework.security.oauth2.server.authorization.OAuth2Authorization;
import org.springframework.security.oauth2.server.authorization.OAuth2AuthorizationService;
import org.springframework.stereotype.Service;
import top.wecoding.core.cache.CacheKey;
import top.wecoding.core.result.PageInfo;
import top.wecoding.iam.common.cache.UserTokenCacheKeyBuilder;
import top.wecoding.iam.common.constant.RedisConstant;
import top.wecoding.iam.common.model.request.TokenInfoPageRequest;
import top.wecoding.iam.common.model.response.TokenInfoResponse;
import top.wecoding.iam.server.service.TokenService;

import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * @author liuyuhui
 * @date 2022/10/6
 * @qq 1515418211
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
    if (accessToken == null || StrUtil.isBlank(accessToken.getToken().getTokenValue())) {
      return true;
    }
    // TODO 清空用户信息
    // 清空access token
    authorizationService.remove(authorization);

    return true;
  }

  @Override
  public PageInfo<TokenInfoResponse> infoPage(TokenInfoPageRequest tokenInfoPageRequest) {
    CacheKey key = new UserTokenCacheKeyBuilder().key(RedisConstant.OAUTH_ACCESS_PREFIX);
    Set<String> keys = redisTemplate.keys(key.getKey());
    if (keys == null || keys.isEmpty()) {
      return PageInfo.empty();
    }

    int total = keys.size();
    int pageSize = tokenInfoPageRequest.getPageSize();
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
