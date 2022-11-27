package top.wecoding.iam.server.util;

import lombok.experimental.UtilityClass;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

/**
 * @author liuyuhui
 */
@UtilityClass
public class PasswordEncoderUtil {

  private static final PasswordEncoder ENCODER = new BCryptPasswordEncoder();

  public static String encode(String rawPassword) {
    return ENCODER.encode(rawPassword);
  }

  public static boolean matches(CharSequence rawPassword, String encodedPassword) {
    return ENCODER.matches(rawPassword, encodedPassword);
  }
}
