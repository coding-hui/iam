package top.wecoding.iam.server.util;

import lombok.experimental.UtilityClass;
import top.wecoding.core.enums.SystemDictEnum;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.core.util.RSAUtil;
import top.wecoding.iam.common.enums.IamErrorCode;

import java.util.Set;
import java.util.stream.Collectors;

/**
 * @author liuyuhui
 * @qq 1515418211
 */
@UtilityClass
public class PasswordUtil {

  private static final String SYMBOL_STRING = "!\"#%$()*+-,./:;<=>?@[\\]^_`{|}~";

  private static final Set<Integer> SYMBOL_CHARACTERS_SET =
      SYMBOL_STRING.chars().boxed().collect(Collectors.toSet());

  public static boolean checkLength(String password) {
    return null != password && 8 <= password.length() && 255 >= password.length();
  }

  public static void checkPwd(String password) {
    AssertUtils.isTrue(checkLength(password), IamErrorCode.PASSWORD_LENGTH_WRONG);

    boolean hasSpecial = false;
    boolean hasNumber = false;
    boolean hasLowerLetter = false;
    boolean hasUpperLetter = false;
    for (char c : password.toCharArray()) {
      if (SYMBOL_CHARACTERS_SET.contains((int) c)) {
        hasSpecial = true;
        continue;
      }
      if ('0' <= c && '9' >= c) {
        hasNumber = true;
        continue;
      }
      if ('a' <= c && 'z' >= c) {
        hasLowerLetter = true;
        continue;
      }
      if ('A' <= c && 'Z' >= c) {
        hasUpperLetter = true;
        continue;
      }
      AssertUtils.error(IamErrorCode.INVALID_PASSWORD, SYMBOL_STRING);
    }
    AssertUtils.isTrue(
        (hasSpecial && hasNumber && (hasLowerLetter || hasUpperLetter)),
        IamErrorCode.WEAK_PASSWORD);
  }

  public static boolean checkContent(String password) {
    if (!checkLength(password)) {
      return false;
    }

    boolean hasSpecial = false;
    boolean hasNumber = false;
    boolean hasLowerLetter = false;
    boolean hasUpperLetter = false;
    boolean hasUnexpected = false;
    for (char c : password.toCharArray()) {
      if (SYMBOL_CHARACTERS_SET.contains((int) c)) {
        hasSpecial = true;
        continue;
      }
      if ('0' <= c && '9' >= c) {
        hasNumber = true;
        continue;
      }
      if ('a' <= c && 'z' >= c) {
        hasLowerLetter = true;
        continue;
      }
      if ('A' <= c && 'Z' >= c) {
        hasUpperLetter = true;
        continue;
      }
      hasUnexpected = true;
      break;
    }
    return (!hasUnexpected && hasSpecial && hasNumber && (hasLowerLetter || hasUpperLetter));
  }

  public static String decrypt(String ciphertext) {
    return RSAUtil.decryptFromBase64(SystemDictEnum.SYSTEM_RSA_PRIVATE_KEY.getValue(), ciphertext);
  }

  public static String encrypt(String context) {
    return RSAUtil.encryptToBase64(SystemDictEnum.SYSTEM_RSA_PUBLIC_KEY.getValue(), context);
  }
}
