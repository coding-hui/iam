package top.wecoding.iam.framework;

import org.springframework.boot.autoconfigure.condition.ConditionalOnWebApplication;
import org.springframework.context.MessageSource;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.support.ReloadableResourceBundleMessageSource;
import org.springframework.web.servlet.config.annotation.WebMvcConfigurer;

import java.util.Locale;

import static org.springframework.boot.autoconfigure.condition.ConditionalOnWebApplication.Type.SERVLET;

/**
 * @author liuyuhui
 * @date 2022/10/4
 * @qq 1515418211
 */
@Configuration
@ConditionalOnWebApplication(type = SERVLET)
public class IAMErrorMessageSourceConfiguration implements WebMvcConfigurer {

  @Bean("iamMessageSource")
  public MessageSource iamMessageSource() {
    ReloadableResourceBundleMessageSource messageSource =
        new ReloadableResourceBundleMessageSource();
    messageSource.addBasenames("classpath:i18n/errors/messages");
    messageSource.addBasenames("classpath:i18n/errors/messages-iam");
    messageSource.addBasenames("classpath:i18n/errors/messages-common");
    messageSource.setDefaultLocale(Locale.CHINA);
    return messageSource;
  }
}
