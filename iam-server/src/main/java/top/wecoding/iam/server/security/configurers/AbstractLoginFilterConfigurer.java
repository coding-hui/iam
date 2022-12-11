package top.wecoding.iam.server.security.configurers;

import org.springframework.security.config.annotation.ObjectPostProcessor;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.web.util.matcher.RequestMatcher;

/**
 * @author liuyuhui
 * @since 0.5
 */
abstract class AbstractLoginFilterConfigurer {

  private final ObjectPostProcessor<Object> objectPostProcessor;

  AbstractLoginFilterConfigurer(ObjectPostProcessor<Object> objectPostProcessor) {
    this.objectPostProcessor = objectPostProcessor;
  }

  abstract void init(HttpSecurity httpSecurity);

  abstract void configure(HttpSecurity httpSecurity);

  abstract RequestMatcher getRequestMatcher();

  protected final <T> T postProcess(T object) {
    return (T) this.objectPostProcessor.postProcess(object);
  }

  protected final ObjectPostProcessor<Object> getObjectPostProcessor() {
    return this.objectPostProcessor;
  }
}
