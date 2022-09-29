package top.wecoding.iam.common.helper;

import cn.hutool.core.convert.Convert;
import cn.hutool.core.util.StrUtil;
import cn.hutool.crypto.digest.BCrypt;
import cn.hutool.extra.spring.SpringUtil;
import top.wecoding.core.constant.StrPool;
import top.wecoding.core.enums.rest.CommonErrorCodeEnum;
import top.wecoding.core.exception.IllegalParameterException;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.core.util.WebUtil;
import top.wecoding.iam.common.provider.ClientDetails;
import top.wecoding.iam.common.provider.ClientDetailsService;

import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.Objects;

import static top.wecoding.core.constant.SecurityConstants.*;

/**
 * @author liuyuhui
 * @date 2022/9/12
 * @qq 1515418211
 */
public class SecurityHelper {

  private static final ClientDetailsService clientDetailsService;

  static {
    clientDetailsService = SpringUtil.getBean(ClientDetailsService.class);
  }

  public static String getClientIdFromHeader() {
    String[] tokens = extractAndDecodeHeader();
    assert tokens.length == 2;
    return tokens[0];
  }

  /**
   * @return [0]: clientId, [1]: password
   */
  public static String[] extractAndDecodeHeader() {
    String header = Objects.requireNonNull(WebUtil.getRequest()).getHeader(BASIC_HEADER_KEY);
    header = Convert.toStr(header).replace(BASIC_HEADER_PREFIX_EXT, BASIC_HEADER_PREFIX);

    AssertUtils.isNotBlank(
        header,
        CommonErrorCodeEnum.COMMON_ERROR,
        "No client authentication information in request header");
    AssertUtils.isTrue(
        header.startsWith(BASIC_HEADER_PREFIX),
        CommonErrorCodeEnum.COMMON_ERROR,
        "Client authentication information does not start with Basic");

    byte[] base64Token =
        header.substring(BASIC_HEADER_PREFIX.length()).getBytes(StandardCharsets.UTF_8);
    byte[] decoded;
    try {
      decoded = Base64.getDecoder().decode(base64Token);
    } catch (Exception e) {
      throw new IllegalParameterException(
          CommonErrorCodeEnum.COMMON_ERROR, "Failed to decode basic authentication token");
    }

    String token = new String(decoded, StandardCharsets.UTF_8);
    int delim = token.indexOf(StrPool.COLON);
    if (delim < 0) {
      throw new IllegalParameterException(
          CommonErrorCodeEnum.COMMON_ERROR, "Invalid basic authentication token");
    }

    // [0]: clientId, [1]: 密码
    return new String[] {token.substring(0, delim), token.substring(delim + 1)};
  }

  public static ClientDetails clientDetails(String clientId) {
    return clientDetailsService.loadClientByClientId(clientId);
  }

  public static boolean validateClient(
      ClientDetails clientDetails, String clientId, String clientSecret) {
    if (clientDetails != null) {
      return StrUtil.equals(clientId, clientDetails.getClientId())
          && BCrypt.checkpw(clientSecret, clientDetails.getClientSecret());
    }
    return false;
  }
}
