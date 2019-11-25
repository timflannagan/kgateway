import { css } from '@emotion/core';
import { colors } from './colors';

export const globalStyles = css`
  html,
  body {
    width: 100vw;
    height: 100vh;
    overflow: auto !important;
  }

  body {
    font-family: 'Proxima Nova', 'Open Sans', 'Helvetica', 'Arial', 'sans-serif';
    margin: 0;
    padding: 0;
    min-height: 100vh;
    min-width: 100vw;
    background: ${colors.januaryGrey};

    .ant-modal-content {
      border-radius: 10px;
      box-shadow: hsla(0, 0%, 0%, 0.1) 0 4px 9px;

      .ant-modal-title {
        font-size: 24px;
        line-height: 26px;
      }
    }

    .ant-popover {
      color: white;
      .ant-popover-title {
        color: white;
        border: none;
      }
      .ant-popover-content {
        min-width: 125px;

        .ant-popover-message-title {
          color: white;
        }
        .ant-popover-inner {
          background: ${colors.novemberGrey};
          border-radius: 2px;

          .ant-popover-inner-content {
            color: white;
          }
        }
        .ant-popover-arrow {
          border-color: ${colors.novemberGrey};
        }
      }
    }

    .ant-select-dropdown {
      .ant-select-dropdown-menu-item {
        display: flex;
        align-items: center;

        svg {
          height: 20px;
          width: 20px;
          min-width: 20px;
          margin-right: 8px;
        }
      }
    }
  }
`;
