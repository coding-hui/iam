package top.wecoding.iam.server.authentication.dao;

import cn.hutool.core.util.StrUtil;
import cn.hutool.extra.servlet.ServletUtil;
import cn.hutool.extra.spring.SpringUtil;
import lombok.SneakyThrows;
import org.springframework.core.Ordered;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.authentication.InternalAuthenticationServiceException;
import org.springframework.security.authentication.UsernamePasswordAuthenticationToken;
import org.springframework.security.authentication.dao.AbstractUserDetailsAuthenticationProvider;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.AuthenticationException;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.core.userdetails.UserDetailsPasswordService;
import org.springframework.security.core.userdetails.UserDetailsService;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.factory.PasswordEncoderFactories;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.oauth2.core.endpoint.OAuth2ParameterNames;
import org.springframework.security.web.authentication.www.BasicAuthenticationConverter;
import org.springframework.util.Assert;
import top.wecoding.core.constant.SecurityConstants;
import top.wecoding.core.util.WebUtils;
import top.wecoding.iam.common.userdetails.IAMUserDetailsService;

import javax.servlet.http.HttpServletRequest;
import java.util.Comparator;
import java.util.Map;
import java.util.Optional;
import java.util.function.Supplier;

/**
 * @author liuyuhui
 * @date 2022/10/3
 * @qq 1515418211
 */
public class IAMDaoAuthenticationProvider extends AbstractUserDetailsAuthenticationProvider {

  /**
   * The plaintext password used to perform PasswordEncoder#matches(CharSequence, String)} on when
   * the user is not found to avoid SEC-2056.
   */
  private static final String USER_NOT_FOUND_PASSWORD = "userNotFoundPassword";

  private static final BasicAuthenticationConverter basicConvert =
      new BasicAuthenticationConverter();

  private PasswordEncoder passwordEncoder;

  /**
   * The password used to perform {@link PasswordEncoder#matches(CharSequence, String)} on when the
   * user is not found to avoid SEC-2056. This is necessary, because some {@link PasswordEncoder}
   * implementations will short circuit if the password is not in a valid format.
   */
  private volatile String userNotFoundEncodedPassword;

  private UserDetailsService userDetailsService;

  private UserDetailsPasswordService userDetailsPasswordService;

  public IAMDaoAuthenticationProvider() {
    setMessageSource(SpringUtil.getBean("securityMessageSource"));
    setPasswordEncoder(PasswordEncoderFactories.createDelegatingPasswordEncoder());
  }

  @Override
  protected void additionalAuthenticationChecks(
      UserDetails userDetails, UsernamePasswordAuthenticationToken authentication)
      throws AuthenticationException {

    // app 模式不用校验密码
    String grantType = WebUtils.getRequest().getParameter(OAuth2ParameterNames.GRANT_TYPE);
    if (StrUtil.equals(SecurityConstants.APP, grantType)) {
      return;
    }

    if (authentication.getCredentials() == null) {
      this.logger.debug("Failed to authenticate since no credentials provided");
      throw new BadCredentialsException(
          this.messages.getMessage(
              "AbstractUserDetailsAuthenticationProvider.badCredentials", "Bad credentials"));
    }
    String presentedPassword = authentication.getCredentials().toString();
    if (!this.passwordEncoder.matches(presentedPassword, userDetails.getPassword())) {
      this.logger.debug("Failed to authenticate since password does not match stored value");
      throw new BadCredentialsException(
          this.messages.getMessage(
              "AbstractUserDetailsAuthenticationProvider.badCredentials", "Bad credentials"));
    }
  }

  @SneakyThrows
  @Override
  protected final UserDetails retrieveUser(
      String username, UsernamePasswordAuthenticationToken authentication) {
    prepareTimingAttackProtection();
    HttpServletRequest request =
        WebUtils.getRequestOptional()
            .orElseThrow(
                (Supplier<Throwable>)
                    () ->
                        new InternalAuthenticationServiceException(
                            "retrieve user failed. web request is empty"));

    Map<String, String> paramMap = ServletUtil.getParamMap(request);
    String grantType = paramMap.get(OAuth2ParameterNames.GRANT_TYPE);
    String clientId = paramMap.get(OAuth2ParameterNames.CLIENT_ID);

    if (StrUtil.isBlank(clientId)) {
      clientId = basicConvert.convert(request).getName();
    }

    Map<String, IAMUserDetailsService> userDetailsServiceMap =
        SpringUtil.getBeansOfType(IAMUserDetailsService.class);

    String finalClientId = clientId;
    Optional<IAMUserDetailsService> optional =
        userDetailsServiceMap.values().stream()
            .filter(service -> service.support(finalClientId, grantType))
            .max(Comparator.comparingInt(Ordered::getOrder));

    if (!optional.isPresent()) {
      throw new InternalAuthenticationServiceException("UserDetailsService error , not register");
    }

    try {
      UserDetails loadedUser = optional.get().loadUserByUsername(username);
      if (loadedUser == null) {
        throw new InternalAuthenticationServiceException(
            "UserDetailsService returned null, which is an interface contract violation");
      }
      return loadedUser;
    } catch (UsernameNotFoundException ex) {
      mitigateAgainstTimingAttack(authentication);
      throw ex;
    } catch (InternalAuthenticationServiceException ex) {
      throw ex;
    } catch (Exception ex) {
      throw new InternalAuthenticationServiceException(ex.getMessage(), ex);
    }
  }

  @Override
  protected Authentication createSuccessAuthentication(
      Object principal, Authentication authentication, UserDetails user) {
    boolean upgradeEncoding =
        this.userDetailsPasswordService != null
            && this.passwordEncoder.upgradeEncoding(user.getPassword());
    if (upgradeEncoding) {
      String presentedPassword = authentication.getCredentials().toString();
      String newPassword = this.passwordEncoder.encode(presentedPassword);
      user = this.userDetailsPasswordService.updatePassword(user, newPassword);
    }
    return super.createSuccessAuthentication(principal, authentication, user);
  }

  private void prepareTimingAttackProtection() {
    if (this.userNotFoundEncodedPassword == null) {
      this.userNotFoundEncodedPassword = this.passwordEncoder.encode(USER_NOT_FOUND_PASSWORD);
    }
  }

  private void mitigateAgainstTimingAttack(UsernamePasswordAuthenticationToken authentication) {
    if (authentication.getCredentials() != null) {
      String presentedPassword = authentication.getCredentials().toString();
      this.passwordEncoder.matches(presentedPassword, this.userNotFoundEncodedPassword);
    }
  }

  protected PasswordEncoder getPasswordEncoder() {
    return this.passwordEncoder;
  }

  /**
   * Sets the PasswordEncoder instance to be used to encode and validate passwords. If not set, the
   * password will be compared using {@link
   * PasswordEncoderFactories#createDelegatingPasswordEncoder()}
   *
   * @param passwordEncoder must be an instance of one of the {@code PasswordEncoder} types.
   */
  public void setPasswordEncoder(PasswordEncoder passwordEncoder) {
    Assert.notNull(passwordEncoder, "passwordEncoder cannot be null");
    this.passwordEncoder = passwordEncoder;
    this.userNotFoundEncodedPassword = null;
  }

  protected UserDetailsService getUserDetailsService() {
    return this.userDetailsService;
  }

  public void setUserDetailsService(UserDetailsService userDetailsService) {
    this.userDetailsService = userDetailsService;
  }

  public void setUserDetailsPasswordService(UserDetailsPasswordService userDetailsPasswordService) {
    this.userDetailsPasswordService = userDetailsPasswordService;
  }
}
