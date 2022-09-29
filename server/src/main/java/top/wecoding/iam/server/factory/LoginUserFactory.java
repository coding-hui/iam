package top.wecoding.iam.server.factory;

import cn.hutool.core.util.ObjectUtil;
import cn.hutool.system.UserInfo;
import lombok.extern.slf4j.Slf4j;
import top.wecoding.core.util.WebUtil;
import top.wecoding.iam.common.model.LoginUser;

import javax.servlet.http.HttpServletRequest;

/**
 * 创建登录用户信息工厂
 *
 * @author liuyuhui
 * @qq 1515418211
 */
@Slf4j
public class LoginUserFactory {

  public static LoginUser createLoginUser(UserInfo userInfo) {
    LoginUser loginUser = new LoginUser();
    HttpServletRequest request = WebUtil.getRequest();
    if (ObjectUtil.isNull(request)) {
      return null;
    }
    // 用户基本信息
    // SysUser sysUser = userInfo.getSysUser();
    // String account = sysUser.getAccount();
    //
    // BeanUtil.copyProperties(sysUser, loginUser);
    //
    // loginUser.setLastLoginIp(IpAddressUtil.getIp(request));
    // loginUser.setLastLoginTime(DateTime.now());
    // loginUser.setLoginLocation(IpAddressUtil.getAddress(request));
    // // loginUser.setLastLoginBrowser(UaUtil.getBrowser(request));
    // // loginUser.setLastLoginOs(UaUtil.getOs(request));
    //
    // Set<Dict> roles =
    //     ObjectUtil.isNull(userInfo.getRoles()) ? new HashSet<>() : userInfo.getRoles();
    // Set<String> roleKeys =
    //     roles.stream().map(dict -> dict.getStr(CommonConstant.CODE)).collect(Collectors.toSet());
    // Set<String> permissions = userInfo.getPermissions();
    // List<Long> dataScopes = userInfo.getDataScopes();
    // // 角色信息
    // loginUser.setRoles(roles);
    // loginUser.setRoleKeys(roleKeys);
    // // 权限信息
    // loginUser.setPermissions(permissions);
    // // 数据范围信息
    // loginUser.setDataScopes(dataScopes);

    return loginUser;
  }
}
