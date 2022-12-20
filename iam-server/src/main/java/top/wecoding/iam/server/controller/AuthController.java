package top.wecoding.iam.server.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;
import org.springframework.web.servlet.ModelAndView;

/**
 * @author liuyuhui
 * @since 0.5
 */
@RestController
public class AuthController {

  @GetMapping("/auth/login")
  public ModelAndView require(
      ModelAndView modelAndView, @RequestParam(required = false) String error) {
    modelAndView.setViewName("ftl/login");
    modelAndView.addObject("error", error);
    return modelAndView;
  }
}
