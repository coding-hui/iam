// /*
//  * Copyright (c) 2022. WeCoding (wecoding@yeah.net).
//  *
//  * Licensed under the GNU LESSER GENERAL PUBLIC LICENSE 3.0;
//  * you may not use this file except in compliance with the License.
//  * You may obtain a copy of the License at
//  *
//  * http://www.gnu.org/licenses/lgpl.html
//  *
//  * Unless required by applicable law or agreed to in writing, software
//  * distributed under the License is distributed on an "AS IS" BASIS,
//  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  * See the License for the specific language governing permissions and
//  * limitations under the License.
//  */
// package top.wecoding.starter.security.service;
//
// import cn.hutool.core.convert.Convert;
// import cn.hutool.core.util.IdUtil;
// import cn.hutool.core.util.StrUtil;
// import io.jsonwebtoken.Claims;
// import lombok.AllArgsConstructor;
// import lombok.extern.slf4j.Slf4j;
// import org.springframework.stereotype.Component;
// import top.wecoding.common.core.constant.SecurityConstants;
// import top.wecoding.common.core.constant.TokenConstant;
// import top.wecoding.common.core.util.HttpServletUtils;
// import top.wecoding.starter.auth.model.AuthInfo;
// import top.wecoding.starter.auth.model.LoginUser;
// import top.wecoding.starter.auth.util.AuthUtil;
// import top.wecoding.starter.jwt.model.JwtPayLoad;
// import top.wecoding.starter.jwt.model.TokenInfo;
// import top.wecoding.starter.jwt.props.JwtProperties;
// import top.wecoding.starter.jwt.util.JwtUtils;
//
// import javax.servlet.http.HttpServletRequest;
// import java.util.Date;
// import java.util.HashMap;
// import java.util.Map;
//
// /**
//  * Token相关服务类
//  *
//  * @author liuyuhui
//  * @qq 1515418211
//  */
// @Slf4j
// @Component
// @AllArgsConstructor
// public class TokenService {
//
//   private final JwtProperties jwtProperties;
//
//   private final LoginUserCache loginUserCache;
//
//   /**
//    * 创建认证 Token
//    *
//    * @param loginUser 登录用户信息
//    * @param expireMillis 过期时间（秒）
//    * @return 认证信息
//    */
//   public AuthInfo createAuthInfo(LoginUser loginUser, Long expireMillis) {
//     // 缓存用户
//     loginUser.setUuid(IdUtil.fastUUID());
//     cacheLoginUser(loginUser);
//
//     JwtPayLoad jwtPayLoad =
//         JwtPayLoad.builder()
//             .uuid(loginUser.getUuid())
//             .userId(loginUser.getUserId())
//             .account(loginUser.getAccount())
//             .clientId(loginUser.getClientId())
//             .realName(loginUser.getRealName())
//             .build();
//     // 构建 token
//     return createAuthInfo(jwtPayLoad, expireMillis);
//   }
//
//   /**
//    * 创建认证 Token
//    *
//    * @param jwtPayLoad 数据声明
//    * @param expireMillis 过期时间（秒）
//    * @return 认证信息
//    */
//   public AuthInfo createAuthInfo(JwtPayLoad jwtPayLoad, Long expireMillis) {
//     if (expireMillis == null || expireMillis <= 0) {
//       expireMillis = jwtProperties.getExpire();
//     }
//
//     Map<String, String> claims = new HashMap<>();
//     claims.put(TokenConstant.TOKEN_TYPE, TokenConstant.ACCESS_TOKEN);
//     claims.put(SecurityConstants.USER_KEY, jwtPayLoad.getUuid());
//     claims.put(SecurityConstants.DETAILS_USER_ID, Convert.toStr(jwtPayLoad.getUserId()));
//     claims.put(SecurityConstants.DETAILS_ACCOUNT, jwtPayLoad.getAccount());
//     claims.put(SecurityConstants.DETAILS_CLIENT_ID, jwtPayLoad.getClientId());
//     claims.put(SecurityConstants.DETAILS_USERNAME, jwtPayLoad.getRealName());
//
//     TokenInfo tokenInfo = JwtUtils.createJWT(claims, expireMillis);
//
//     return AuthInfo.builder()
//         .accessToken(tokenInfo.getToken())
//         .expireMillis(tokenInfo.getExpiresIn())
//         .expiration(tokenInfo.getExpiration())
//         .tokenType(TokenConstant.ACCESS_TOKEN)
//         .refreshToken(createRefreshToken(jwtPayLoad).getToken())
//         .uuid(jwtPayLoad.getUuid())
//         .userId(jwtPayLoad.getUserId())
//         .account(jwtPayLoad.getAccount())
//         .realName(jwtPayLoad.getRealName())
//         .clientId(jwtPayLoad.getClientId())
//         .license(SecurityConstants.PROJECT_LICENSE)
//         .build();
//   }
//
//   /**
//    * 创建refreshToken
//    *
//    * @param jwtPayLoad 数据声明
//    * @return refreshToken
//    */
//   private TokenInfo createRefreshToken(JwtPayLoad jwtPayLoad) {
//     Map<String, String> claims = new HashMap<>(16);
//     claims.put(TokenConstant.TOKEN_TYPE, TokenConstant.REFRESH_TOKEN);
//     claims.put(SecurityConstants.USER_KEY, jwtPayLoad.getUuid());
//     claims.put(SecurityConstants.DETAILS_USER_ID, Convert.toStr(jwtPayLoad.getUserId()));
//     claims.put(SecurityConstants.DETAILS_ACCOUNT, jwtPayLoad.getAccount());
//     claims.put(SecurityConstants.DETAILS_CLIENT_ID, jwtPayLoad.getClientId());
//     return JwtUtils.createJWT(claims, jwtProperties.getRefreshExpire());
//   }
//
//   /**
//    * 缓存登录用户信息
//    *
//    * @param loginUser 登录用户信息
//    */
//   public void cacheLoginUser(LoginUser loginUser) {
//     loginUserCache.set(loginUser.getUuid(), loginUser, jwtProperties.getRefreshExpire());
//   }
//
//   /**
//    * 删除登录用户信息
//    *
//    * @param token 令牌
//    */
//   public void removeLoginUser(String token) {
//     if (StrUtil.isBlank(token)) {
//       return;
//     }
//     loginUserCache.del(JwtUtils.getUserKey(token));
//   }
//
//   /**
//    * 获取用户身份信息
//    *
//    * @return 用户信息
//    */
//   public LoginUser getLoginUser() {
//     return getLoginUser(HttpServletUtils.getRequest());
//   }
//
//   /**
//    * 获取用户身份信息
//    *
//    * @return 用户信息
//    */
//   public LoginUser getLoginUser(HttpServletRequest request) {
//     return getLoginUser(AuthUtil.getToken(request));
//   }
//
//   /**
//    * 获取用户身份信息，取不到返回null
//    *
//    * @return 用户信息
//    */
//   public LoginUser getLoginUser(String token) {
//     try {
//       if (StrUtil.isNotEmpty(token)) {
//         return loginUserCache.get(JwtUtils.getUserKey(token));
//       }
//     } catch (Exception e) {
//       log.warn(" >>> 获取登录用户失败. Thread:{}", Thread.currentThread());
//     }
//     return null;
//   }
//
//   /**
//    * 解析 Token
//    *
//    * @param token token
//    * @return 用户信息
//    */
//   public AuthInfo getAuthInfo(String token) {
//     Claims claims = JwtUtils.parseToken(token);
//     Date expiration = claims.getExpiration();
//     return AuthInfo.builder()
//         .accessToken(token)
//         .tokenType(JwtUtils.getValue(claims, TokenConstant.TOKEN_TYPE))
//         .uuid(JwtUtils.getUserKey(claims))
//         .userId(JwtUtils.getUserId(claims))
//         .account(JwtUtils.getUserAccount(claims))
//         .realName(JwtUtils.getUserRealName(claims))
//         .clientId(JwtUtils.getClientId(claims))
//         .expireMillis(Convert.toLong(expiration, 0L))
//         .build();
//   }
// }
