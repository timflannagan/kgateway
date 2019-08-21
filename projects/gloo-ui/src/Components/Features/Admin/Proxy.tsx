import * as React from 'react';
/** @jsx jsx */
import { jsx } from '@emotion/core';

import styled from '@emotion/styled/macro';
import { withRouter, RouteComponentProps } from 'react-router';
import { colors, healthConstants } from 'Styles';
import { useGetProxiesList } from 'Api/v2/useProxyClientV2';
import { ProxyDetails } from 'proto/github.com/solo-io/solo-projects/projects/grpcserver/api/v1/proxy_pb';
import { SectionCard } from 'Components/Common/SectionCard';
import { ReactComponent as ProxyLogo } from 'assets/proxy-icon.svg';
import { FileDownloadLink } from 'Components/Common/FileDownloadLink';
import { YamlDisplayer } from 'Components/Common/DisplayOnly/YamlDisplayer';
import { useSelector } from 'react-redux';
import { AppState } from 'store';

const InsideHeader = styled.div`
  display: flex;
  justify-content: space-between;
  font-size: 18px;
  line-height: 22px;
  margin-bottom: 18px;
  color: ${colors.novemberGrey};
`;

const ProxyLogoFullSize = styled(ProxyLogo)`
  width: 33px !important;
  max-height: none !important;
`;

interface Props {}

export const Proxys = (props: Props) => {
  const [gatewaysOpen, setGatewaysOpen] = React.useState<boolean[]>([]);

  const namespacesList = useSelector(
    (state: AppState) => state.config.namespacesList
  );
  const {
    data,
    loading,
    error,
    setNewVariables,
    dataObj: proxyObj
  } = useGetProxiesList({
    namespaces: namespacesList
  });
  const [allProxies, setAllProxies] = React.useState<ProxyDetails.AsObject[]>(
    []
  );

  React.useEffect(() => {
    if (!!data) {
      const newProxies = data
        .toObject()
        .proxyDetailsList.filter(proxy => !!proxy.proxy);
      setAllProxies(newProxies);
    }
  }, [loading]);

  if (!data || (!data && loading)) {
    return <div>Loading...</div>;
  }

  const toggleExpansion = (indexToggled: number) => {
    setGatewaysOpen(
      gatewaysOpen.map((isOpen, ind) => {
        if (ind !== indexToggled) {
          return false;
        }

        return !isOpen;
      })
    );
  };

  return (
    <React.Fragment>
      {allProxies.map((proxy, ind) => {
        return (
          <SectionCard
            key={proxy.proxy!.metadata!.name + ind}
            cardName={proxy.proxy!.metadata!.name}
            logoIcon={<ProxyLogoFullSize />}
            headerSecondaryInformation={[
              {
                title: 'Namespace',
                value: proxy.proxy!.metadata!.namespace
              },
              {
                title: 'BindPort',
                value: proxy.proxy!.listenersList[0].bindPort.toString()
              }
            ]}
            health={
              proxy.proxy!.status
                ? proxy.proxy!.status!.state
                : healthConstants.Pending.value
            }
            healthMessage={'Proxy Status'}>
            <InsideHeader>
              <div>Code Log (Read Only)</div>{' '}
              {!!proxy.raw && (
                <FileDownloadLink
                  fileName={proxy.raw.fileName}
                  fileContent={proxy.raw.content}
                />
              )}
            </InsideHeader>
            {!!proxy.raw && <YamlDisplayer content={proxy.raw.content} />}
          </SectionCard>
        );
      })}
    </React.Fragment>
  );
};
