package top.wecoding.iam.server.service.impl;

import com.baomidou.mybatisplus.extension.plugins.pagination.Page;
import com.baomidou.mybatisplus.extension.service.impl.ServiceImpl;
import lombok.RequiredArgsConstructor;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.stereotype.Service;
import top.wecoding.core.result.PageInfo;
import top.wecoding.core.util.AssertUtils;
import top.wecoding.core.util.PageUtil;
import top.wecoding.iam.common.constant.RedisConstant;
import top.wecoding.iam.common.enums.IamErrorCode;
import top.wecoding.iam.common.model.request.CreateOauth2ClientPageRequest;
import top.wecoding.iam.common.model.request.CreateOauth2ClientRequest;
import top.wecoding.iam.common.model.request.UpdateOauth2ClientRequest;
import top.wecoding.iam.common.model.response.Oauth2ClientInfoResponse;
import top.wecoding.iam.server.entity.Oauth2Client;
import top.wecoding.iam.server.mapper.Oauth2ClientMapper;
import top.wecoding.iam.server.service.Oauth2ClientService;
import top.wecoding.iam.server.util.Oauth2ClientUtil;

/**
 * @author liuyuhui
 * @date 2022/10/5
 * @qq 1515418211
 */
@Service
@RequiredArgsConstructor
public class Oauth2ClientServiceImpl extends ServiceImpl<Oauth2ClientMapper, Oauth2Client>
    implements Oauth2ClientService {

  private final Oauth2ClientMapper clientMapper;

  @Override
  public Oauth2ClientInfoResponse getInfo(String clientId) {
    Oauth2Client client = clientMapper.getByClientId(clientId);

    AssertUtils.isNotNull(client, IamErrorCode.CLIENT_DOES_NOT_EXIST);

    return Oauth2ClientUtil.toOauth2ClientInfoResponse(client);
  }

  @Override
  public void create(CreateOauth2ClientRequest createOauth2ClientRequest) {
    Oauth2Client client = Oauth2ClientUtil.ofOauth2Client(createOauth2ClientRequest);

    Oauth2Client oldClient = clientMapper.getByClientId(client.getClientId());

    AssertUtils.isNull(oldClient, IamErrorCode.CLIENT_ID_ALREADY_EXISTS);

    AssertUtils.isTrue(clientMapper.insert(client) > 0, IamErrorCode.CLIENT_ADD_FAILED);
  }

  @Override
  @CacheEvict(value = RedisConstant.CLIENT_DETAILS_KEY, key = "#updateOauth2ClientRequest.clientId")
  public void update(UpdateOauth2ClientRequest updateOauth2ClientRequest) {
    String id = updateOauth2ClientRequest.getId();

    AssertUtils.isNotNull(clientMapper.getById(id), IamErrorCode.CLIENT_DOES_NOT_EXIST);

    Oauth2Client client = Oauth2ClientUtil.ofOauth2Client(updateOauth2ClientRequest);

    AssertUtils.isTrue(clientMapper.updateById(client) > 0, IamErrorCode.CLIENT_UPDATE_FAILED);
  }

  @Override
  @CacheEvict(value = RedisConstant.CLIENT_DETAILS_KEY, key = "#clientId")
  public void delete(String clientId) {
    Oauth2Client client = clientMapper.getByClientId(clientId);

    AssertUtils.isNotNull(client, IamErrorCode.CLIENT_DOES_NOT_EXIST);

    AssertUtils.isTrue(
        clientMapper.deleteById(client.getId()) > 0, IamErrorCode.CLIENT_DELETE_FAILED);
  }

  @Override
  public PageInfo<Oauth2ClientInfoResponse> infoPage(
      CreateOauth2ClientPageRequest clientPageRequest) {
    Page<Oauth2Client> pageResult =
        clientMapper.page(PageUtil.getPageFromRequest(clientPageRequest), clientPageRequest);
    return PageInfo.of(pageResult.getRecords(), clientPageRequest, pageResult.getTotal())
        .map(Oauth2ClientUtil::toOauth2ClientInfoResponse);
  }
}
