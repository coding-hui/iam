package top.wecoding.iam.server.entity;

import com.baomidou.mybatisplus.annotation.TableField;
import com.baomidou.mybatisplus.annotation.TableId;
import com.baomidou.mybatisplus.annotation.TableName;
import java.util.Date;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.EqualsAndHashCode;
import lombok.NoArgsConstructor;
import top.wecoding.mybatis.base.BaseEntity;

/**
 * 系统用户表
 *
 * @author liuyuhui
 */
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
@TableName("iam_user")
@EqualsAndHashCode(callSuper = true)
public class User extends BaseEntity {

  @TableId private String id;

  private String tenantId;

  @TableField("pwd")
  private String password;

  @TableField("def_pwd")
  private Boolean defaultPwd;

  private Integer userState;

  private String lastLoginIp;

  private Date lastLoginTime;

  private Integer loginCount;
}
